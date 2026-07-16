# IAS Platform — Backend

Go-based IoT data ingest and automation backend serving the IAS Platform front-end.

## Stack

- **Go** 1.26 — HTTP server (stdlib `net/http`), worker pool, MQTT client
- **PostgreSQL** — operational ledger (devices, dashboards, raw ingest records)
- **InfluxDB** 2.7 — time-series storage for decoded/processed device data
- **Redis** — cache layer

## Quick Start (Local Development)

### Prerequisites

- Go 1.26+
- PostgreSQL, InfluxDB 2, and Redis running locally (or via [Docker Compose](../../docker-compose.yaml))
- Node.js 20+ (for the front-end; front-end is served separately in dev)

### 1. Configure environment

```bash
cp example.env .env
# Edit .env with your database credentials, ports, and feature flags
```

### 2. Build & Run

```bash
make dev        # build the binary
make run-dev    # or run without building (go run .)
```

The server listens on `HTTP_SERVER_PORT` (default from `example.env`: `8080`).

### Makefile targets

| Target | Description |
|---|---|
| `dev` | Build the Go binary |
| `run-dev` | Run with `go run .` |
| `clean` | Remove the built binary |
| `clean-build` | Alias for `clean` |

## Environment Variables

See `example.env` for a complete reference. Key variables:

| Variable | Default | Description |
|---|---|---|
| `HTTP_SERVER_PORT` | `8080` | HTTP listen port |
| `HTTP_SERVER_AUTOSTART` | `true` | Start the HTTP server |
| `POSTGRES_HOST` | `localhost` | PostgreSQL host |
| `REDIS_HOST` | `localhost` | Redis host |
| `INFLUXDB_URL` | — | InfluxDB server URL |
| `IAS_HC_BACKEND_ENABLE` | `true` | Enable API routes |
| `IAS_ENABLE_EXTENSION` | `false` | Enable extension system |
| `MQTT_ENABLED` | `true` | Enable MQTT sensor ingest |
| `WORKER_ENABLED` | `true` | Enable background worker (decodes raw ingest → InfluxDB) |
| `AUTH_ENABLED` | `false` | Enable LDAP authentication |

## Docker Deployment

Docker artifacts live at the [repository root](../../). See the root `README.md` or run:

```bash
docker compose up -d
```
