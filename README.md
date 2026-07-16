# IAS Platform

Industrial IoT automation and data ingest platform. Ingests device telemetry via MQTT, decodes payloads through configurable extensions, stores time-series data in InfluxDB, and surfaces it through dashboards in a Vue 3 SPA.

## Architecture

```
┌─────────────────────────────────────────────────┐
│                    Browser                       │
└──────────────────────┬──────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────┐
│  Nginx (:80)                                    │
│    /         → SPA static files (Vue 3)         │
│    /api/*    → Go backend (:9090)               │
│    /media/*  → Media server (:8080)             │
└──────────────────────┬──────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         ▼             ▼             ▼
   ┌──────────┐ ┌──────────┐ ┌──────────┐
   │PostgreSQL│ │ InfluxDB │ │  Redis   │
   │ ledger   │ │ 2.7 TSDB │ │  cache   │
   └──────────┘ └──────────┘ └──────────┘
```

### Data flow

1. **MQTT** — devices publish telemetry to the MQTT broker
2. **Worker** — Go worker pool polls raw ingest records from PostgreSQL, runs them through JS decoders (extensions), and writes decoded payloads to **InfluxDB**
3. **API** — the Go HTTP server serves dashboard metrics, device management, ingest logs, and device profiles
4. **Frontend** — Vue 3 SPA with drag-and-drop dashboards, line/bar chart widgets, and extensions management

## Project Structure

```
ias-platform/
├── apps/
│   ├── backend/          # Go API + ingest backend
│   │   ├── db/           # PostgreSQL & InfluxDB clients
│   │   ├── extension/    # Extension subprocess manager
│   │   ├── ingest/       # HTTP handlers, MQTT client, worker
│   │   └── main.go
│   └── frontend/         # Vue 3 SPA (Vite)
│       └── src/
│           ├── api/      # API client modules
│           ├── components/widgets/
│           ├── components/dashboards/
│           ├── views/    # Page components
│           └── router/
├── docker-compose.yaml
├── Dockerfile.backend
├── Dockerfile.frontend
├── nginx.conf
└── .env.docker
```

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3, Vite 7, PrimeVue 4, ECharts 6, Vue Router |
| Backend | Go 1.26 (stdlib `net/http`) |
| Operational DB | PostgreSQL 16 |
| Time-series DB | InfluxDB 2.7 |
| Cache | Redis 7 |
| Message broker | MQTT (Eclipse Paho client) |
| Auth | LDAP (optional) |
| Deployment | Docker Compose, Nginx |

## Quick Start (Local Development)

### Prerequisites

- Go 1.26+
- Node.js 20+
- PostgreSQL 16, InfluxDB 2.7, and Redis 7 (install locally or use Docker)

### 1. Start infrastructure services

```bash
docker compose up -d postgres influxdb redis
```

### 2. Configure backend

```bash
cd apps/backend
cp example.env .env
# Edit .env: set POSTGRES_HOST=localhost, REDIS_HOST=localhost, etc.
make run-dev
```

### 3. Start frontend

```bash
cd apps/frontend
npm install
npm run dev
```

Open **http://localhost:5173**.

## Docker Deployment

For a fully containerized production deployment:

### 1. Configure

```bash
cp apps/backend/example.env apps/backend/.env   # local dev
# .env.docker is pre-configured for Docker service names
```

### 2. Build & Start

```bash
docker compose up -d
```

This starts:

| Service | Container | Port |
|---|---|---|
| **frontend** | ias-frontend | `80` (public) |
| **backend** | ias-backend | `9090` (internal) |
| **postgres** | ias-postgres | `5432` (internal) |
| **influxdb** | ias-influxdb | `8086` (internal) |
| **redis** | ias-redis | `6379` (internal) |

Open **http://localhost**.

### 3. Stop

```bash
docker compose down
```

### 4. Rebuild after code changes

```bash
docker compose up -d --build
```

### Docker files

| File | Purpose |
|---|---|
| `docker-compose.yaml` | Service orchestration |
| `Dockerfile.backend` | Go backend (multi-stage: build → alpine runtime) |
| `Dockerfile.frontend` | Vue 3 frontend (multi-stage: node build → nginx serve) |
| `nginx.conf` | SPA fallback + `/api` and `/media` reverse proxy |
| `.env.docker` | Docker-specific environment variables |

## Environment Configuration

| File | Use |
|---|---|
| `apps/backend/example.env` | Template for local development |
| `apps/backend/.env` | Local dev config (git-ignored) |
| `.env.docker` | Docker deployment config |

Key backend variables (see `example.env` for full list):

- `HTTP_SERVER_PORT` — API listen port
- `POSTGRES_HOST` / `PORT` / `USER` / `PASSWORD` / `DB`
- `REDIS_HOST` / `PORT`
- `INFLUXDB_URL` / `TOKEN` / `ORG` / `BUCKET`
- `IAS_HC_BACKEND_ENABLE` — enable API routes
- `IAS_ENABLE_EXTENSION` — enable extension subprocess system
- `MQTT_ENABLED` — enable MQTT ingest
- `WORKER_ENABLED` — enable background data processor
- `AUTH_ENABLED` — enable LDAP authentication
