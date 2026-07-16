package influx

import (
	"context"
	"log/slog"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
	Bucket   string
)

func InitClient() error {
	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")
	org := os.Getenv("INFLUXDB_ORG")
	Bucket = os.Getenv("INFLUXDB_BUCKET")

	if url == "" || token == "" || org == "" || Bucket == "" {
		slog.Error("InfluxDB environment variables missing: INFLUXDB_URL, INFLUXDB_TOKEN, INFLUXDB_ORG, INFLUXDB_BUCKET must all be set")
		return nil
	}

	client = influxdb2.NewClient(url, token)

	health, err := client.Health(context.Background())
	if err != nil {
		slog.Error("Failed to connect to InfluxDB", "error", err)
		client.Close()
		return err
	}

	slog.Info("InfluxDB client initialized",
		"status", health.Status,
		"version", *health.Version,
		"org", org,
		"bucket", Bucket,
	)

	writeAPI = client.WriteAPIBlocking(org, Bucket)
	queryAPI = client.QueryAPI(org)

	return nil
}

func CloseClient() {
	if client != nil {
		client.Close()
		slog.Info("InfluxDB client closed")
	}
}
