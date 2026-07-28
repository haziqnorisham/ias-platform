# IAS Platform

Industrial IoT automation and data ingest platform. Ingests device telemetry via MQTT, decodes payloads through configurable extensions, stores time-series data in InfluxDB, and surfaces it through dashboards in a Vue 3 SPA. Ships with a pluggable extension system (local subprocess or remote HTTP) including ONVIF CCTV camera management with MediaMTX streaming, Telegram alerting, and surveillance event processing.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                         Browser                         │
└──────────────────────────┬──────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│  Nginx (:80)                                            │
│    /      → SPA static files (Vue 3)                    │
│    /api/* → Go backend (:9090)                          │
│             (backend reverse-proxies                    │
│              /api/extensions/{name}/* to extensions)    │
└──────┬──────────────────────────────┬───────────────────┘
       │                              │
 ┌─────┼───────┐              ┌───────┴────────┐
 ▼     ▼       ▼              ▼                ▼
┌──────────┐┌────────┐┌───────┐┌────────────┐┌───────────┐
│PostgreSQL││InfluxDB││ Redis ││ Extensions ││ MediaMTX  │
│ ledger   ││ TSDB   ││ cache ││ local /    ││ RTSP→HLS/ │
│ 18       ││ 2.7    ││ 7     ││ remote     ││ WebRTC    │
└──────────┘└────────┘└───────┘└─────┬──────┘└───────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                ▼                ▼
             ┌────────────┐  ┌────────────┐  ┌─────────────────┐
             │ cctv-onvif │  │  telegram  │  │mbpj-surveillance│
             │  (Node.js) │  │    (Go)    │  │    (Python)     │
             └────────────┘  └────────────┘  └─────────────────┘
```

### Data flow

1. **MQTT** — devices publish telemetry to the MQTT broker
2. **Worker** — Go worker pool polls raw ingest records from PostgreSQL, decodes payloads with JavaScript decoder scripts from device profiles (goja runtime), and writes decoded data to **InfluxDB**
3. **API** — the Go HTTP server serves dashboard metrics, device management, ingest logs, device profiles, and auth
4. **Extensions** — child processes (or remote HTTP services) discovered from `extension.json` manifests; the backend reverse-proxies `/api/extensions/{name}/*` to each extension's HTTP server
5. **Frontend** — Vue 3 SPA with drag-and-drop dashboards, line/bar chart widgets, extension widgets, and an extensions hub

## Project Structure

```
ias-platform/
├── apps/
│   ├── backend/              # Go API + ingest backend
│   │   ├── db/               # PostgreSQL & InfluxDB clients
│   │   ├── docs/             # Extension developer guide, TO-DO
│   │   ├── extension/        # Extension manager (local + remote)
│   │   ├── extensions/       # Bundled extensions
│   │   │   ├── cctv-onvif/       # ONVIF discovery + MediaMTX streams (Node.js)
│   │   │   ├── mbpj-surveillance/# Surveillance events + Telegram alerts (Python)
│   │   │   └── telegram/         # Telegram messaging (Go)
│   │   ├── ingest/           # HTTP handlers, MQTT client, auth
│   │   ├── worker/           # Background ingest processing pool
│   │   └── main.go
│   └── frontend/             # Vue 3 SPA (Vite)
│       └── src/
│           ├── api/          # API client modules
│           ├── components/
│           │   ├── widgets/      # MetricCard, WidgetWrapper, ExtensionWidget
│           │   ├── dashboards/   # DashboardEditor, chart widgets
│           │   ├── devices/      # Device cards & dialogs
│           │   └── ingest/       # Ingest data dialogs
│           ├── composables/  # useAuth
│           ├── views/        # 15 page components
│           └── router/
├── docs/
│   └── extension-widget-factory.md  # Dashboard widget plugin architecture
├── Design.md                 # Complete design system documentation
├── docker-compose.yaml
├── Dockerfile.backend
├── Dockerfile.frontend
├── nginx.conf
└── .env.docker
```

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Vue 3, Vite 7, PrimeVue 4, ECharts 6, vue3-grid-layout, Vue Router |
| Backend | Go 1.26 (stdlib `net/http`) |
| Operational DB | PostgreSQL 18 |
| Time-series DB | InfluxDB 2.7 |
| Cache | Redis 7 |
| Message broker | MQTT (Eclipse Paho client) |
| Auth | LDAP (optional, LLDAP bundled in Docker) |
| Media streaming | MediaMTX (RTSP → HLS/WebRTC) |
| Extensions | Node.js (Express), Python (Flask), Go |
| Deployment | Docker Compose, Nginx |

## Extensions

Extensions are standalone services that plug into the platform. Each has an `extension.json` manifest and exposes `/health` and `/execute` HTTP endpoints. The backend loads manifests from the extensions directory, spawns local extensions as child processes (port announced via `IAS_EXTENSION_PORT=`), and reverse-proxies UI/API calls to them.

### Bundled extensions

| Extension | Runtime | Purpose |
|---|---|---|
| **cctv-onvif** | Node.js | ONVIF camera discovery, credential management, MediaMTX stream registration (RTSP → HLS), saved cameras, dashboard widgets (saved-cameras, live-feed) |
| **mbpj-surveillance** | Python (Flask) | Surveillance event intake, Telegram alerts with pre-roll video clips |
| **telegram** | Go | Telegram messaging extension |

### Extension types

- **`local`** (default) — backend spawns the manifest `command` as a child process
- **`remote`** — backend proxies to an externally hosted `url` (no process management)

### Dashboard widgets

Extensions can contribute widgets to the dashboard grid via the **widget factory pattern** — the extension serves a `widget.js` that registers factories with the host app, which mounts them into Vue-managed containers. See [`docs/extension-widget-factory.md`](docs/extension-widget-factory.md) for the full architecture.

Extension UI pages are served as custom elements (`component.js`) and rendered in the Extensions Hub (`/extensions/{name}`).

## Quick Start (Local Development)

### Prerequisites

- Go 1.26+
- Node.js 20+
- PostgreSQL, InfluxDB 2.7, and Redis 7 (install locally or use Docker)

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

### 4. (Optional) CCTV streaming

For the cctv-onvif extension's live streams, run MediaMTX:

```bash
docker compose -f apps/backend/extensions/cctv-onvif/docker-compose.mediamtx.yaml up -d
```

Set `MEDIAMTX_API_URL` in the extension's environment (default API port `9997`).

## Docker Deployment

For a fully containerized deployment:

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

| Service | Container | Port | Purpose |
|---|---|---|---|
| **frontend** | ias-frontend | `8080` (public) | Vue 3 SPA via Nginx |
| **backend** | ias-backend | `9090` (internal) | Go API |
| **postgres** | ias-postgres | `5432` | Operational DB |
| **influxdb** | ias-influxdb | `8086` | Time-series DB |
| **redis** | ias-redis | internal | Cache |
| **lldap** | ias-lldap | `17170`, `3890` | LDAP auth server |
| **adminer** | ias_adminer | `7070` | DB admin UI |
| **dozzle** | — | `6060` | Container log viewer |

Extensions are mounted into the backend container at `/app/extensions` (via `./app_data/backend/extensions`).

Open **http://localhost:8080**.

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
| `nginx.conf` | SPA fallback + `/api` reverse proxy |
| `.env.docker` | Docker-specific environment variables |
| `apps/backend/extensions/cctv-onvif/docker-compose.mediamtx.yaml` | MediaMTX streaming service (optional) |

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
- `IAS_ENABLE_EXTENSION` — enable extension system
- `MQTT_ENABLED` / `MQTT_BROKER_URL` / `MQTT_TOPIC` — MQTT ingest
- `WORKER_ENABLED` / `WORKER_COUNT` / `WORKER_BATCH_SIZE` — background data processor
- `AUTH_ENABLED` — enable LDAP authentication
- `LDAP_URL` / `LDAP_BASE_DN` / `LDAP_BIND_DN` / `LDAP_USER_FILTER` / `LDAP_ADMIN_GROUP` — LDAP config

## Documentation

| Document | Description |
|---|---|
| [`Design.md`](Design.md) | Design system — colors, typography, components, page layouts |
| [`docs/extension-widget-factory.md`](docs/extension-widget-factory.md) | Dashboard widget plugin architecture |
| [`apps/backend/docs/EXTENSION_DEVELOPER_GUIDE.md`](apps/backend/docs/EXTENSION_DEVELOPER_GUIDE.md) | How to build an extension |
| [`apps/backend/README.md`](apps/backend/README.md) | Backend-specific details |
| [`apps/frontend/README.md`](apps/frontend/README.md) | Frontend-specific details |
| [`apps/frontend/AGENTS.md`](apps/frontend/AGENTS.md) | Frontend coding conventions |
