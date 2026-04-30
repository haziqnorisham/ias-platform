package pg

type Hc_Device_Metadata struct {
	DeviceId    string  `db:"id"`
	DeviceEui   string  `db:"device_eui"`
	DeviceName  string  `db:"device_name"`
	Latitude    float64 `db:"latitude"`
	Longitude   float64 `db:"longitude"`
	CreatedDate string  `db:"created_date"`
}

func (p *PostgresStorage) CreateHcSchemaIfNotExists() error {
	query := `CREATE TABLE IF NOT EXISTS hc_device_metadata (
		id SERIAL PRIMARY KEY,
		device_eui VARCHAR(50) UNIQUE NOT NULL,
		device_name VARCHAR(100) NOT NULL,
		latitude DECIMAL(10, 8),
		longitude DECIMAL(11, 8),
		created_date DATE DEFAULT CURRENT_DATE
	);`

	_, err := p.DB.Exec(query)
	return err
}

func (p *PostgresStorage) GetAllDevices() ([]Hc_Device_Metadata, error) {
	query := `SELECT id, device_eui, device_name, latitude, longitude, created_date FROM hc_device_metadata;`
	rows, err := p.DB.Query(query)
	if err != nil {
		defer rows.Close()
		return nil, err
	}
	defer rows.Close()

	var devices []Hc_Device_Metadata

	for rows.Next() {
		var hc_device Hc_Device_Metadata
		err := rows.Scan(
			&hc_device.DeviceId,
			&hc_device.DeviceEui,
			&hc_device.DeviceName,
			&hc_device.Latitude,
			&hc_device.Longitude,
			&hc_device.CreatedDate,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, hc_device)
	}

	return devices, nil
}
