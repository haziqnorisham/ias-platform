package http

import (
	"encoding/json"
	ias_pg "ias/automation/db/pg"
	"log/slog"
	"net/http"
	"strings"
)

// Non Http Utility functions related to HC schema and handlers
func SetupHcSchema() error {
	slog.Info("Creating HC Schema", "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
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

// GetRawIngest handles POST /api/get_raw_ingest
// Accepts JSON body: { "limit": 100, "offset": 0, "sort_by_message_id": "desc", "status": "" }
// Paginates through hc_raw_ingest records with optional status filter and sort order.
func GetRawIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse optional JSON body for filters
	type requestBody struct {
		Limit       int    `json:"limit"`
		Offset      int    `json:"offset"`
		SortByMsgID string `json:"sort_by_message_id"` // "asc" or "desc"
		Status      string `json:"status"`             // "processed", "unprocessed", "reprocess", or ""
	}
	body := requestBody{
		Limit:       100,
		Offset:      0,
		SortByMsgID: "desc",
		Status:      "",
	}

	// Try to decode body; ignore decode errors and use defaults
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Warn("Failed to decode request body, using defaults", "error", err, "process", "hc_handler_main")
		}
	}

	// Validate sort direction (case-insensitive)
	body.SortByMsgID = strings.ToLower(body.SortByMsgID)
	if body.SortByMsgID != "asc" && body.SortByMsgID != "desc" {
		body.SortByMsgID = "desc"
	}

	// Validate status filter (case-insensitive)
	body.Status = strings.ToLower(body.Status)
	if body.Status != "" && body.Status != ias_pg.IngestStatusProcessed && body.Status != ias_pg.IngestStatusUnprocessed && body.Status != ias_pg.IngestStatusReprocess {
		http.Error(w, `{"error":"invalid status, must be 'processed', 'unprocessed', or 'reprocess'"}`, http.StatusBadRequest)
		return
	}

	slog.Info("Querying raw ingest records",
		"limit", body.Limit,
		"offset", body.Offset,
		"sort_by_message_id", body.SortByMsgID,
		"status", body.Status,
		"process", "hc_handler_main",
	)

	ias_db := ias_pg.NewPostgresStorage(nil)

	records, err := ias_db.QueryRawIngest(body.Limit, body.Offset, body.SortByMsgID, body.Status)
	if err != nil {
		slog.Error("Failed to query raw ingest", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to query raw ingest"}`, http.StatusInternalServerError)
		return
	}

	// Get total count of matching records (ignoring pagination)
	total, err := ias_db.CountRawIngest(body.Status)
	if err != nil {
		slog.Error("Failed to count raw ingest records", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to count raw ingest records"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("Raw ingest records retrieved",
		"count", len(records),
		"total", total,
		"process", "hc_handler_main",
	)

	// Return empty array instead of null when no records
	if records == nil {
		records = []ias_pg.HcRawIngest{}
	}

	response := struct {
		Total   int                  `json:"total"`
		Records []ias_pg.HcRawIngest `json:"records"`
	}{
		Total:   total,
		Records: records,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetDeviceProfiles handles POST /api/get_device_profiles
func GetDeviceProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slog.Info("Retrieving all device profiles", "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	profiles, err := ias_db.GetAllDeviceProfiles()
	if err != nil {
		slog.Error("Failed to retrieve device profiles", "error", err)
		http.Error(w, `{"error":"failed to retrieve device profiles"}`, http.StatusInternalServerError)
		return
	}
	if profiles == nil {
		profiles = []ias_pg.HcDeviceProfile{}
	}
	jsonData, err := json.Marshal(profiles)
	if err != nil {
		slog.Error("Failed to marshal device profiles", "error", err)
		http.Error(w, `{"error":"failed to marshal device profiles"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetDeviceProfileByID handles POST /api/get_device_profile
func GetDeviceProfileByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProfileID int `json:"profile_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Retrieving device profile by ID", "profile_id", req.ProfileID, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	profile, err := ias_db.GetDeviceProfileByID(req.ProfileID)
	if err != nil {
		slog.Error("Failed to retrieve device profile", "profile_id", req.ProfileID, "error", err)
		http.Error(w, `{"error":"profile not found"}`, http.StatusNotFound)
		return
	}
	jsonData, err := json.Marshal(profile)
	if err != nil {
		slog.Error("Failed to marshal device profile", "error", err)
		http.Error(w, `{"error":"failed to marshal device profile"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetDeviceByIDHandler handles POST /api/get_device
func GetDeviceByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.DeviceID == "" {
		http.Error(w, `{"error":"device_id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Retrieving device by ID", "device_id", req.DeviceID, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	device, err := ias_db.GetDeviceByID(req.DeviceID)
	if err != nil {
		slog.Error("Failed to retrieve device", "device_id", req.DeviceID, "error", err)
		http.Error(w, `{"error":"device not found"}`, http.StatusNotFound)
		return
	}
	jsonData, err := json.Marshal(device)
	if err != nil {
		slog.Error("Failed to marshal device", "error", err)
		http.Error(w, `{"error":"failed to marshal device"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateDevice handles POST /api/create_device
func CreateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var device ias_pg.HcDevice
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if device.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Creating device", "device_name", device.Name, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.InsertDevice(device); err != nil {
		slog.Error("Failed to create device", "error", err)
		http.Error(w, `{"error":"failed to create device"}`, http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(map[string]string{"message": "device created successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// UpdateDeviceHandler handles POST /api/update_device
func UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var device ias_pg.HcDevice
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if device.Id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Updating device", "device_id", device.Id, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.UpdateDevice(device); err != nil {
		slog.Error("Failed to update device", "error", err)
		http.Error(w, `{"error":"failed to update device"}`, http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(map[string]string{"message": "device updated successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteDeviceHandler handles POST /api/delete_device
func DeleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		DeviceID string `json:"device_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.DeviceID == "" {
		http.Error(w, `{"error":"device_id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Deleting device", "device_id", req.DeviceID, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.DeleteDevice(req.DeviceID); err != nil {
		slog.Error("Failed to delete device", "error", err)
		http.Error(w, `{"error":"failed to delete device"}`, http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(map[string]string{"message": "device deleted successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateDeviceProfile handles POST /api/create_device_profile
func CreateDeviceProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var profile ias_pg.HcDeviceProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if profile.ProfileName == "" {
		http.Error(w, `{"error":"profile_name is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Creating device profile", "profile_name", profile.ProfileName, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	id, err := ias_db.InsertDeviceProfile(profile)
	if err != nil {
		slog.Error("Failed to create device profile", "error", err)
		http.Error(w, `{"error":"failed to create device profile"}`, http.StatusInternalServerError)
		return
	}
	resp := struct {
		ProfileID int    `json:"profile_id"`
		Message   string `json:"message"`
	}{ProfileID: id, Message: "device profile created successfully"}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// UpdateDeviceProfile handles POST /api/update_device_profile
func UpdateDeviceProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var profile ias_pg.HcDeviceProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if profile.ProfileID == 0 {
		http.Error(w, `{"error":"profile_id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Updating device profile", "profile_id", profile.ProfileID, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.UpdateDeviceProfile(profile); err != nil {
		slog.Error("Failed to update device profile", "error", err)
		http.Error(w, `{"error":"failed to update device profile"}`, http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(map[string]string{"message": "device profile updated successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteDeviceProfile handles POST /api/delete_device_profile
func DeleteDeviceProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProfileID int `json:"profile_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ProfileID == 0 {
		http.Error(w, `{"error":"profile_id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Deleting device profile", "profile_id", req.ProfileID, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.DeleteDeviceProfile(req.ProfileID); err != nil {
		slog.Error("Failed to delete device profile", "error", err)
		http.Error(w, `{"error":"failed to delete device profile"}`, http.StatusInternalServerError)
		return
	}
	jsonData, _ := json.Marshal(map[string]string{"message": "device profile deleted successfully"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetProcessedData handles POST /api/get_processed_data
func GetProcessedData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type requestBody struct {
		Limit         int    `json:"limit"`
		Offset        int    `json:"offset"`
		SortByID      string `json:"sort_by_id"`
		Success       *bool  `json:"success"`
		DeviceID      string `json:"device_id"`
		RawMessageID  *int64 `json:"raw_message_id"`
	}
	body := requestBody{
		Limit:    100,
		Offset:   0,
		SortByID: "desc",
	}

	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			slog.Warn("Failed to decode request body, using defaults", "error", err, "process", "hc_handler_main")
		}
	}

	body.SortByID = strings.ToLower(body.SortByID)
	if body.SortByID != "asc" && body.SortByID != "desc" {
		body.SortByID = "desc"
	}

	slog.Info("Querying processed data records",
		"limit", body.Limit,
		"offset", body.Offset,
		"sort_by_id", body.SortByID,
		"success", body.Success,
		"device_id", body.DeviceID,
		"raw_message_id", body.RawMessageID,
		"process", "hc_handler_main",
	)

	ias_db := ias_pg.NewPostgresStorage(nil)

	records, err := ias_db.QueryProcessedData(body.Limit, body.Offset, body.SortByID, body.Success, body.DeviceID, body.RawMessageID)
	if err != nil {
		slog.Error("Failed to query processed data", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to query processed data"}`, http.StatusInternalServerError)
		return
	}

	total, err := ias_db.CountProcessedData(body.Success, body.DeviceID, body.RawMessageID)
	if err != nil {
		slog.Error("Failed to count processed data records", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to count processed data records"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("Processed data records retrieved",
		"count", len(records),
		"total", total,
		"process", "hc_handler_main",
	)

	if records == nil {
		records = []ias_pg.HcProcessedData{}
	}

	response := struct {
		Total   int                     `json:"total"`
		Records []ias_pg.HcProcessedData `json:"records"`
	}{
		Total:   total,
		Records: records,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		slog.Error("Failed to marshal response", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to marshal response"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// ReprocessRawIngest handles POST /api/reprocess_raw_ingest
func ReprocessRawIngest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MessageIDs []int64 `json:"message_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if len(req.MessageIDs) == 0 {
		http.Error(w, `{"error":"message_ids is required"}`, http.StatusBadRequest)
		return
	}

	slog.Info("Reprocessing raw ingest records",
		"count", len(req.MessageIDs),
		"message_ids", req.MessageIDs,
		"process", "hc_handler_main",
	)

	ias_db := ias_pg.NewPostgresStorage(nil)
	affected, err := ias_db.BatchUpdateRawIngestStatus(req.MessageIDs, ias_pg.IngestStatusReprocess)
	if err != nil {
		slog.Error("Failed to reprocess raw ingest records", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to reprocess raw ingest records"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("Raw ingest records marked for reprocessing",
		"affected", affected,
		"process", "hc_handler_main",
	)

	response := struct {
		Affected int64  `json:"affected"`
		Message  string `json:"message"`
	}{
		Affected: affected,
		Message:  "records marked for reprocessing",
	}

	jsonData, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
