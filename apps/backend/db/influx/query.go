package influx

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type ProcessedPoint struct {
	DeviceID     string                 `json:"device_id"`
	ProfileID    int                    `json:"profile_id"`
	RawMessageID int64                  `json:"raw_message_id"`
	MeasuredAt   time.Time              `json:"measured_at"`
	Payload      map[string]interface{} `json:"payload"`
}

func QueryProcessedData(limit int, offset int, sortDesc bool, deviceID string, rawMessageID *int64) ([]ProcessedPoint, error) {
	if limit <= 0 {
		limit = 100
	}

	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: 1970-01-01T00:00:00Z)
			|> filter(fn: (r) => r._measurement == "processed_data")
	`, Bucket)

	if deviceID != "" {
		flux += fmt.Sprintf(`
			|> filter(fn: (r) => r.device_id == "%s")
		`, deviceID)
	}

	flux += `
		|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
	`

	if rawMessageID != nil {
		flux += fmt.Sprintf(`
			|> filter(fn: (r) => r.raw_message_id == %d)
		`, *rawMessageID)
	}

	flux += fmt.Sprintf(`
		|> sort(columns: ["_time"], desc: %s)
		|> limit(n: %d, offset: %d)
	`, strconv.FormatBool(sortDesc), limit, offset)

	return queryProcessedPoints(flux)
}

func QueryDeviceHistory(deviceID string, limit int, startTime time.Time) ([]ProcessedPoint, error) {
	if limit <= 0 {
		limit = 1000
	}

	rangeStart := "1970-01-01T00:00:00Z"
	if !startTime.IsZero() {
		rangeStart = startTime.Format(time.RFC3339)
	}

	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s)
			|> filter(fn: (r) => r._measurement == "processed_data" and r.device_id == "%s")
			|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
			|> sort(columns: ["_time"], desc: true)
			|> limit(n: %d)
			|> sort(columns: ["_time"], desc: false)
	`, Bucket, rangeStart, deviceID, limit)

	return queryProcessedPoints(flux)
}

func QueryLatestByDeviceIDs(deviceIDs []string) (map[string]ProcessedPoint, error) {
	if len(deviceIDs) == 0 {
		return map[string]ProcessedPoint{}, nil
	}

	deviceSet := make([]string, len(deviceIDs))
	for i, id := range deviceIDs {
		deviceSet[i] = fmt.Sprintf(`"%s"`, id)
	}

	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: 1970-01-01T00:00:00Z)
			|> filter(fn: (r) => r._measurement == "processed_data")
			|> filter(fn: (r) => contains(value: r.device_id, set: [%s]))
			|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
			|> group(columns: ["device_id"])
			|> sort(columns: ["_time"], desc: true)
			|> limit(n: 1)
			|> group()
	`, Bucket, strings.Join(deviceSet, ", "))

	return queryProcessedPointsMap(flux)
}

func queryProcessedPoints(flux string) ([]ProcessedPoint, error) {
	slog.Debug("Executing InfluxDB query", "flux", flux)

	result, err := queryAPI.Query(context.Background(), flux)
	if err != nil {
		return nil, fmt.Errorf("influx query failed: %w", err)
	}

	var points []ProcessedPoint
	for result.Next() {
		record := result.Record()
		pt := rowToProcessedPoint(record.Values())
		if pt != nil {
			points = append(points, *pt)
		}
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("influx query iteration error: %w", result.Err())
	}

	if points == nil {
		points = []ProcessedPoint{}
	}

	return points, nil
}

func queryProcessedPointsMap(flux string) (map[string]ProcessedPoint, error) {
	slog.Debug("Executing InfluxDB query", "flux", flux)

	result, err := queryAPI.Query(context.Background(), flux)
	if err != nil {
		return nil, fmt.Errorf("influx query failed: %w", err)
	}

	points := make(map[string]ProcessedPoint)
	for result.Next() {
		record := result.Record()
		pt := rowToProcessedPoint(record.Values())
		if pt != nil {
			points[pt.DeviceID] = *pt
		}
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("influx query iteration error: %w", result.Err())
	}

	return points, nil
}

func rowToProcessedPoint(values map[string]interface{}) *ProcessedPoint {
	deviceID, _ := values["device_id"].(string)
	if deviceID == "" {
		return nil
	}

	rawJSON, ok := values["raw_json"].(string)
	if !ok {
		return nil
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
		payload = make(map[string]interface{})
	}

	measuredAt, _ := values["_time"].(time.Time)

	rawMsgID := int64(0)
	switch v := values["raw_message_id"].(type) {
	case float64:
		rawMsgID = int64(v)
	case int64:
		rawMsgID = v
	case string:
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			rawMsgID = id
		}
	}

	profileID := 0
	switch v := values["profile_id"].(type) {
	case float64:
		profileID = int(v)
	case int64:
		profileID = int(v)
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			profileID = id
		}
	}

	return &ProcessedPoint{
		DeviceID:     deviceID,
		ProfileID:    profileID,
		RawMessageID: rawMsgID,
		MeasuredAt:   measuredAt,
		Payload:      payload,
	}
}
