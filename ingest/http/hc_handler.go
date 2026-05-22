package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	ias_pg "ias/automation/db/pg"
	"log/slog"
	"net/http"
	"os"
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

func saveDeviceProfileImage(profileID int, imageBase64 string) error {
	if imageBase64 == "" {
		return nil
	}
	if err := os.MkdirAll("public", 0755); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return err
	}
	return os.WriteFile(fmt.Sprintf("public/device_profile_%d.png", profileID), data, 0644)
}

func getDeviceProfileImageURL(profileID int) *string {
	path := fmt.Sprintf("public/device_profile_%d.png", profileID)
	if _, err := os.Stat(path); err == nil {
		url := fmt.Sprintf("/api/image/device_profile_%d.png", profileID)
		return &url
	}
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
		MessageID   *int64 `json:"message_id"`         // optional, filter by specific message_id
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
		"message_id", body.MessageID,
		"process", "hc_handler_main",
	)

	ias_db := ias_pg.NewPostgresStorage(nil)

	records, err := ias_db.QueryRawIngest(body.Limit, body.Offset, body.SortByMsgID, body.Status, body.MessageID)
	if err != nil {
		slog.Error("Failed to query raw ingest", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to query raw ingest"}`, http.StatusInternalServerError)
		return
	}

	// Get total count of matching records (ignoring pagination)
	total, err := ias_db.CountRawIngest(body.Status, body.MessageID)
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
	type profileWithImage struct {
		ProfileID              int     `json:"profile_id"`
		ProfileName            string  `json:"profile_name"`
		Manufacturer           string  `json:"manufacturer"`
		ModelNumber            string  `json:"model_number"`
		CommunicationsProtocol string  `json:"communications_protocol"`
		Decoder                string  `json:"decoder"`
		ImageURL               *string `json:"image_url"`
	}
	results := make([]profileWithImage, len(profiles))
	for i, p := range profiles {
		results[i] = profileWithImage{
			ProfileID:              p.ProfileID,
			ProfileName:            p.ProfileName,
			Manufacturer:           p.Manufacturer,
			ModelNumber:            p.ModelNumber,
			CommunicationsProtocol: p.CommunicationsProtocol,
			Decoder:                p.Decoder,
			ImageURL:               getDeviceProfileImageURL(p.ProfileID),
		}
	}
	jsonData, err := json.Marshal(results)
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
	resp := struct {
		ProfileID              int     `json:"profile_id"`
		ProfileName            string  `json:"profile_name"`
		Manufacturer           string  `json:"manufacturer"`
		ModelNumber            string  `json:"model_number"`
		CommunicationsProtocol string  `json:"communications_protocol"`
		Decoder                string  `json:"decoder"`
		ImageURL               *string `json:"image_url"`
	}{
		ProfileID:              profile.ProfileID,
		ProfileName:            profile.ProfileName,
		Manufacturer:           profile.Manufacturer,
		ModelNumber:            profile.ModelNumber,
		CommunicationsProtocol: profile.CommunicationsProtocol,
		Decoder:                profile.Decoder,
		ImageURL:               getDeviceProfileImageURL(profile.ProfileID),
	}
	jsonData, err := json.Marshal(resp)
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
	var req struct {
		ProfileName            string `json:"profile_name"`
		Manufacturer           string `json:"manufacturer"`
		ModelNumber            string `json:"model_number"`
		CommunicationsProtocol string `json:"communications_protocol"`
		Decoder                string `json:"decoder"`
		ImageBase64            string `json:"image_base64"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ProfileName == "" {
		http.Error(w, `{"error":"profile_name is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Creating device profile", "profile_name", req.ProfileName, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	profile := ias_pg.HcDeviceProfile{
		ProfileName:            req.ProfileName,
		Manufacturer:           req.Manufacturer,
		ModelNumber:            req.ModelNumber,
		CommunicationsProtocol: req.CommunicationsProtocol,
		Decoder:                req.Decoder,
	}
	id, err := ias_db.InsertDeviceProfile(profile)
	if err != nil {
		slog.Error("Failed to create device profile", "error", err)
		http.Error(w, `{"error":"failed to create device profile"}`, http.StatusInternalServerError)
		return
	}
	if err := saveDeviceProfileImage(id, req.ImageBase64); err != nil {
		slog.Error("Failed to save device profile image", "profile_id", id, "error", err)
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
	var req struct {
		ProfileID              int    `json:"profile_id"`
		ProfileName            string `json:"profile_name"`
		Manufacturer           string `json:"manufacturer"`
		ModelNumber            string `json:"model_number"`
		CommunicationsProtocol string `json:"communications_protocol"`
		Decoder                string `json:"decoder"`
		ImageBase64            string `json:"image_base64"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.ProfileID == 0 {
		http.Error(w, `{"error":"profile_id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Updating device profile", "profile_id", req.ProfileID, "process", "hc_handler_main")
	if req.ImageBase64 != "" {
		if err := saveDeviceProfileImage(req.ProfileID, req.ImageBase64); err != nil {
			slog.Error("Failed to save device profile image", "profile_id", req.ProfileID, "error", err)
		}
	}
	ias_db := ias_pg.NewPostgresStorage(nil)
	profile := ias_pg.HcDeviceProfile{
		ProfileID:              req.ProfileID,
		ProfileName:            req.ProfileName,
		Manufacturer:           req.Manufacturer,
		ModelNumber:            req.ModelNumber,
		CommunicationsProtocol: req.CommunicationsProtocol,
		Decoder:                req.Decoder,
	}
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
	os.Remove(fmt.Sprintf("public/device_profile_%d.png", req.ProfileID))
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
		Limit        int    `json:"limit"`
		Offset       int    `json:"offset"`
		SortByID     string `json:"sort_by_id"`
		Success      *bool  `json:"success"`
		DeviceID     string `json:"device_id"`
		RawMessageID *int64 `json:"raw_message_id"`
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
		Total   int                      `json:"total"`
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

// GetServerConfig handles POST /api/get_server_config
func GetServerConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	slog.Info("Retrieving server configuration", "process", "hc_handler_main")

	data, err := os.ReadFile(".env")
	if err != nil {
		slog.Error("Failed to read .env file", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to read configuration"}`, http.StatusInternalServerError)
		return
	}

	config := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"'`)

		upper := strings.ToUpper(key)
		if strings.Contains(upper, "PASSWORD") || strings.Contains(upper, "TOKEN") || strings.Contains(upper, "SECRET") {
			val = "***"
		}
		config[key] = val
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		slog.Error("Failed to marshal server config", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to marshal configuration"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// SaveDashboard handles POST /api/save_dashboard
// Body: { id?, name, layout_json } — creates if no id, updates if id present
func SaveDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var dashboard ias_pg.HcDashboard
	if err := json.NewDecoder(r.Body).Decode(&dashboard); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if dashboard.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}
	ias_db := ias_pg.NewPostgresStorage(nil)
	var result ias_pg.HcDashboard
	var err error
	if dashboard.Id == 0 {
		slog.Info("Creating dashboard", "name", dashboard.Name, "process", "hc_handler_main")
		result, err = ias_db.InsertDashboard(dashboard)
	} else {
		slog.Info("Updating dashboard", "id", dashboard.Id, "process", "hc_handler_main")
		result, err = ias_db.UpdateDashboard(dashboard)
	}
	if err != nil {
		slog.Error("Failed to save dashboard", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to save dashboard"}`, http.StatusInternalServerError)
		return
	}
	resp := struct {
		Id         int    `json:"id"`
		Name       string `json:"name"`
		LayoutJSON string `json:"layout_json"`
	}{
		Id:         result.Id,
		Name:       result.Name,
		LayoutJSON: result.LayoutJSON,
	}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetDashboards handles POST /api/get_dashboards
// Returns summary list: Array of { id, name, updated_at }
func GetDashboards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slog.Info("Retrieving all dashboard summaries", "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	dashboards, err := ias_db.GetAllDashboardSummaries()
	if err != nil {
		slog.Error("Failed to retrieve dashboards", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to retrieve dashboards"}`, http.StatusInternalServerError)
		return
	}
	if dashboards == nil {
		dashboards = []ias_pg.HcDashboard{}
	}
	summaries := make([]map[string]interface{}, len(dashboards))
	for i, d := range dashboards {
		summaries[i] = map[string]interface{}{
			"id":         d.Id,
			"name":       d.Name,
			"updated_at": d.UpdatedAt,
		}
	}
	jsonData, err := json.Marshal(summaries)
	if err != nil {
		slog.Error("Failed to marshal dashboards", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to marshal dashboards"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetDashboard handles POST /api/get_dashboard
// Body: { id }, returns: { id, name, layout_json }
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Id int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Id == 0 {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Retrieving dashboard by ID", "id", req.Id, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	dashboard, err := ias_db.GetDashboardByID(req.Id)
	if err != nil {
		slog.Error("Failed to retrieve dashboard", "id", req.Id, "error", err)
		http.Error(w, `{"error":"dashboard not found"}`, http.StatusNotFound)
		return
	}
	resp := struct {
		Id         int    `json:"id"`
		Name       string `json:"name"`
		LayoutJSON string `json:"layout_json"`
	}{
		Id:         dashboard.Id,
		Name:       dashboard.Name,
		LayoutJSON: dashboard.LayoutJSON,
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal dashboard", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to marshal dashboard"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteDashboardHandler handles POST /api/delete_dashboard
// Body: { id }, returns: { id, status: "deleted" }
func DeleteDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Id int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Id == 0 {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}
	slog.Info("Deleting dashboard", "id", req.Id, "process", "hc_handler_main")
	ias_db := ias_pg.NewPostgresStorage(nil)
	if err := ias_db.DeleteDashboard(req.Id); err != nil {
		slog.Error("Failed to delete dashboard", "error", err, "process", "hc_handler_main")
		http.Error(w, `{"error":"failed to delete dashboard"}`, http.StatusInternalServerError)
		return
	}
	resp := struct {
		Id     int    `json:"id"`
		Status string `json:"status"`
	}{
		Id:     req.Id,
		Status: "deleted",
	}
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
