# Data Flow — IAS Automation

## Overview

Data flows through four stages: **ingest** → **decode** → **store** → **serve**. Postgres holds the operational ledger (raw payloads, device config, processing status). InfluxDB holds the decoded time-series values that the frontend queries.

```
MQTT Broker ──► ingest/mqtt ──► hc_raw_ingest (Postgres)
                                     │
                              worker/worker.go
                            (poll + Goja JS decode)
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                                 ▼
        InfluxDB processed_data          hc_ingest_summary (Postgres)
        (decoded time-series)            (success/error ledger)
                    │
          ingest/http handlers
          (Flux queries via db/influx)
                    │
              Frontend UI
```

---

## Stage 1: Ingest

### 1a. MQTT message arrives
**File:** `ingest/mqtt/hc_handler.go:27`

The MQTT client (`ingest/mqtt/client.go`) connects to the broker using env vars (`MQTT_BROKER_URL`, `MQTT_CLIENT_ID`, `MQTT_TOPIC`). When a message arrives on a matching topic (e.g. `/sensor/+/device-id/up`), the `HcDbHandler` callback fires.

**What it does:**
1. Extracts `device_id` from the MQTT topic (using the `{device_id}` placeholder defined in `MQTT_TOPIC`).
2. Calls `db.InsertRawIngest(topic, payload, &deviceID, "mqtt", "unprocessed")` → writes one row to `hc_raw_ingest` with `status = 'unprocessed'`.

### 1b. PostgreSQL raw ingest table
**File:** `db/pg/hc_pg_collections.go:92`

```sql
hc_raw_ingest
├── message_id    BIGSERIAL PK
├── topic         VARCHAR(255)   -- MQTT topic the message arrived on
├── payload       TEXT           -- raw JSON as a string
├── device_id     VARCHAR(50)    -- extracted from topic
├── ingest_method VARCHAR(20)    -- 'mqtt' (default)
├── status        VARCHAR(20)    -- 'unprocessed' → 'processed' | 'error'
└── received_at   TIMESTAMPTZ
```

---

## Stage 2: Decode (Worker)

### 2a. Poll loop
**File:** `worker/worker.go:83`

The scheduler polls Postgres on a configurable interval (`WORKER_POLL_INTERVAL`, default 5s). It fetches up to `WORKER_BATCH_SIZE` rows where `status IN ('unprocessed', 'reprocess')` via `db.GetUnprocessedIngestBatch()`.

Each row is dispatched to a pool of `WORKER_COUNT` concurrent goroutines.

### 2b. processRecord
**File:** `worker/worker.go:138`

For each raw ingest record:

1. **Resolve device** — `db.GetDeviceByID(deviceID)`. Fails if device doesn't exist.
2. **Resolve profile** — `db.GetDeviceProfileByID(device.ProfileID)`. Fails if device has no profile assigned.
3. **Get decoder script** — `profile.Decoder` is a JavaScript string stored in `hc_device_profiles`.
4. **Run decoder** — `decodePayload(decoderScript, rawPayload)`:
   - Creates a Goja JavaScript VM.
   - Runs the decoder script string.
   - Calls the `decode(payload)` function with the raw payload string.
   - JSON-marshals the return value.
5. **Write to InfluxDB** (if decode succeeded) — `influx.WriteProcessedPoint()`.
6. **Upsert summary** — `db.UpsertIngestSummary()` → `hc_ingest_summary`.
7. **Update status** — `db.UpdateRawIngestStatus()` sets status to `'processed'` or `'error'`.

**Failure modes:**
| Failure | Status | Retry? |
|---|---|---|
| Null device_id | `error` | No — manual reprocess via UI |
| Device not found | `error` | No |
| No profile assigned | `error` | No |
| Profile not found | `error` | No |
| No decoder script | `error` | No |
| Decoder JS throws | `error` | No |
| InfluxDB write fails | `error` | No — manual reprocess via UI |

All `error` records can be reprocessed via the `/api/reprocess_raw_ingest` endpoint, which sets status back to `'reprocess'` — the poll loop then picks them up again.

---

## Stage 3: Store

### 3a. InfluxDB: processed_data measurement
**File:** `db/influx/points.go:55`

Every successful decode writes one point to the `processed_data` measurement:

```
Measurement: processed_data
┌─────────────────────────────────────────────────────────────┐
│ Tags (indexed, low-cardinality)                             │
├──────────────────┬──────────────────────────────────────────┤
│ device_id        │ "sensor-abc123"                          │
│ profile_id       │ "3" (stored as string in InfluxDB)      │
├──────────────────┼──────────────────────────────────────────┤
│ Fields (not indexed)                                        │
├──────────────────┬──────────────────────────────────────────┤
│ temperature      │ 25.5 (float64)                           │
│ humidity         │ 62.0 (float64)                           │
│ status           │ "ok" (string)                            │
│ meta.battery     │ 3.7 (float64, one-level dot flattening)  │
│ raw_json         │ '{"temperature":25.5,...}' (verbatim)    │
│ raw_message_id   │ 1503 (int64, link to hc_raw_ingest)     │
├──────────────────┼──────────────────────────────────────────┤
│ Timestamp        │ processing time (time.Now())             │
└──────────────────┴──────────────────────────────────────────┘
```

**Flattening rules** (`FlattenJSON` in `db/influx/points.go:12`):
- Top-level numbers → `float64` fields (avoids type-conflict write errors across points).
- Top-level booleans/strings → kept as-is.
- One level of nested objects → dot-notation fields (e.g. `{"meta":{"battery":3.7}}` → field `meta.battery=3.7`).
- Arrays, deeper nesting, nulls → skipped. Full structure is preserved in `raw_json`.
- `raw_json` stores the complete decoded JSON string verbatim — this is the **authoritative source** for reconstruction.

### 3b. PostgreSQL: hc_ingest_summary (deep dive)

**File:** `db/pg/hc_pg_collections.go:107` (schema), `:382` (upsert), `worker/worker.go:194` (call site)

#### Schema

```sql
CREATE TABLE IF NOT EXISTS hc_ingest_summary (
    id              SERIAL PRIMARY KEY,
    raw_ingest_id   BIGINT UNIQUE REFERENCES hc_raw_ingest(message_id),
    device_id       VARCHAR(50),
    profile_id      INT,
    success         BOOLEAN NOT NULL DEFAULT true,
    error_message   TEXT DEFAULT '',
    process_count   INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ingest_summary_device
    ON hc_ingest_summary (device_id);
```

| Column | Type | Purpose |
|---|---|---|
| `id` | SERIAL PK | Internal row identifier |
| `raw_ingest_id` | BIGINT UNIQUE | 1:1 link to `hc_raw_ingest.message_id`. The `UNIQUE` constraint ensures exactly one summary row per raw ingest record. |
| `device_id` | VARCHAR(50) | Which device this payload was processed for (extracted from MQTT topic during ingest) |
| `profile_id` | INT | Which device profile's decoder was used |
| `success` | BOOLEAN | Did decoding + InfluxDB write succeed? |
| `error_message` | TEXT | If `success=false`, the specific failure reason (e.g. "device has no profile assigned", "decoder error: ...", "influxdb write error: ...") |
| `process_count` | INT | How many times this raw ingest has been processed. Starts at 1 on first attempt, increments by 1 on each reprocess. |
| `created_at` | TIMESTAMPTZ | When the summary row was first created |
| `updated_at` | TIMESTAMPTZ | When the summary row was last updated (reprocessed) |

#### Role in the system

`hc_ingest_summary` is a **write-only processing ledger**. It records the outcome of every processing attempt — both successes and failures — in a single compact row per raw ingest. It is the authoritative source for answering "what happened to this specific raw message?"

**It is currently not read by any API endpoint.** After the InfluxDB migration, the old join-based query (`GetSuccessfulIngestByDeviceID` — which joined `hc_raw_ingest` → `hc_ingest_summary` → `hc_processed_data`) was removed. The index on `device_id` is forward-looking: it enables future queries like "show me all processing failures for device X" or "what error rate does device Y have" without a full table scan.

#### Write path (the only interaction)

Every time the worker processes a raw ingest record, regardless of success or failure, it upserts one row into `hc_ingest_summary`:

**Call site:** `worker/worker.go:194`
```go
if err := db.UpsertIngestSummary(raw.MessageID, deviceID, profileID, success, errorMsg); err != nil {
    log.Error("Failed to upsert ingest summary", "error", err)
}
```

**UPSERT logic** (`db/pg/hc_pg_collections.go:382`):

```sql
INSERT INTO hc_ingest_summary
    (raw_ingest_id, device_id, profile_id, success, error_message, process_count)
VALUES ($1, $2, $3, $4, $5, 1)
ON CONFLICT (raw_ingest_id)
DO UPDATE SET
    device_id     = $2,                -- overwrite with latest
    profile_id    = $3,                -- overwrite with latest
    success       = $4,                -- overwrite with latest outcome
    error_message = $5,                -- overwrite with latest
    process_count = hc_ingest_summary.process_count + 1,  -- increment
    updated_at    = NOW();
```

**Key behaviour:**
- **First attempt:** `INSERT` a new row with `process_count = 1`.
- **Subsequent reprocess** (e.g. after manual `reprocess` via the UI, or if `WORKER_PROCESS_ORDER` is `desc` and the worker picks up `error`/`reprocess` statuses): `ON CONFLICT` triggers the `UPDATE`, overwriting `success`/`error_message` with the latest outcome and incrementing `process_count`.
- **device_id and profile_id** are always overwritten — they are immutable for a given raw ingest (the raw record's `device_id` and the device's assigned `profile_id` don't change between attempts), but writing them on every upsert keeps the logic simple and idempotent.

#### Timeline of one record's lifecycle through hc_ingest_summary

```
Time  Event                                      hc_ingest_summary row
────  ───────────────────────────────────────     ───────────────────────────────────────────
 T0   MQTT message arrives
      InsertRawIngest → hc_raw_ingest
      (message_id=1500, status='unprocessed')
                                                  (no row yet)

 T1   Worker polls, picks up message_id=1500
      Decoder succeeds, InfluxDB write OK
      UpsertIngestSummary(1500, ...)              INSERT → success=true, process_count=1,
                                                  error_message='', updated_at=T1

 T2   Admin clicks "reprocess" in UI
      BatchUpdateRawIngestStatus → status='reprocess'

 T3   Worker polls, picks up message_id=1500 again
      Decoder succeeds, InfluxDB write OK
      UpsertIngestSummary(1500, ...)              UPDATE → success=true, process_count=2,
                                                  updated_at=T3

 T4   Admin clicks "reprocess" again
      This time decoder script was changed,
      decode throws an error
      UpsertIngestSummary(1500, ...)              UPDATE → success=false,
                                                  error_message='decoder error: ...',
                                                  process_count=3, updated_at=T4
```

#### What hc_ingest_summary does NOT replace

- **`hc_raw_ingest.status`** still drives the worker poll loop (`WHERE status IN ('unprocessed', 'reprocess')`). The summary is a ledger of *outcomes*; the raw record's status is the *dispatch signal*.
- **InfluxDB `processed_data`** is the source of truth for decoded values. The summary does not store payload data — only metadata about whether processing succeeded and how many times it was attempted.

#### Operational queries (manual, not automated in code)

```sql
-- Count processing failures by device in the last 24 hours
SELECT device_id, COUNT(*) AS failures
FROM hc_ingest_summary
WHERE success = false AND updated_at > NOW() - INTERVAL '24 hours'
GROUP BY device_id;

-- Find raw ingest records that were processed > 3 times (reprocessing loops?)
SELECT raw_ingest_id, process_count, success, error_message
FROM hc_ingest_summary
WHERE process_count > 3
ORDER BY process_count DESC;

-- List the most recent processing errors
SELECT s.raw_ingest_id, s.device_id, s.error_message, s.updated_at,
       r.topic, r.received_at
FROM hc_ingest_summary s
JOIN hc_raw_ingest r ON r.message_id = s.raw_ingest_id
WHERE s.success = false
ORDER BY s.updated_at DESC
LIMIT 50;
```

---

## Stage 4: Serve (HTTP API → Frontend)

All three processed-data endpoints run Flux queries against InfluxDB, not Postgres. The query logic lives in `db/influx/query.go`.

### 4a. Dashboard cards: `POST /api/get_dashboard_metric`
**File:** `ingest/http/hc_handler.go:960`

The performance-critical path. For each "card" metric, the frontend sends a `deviceID` and `column_name` (dot-separated path into the decoded JSON, e.g. `temperature` or `meta.battery`).

**Backend flow:**
1. `influx.QueryLatestByDeviceIDs(deviceIDs)` runs this Flux query:
   ```
   from(bucket) |> range(start: -100y)
     |> filter(r._measurement == "processed_data")
     |> filter(r => contains(value: r.device_id, set: ["id1","id2",...]))
     |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
     |> group(columns: ["device_id"])
     |> sort(columns: ["_time"], desc: true)
     |> limit(n: 1)
   ```
2. Returns `map[string]ProcessedPoint` — one per device, each with `Payload` (parsed `map[string]interface{}` from `raw_json`).
3. `extractValue(latest.Payload, columnName)` traverses the parsed map via dot-separated keys.
4. Response key is the exact same shape as before the migration:
   ```json
   { "metrics": [{ "type": "card", "deviceID": "...", "column_name": "...", "value": 25.5, "processed_at": "..." }] }
   ```

**For barchart/linechart metrics**, the handler calls `influx.QueryDeviceHistory(deviceID, 100)` instead, returning the most recent 100 points ordered chronologically.

### 4b. Browse: `POST /api/get_processed_data`
**File:** `ingest/http/hc_handler.go:574`

Admin/debug endpoint. Paginates through all processed points with optional device filter.

**Backend flow:**
1. `influx.QueryProcessedData(limit, offset, sortDesc, deviceID, rawMessageID)` runs Flux with optional tag filters, `pivot()`, `sort()`, and `limit()`.
2. Returns `[]ProcessedPoint`.
3. Response shape (breaking change — `total` removed):
   ```json
   { "records": [{ "device_id": "...", "profile_id": 3, "raw_message_id": 1503, "measured_at": "...", "payload": { "temperature": 25.5, ... } }] }
   ```
4. If `success=false` is passed in the request, returns `{"records":[]}` immediately (failed decodes have no InfluxDB points).

### 4c. Device history: `POST /api/get_device_successful_ingest`
**File:** `ingest/http/hc_handler.go:535`

Returns recent processed data for a single device (successful decodes only).

**Backend flow:**
1. `influx.QueryDeviceHistory(deviceID, limit)` runs Flux filtered by `device_id`, sorted desc, limited.
2. Returns `[]ProcessedPoint`.

**Breaking change:** Response shape changed from the old join-based `HcDeviceIngestSummary` (which included raw MQTT topic/payload/ingest_method) to the simpler `ProcessedPoint` (device_id, profile_id, raw_message_id, measured_at, payload). Raw MQTT data can still be fetched via `/api/get_raw_ingest` filtered by `message_id` if needed.

---

## Read Path Summary

| Frontend use case | Endpoint | InfluxDB function | Flux strategy |
|---|---|---|---|
| Dashboard card (latest value) | `/api/get_dashboard_metric` | `QueryLatestByDeviceIDs` | `group → sort desc → limit 1` per device |
| Dashboard chart (time series) | `/api/get_dashboard_metric` | `QueryDeviceHistory` | filter by device, sort asc, limit 100 |
| Admin browse page | `/api/get_processed_data` | `QueryProcessedData` | optional device filter, sort, limit/offset |
| Device history view | `/api/get_device_successful_ingest` | `QueryDeviceHistory` | filter by device, sort desc, limit N |

---

## Data Storage Ownership

| Data | Primary store | Backup/query store |
|---|---|---|
| Raw MQTT payloads | Postgres `hc_raw_ingest` | — |
| Device definitions | Postgres `hc_devices` | — |
| Device profiles / decoders | Postgres `hc_device_profiles` | — |
| Dashboards | Postgres `hc_dashboards` | — |
| Processing status ledger | Postgres `hc_ingest_summary` | — |
| Decoded time-series values | InfluxDB `processed_data` | `raw_json` fallback field |

---

## Configuration Reference

| Env var | Purpose |
|---|---|
| `INFLUXDB_URL` | InfluxDB server URL |
| `INFLUXDB_TOKEN` | Auth token (admin token for dev, scoped token for prod) |
| `INFLUXDB_ORG` | InfluxDB organisation name |
| `INFLUXDB_BUCKET` | InfluxDB bucket name |
| `MQTT_BROKER_URL` | MQTT broker address |
| `MQTT_TOPIC` | Topic pattern with `{device_id}` placeholder |
| `WORKER_ENABLED` | Enable/disable the processing worker |
| `WORKER_COUNT` | Number of concurrent decoder goroutines |
| `WORKER_POLL_INTERVAL` | How often to check for unprocessed records |
| `WORKER_BATCH_SIZE` | Max records per poll |
| `WORKER_PROCESS_ORDER` | `asc` (FIFO) or `desc` (LIFO) |
| `IAS_HC_BACKEND_ENABLE` | Master switch for HC backend (gates both Postgres schema init and InfluxDB init) |
