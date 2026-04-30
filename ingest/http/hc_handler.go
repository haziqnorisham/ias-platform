package http

import (
	"encoding/json"
	ias_pg "ias/automation/db/pg"
	"log/slog"
	"net/http"
)

// Non Http Utility functions related to HC schema and handlers
func SetupHcSchema() error {
	slog.Info("Creating HC Schema", "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	defer ias_db.DB.Close()
	err := ias_db.CreateHcSchemaIfNotExists()
	if err != nil {
		slog.Error("Failed to create HC schema", "error", err)
		return err
	}
	slog.Info("HC schema created successfully", "process", "hc_handler_main")
	return nil
}

// HTTP Handlers related to HC schema
func GetAllDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slog.Info("Retrieving all devices from HC schema", "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	defer ias_db.DB.Close()
	devices, err := ias_db.GetAllDevices()
	if err != nil {
		slog.Error("Failed to retrieve devices", "error", err)
		return
	}
	slog.Info("Devices retrieved successfully", "process", "hc_handler_main")
	jsonData, err := json.Marshal(devices)
	if err != nil {
		slog.Error("Failed to marshal devices to JSON", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
