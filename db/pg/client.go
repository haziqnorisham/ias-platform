package pg

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	DB *sql.DB
}

type PpjTreeSensor struct {
	DeviceEUI              string  `db:"device_eui"`
	DeviceName             string  `db:"device_name"`
	Latitude               float64 `db:"latitude"`
	Longitude              float64 `db:"longitude"`
	DisplacementAlertAngle *int32  `db:"displacement_alert_angle"`
	Note                   *string `db:"note"`
	TreeID                 *string `db:"tree_id"`
	ID                     *int32  `db:"id"`
	GatewayID              *string `db:"gateway_id"`
	GatewayName            *string `db:"gateway_name"`
	Description            *string `db:"description"`
	GatewayEUI             *string `db:"gateway_eui"`
	NetworkServer          *string `db:"network_server"`
	GatewayModel           *string `db:"gateway_model"`
	LastSeen               *string `db:"last_seen"`
	IsActive               *bool   `db:"is_active"`
	CreatedAt              *string `db:"created_at"`
	UpdatedAt              *string `db:"updated_at"`
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	db, err := sql.Open("pgx", "postgres://"+os.Getenv("POSTGRES_USER")+":"+os.Getenv("POSTGRES_PASSWORD")+"@"+os.Getenv("POSTGRES_HOST")+":"+os.Getenv("POSTGRES_PORT")+"/"+os.Getenv("POSTGRES_DB")+"?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return &PostgresStorage{DB: db}
}
func (p *PostgresStorage) QueryData(query string) ([]PpjTreeSensor, error) {
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []PpjTreeSensor

	for rows.Next() {
		var sensor PpjTreeSensor
		err := rows.Scan(
			&sensor.DeviceEUI,
			&sensor.DeviceName,
			&sensor.Latitude,
			&sensor.Longitude,
			&sensor.DisplacementAlertAngle,
			&sensor.Note,
			&sensor.TreeID,
			&sensor.ID,
			&sensor.GatewayID,
			&sensor.GatewayName,
			&sensor.Description,
			&sensor.GatewayEUI,
			&sensor.NetworkServer,
			&sensor.GatewayModel,
			&sensor.LastSeen,
			&sensor.IsActive,
			&sensor.CreatedAt,
			&sensor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}
