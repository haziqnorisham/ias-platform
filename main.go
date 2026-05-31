package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	ias_pg "ias/automation/db/pg"
	ingest_http "ias/automation/ingest/http"
	ingest_mqtt "ias/automation/ingest/mqtt"
	"ias/automation/worker"

	"github.com/joho/godotenv"
)

func main() {
	initLogger()
	loadEnv()
	initSharedPool()
	defer ias_pg.CloseSharedPool()
	setupHCBackendIfEnabled()
	startHTTPServerIfEnabled()
	startMQTTIfEnabled()
	sched := startWorkerIfEnabled()
	waitForShutdown(sched)
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	slog.Info("Application is starting")
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}
	slog.Info("Environment variables loaded", "process", "main")
}

func initSharedPool() {
	slog.Info("Initializing PostgreSQL shared connection pool", "process", "main")
	if err := ias_pg.InitSharedPool(); err != nil {
		slog.Error("Failed to initialize PostgreSQL shared pool", "error", err)
		os.Exit(1)
	}
	slog.Info("PostgreSQL shared connection pool initialized", "process", "main")
}

func setupHCBackendIfEnabled() {
	if os.Getenv("IAS_HC_BACKEND_ENABLE") != "true" {
		slog.Info("IAS HC Backend Server is disabled", "process", "main")
		return
	}
	slog.Info("IAS HC Backend Server is enabled", "process", "main")
	if err := ingest_http.SetupHcSchema(); err != nil {
		slog.Error("Failed to setup HC schema", "error", err)
		os.Exit(1)
	}
}

func startHTTPServerIfEnabled() {
	if os.Getenv("HTTP_SERVER_AUTOSTART") != "true" {
		return
	}
	slog.Info("Starting HTTP server", "process", "main")
	ingest_http.SetupRoutes()
	ingest_http.StartServer()
}

func startMQTTIfEnabled() {
	if os.Getenv("MQTT_ENABLED") != "true" {
		return
	}
	slog.Info("MQTT sensor monitoring is enabled, connecting to broker", "process", "main")

	var mqttHandlers []ingest_mqtt.MessageHandler
	if os.Getenv("IAS_HC_BACKEND_ENABLE") == "true" {
		slog.Info("HC raw ingest handler attached to MQTT subscription", "process", "main")
		mqttHandlers = append(mqttHandlers, ingest_mqtt.HcDbHandler())
	}

	if err := ingest_mqtt.ConnectAndSubscribe(mqttHandlers...); err != nil {
		slog.Error("Failed to start MQTT client", "error", err, "process", "main")
		os.Exit(1)
	}
	slog.Info("MQTT client connected and subscribed", "process", "main")
}

func startWorkerIfEnabled() *worker.Scheduler {
	if os.Getenv("WORKER_ENABLED") != "true" {
		return nil
	}
	slog.Info("Job scheduler is enabled, starting worker", "process", "main")
	sched := worker.NewScheduler()
	sched.Start()
	return sched
}

func waitForShutdown(sched *worker.Scheduler) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	slog.Info("Received signal, shutting down", "signal", sig.String())

	if sched != nil {
		slog.Info("Stopping job scheduler", "process", "main")
		sched.Stop()
	}

	if ingest_http.IsRunning {
		slog.Info("Shutting down HTTP server", "process", "main")
		ingest_http.StopServer()
	}

	if ingest_mqtt.IsRunning() {
		slog.Info("Shutting down MQTT client", "process", "main")
		ingest_mqtt.StopClient()
	}

	slog.Info("Application shut down gracefully", "process", "main")
}
