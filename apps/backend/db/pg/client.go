package pg

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var sharedPool *sql.DB

func InitSharedPool() error {
	dsn := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT") + "/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return err
	}
	sharedPool = db
	slog.Info("PostgreSQL shared connection pool initialized", "process", "pg_client")
	return nil
}

func CloseSharedPool() {
	if sharedPool != nil {
		sharedPool.Close()
		slog.Info("PostgreSQL shared connection pool closed", "process", "pg_client")
	}
}

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	if db != nil {
		return &PostgresStorage{DB: db}
	}
	return &PostgresStorage{DB: sharedPool}
}
