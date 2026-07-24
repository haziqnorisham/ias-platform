# Extension Widget Factory

## Overview

The extension widget factory pattern allows extensions to provide dashboard widgets that render inside the main application's dashboard grid **without custom elements or Shadow DOM**. This eliminates VDOM conflicts with Vue that would otherwise cause `insertBefore` null-parent errors.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  Extension (cctv-onvif/widget.mjs)                           │
│                                                              │
│  window.__ias_registerWidget('cctv-onvif:saved-cameras', {   │
│    label: 'Saved Cameras',                                   │
│    create(container, config) → { mount, update, unmount }    │
│  })                                                          │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Main App (main.js)                                          │
│                                                              │
│  window.__ias_registerWidget = (name, factory) → registry    │
│  window.__ias_getWidget = (name) → factory                   │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Vue Wrapper (components/widgets/ExtensionWidget.vue)        │
│                                                              │
│  onMounted → factory.create(hostEl, config) → instance.mount │
│  onUnmounted → instance.unmount()                            │
│  watch(config) → instance.update(config)                     │
└──────────────────────────┬───────────────────────────────────┘
                           │
                           ▼
┌──────────────────────────────────────────────────────────────┐
│  Dashboard (DashboardView.vue)                               │
│                                                              │
│  <ExtensionWidget                                            │
│    :widget-key="cctv-onvif:saved-cameras"                    │
│    :config="{ streamName: 'MsMediaProfile2' }"               │
│  />                                                          │
└──────────────────────────────────────────────────────────────┘
```

## Why Not Web Components?

Web components (via `customElements.define()`) use Shadow DOM, which Vue's VDOM reconciler cannot read or manage. This causes Vue to try removing/moving children it doesn't know about, resulting in:

```
TypeError: null is not an object (evaluating 'parent.insertBefore')
```

The widget factory pattern avoids this entirely:
- Vue owns a plain `<div>` element in its VDOM
- The factory renders into that div using vanilla DOM APIs
- Vue never touches the widget's content — it only manages the empty container div

## Widget Factory API

### `window.__ias_registerWidget(name, factory)`

Registers a widget factory under a unique name.

- **`name`** `string` — Unique widget key, conventionally `extensionName:widgetId` (e.g. `cctv-onvif:saved-cameras`)
- **`factory`** `object` — Widget definition:
  - **`label`** `string` — Human-readable label shown in the dashboard editor
  - **`defaultW`** `number` (optional) — Default grid width, defaults to `4`
  - **`defaultH`** `number` (optional) — Default grid height, defaults to `4`
  - **`create(container, config)`** — Factory function that returns a widget instance:

### `create(container, config)` → Widget Instance

Called by the Vue wrapper when mounting the widget.

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | `HTMLElement` | The empty `<div>` to render into. Owned by Vue's VDOM. |
| `config` | `object` | Widget configuration from the dashboard layout. |

Returns an object with lifecycle methods:

| Method | When called | Purpose |
|--------|------------|---------|
| `mount()` | On first render | Initial rendering into the container |
| `update(config)` | Config changes | Re-render with new config (optional) |
| `unmount()` | Widget removed | Clean up (clear intervals, remove listeners, etc.) |

## Example: Saved Cameras Widget

```js
// widget.mjs
const BASE = '/api/extensions/cctv-onvif'

window.__ias_registerWidget('cctv-onvif:saved-cameras', {
  label: 'Saved Cameras',
  defaultW: 4,
  defaultH: 4,

  create(container) {
    let interval = null

    async function fetch() {
      const resp = await fetch(`${BASE}/execute`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action: 'list-cameras' }),
      })
      const json = await resp.json()
      render(json.data?.cameras || [])
    }

    function render(cameras) {
      container.innerHTML = cameras.length
        ? cameras.map(c => `<div>${c.name || c.ip}</div>`).join('')
        : '<div>No cameras</div>'
    }

    fetch()
    interval = setInterval(fetch, 5000)

    return {
      unmount() {
        clearInterval(interval)
        container.innerHTML = ''
      },
    }
  },
})
```

## Declaring Widgets in `/health`

Extensions declare available widgets in their `/health` endpoint:

```json
{
  "widgets": [
    { "id": "saved-cameras", "label": "Saved Cameras", "w": 4, "h": 4 },
    { "id": "live-feed", "label": "Live Feed — MsMediaProfile2", "w": 6, "h": 4, "config": { "streamName": "MsMediaProfile2" } }
  ]
}
```

The `config` field is passed to `create(container, config)` at mount time and persisted in the dashboard's `layout_json`.

## Extension Layout JSON

Each widget stores its config in the dashboard layout:

```json
{
  "i": "5",
  "x": 0, "y": 6,
  "w": 4, "h": 4,
  "type": "extension",
  "extensionName": "cctv-onvif",
  "widgetId": "live-feed",
  "widgetLabel": "Live Feed — MsMediaProfile2",
  "widgetConfig": { "streamName": "MsMediaProfile2" }
}
```
