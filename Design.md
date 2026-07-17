# IAS Platform — Design System

## Product Identity

**IAS Platform** is an industrial IoT automation and data ingest platform. It ingests device telemetry via MQTT, decodes payloads through configurable extensions, stores time-series data in InfluxDB, and surfaces it through drag-and-drop dashboards in a Vue 3 SPA. The product shape is a **dashboard platform** — the data pipeline feeds a canvas of widgets that users arrange, configure, and monitor.

| Token | Value |
|---|---|
| App Title | IAS Health Center |
| Package | `ias_hc` |
| Version | ALPHA v0.01 |
| Tenant | DEV-ENV |
| Footer | (c) 2026 Camart Sdn. Bhd. |

**Logo**: `/bitmap.png` — 180px wide on login, full-width in sidenav.

## Architecture as Page Shape

The app is structured around a 15-route SPA with a persistent 16rem sidebar. The dominant feature is the **Home dashboard** — a 24-column drag-and-drop grid of metric cards, chart widgets, and item lists. This grid pattern repeats on the DashboardEditor (read-write) and DashboardView (read-only, live-polling).

```
┌──────────┬──────────────────────────────────────────┐
│ SideNav  │ Breadcrumb                               │
│ 16rem    │ ┌──────────────────────────────────────┐ │
│          │ │ <router-view />                      │ │
│ Home     │ │                                      │ │
│ Dashbrds │ │  GridLayout (24 cols × 50px rows)    │ │
│ Devices  │ │  ┌──────────┐ ┌──────────┐          │ │
│ Profiles │ │  │  Metric  │ │  Metric  │          │ │
│ Ingest   │ │  │  Card    │ │  Card    │          │ │
│ Extens.  │ │  └──────────┘ └──────────┘          │ │
│ Settings │ │  ┌──────────────────────────────────┐│ │
│ About    │ │  │  List Widget (extensions / dash) ││ │
│          │ │  └──────────────────────────────────┘│ │
│ [Profile]│ │                                      │ │
│ [logout] │ │                                      │ │
└──────────┴──────────────────────────────────────────┘
         Footer: (c) 2026 Camart Sdn. Bhd.
```

## Navigation

11 fixed items in the sidebar, each with a PrimeIcon. The home icon and chart-bar icon signal the app's primary purpose: monitoring device data through metrics and charts.

| # | Label | Icon | Route |
|---|---|---|---|
| 1 | Home | `pi pi-home` | `/` |
| 2 | IAS AI (Preview) | `pi pi-sparkles` | `/ai` |
| 3 | Dashboards | `pi pi-chart-bar` | `/dashboards` |
| 4 | Devices | `pi pi-microchip` | `/devices` |
| 5 | Device Profiles | `pi pi-wrench` | `/device-profiles` |
| 6 | Data Browser | `pi pi-database` | `/data-browser` |
| 7 | Ingest Logs | `pi pi-chevron-circle-down` | `/ingest-logs` |
| 8 | Settings | `pi pi-cog` | `/settings` |
| 9 | Extensions | `pi pi-bolt` | `/extensions` |
| 10 | Diagnostics | `pi pi-wave-pulse` | `/diagnostics` |
| 11 | About | `pi pi-info-circle` | `/about` |

Auth is LDAP-based. Login page appears standalone (no sidebar) with the logo, username/password fields, and a "Sign In" heading.

## Widget Type System

The dashboard editor supports 5 widget types, each stored as a JSON object in `layout_json`. Two additional `*list` types exist on the Home page only.

| Type | Properties | Visual |
|---|---|---|
| `card` | `cardTitle`, `cardValue`, `config.query` | `MetricCard` — large centered number |
| `barchart` | `chartTitle`, `config.query.y_axis` | `BarChartWidget` — ECharts bar |
| `linechart` | `lineChartTitle`, `config.query.y_axis` | `LineChartWidget` — ECharts line |
| `table` | `tableTitle` | `TableWidget` — PrimeVue DataTable |
| `text` | `textTitle`, `textContent` | `TextWidget` — plain text |
| `extensionslist` | (auto-fetched) | Scrollable list with Open buttons |
| `dashboardslist` | (auto-fetched) | Scrollable list with Open buttons |

**Chart data binding**: Widgets with `config.query.deviceID` poll `POST /api/get_dashboard_metric` every N seconds. The backend queries InfluxDB via Flux and returns `data_points: [{x, y, processed_at}]`. The X axis uses `type: 'time'` in ECharts (auto-formatted labels). Y axis is a dot-separated payload path (e.g. `object.temperature`). Cards receive the latest scalar value via `column_name`.

**Dashboard controls**: Above the grid, a controls panel (`#1a1a1a` bg, `#2a2a2e` border, 8px radius) holds the time range selector (5m through 1yr + custom date picker), a refresh interval dropdown (5s/10s/30s/60s/Off), and a countdown timer with an SVG progress ring (`#48897b` stroke on `#2a2a2e` track).

## Color System

### Primary Brand
| Token | Hex |
|---|---|
| Primary | `#48897b` |
| Primary light | `#5fa89a` |
| Primary lightest | `#7dc4b5` |

### Backgrounds (dark greys, darkest → lightest)
| Token | Hex | Usage |
|---|---|---|
| Canvas outer | `#0e0e10` | Login page, editor container |
| Canvas inner | `#0a0a0c` | Editor canvas, textarea |
| SideNav | `#18181B` | Left sidebar |
| Widget body | `#1a1a1a` | Widget content, card backgrounds |
| Toolbar | `#141416` | Dashboard editor header |
| Widget header | `#202024` | Widget/card headers, tooltips |
| Nav hover | `#3a3a3e` | SideNav item hover |
| Nav active | `rgba(255,255,255,0.08)` | Active nav state |

### Borders
| Token | Hex | Usage |
|---|---|---|
| Primary border | `#212121` | SideNav edges, dialogs, login card |
| Secondary border | `#2a2a2e` | Widget borders, cards, dividers, chart axes |
| Hover border | `#3a3a3e` | Card hover, image preview |

### Text
| Role | Hex |
|---|---|
| Primary text | `#e0e0e0` |
| Secondary / labels | `#a0a0a0` |
| Muted (widget titles) | `#aaa` |
| Subtle (subtitle) | `#888` / `#777` |
| Metric value | `#6d6d6d` |
| Disabled / placeholder | `#666` |
| Dim (empty states) | `#555` |
| Extra dim | `#444` |
| Light (on dark) | `#ccc` |

### Semantic
| State | Color |
|---|---|
| Success | `#4CAF50` icon |
| Active chip | `#2d5a4e` bg, `#4cff88` dot |
| Error | `#f44336` / `#e57373` |
| Link | `#64b5f6` → `#90caf9` (hover) |
| Warning badge | `rgba(240,173,78,0.85)` |

### Chart Colors
Chart grid lines: `rgba(255,255,255,0.06)` dashed. Axis labels: `#6d6d6d`. Tooltip: `#202024` bg, `#ffffff` text, `#2a2a2e` border. Bar gradient: `#5fa89a` → `#48897b`. Line: `#48897b` width 2, with area gradient `rgba(72,137,123,0.25)` → `rgba(72,137,123,0.02)`.

## Typography

**Font**: Space Grotesk, weights 300–700, loaded from Google Fonts. Monospace for code: SF Mono / Fira Code / Cascadia Code. Base HTML font size: **12px**.

| Token | px | Usage |
|---|---|---|
| `--font-size-2xs` | 7.8 | Meta labels, nav badge |
| `--font-size-xs` | 8.6 | Section labels, drag handle, profile role |
| `--font-size-sm` | 9.8 | Body text, form labels, code |
| `--font-size-md` | 10.8 | Default text, nav items |
| `--font-size-lg` | 12.6 | Section titles, card icons |
| `--font-size-xl` | 15 | Page titles, bold device names |
| `--font-size-2xl` | 24 | Metric values, welcome screen headings |
| `--font-size-3xl` | 36 | Empty state icons, modal icons |

## Spacing

| Element | Value |
|---|---|
| SideNav width | 16rem (256px) |
| Main content offset | left: 16rem |
| Page padding | 40px horizontal |
| Grid page padding | 8px |
| Grid columns × height | 24 cols × 50px rows |
| Grid margin | 12px |
| Widget header padding | 4px 8px |
| MetricCard padding | 12px |
| Card / widget gap | 10px–24px |
| Dialog padding | 1.25rem–1.5rem |
| Toolbar margin-bottom | 1rem |
| Section margin-bottom | 0.75rem |

## Border Radius

| Value | Usage |
|---|---|
| 0px | SideNav, nav items |
| 4px | Code blocks, drag handles, scrollbar |
| 6px | Image preview, mappings container |
| 8px | Grid items, controls panel, widgets |
| 10px | Editor container |
| 12px | Breadcrumb, cards, main content |
| 50% | Chat avatars, send button |

## Component Library

### Core Primitives
- **Panel** — `title`, `subtitle` — wraps PrimeVue Card
- **WidgetWrapper** — `title`, `icon`, `variant` (`default`|`flat`|`transparent`), `draggable`, `bodyPadding` — universal widget shell with CSS variable theming
- **MetricCard** — `value`, `title`, `hideTitle`, `loading`, `error` — centered single-value display, `#6d6d6d` value color
- **ConfirmDialog** — reusable delete confirmation, 420px, PrimeVue Dialog with danger severity

### Dashboard
- **DashboardEditor** — full editor with canvas, toolbar, widget management, persistence to `POST /api/save_dashboard`
- **EditorWidget** — individual editable widget with double-click title editing and DataSourceConfig dialog
- **DataSourceConfig** — Dynamic Data toggle, Device selector, Column Name field
- **BarChartWidget** / **LineChartWidget** — ECharts charts with `dataPoints`, `loading`, `error` props, reactive `setOption`, loading spinner and error icon overlays

### Charts (Legacy)
- **BarChart** — Chart.js bar with sample data
- **eSampleChart** — ECharts bar with sample data
- **BaseCard** — chart wrapper (legacy, unused in dashboards)

## State Patterns

All state is local `ref()`/`reactive()` per component — no Pinia, no Vuex. API calls use `apiFetch()` from `src/api/index.js` with relative URLs (`BASE_URL = ''`). Auth state lives in `useAuth` composable. Toast notifications at `top-center`. Dashboard metrics poll every N seconds via `setInterval`/`clearInterval` with `startPolling()`/`stopPolling()` helpers.

## Empty / Loading / Error States

Every data-fetching view implements three states:
1. **Loading** — `ProgressSpinner` centered (32px or 40px) with descriptive text
2. **Empty** — icon (`pi pi-box`, `pi pi-chart-bar`) + message in `#555`, flex-column centered
3. **Error** — `pi pi-exclamation-triangle` icon in `#e57373` + error message

Chart widgets additionally show overlay spinners on individual widget loading and error icons on widget error.

## Form Dialog Pattern

All dialogs share a consistent structure:
- Header: `border-bottom: 1px solid #212121`, `padding: 1.25rem 1.5rem`
- Content: `padding: 1.5rem`, grid layout with `gap: 1.25rem`
- Labels: Space Grotesk, `--font-size-sm`, `font-weight: 600`, `color: #a0a0a0`
- Required fields: red asterisk `color: #f44336`
- Footer: Cancel (text) left + Save (primary) right

## Transitions

| Element | Duration | Easing |
|---|---|---|
| Grid item hover ring | 0.15s | ease |
| List item hover bg | 0.15s | — |
| Card border hover | 0.15s | — |
| Link color | 0.15s | — |
| Drag handle color | 0.15s | ease |
| Metric card opacity (loading) | 0.2s | ease |
| Edit button opacity | 0.15s | ease |
| Payload cell bg | 0.2s | ease |

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend framework | Vue 3 (Composition API, `<script setup>`) |
| Build tool | Vite 7 |
| UI library | PrimeVue 4 (Aura dark theme) |
| Icons | PrimeIcons 7 |
| Charts | ECharts 6 (dashboards), Chart.js 4 (legacy) |
| Grid layout | vue3-grid-layout 1.0 |
| Routing | Vue Router 4 |
| Auth | LDAP (go-ldap) |
| Backend | Go 1.26 (`net/http`) |
| Operational DB | PostgreSQL 16 |
| Time-series DB | InfluxDB 2.7 |
| Cache | Redis 7 |
| Message broker | MQTT (Eclipse Paho) |
| Deployment | Docker Compose + Nginx |
