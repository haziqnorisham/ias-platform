package influx

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func FlattenJSON(raw string) map[string]interface{} {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &root); err != nil {
		return nil
	}

	result := make(map[string]interface{})
	for key, val := range root {
		if val == nil {
			continue
		}
		m, ok := val.(map[string]interface{})
		if ok {
			for nestedKey, nestedVal := range m {
				if nestedVal == nil {
					continue
				}
				switch v := nestedVal.(type) {
				case float64:
					result[fmt.Sprintf("%s.%s", key, nestedKey)] = v
				case bool:
					result[fmt.Sprintf("%s.%s", key, nestedKey)] = v
				case string:
					result[fmt.Sprintf("%s.%s", key, nestedKey)] = v
				}
			}
			continue
		}
		switch v := val.(type) {
		case float64:
			result[key] = v
		case bool:
			result[key] = v
		case string:
			result[key] = v
		}
	}
	return result
}

func WriteProcessedPoint(deviceID string, profileID int, rawMessageID int64, rawJSON string, measuredAt time.Time) error {
	p := influxdb2.NewPointWithMeasurement("processed_data").
		AddTag("device_id", deviceID).
		AddTag("profile_id", fmt.Sprintf("%d", profileID)).
		AddField("raw_json", rawJSON).
		AddField("raw_message_id", rawMessageID)

	flattened := FlattenJSON(rawJSON)
	for key, val := range flattened {
		p.AddField(key, val)
	}

	if measuredAt.IsZero() {
		measuredAt = time.Now()
	}
	p.SetTime(measuredAt)

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		slog.Error("Failed to write point to InfluxDB",
			"device_id", deviceID,
			"raw_message_id", rawMessageID,
			"error", err,
		)
		return err
	}

	return nil
}
