package mqtt

import (
	"log/slog"
	"os"
	"strings"

	ias_pg "ias/automation/db/pg"
)

// extractDeviceID extracts a device identifier from an MQTT topic by aligning
// the incoming topic with the subscription topic template from MQTT_TOPIC env var.
// The template uses {device_id} as a placeholder (e.g. "sensors/{device_id}").
// Returns nil if the device ID position cannot be inferred or is empty.
func extractDeviceID(topic string) *string {
	subTopic := os.Getenv("MQTT_TOPIC")
	if subTopic == "" {
		return nil
	}

	templateParts := strings.Split(strings.Trim(subTopic, "/"), "/")
	var deviceIdx int = -1
	for i, part := range templateParts {
		if part == "{device_id}" {
			deviceIdx = i
			break
		}
	}

	if deviceIdx == -1 {
		return nil
	}

	parts := strings.Split(strings.Trim(topic, "/"), "/")
	if len(parts) <= deviceIdx {
		return nil
	}

	deviceID := parts[deviceIdx]
	if deviceID == "" || deviceID == "+" {
		return nil
	}

	return &deviceID
}

// HcDbHandler creates a MessageHandler that stores every incoming MQTT payload
// into the hc_raw_ingest PostgreSQL table with ingest_method="mqtt" and status="unprocessed".
func HcDbHandler() MessageHandler {
	return func(topic string, payload []byte) {
		deviceID := extractDeviceID(topic)

		db := ias_pg.NewPostgresStorage(nil)

		if err := db.InsertRawIngest(topic, payload, deviceID, "mqtt", ias_pg.IngestStatusUnprocessed); err != nil {
			slog.Error("Failed to store raw ingest",
				"topic", topic,
				"device_id", deviceID,
				"error", err,
				"process", "mqtt_hc_handler",
			)
			return
		}

		slog.Info("Raw ingest stored",
			"topic", topic,
			"device_id", deviceID,
			"payload_length", len(payload),
			"process", "mqtt_hc_handler",
		)
	}
}
