package pg

import (
	"strconv"
	"strings"
	"time"
)

type HcDevice struct {
	Id              int     `db:"id"`
	Name            string  `db:"name"`
	Description     *string `db:"description"`
	ProfileID       *int    `db:"profile_id"`
	Status          string  `db:"status"`
	LocationLabel   *string `db:"location_label"`
	CreatedDate     string  `db:"created_date"`
	LastUpdatedDate string  `db:"last_updated_date"`
}

// HcDeviceProfile represents a device profile stored in PostgreSQL.
type HcDeviceProfile struct {
	ProfileID              int    `db:"profile_id"`
	ProfileName            string `db:"profile_name"`
	Manufacturer           string `db:"manufacturer"`
	ModelNumber            string `db:"model_number"`
	CommunicationsProtocol string `db:"communications_protocol"`
	Decoder                string `db:"decoder"`
}

// Ingest status constants.
const (
	IngestStatusUnprocessed = "unprocessed"
	IngestStatusProcessed   = "processed"
	IngestStatusError       = "error"
)

// HcProcessedData represents a processed data record stored in PostgreSQL.
type HcProcessedData struct {
	ID               int64     `db:"id"`
	RawMessageID     int64     `db:"raw_message_id"`
	DeviceID         string    `db:"device_id"`
	ProfileID        int       `db:"profile_id"`
	ProcessedPayload string    `db:"processed_payload"`
	Success          bool      `db:"success"`
	ErrorMessage     string    `db:"error_message"`
	ProcessedAt      time.Time `db:"processed_at"`
}

// HcRawIngest represents a raw ingest message stored in PostgreSQL.
type HcRawIngest struct {
	MessageID    int64     `db:"message_id"`
	Topic        string    `db:"topic"`
	Payload      string    `db:"payload"`
	DeviceID     *string   `db:"device_id"`
	IngestMethod string    `db:"ingest_method"`
	Status       string    `db:"status"`
	ReceivedAt   time.Time `db:"received_at"`
}

func (p *PostgresStorage) CreateHcSchemaIfNotExists() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS hc_devices (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT DEFAULT '',
			profile_id INT REFERENCES hc_device_profiles(profile_id),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			location_label VARCHAR(255) DEFAULT '',
			created_date DATE DEFAULT CURRENT_DATE,
			last_updated_date DATE DEFAULT CURRENT_DATE
		);`,
		`CREATE TABLE IF NOT EXISTS hc_raw_ingest (
			message_id BIGSERIAL PRIMARY KEY,
			topic VARCHAR(255) NOT NULL,
			payload TEXT NOT NULL,
			device_id VARCHAR(50),
			ingest_method VARCHAR(20) NOT NULL DEFAULT 'mqtt',
			status VARCHAR(20) NOT NULL DEFAULT 'unprocessed',
			received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS hc_device_profiles (
			profile_id SERIAL PRIMARY KEY,
			profile_name VARCHAR(100) NOT NULL,
			manufacturer VARCHAR(100) NOT NULL DEFAULT '',
			model_number VARCHAR(100) NOT NULL DEFAULT '',
			communications_protocol VARCHAR(50) NOT NULL DEFAULT '',
			decoder TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS hc_processed_data (
			id BIGSERIAL PRIMARY KEY,
			raw_message_id BIGINT REFERENCES hc_raw_ingest(message_id),
			device_id VARCHAR(50),
			profile_id INT,
			processed_payload JSONB DEFAULT '{}',
			success BOOLEAN NOT NULL DEFAULT true,
			error_message TEXT DEFAULT '',
			processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
	}

	for _, q := range queries {
		if _, err := p.DB.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

// InsertRawIngest stores a raw ingest message into the hc_raw_ingest table.
func (p *PostgresStorage) InsertRawIngest(topic string, payload []byte, deviceID *string, ingestMethod string, status string) error {
	query := `INSERT INTO hc_raw_ingest (topic, payload, device_id, ingest_method, status) VALUES ($1, $2, $3, $4, $5);`
	_, err := p.DB.Exec(query, topic, string(payload), deviceID, ingestMethod, status)
	return err
}

// QueryRawIngest retrieves raw ingest records with pagination, optional status filter, and sort order.
// Parameters:
//   - limit: max records to return (default 100)
//   - offset: number of records to skip (for pagination)
//   - sortByMsgID: "asc" or "desc" ordering by message_id (default "desc")
//   - status: "" for all, or "processed"/"unprocessed" to filter
func (p *PostgresStorage) QueryRawIngest(limit int, offset int, sortByMsgID string, status string) ([]HcRawIngest, error) {
	if limit <= 0 {
		limit = 100
	}
	// Normalise sort direction to lowercase for case-insensitive matching
	sortByMsgID = strings.ToLower(sortByMsgID)
	if sortByMsgID != "asc" {
		sortByMsgID = "desc"
	}

	var query string
	var args []interface{}
	argIdx := 1

	if status != "" {
		query = `SELECT message_id, topic, payload, device_id, ingest_method, status, received_at FROM hc_raw_ingest WHERE status = $` + strconv.Itoa(argIdx) + ` ORDER BY message_id ` + sortByMsgID + ` LIMIT $` + strconv.Itoa(argIdx+1) + ` OFFSET $` + strconv.Itoa(argIdx+2) + `;`
		args = append(args, status, limit, offset)
	} else {
		query = `SELECT message_id, topic, payload, device_id, ingest_method, status, received_at FROM hc_raw_ingest ORDER BY message_id ` + sortByMsgID + ` LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1) + `;`
		args = append(args, limit, offset)
	}

	rows, err := p.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []HcRawIngest
	for rows.Next() {
		var r HcRawIngest
		if err := rows.Scan(&r.MessageID, &r.Topic, &r.Payload, &r.DeviceID, &r.IngestMethod, &r.Status, &r.ReceivedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

// CountRawIngest returns the total number of records matching the optional status filter.
func (p *PostgresStorage) CountRawIngest(status string) (int, error) {
	var query string
	var args []interface{}
	if status != "" {
		query = `SELECT COUNT(*) FROM hc_raw_ingest WHERE status = $1;`
		args = append(args, status)
	} else {
		query = `SELECT COUNT(*) FROM hc_raw_ingest;`
	}

	var total int
	if err := p.DB.QueryRow(query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (p *PostgresStorage) GetAllDevices() ([]HcDevice, error) {
	query := `SELECT id, name, description, profile_id, status, location_label, created_date, last_updated_date FROM hc_devices;`
	rows, err := p.DB.Query(query)
	if err != nil {
		defer rows.Close()
		return nil, err
	}
	defer rows.Close()

	var devices []HcDevice

	for rows.Next() {
		var hc_device HcDevice
		err := rows.Scan(
			&hc_device.Id,
			&hc_device.Name,
			&hc_device.Description,
			&hc_device.ProfileID,
			&hc_device.Status,
			&hc_device.LocationLabel,
			&hc_device.CreatedDate,
			&hc_device.LastUpdatedDate,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, hc_device)
	}

	return devices, nil
}

// GetDeviceByID retrieves a single device by its id (serial primary key).
func (p *PostgresStorage) GetDeviceByID(deviceID string) (*HcDevice, error) {
	query := `SELECT id, name, description, profile_id, status, location_label, created_date, last_updated_date FROM hc_devices WHERE id = $1;`
	var d HcDevice
	err := p.DB.QueryRow(query, deviceID).Scan(
		&d.Id,
		&d.Name,
		&d.Description,
		&d.ProfileID,
		&d.Status,
		&d.LocationLabel,
		&d.CreatedDate,
		&d.LastUpdatedDate,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// InsertDevice stores a new device into hc_devices.
func (p *PostgresStorage) InsertDevice(device HcDevice) error {
	query := `INSERT INTO hc_devices (id, name, description, profile_id, status, location_label)
		VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := p.DB.Exec(query, device.Id, device.Name, device.Description, device.ProfileID, device.Status, device.LocationLabel)
	return err
}

// UpdateDevice updates an existing device identified by id.
func (p *PostgresStorage) UpdateDevice(device HcDevice) error {
	query := `UPDATE hc_devices SET name=$1, description=$2, profile_id=$3, status=$4, location_label=$5, last_updated_date=CURRENT_DATE WHERE id=$6;`
	_, err := p.DB.Exec(query, device.Name, device.Description, device.ProfileID, device.Status, device.LocationLabel, device.Id)
	return err
}

// DeleteDevice removes a device by its id.
func (p *PostgresStorage) DeleteDevice(deviceID string) error {
	query := `DELETE FROM hc_devices WHERE id = $1;`
	_, err := p.DB.Exec(query, deviceID)
	return err
}

// InsertDeviceProfile stores a new device profile into hc_device_profiles.
func (p *PostgresStorage) InsertDeviceProfile(profile HcDeviceProfile) (int, error) {
	query := `INSERT INTO hc_device_profiles (profile_name, manufacturer, model_number, communications_protocol, decoder)
		VALUES ($1, $2, $3, $4, $5) RETURNING profile_id;`
	var id int
	err := p.DB.QueryRow(query, profile.ProfileName, profile.Manufacturer, profile.ModelNumber, profile.CommunicationsProtocol, profile.Decoder).Scan(&id)
	return id, err
}

// GetAllDeviceProfiles retrieves all device profiles.
func (p *PostgresStorage) GetAllDeviceProfiles() ([]HcDeviceProfile, error) {
	query := `SELECT profile_id, profile_name, manufacturer, model_number, communications_protocol, decoder FROM hc_device_profiles ORDER BY profile_id;`
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []HcDeviceProfile
	for rows.Next() {
		var pr HcDeviceProfile
		if err := rows.Scan(&pr.ProfileID, &pr.ProfileName, &pr.Manufacturer, &pr.ModelNumber, &pr.CommunicationsProtocol, &pr.Decoder); err != nil {
			return nil, err
		}
		profiles = append(profiles, pr)
	}
	return profiles, rows.Err()
}

// GetDeviceProfileByID retrieves a single device profile by its profile_id.
func (p *PostgresStorage) GetDeviceProfileByID(profileID int) (*HcDeviceProfile, error) {
	query := `SELECT profile_id, profile_name, manufacturer, model_number, communications_protocol, decoder FROM hc_device_profiles WHERE profile_id = $1;`
	var pr HcDeviceProfile
	err := p.DB.QueryRow(query, profileID).Scan(&pr.ProfileID, &pr.ProfileName, &pr.Manufacturer, &pr.ModelNumber, &pr.CommunicationsProtocol, &pr.Decoder)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// UpdateDeviceProfile updates an existing device profile.
func (p *PostgresStorage) UpdateDeviceProfile(profile HcDeviceProfile) error {
	query := `UPDATE hc_device_profiles SET profile_name=$1, manufacturer=$2, model_number=$3, communications_protocol=$4, decoder=$5 WHERE profile_id=$6;`
	_, err := p.DB.Exec(query, profile.ProfileName, profile.Manufacturer, profile.ModelNumber, profile.CommunicationsProtocol, profile.Decoder, profile.ProfileID)
	return err
}

// DeleteDeviceProfile removes a device profile by its profile_id.
func (p *PostgresStorage) DeleteDeviceProfile(profileID int) error {
	query := `DELETE FROM hc_device_profiles WHERE profile_id = $1;`
	_, err := p.DB.Exec(query, profileID)
	return err
}

// GetUnprocessedIngestBatch fetches up to limit unprocessed hc_raw_ingest records
// ordered by message_id ascending (oldest first).
func (p *PostgresStorage) GetUnprocessedIngestBatch(limit int) ([]HcRawIngest, error) {
	query := `SELECT message_id, topic, payload, device_id, ingest_method, status, received_at FROM hc_raw_ingest WHERE status = $1 ORDER BY message_id ASC LIMIT $2;`
	rows, err := p.DB.Query(query, IngestStatusUnprocessed, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []HcRawIngest
	for rows.Next() {
		var r HcRawIngest
		if err := rows.Scan(&r.MessageID, &r.Topic, &r.Payload, &r.DeviceID, &r.IngestMethod, &r.Status, &r.ReceivedAt); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, rows.Err()
}

// UpdateRawIngestStatus updates the status of a raw ingest record.
func (p *PostgresStorage) UpdateRawIngestStatus(messageID int64, status string) error {
	query := `UPDATE hc_raw_ingest SET status = $1 WHERE message_id = $2;`
	_, err := p.DB.Exec(query, status, messageID)
	return err
}

// InsertProcessedData stores a processed data record.
func (p *PostgresStorage) InsertProcessedData(data HcProcessedData) error {
	query := `INSERT INTO hc_processed_data (raw_message_id, device_id, profile_id, processed_payload, success, error_message) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := p.DB.Exec(query, data.RawMessageID, data.DeviceID, data.ProfileID, data.ProcessedPayload, data.Success, data.ErrorMessage)
	return err
}
