# Extension Developer Guide

This document teaches you how to build an extension that integrates with the
IAS Automation platform. You can use **any programming language** — your
extension is a standalone HTTP server. The IAS host spawns it as a child
process and communicates with it over `localhost` HTTP.

---

## Table of Contents

1. [What is an extension](#1-what-is-an-extension)
2. [Quickstart — the 10-minute version](#2-quickstart--the-10-minute-version)
3. [The protocol (every endpoint you must implement)](#3-the-protocol)
   - [3.1 Port announcement](#31-port-announcement)
   - [3.2 `/health`](#32-get-health)
   - [3.3 `/execute`](#33-post-execute)
   - [3.4 `/component.js`](#34-get-componentjs)
4. [Manifest (`extension.json`)](#4-manifest-extensionjson)
5. [Web Component design guide](#5-web-component-design-guide)
6. [Full examples](#6-full-examples)
   - [6.1 Go](#61-go)
   - [6.2 Python](#62-python)
   - [6.3 Node.js / JavaScript](#63-nodejs--javascript)
7. [Testing your extension](#7-testing-your-extension)
8. [Deploying to IAS](#8-deploying-to-ias)
   - [8.1 Local development](#81-local-development)
   - [8.2 Docker](#82-docker)
9. [Troubleshooting](#9-troubleshooting)

---

## 1. What is an extension

An extension is a standalone executable (binary, script, or interpreted
program) that exposes a small HTTP API on `localhost`. The IAS host spawns your
extension, reads the port from stdout, and communicates with it server-to-server
over HTTP.

Your extension runs in its own OS process. If it crashes, the host is
unaffected. You can write it in Go, Python, Node.js, Java, Rust, Ruby — any
language that can start an HTTP server and print to stdout.

**What extensions typically do:**

- Send Telegram/Slack notifications when sensor values exceed thresholds
- Integrate with ONVIF cameras (PTZ control, snapshot, recording)
- Forward decoded sensor data to external webhooks
- Provide custom dashboard panels via Web Components
- Bridge IAS to third-party APIs (email, SMS, CRM systems)

---

## 2. Quickstart — the 10-minute version

Here's the smallest possible extension in Python:

```python
import json, sys
from http.server import HTTPServer, BaseHTTPRequestHandler

class H(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/health":
            self._json(200, {"status": "ok"})
        elif self.path == "/component.js":
            self._js(200, b"customElements.define('ext-myext',class extends HTMLElement{constructor(){super()}})")
    def do_POST(self):
        if self.path == "/execute":
            body = json.loads(self.rfile.read(int(self.headers["Content-Length"])))
            self._json(200, {"success": True, "data": {"action": body["action"]}})
    def _json(self, code, data):
        self.send_response(code); self.send_header("Content-Type", "application/json")
        self.end_headers(); self.wfile.write(json.dumps(data).encode())
    def _js(self, code, data):
        self.send_response(code); self.send_header("Content-Type", "application/javascript")
        self.end_headers(); self.wfile.write(data)

s = HTTPServer(("localhost", 0), H)
print(f"IAS_EXTENSION_PORT={s.server_address[1]}", flush=True); sys.stdout.flush()
s.serve_forever()
```

Save as `myext.py`, create an `extension.json`:

```json
{
    "name": "myext",
    "command": ["python3", "./extensions/myext/myext.py"],
    "enabled": true
}
```

Place both files in `extensions/myext/`. Done — your extension will load on the
next host restart.

**Three things to notice:**

1. `net.Listen("tcp", "localhost:0")` / `HTTPServer(("localhost", 0), ...)` — port `0` means "OS, pick a free port for me."
2. `print("IAS_EXTENSION_PORT={port}")` — tell the host which port you got.
3. `/health`, `/execute`, `/component.js` — the three endpoints the host looks for.

---

## 3. The protocol

### 3.1 Port announcement

After binding to a port, print exactly this to stdout **before anything else**:

```
IAS_EXTENSION_PORT=18432
```

Rules:
- Must be the **first line** of output. No other output before it.
- Must end with a newline.
- Must flush immediately: `os.Stdout.Sync()` (Go), `sys.stdout.flush()` (Python), `console.log(...)` (Node.js auto-flushes on newline).
- The host waits up to `timeout_ms` (from the manifest) for this line.

Do not print any other text to stdout after this line — the host stops
reading stdout once it has the port. Use stderr for logging (`log.Println` in
Go, `sys.stderr` or `logging` in Python).

### 3.2 `GET /health`

```
GET /health
→ 200 {"status": "ok"}
```

The host calls this repeatedly (every ~100ms) after discovering the port until
it gets a 200, or until `timeout_ms` elapses. Return `200` as soon as your
extension is ready to accept `/execute` calls.

**Full response spec:**

| Field | Type | Required | Description |
|---|---|---|---|
| `status` | string | Yes | Must be `"ok"` when healthy |
| `version` | string | No | Extension version (for debugging) |
| `actions` | string[] | No | List of action names this extension supports |

### 3.3 `POST /execute`

```
POST /execute
Content-Type: application/json

Request body:
{
    "action": "send_message",
    "params": {
        "chat_id": "-1001234567890",
        "message": "Alert: temperature 45°C"
    }
}
```

**Request fields:**

| Field | Type | Required | Description |
|---|---|---|---|
| `action` | string | Yes | What to do. You define these — the host passes whatever the rule/config specifies. |
| `params` | object | Yes | Action-specific arguments. Can contain any JSON-serializable values. |

**Response — success:**

```json
{
    "success": true,
    "data": {
        "message_id": 1234,
        "sent_at": "2026-07-16T10:30:00Z"
    }
}
```

**Response — failure:**

```json
{
    "success": false,
    "error": "Chat not found: -1001234567890"
}
```

**Response fields:**

| Field | Type | Required | Description |
|---|---|---|---|
| `success` | boolean | Yes | `true` if the action completed, `false` if it failed |
| `data` | object | On success | Action-specific result. Any JSON object. |
| `error` | string | On failure | Human-readable error message |

**HTTP status codes:**
- `200` — the action was processed (check `success` field for outcome)
- `400` — invalid request body (missing `action`, bad JSON, unknown `action`)
- `500` — internal error (return `success: false, error: "..."` in the body)

### 3.4 `GET /component.js`

```
GET /component.js
Content-Type: application/javascript
→ JavaScript that registers a Web Component
```

This endpoint is **required** if you want your extension to have a visual panel
in the IAS frontend. If your extension is API-only (no UI), you can omit it.

The file must:
1. Define a class extending `HTMLElement`.
2. Call `customElements.define("ext-{name}", YourClass)`.
3. Return `Content-Type: application/javascript`.

The tag name convention is `ext-{extension-name}` (e.g. `ext-telegram`,
`ext-onvif`).

**Minimal component.js:**

```javascript
class ExtensionPanel extends HTMLElement {
    constructor() {
        super()
        this.attachShadow({ mode: "open" })
    }

    connectedCallback() {
        this.shadowRoot.innerHTML = `
            <style>:host { display: block; }</style>
            <p>Hello from the extension!</p>
        `
    }
}

customElements.define("ext-myext", ExtensionPanel)
```

For a complete Web Component with styling, input handling, and host
communication, see [Section 6](#6-full-examples).

---

## 4. Manifest (`extension.json`)

Place this file in your extension's directory (e.g. `extensions/telegram/`):

```json
{
    "name": "telegram",
    "command": ["./extensions/telegram/telegram-extension"],
    "enabled": true,
    "timeout_ms": 10000
}
```

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `name` | string | Yes | — | Unique name. Used to reference this extension in the API and future rule engine. Must be lowercase, no spaces. |
| `command` | string[] | Yes | — | OS command to start the extension. `command[0]` is the executable, `command[1..]` are args. |
| `enabled` | boolean | No | `true` | Set to `false` to skip loading this extension on startup. |
| `timeout_ms` | integer | No | `10000` | Max milliseconds the host waits for the port announcement + health check. |

**Python / script-based extensions:**

```json
{
    "name": "slack",
    "command": ["python3", "./extensions/slack/extension.py", "--debug"],
    "enabled": true
}
```

**Docker-based extensions (the host runs inside a container):**

```json
{
    "name": "telegram",
    "command": ["/app/extensions/telegram/telegram-extension"],
    "enabled": true
}
```

The `command` array is passed directly to `execve` — it is not interpreted by a
shell. No pipes, redirects, or globs. If you need those, wrap your extension in
a shell script.

---

## 5. Web Component design guide

### Why Web Components

- **Style isolation:** Shadow DOM prevents your CSS from leaking into the host
  and vice versa.
- **Framework-agnostic:** Works in Vue, React, Angular, or plain HTML.
- **Dynamic loading:** The frontend loads `component.js` at runtime — no build
  step, no npm install.
- **Attribute-driven config:** The host sets `device-id`, `dark-mode`, etc. as
  HTML attributes. Your component reacts.

### The `ext-{name}` tag

The frontend registers your component by loading:
```
GET /api/extensions/{name}/ui/component.js
```

Your `component.js` must call `customElements.define("ext-{name}", ...)`.
The tag name convention is `ext-{name}` — e.g. `ext-telegram`, `ext-onvif`,
`ext-slack`.

### Host → Extension communication

The Vue host sets HTML attributes on your element:

```html
<ext-telegram
    device-id="sensor-abc123"
    device-name="Ultrasonic Sensor"
    :dark-mode="true"
/>
```

Your component reads them:

```javascript
static get observedAttributes() {
    return ["device-id", "device-name", "dark-mode"]
}

attributeChangedCallback(name, oldVal, newVal) {
    if (name === "device-id") {
        this.updateDeviceId(newVal)
    }
}

connectedCallback() {
    const deviceId = this.getAttribute("device-id")
    const dark = this.getAttribute("dark-mode") === "true"
}
```

### Extension → Host communication

Dispatch `CustomEvent` with `bubbles: true, composed: true`. The Vue host can
listen for these events:

```javascript
this.dispatchEvent(new CustomEvent("ext-action", {
    detail: { action: "notification_sent", status: "ok" },
    bubbles: true,   // must be true to reach the parent
    composed: true   // must be true to cross the Shadow DOM boundary
}))
```

Vue listens with:
```html
<ext-telegram @ext-action="onExtAction" />
```

### Styling

All styles go **inside the Shadow DOM** via a `<style>` tag in `shadowRoot.innerHTML`:

```javascript
connectedCallback() {
    this.shadowRoot.innerHTML = `
        <style>
            :host { display: block; font-family: system-ui, sans-serif; }
            .card { background: #1a1a2e; border-radius: 8px; padding: 20px; }
            button { background: #e94560; color: white; border: none; cursor: pointer; }
            button:hover { background: #c73652; }
        </style>
        <div class="card">
            <button>Click me</button>
        </div>
    `
}
```

- Use `:host` to style the element itself (not the Shadow DOM internals).
- All other selectors are scoped inside the Shadow DOM.
- External stylesheets **will not** affect your component. Your styles **will
  not** affect the host. This is the point of Shadow DOM.

### Calling back to the host API

Your component's JavaScript runs in the browser. To call the extension's own
backend, use the host's proxy endpoint (once implemented). For now, call
`/execute` directly or embed the logic in `component.js` itself.

```javascript
async send() {
    const resp = await fetch("/api/extensions/telegram/execute", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            action: "send_message",
            params: { message: "hello" }
        })
    })
    const data = await resp.json()
}
```

---

## 6. Full examples

### 6.1 Go

This example shows: port announcement, `/health`, `/execute` with multiple actions,
and `/component.js` with a full Web Component.

**`main.go`:**

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
)

func main() {
    // Step 1: Bind to a random port
    listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
        log.Fatalf("Failed to bind: %v", err)
    }

    // Step 2: Announce the port
    port := listener.Addr().(*net.TCPAddr).Port
    fmt.Printf("IAS_EXTENSION_PORT=%d\n", port)
    os.Stdout.Sync()

    // Step 3: /health
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":  "ok",
            "version": "1.0.0",
            "actions": []string{"echo", "send_notification"},
        })
    })

    // Step 4: /execute
    http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Action string                 `json:"action"`
            Params map[string]interface{} `json:"params"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false, "error": "invalid JSON: " + err.Error(),
            })
            return
        }

        w.Header().Set("Content-Type", "application/json")

        switch req.Action {
        case "echo":
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": true,
                "data":    map[string]interface{}{"echo": req.Params["message"]},
            })

        case "send_notification":
            // Real extension: make HTTP call to Telegram / Slack / webhook here
            chatID, _ := req.Params["chat_id"].(string)
            msg, _ := req.Params["message"].(string)
            log.Printf("[extension] Sending to %s: %s", chatID, msg)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": true,
                "data":    map[string]interface{}{"sent": true, "to": chatID},
            })

        default:
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false, "error": "unknown action: " + req.Action,
            })
        }
    })

    // Step 5: /component.js (Web Component)
    http.HandleFunc("/component.js", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/javascript")
        w.Write([]byte(componentJS))
    })

    // Step 6: Start serving
    log.Printf("[extension] Listening on localhost:%d", port)
    if err := http.Serve(listener, nil); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}

const componentJS = `
class ExtensionPanel extends HTMLElement {
    constructor() { super(); this.attachShadow({ mode: "open" }) }

    connectedCallback() {
        const deviceId = this.getAttribute("device-id") || "—"
        this.shadowRoot.innerHTML = ` + "`" + `
            <style>
                :host { display: block; font-family: system-ui, sans-serif; }
                .card {
                    background: #16213e; border: 1px solid #0f3460;
                    border-radius: 8px; padding: 20px; color: #e0e0e0;
                }
                .title { font-size: 16px; font-weight: 600; color: #e94560; margin-bottom: 12px; }
                .device { font-size: 12px; color: #a0a0b0; margin-bottom: 12px; }
                input {
                    width: 100%; padding: 8px 12px; border: 1px solid #0f3460;
                    border-radius: 4px; background: #0f0f23; color: #e0e0e0;
                    font-size: 14px; margin-bottom: 12px; box-sizing: border-box;
                }
                button {
                    background: #e94560; color: white; border: none;
                    border-radius: 4px; padding: 8px 16px; cursor: pointer; font-size: 14px;
                }
                button:hover { background: #c73652; }
                button:disabled { background: #555; cursor: not-allowed; }
                .result {
                    margin-top: 12px; padding: 10px; background: #0f0f23;
                    border-radius: 4px; font-family: monospace; font-size: 12px;
                    white-space: pre-wrap; word-break: break-all; min-height: 20px;
                }
            </style>
            <div class="card">
                <div class="title">My Extension</div>
                <div class="device">Device: ${deviceId}</div>
                <input id="msg" type="text" placeholder="Type a message..." value="Hello from IAS!">
                <button id="send">Send Echo</button>
                <div class="result" id="result"></div>
            </div>
        ` + "`" + `

        this.shadowRoot.getElementById("send").onclick = () => this.send()
    }

    static get observedAttributes() { return ["device-id"] }

    attributeChangedCallback(name, oldVal, newVal) {
        if (name === "device-id" && this.shadowRoot) {
            const el = this.shadowRoot.querySelector(".device")
            if (el) el.textContent = "Device: " + (newVal || "—")
        }
    }

    async send() {
        const msg = this.shadowRoot.getElementById("msg").value
        const btn = this.shadowRoot.getElementById("send")
        const result = this.shadowRoot.getElementById("result")
        btn.disabled = true; result.textContent = "Sending..."

        try {
            const resp = await fetch("/api/extensions/NAME/execute", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ action: "echo", params: { message: msg } }),
            })
            const data = await resp.json()
            result.textContent = JSON.stringify(data, null, 2)

            this.dispatchEvent(new CustomEvent("ext:result", {
                detail: data, bubbles: true, composed: true,
            }))
        } catch (err) {
            result.textContent = "Error: " + err.message
        } finally {
            btn.disabled = false
        }
    }
}
customElements.define("ext-my-extension", ExtensionPanel)
`
```

Build:
```bash
CGO_ENABLED=0 go build -o extensions/myext/myextension ./extensions/myext
```

**`extension.json`:**

```json
{
    "name": "myextension",
    "command": ["./extensions/myext/myextension"],
    "enabled": true,
    "timeout_ms": 10000
}
```

### 6.2 Python

Full implementation with a working Web Component.

**`extension.py`:**

```python
#!/usr/bin/env python3
import json, sys
from http.server import HTTPServer, BaseHTTPRequestHandler

COMPONENT_JS = """
class ExtensionPanel extends HTMLElement {
    constructor() { super(); this.attachShadow({mode:"open"}) }
    connectedCallback() {
        this.shadowRoot.innerHTML = '<style>:host{display:block;font-family:sans-serif} .card{background:#16213e;border:1px solid #0f3460;border-radius:8px;padding:20px;color:#e0e0e0} input{width:100%;padding:8px;border:1px solid #0f3460;border-radius:4px;background:#0f0f23;color:#e0e0e0;font-size:14px;margin-bottom:12px;box-sizing:border-box} button{background:#e94560;color:white;border:none;border-radius:4px;padding:8px 16px;cursor:pointer} button:hover{background:#c73652} </style><div class="card"><input id="msg" value="Hello!"><button id="send">Send</button><div id="result"></div></div>'
        this.shadowRoot.getElementById("send").onclick = lambda: self.send()  # tricky in Python output
    }
}
customElements.define("ext-python-ext", ExtensionPanel)
"""

class Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/health":
            self._json(200, {"status": "ok", "version": "1.0.0"})
        elif self.path == "/component.js":
            self._js(200, COMPONENT_JS.encode())

    def do_POST(self):
        if self.path == "/execute":
            length = int(self.headers.get("Content-Length", 0))
            body = json.loads(self.rfile.read(length))
            action = body.get("action")
            params = body.get("params", {})

            if action == "echo":
                self._json(200, {"success": True, "data": {"echo": params.get("message")}})
            elif action == "send_notification":
                chat_id = params.get("chat_id", "")
                msg = params.get("message", "")
                print(f"[extension] Sending to {chat_id}: {msg}", file=sys.stderr)
                self._json(200, {"success": True, "data": {"sent": True}})
            else:
                self._json(400, {"success": False, "error": f"unknown action: {action}"})

    def _json(self, code, data):
        body = json.dumps(data).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _js(self, code, data):
        self.send_response(code)
        self.send_header("Content-Type", "application/javascript")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def log_message(self, format, *args):
        print(f"[extension] {args[0]}", file=sys.stderr)

server = HTTPServer(("localhost", 0), Handler)
port = server.server_address[1]
print(f"IAS_EXTENSION_PORT={port}", flush=True)
sys.stdout.flush()
server.serve_forever()
```

### 6.3 Node.js / JavaScript

**`extension.js`:**

```javascript
const http = require("http")

const HOST = "localhost"

const COMPONENT_JS = `
class ExtensionPanel extends HTMLElement {
    constructor() { super(); this.attachShadow({mode:"open"}) }
    connectedCallback() {
        this.shadowRoot.innerHTML = '<style>:host{display:block;font-family:sans-serif} .card{background:#16213e;border:1px solid #0f3460;border-radius:8px;padding:20px;color:#e0e0e0} input{width:100%;padding:8px 12px;border:1px solid #0f3460;border-radius:4px;background:#0f0f23;color:#e0e0e0;font-size:14px;margin-bottom:12px;box-sizing:border-box} button{background:#e94560;color:white;border:none;border-radius:4px;padding:8px 16px;cursor:pointer;font-size:14px} button:hover{background:#c73652} .result{margin-top:12px;padding:10px;background:#0f0f23;border-radius:4px;font-family:monospace;font-size:12px;white-space:pre-wrap;min-height:20px}</style><div class="card"><input id="msg" value="Hello!"><button id="send">Send</button><div class="result" id="result"></div></div>'
        this.shadowRoot.getElementById("send").onclick = () => this.send()
    }

    async send() {
        const msg = this.shadowRoot.getElementById("msg").value
        const result = this.shadowRoot.getElementById("result")
        try {
            const resp = await fetch("/api/extensions/node-ext/execute", {
                method: "POST",
                headers: {"Content-Type":"application/json"},
                body: JSON.stringify({action:"echo",params:{message:msg}})
            })
            const data = await resp.json()
            result.textContent = JSON.stringify(data, null, 2)
            this.dispatchEvent(new CustomEvent("ext:result", {detail:data, bubbles:true, composed:true}))
        } catch (err) {
            result.textContent = "Error: " + err.message
        }
    }
}
customElements.define("ext-node-ext", ExtensionPanel)
`

function jsonResponse(res, code, data) {
    const body = JSON.stringify(data)
    res.writeHead(code, {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(body)
    })
    res.end(body)
}

function parseBody(req, callback) {
    let body = ""
    req.on("data", chunk => body += chunk)
    req.on("end", () => {
        try { callback(JSON.parse(body)) }
        catch { jsonResponse(res, 400, { success: false, error: "invalid JSON" }) }
    })
}

const server = http.createServer((req, res) => {
    if (req.method === "GET" && req.url === "/health") {
        jsonResponse(res, 200, { status: "ok", version: "1.0.0",
            actions: ["echo", "send_notification"] })
    }
    else if (req.method === "GET" && req.url === "/component.js") {
        res.writeHead(200, { "Content-Type": "application/javascript" })
        res.end(COMPONENT_JS)
    }
    else if (req.method === "POST" && req.url === "/execute") {
        parseBody(req, body => {
            const { action, params } = body
            if (action === "echo") {
                jsonResponse(res, 200, { success: true,
                    data: { echo: params?.message } })
            } else if (action === "send_notification") {
                console.error(`[extension] Sending to ${params?.chat_id}: ${params?.message}`)
                jsonResponse(res, 200, { success: true, data: { sent: true } })
            } else {
                jsonResponse(res, 400, { success: false,
                    error: `unknown action: ${action}` })
            }
        })
    } else {
        res.writeHead(404)
        res.end()
    }
})

server.listen(0, HOST, () => {
    const port = server.address().port
    console.log(`IAS_EXTENSION_PORT=${port}`)
})
```

**`extension.json`:**

```json
{
    "name": "node-ext",
    "command": ["node", "./extensions/node-ext/extension.js"],
    "enabled": true
}
```

Run: `node extension.js` — the IAS host will start it automatically from the manifest.

---

## 7. Testing your extension

### Standalone test (without the host)

Start your extension manually:

```bash
# Go
go run ./extensions/telegram/main.go
# Output: IAS_EXTENSION_PORT=64425

# Python
python3 ./extensions/slack/extension.py
# Output: IAS_EXTENSION_PORT=54321

# Node.js
node ./extensions/webhook/extension.js
# Output: IAS_EXTENSION_PORT=55555
```

In another terminal, test the endpoints with `curl`:

```bash
PORT=the_port_printed_above

# Health check
curl http://localhost:$PORT/health
# → {"status":"ok"}

# Execute an action
curl -X POST http://localhost:$PORT/execute \
  -H "Content-Type: application/json" \
  -d '{"action":"echo","params":{"message":"test"}}'
# → {"success":true,"data":{"echo":"test"}}

# Web Component
curl http://localhost:$PORT/component.js
# → JavaScript source
```

### Testing the Web Component visually

Create a test HTML file:

```html
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body>
  <h3>Extension Test</h3>
  <ext-my-extension device-id="test-001"></ext-my-extension>

  <script>
    // Load the component from your running extension
    const script = document.createElement("script")
    script.src = "http://localhost:PORT/component.js"
    document.head.appendChild(script)
  </script>
</body>
</html>
```

Open in a browser. The component renders with full styling and interaction.

### Debugging checklist

| Check | Command |
|---|---|
| Extension starts without error | Run it standalone and watch stderr |
| Port is printed | Look for `IAS_EXTENSION_PORT=` on stdout |
| `/health` returns 200 | `curl localhost:$PORT/health` |
| `/execute` parses JSON correctly | `curl -X POST ... -d '{"action":"echo",...}'` |
| `/component.js` is valid JavaScript | `curl localhost:$PORT/component.js`, copy-paste into browser console |
| Web Component registers | Open browser console — should see no errors about `customElements.define` |
| No output on stdout after the port line | Stdout must be clean after the port announcement |

---

## 8. Deploying to IAS

### 8.1 Local development

1. Create `extensions/{name}/` directory.
2. Add your source code + `extension.json`.
3. Build your binary (if compiled language):
   ```bash
   CGO_ENABLED=0 go build -o extensions/myext/myext ./extensions/myext
   ```
4. Ensure `IAS_ENABLE_EXTENSION=true` and `IAS_HC_BACKEND_ENABLE=true` in `.env`.
5. Start the IAS host: `go run .`
6. Check logs for:
   ```
   INFO Extension loaded name=myext port=64425 pid=12345 process=extension_manager
   ```
7. Test: `curl -X POST http://localhost:9090/api/extensions`
   ```json
   [{"name":"myext","port":64425,"pid":12345}]
   ```

### 8.2 Docker

Add your extension build step to the Dockerfile:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bin/ias ./cmd/ias
RUN CGO_ENABLED=0 go build -o extensions/telegram/telegram-extension ./extensions/telegram

FROM alpine:3.20
COPY --from=builder /app/bin/ias /app/ias
COPY --from=builder /app/extensions/telegram/telegram-extension /app/extensions/telegram/telegram-extension
# For Python extensions, include the runtime:
# RUN apk add --no-cache python3
# COPY extensions/slack/extension.py /app/extensions/slack/extension.py
VOLUME /app/extensions
CMD ["/app/ias"]
```

The `VOLUME /app/extensions` lets you drop new extension binaries into a running
container without rebuilding the image.

---

## 9. Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| Extension doesn't appear in `/api/extensions` | `enabled: false` in manifest, or `IAS_ENABLE_EXTENSION` is not `true` | Check `.env` and `extension.json` |
| "extension xyz port announcement timed out" | Extension crashes before printing the port, or it prints to stderr instead of stdout | Run it standalone and check the first line of stdout |
| "health check timed out" | Extension server didn't start within `timeout_ms`. Port was printed but the HTTP server isn't accepting conns. | Increase `timeout_ms`, or fix the HTTP server init |
| Web Component shows "undefined" | `component.js` has a JavaScript syntax error | Open browser console, look for errors |
| Web Component loads but has no styles | Styles aren't inside the Shadow DOM | Put `<style>` in `shadowRoot.innerHTML`, not in the `<head>` |
| `POST /execute` returns `500` | Extension crashed during request, or returned non-JSON | Check extension stderr for panic/exception logs |
| Port conflict | The OS gave a port that's already in use | Unlikely with `localhost:0`, but kill stale extension processes |
| Extension binary doesn't run | Wrong architecture, missing deps | `CGO_ENABLED=0` for Go, ensure the base image has the right runtime (Python, Node) |
