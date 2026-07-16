package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"

	ias_pg "ias/automation/db/pg"
	ias_influx "ias/automation/db/influx"

	"github.com/dop251/goja"
)

type Scheduler struct {
	cancel       context.CancelFunc
	workerCount  int
	batchSize    int
	pollInterval time.Duration
	processOrder string
	jobs         chan ias_pg.HcRawIngest
	done         chan struct{}
}

func NewScheduler() *Scheduler {
	workerCount, _ := strconv.Atoi(getEnvOrDefault("WORKER_COUNT", "2"))
	batchSize, _ := strconv.Atoi(getEnvOrDefault("WORKER_BATCH_SIZE", "20"))
	pollInterval, err := time.ParseDuration(getEnvOrDefault("WORKER_POLL_INTERVAL", "5s"))
	if err != nil {
		pollInterval = 5 * time.Second
	}
	processOrder := getEnvOrDefault("WORKER_PROCESS_ORDER", "asc")
	if processOrder != "asc" {
		processOrder = "desc"
	}

	return &Scheduler{
		workerCount:  workerCount,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		processOrder: processOrder,
		jobs:         make(chan ias_pg.HcRawIngest, batchSize),
		done:         make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	if s.workerCount <= 0 {
		slog.Warn("Worker count is 0, scheduler will not process records",
			"worker_count", s.workerCount,
			"process", "worker",
		)
	}

	for i := 0; i < s.workerCount; i++ {
		go s.worker(ctx, i)
	}

	go s.poll(ctx)

	slog.Info("Job scheduler started",
		"worker_count", s.workerCount,
		"batch_size", s.batchSize,
		"poll_interval", s.pollInterval.String(),
		"process", "worker",
	)
}

func (s *Scheduler) Stop() {
	slog.Info("Stopping job scheduler", "process", "worker")
	if s.cancel != nil {
		s.cancel()
	}
	<-s.done
	slog.Info("Job scheduler stopped", "process", "worker")
}

func (s *Scheduler) poll(ctx context.Context) {
	defer close(s.done)

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	dispatch := func() {
		db := ias_pg.NewPostgresStorage(nil)
		records, err := db.GetUnprocessedIngestBatch(s.batchSize, s.processOrder)
		if err != nil {
			slog.Error("Failed to fetch unprocessed ingest batch", "error", err, "process", "worker")
			return
		}
		if len(records) == 0 {
			return
		}
		slog.Debug("Fetched unprocessed records", "count", len(records), "process", "worker")

		for _, record := range records {
			select {
			case s.jobs <- record:
			case <-ctx.Done():
				return
			}
		}
	}

	dispatch()

	for {
		select {
		case <-ticker.C:
			dispatch()
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) worker(ctx context.Context, id int) {
	slog.Info("Worker started", "worker_id", id, "process", "worker")
	db := ias_pg.NewPostgresStorage(nil)

	for {
		select {
		case record := <-s.jobs:
			processRecord(db, record)
		case <-ctx.Done():
			slog.Info("Worker stopped", "worker_id", id, "process", "worker")
			return
		}
	}
}

func processRecord(db *ias_pg.PostgresStorage, raw ias_pg.HcRawIngest) {
	log := slog.With("message_id", raw.MessageID, "device_id", raw.DeviceID, "process", "worker")

	var deviceID string
	var profileID int
	var processedPayload string
	var success bool
	var errorMsg string

	if raw.DeviceID == nil {
		deviceID = ""
		success = false
		errorMsg = "device_id is null in raw ingest record"
	} else {
		deviceID = *raw.DeviceID
		device, err := db.GetDeviceByID(deviceID)
		if err != nil {
			success = false
			errorMsg = "device not found: " + err.Error()
		} else if device.ProfileID == nil {
			success = false
			errorMsg = "device has no profile assigned"
		} else {
			profileID = *device.ProfileID
			profile, err := db.GetDeviceProfileByID(profileID)
			if err != nil {
				success = false
				errorMsg = "profile not found: " + err.Error()
			} else if profile.Decoder == "" {
				success = false
				errorMsg = "profile has no decoder script"
			} else {
				result, err := decodePayload(profile.Decoder, raw.Payload)
				if err != nil {
					success = false
					errorMsg = "decoder error: " + err.Error()
				} else {
					success = true
					processedPayload = result
				}
			}
		}
	}

	if processedPayload == "" {
		processedPayload = "{}"
	}

	if success {
		if err := ias_influx.WriteProcessedPoint(deviceID, profileID, raw.MessageID, processedPayload, time.Time{}); err != nil {
			log.Error("Failed to write InfluxDB point, marking as error", "error", err)
			success = false
			errorMsg = "influxdb write error: " + err.Error()
		}
	}

	if err := db.UpsertIngestSummary(raw.MessageID, deviceID, profileID, success, errorMsg); err != nil {
		log.Error("Failed to upsert ingest summary", "error", err)
	}

	status := ias_pg.IngestStatusProcessed
	if !success {
		status = ias_pg.IngestStatusError
	}

	if err := db.UpdateRawIngestStatus(raw.MessageID, status); err != nil {
		log.Error("Failed to update raw ingest status", "error", err)
		return
	}

	if success {
		log.Info("Successfully processed ingest record",
			"profile_id", profileID,
		)
	} else {
		log.Warn("Failed to process ingest record",
			"error", errorMsg,
		)
	}
}

func decodePayload(decoderScript string, rawPayload string) (string, error) {
	vm := goja.New()

	if _, err := vm.RunString(decoderScript); err != nil {
		return "", err
	}

	decodeFn, ok := goja.AssertFunction(vm.Get("decode"))
	if !ok {
		return "", errors.New("decode function not found in decoder script")
	}

	result, err := decodeFn(goja.Undefined(), vm.ToValue(rawPayload))
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(result.Export())
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
