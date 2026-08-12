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

// SamplePoint is a single downsampled value (e.g. the mean of a field over a time window).
type SamplePoint struct {
	Time  time.Time
	Value interface{}
}

// fluxRangeStart renders the Flux range()-compatible start for a query.
func fluxRangeStart(startTime time.Time) string {
	if startTime.IsZero() {
		return "1970-01-01T00:00:00Z"
	}
	return startTime.Format(time.RFC3339)
}

// fluxDuration renders a Go duration as a Flux duration literal.
func fluxDuration(d time.Duration) string {
	if d <= 0 {
		d = time.Second
	}
	if d%(24*time.Hour) == 0 {
		return fmt.Sprintf("%dd", d/(24*time.Hour))
	}
	if d%time.Hour == 0 {
		return fmt.Sprintf("%dh", d/time.Hour)
	}
	if d%time.Minute == 0 {
		return fmt.Sprintf("%dm", d/time.Minute)
	}
	if d%time.Second == 0 {
		return fmt.Sprintf("%ds", d/time.Second)
	}
	return fmt.Sprintf("%dms", d/time.Millisecond)
}

// DownsampleWindow computes the aggregation window size needed to keep the total resulting
// points at or under targetPoints for a given time span. The window is rounded up to a "nice"
// unit so blocks stay stable across polls, with a 1s floor.
func DownsampleWindow(span time.Duration, targetPoints int) time.Duration {
	if targetPoints <= 0 {
		targetPoints = 1000
	}
	w := span / time.Duration(targetPoints)
	if w <= 0 {
		w = time.Second
	}
	if w < time.Second {
		w = time.Second
	}
	nice := []time.Duration{
		time.Second,
		5 * time.Second,
		10 * time.Second,
		30 * time.Second,
		time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		time.Hour,
		3 * time.Hour,
		6 * time.Hour,
		12 * time.Hour,
		24 * time.Hour,
		2 * 24 * time.Hour,
		7 * 24 * time.Hour,
	}
	for _, n := range nice {
		if w <= n {
			return n
		}
	}
	return w
}

// QueryDeviceHistoryWindowed returns one representative point per fixed time window covering the
// whole requested range (using the last point recorded in each window), instead of dropping data
// beyond a hard limit. Points are returned in ascending time order.
func QueryDeviceHistoryWindowed(deviceID string, startTime time.Time, every time.Duration) ([]ProcessedPoint, error) {
	if every <= 0 {
		every = time.Hour
	}

	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s)
			|> filter(fn: (r) => r._measurement == "processed_data" and r.device_id == "%s")
			|> sort(columns: ["_time"], desc: false)
			|> aggregateWindow(every: %s, fn: last, createEmpty: false)
			|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
			|> group()
			|> sort(columns: ["_time"], desc: false)
	`, Bucket, fluxRangeStart(startTime), deviceID, fluxDuration(every))

	return queryProcessedPoints(flux)
}

// QueryFieldMean aggregates a single numeric field to one mean value per fixed time window during
// the requested range. It errors when the field is not numeric (e.g. a string/boolean), so callers
// can fall back to QueryDeviceHistoryWindowed.
func QueryFieldMean(deviceID string, field string, startTime time.Time, every time.Duration) ([]SamplePoint, error) {
	if every <= 0 {
		every = time.Hour
	}

	flux := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s)
			|> filter(fn: (r) => r._measurement == "processed_data" and r.device_id == "%s" and r._field == "%s")
			|> aggregateWindow(every: %s, fn: mean, createEmpty: false)
	`, Bucket, fluxRangeStart(startTime), deviceID, field, fluxDuration(every))

	slog.Debug("Executing InfluxDB mean query", "flux", flux)

	result, err := queryAPI.Query(context.Background(), flux)
	if err != nil {
		return nil, fmt.Errorf("influx mean query failed: %w", err)
	}

	var points []SamplePoint
	for result.Next() {
		rec := result.Record()
		points = append(points, SamplePoint{
			Time:  rec.Time(),
			Value: rec.Value(),
		})
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("influx mean query iteration error: %w", result.Err())
	}

	if points == nil {
		points = []SamplePoint{}
	}

	return points, nil
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
