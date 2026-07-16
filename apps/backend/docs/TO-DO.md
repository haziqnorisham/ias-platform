# TO-DO: Dynamic Integration Module System (hashicorp/go-plugin)

## Overview

Each integration module (Telegram, ONVIF NVR, Slack, webhooks, etc.) runs as
its own OS process. The IAS host spawns and manages them via
`hashicorp/go-plugin`. Modules communicate over Unix sockets using Go's
`net/rpc`.

```
                        net/rpc over
┌──────────────────┐    localhost TCP    ┌──────────────────┐
│   IAS host       │◄──────────────────►│ telegram-module   │
│   (main process) │                    │  (subprocess)     │
│                  │                    └──────────────────┘
│  JS decoder      │
│  (untouched)     │    net/rpc over    ┌──────────────────┐
│  ─────────────── │◄──────────────────►│ onvif-module      │
│  MQTT → Goja     │                    │  (subprocess)     │
│  → InfluxDB      │                    └──────────────────┘
└──────────────────┘
```

**The JS decoder pipeline is untouched.** This system is for external
integrations only.

## Current state (implemented)

| Component | File | Status |
|---|---|---|
| Shared interface | `plugin/interface.go` | Done |
| Module manager | `plugin/manager.go` | Done |
| Telegram example module | `integrations/telegram/main.go` | Done (echo placeholder) |
| Module manifest | `integrations/telegram/module.json` | Done |
| Host wiring | `main.go` → `initModulesIfEnabled()` | Done |
| Module list endpoint | `POST /api/modules` | Done |

## How it works

### Interface (`plugin/interface.go`)

```go
type Integration interface {
    Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
    HealthCheck() error
}
```

Generic `map` in → `map` out. The module defines what keys it expects and what
it returns. No protobuf, no codegen — just `net/rpc`.

### Module manager (`plugin/manager.go`)

- `InitGlobal(modulesDir)` — scans `integrations/*/module.json`, spawns each
  binary, handshake, health check, registers in the registry.
- `GetGlobal()` — returns the singleton manager (used by HTTP handlers).
- `ShutdownGlobal()` — kills all subprocesses (called on SIGTERM).
- `Load` / `Unload` / `Call` / `List` on the manager struct.

### Adding a new module

1. Create `integrations/{name}/main.go` — implements `plugin.Integration`.
2. Create `integrations/{name}/module.json`:
   ```json
   { "name": "slack", "binary": "./integrations/slack/slack-module", "enabled": true }
   ```
3. Build it: `CGO_ENABLED=0 go build -o integrations/slack/slack-module ./integrations/slack`
4. Restart the host (or use the future `POST /api/modules/load` endpoint).

### HTTP API

| Endpoint | Method | Returns |
|---|---|---|
| `/api/modules` | POST | `["telegram", "slack"]` — list of loaded module names |

## Next steps (to implement)

1. **`POST /api/modules/load`** — load a module at runtime from a binary path.
2. **`POST /api/modules/unload`** — kill a module by name.
3. **`POST /api/modules/call`** — manually invoke a module (test/debug):
   ```json
   { "module": "telegram", "input": { "chat_id": "...", "message": "test" } }
   ```
   → `{ "ok": true, ... }`
4. **Rule engine (`plugin/trigger.go`)** — evaluate conditions on decoded
   payloads and fire modules automatically. Needs new Postgres tables
   `integration_rules` + `integration_actions`. Fire-and-forget goroutine hook
   in `worker/worker.go` after successful InfluxDB write.
5. **Real Telegram module** — replace the echo placeholder with actual Bot API
   calls (read `TELEGRAM_BOT_TOKEN` from env, POST to `api.telegram.org`).
6. **ONVIF module** — WS-Discovery + RTSP camera control. Longer-lived
   connections, may need gRPC streaming mode.

## Build

```bash
# Build the host
go build -o bin/ias ./cmd/ias

# Build a module
CGO_ENABLED=0 go build -o integrations/telegram/telegram-module ./integrations/telegram

# Build everything
go build ./...
```
