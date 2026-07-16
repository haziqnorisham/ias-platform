# IAS Platform — Frontend

Vue 3 SPA front-end for the IAS Platform, built with Vite.

## Stack

- **Vue 3** — Composition API (`<script setup>`)
- **Vite 7** — dev server + production bundler
- **PrimeVue 4** — UI component library (Aura dark theme)
- **ECharts 6** — charting (dashboards)
- **Vue Router** — SPA routing

## Quick Start

### Prerequisites

- Node.js 20.19+ or 22.12+
- Backend services running (see [backend README](../backend/README.md))

### Setup

```sh
npm install
```

### Development

```sh
npm run dev
```

The Vite dev server starts at **http://localhost:5173** and proxies:
- `/api` → `http://localhost:9090` (main IAS backend)
- `/media` → `http://localhost:8080` (media server, prefix stripped)

Configure proxy targets in `vite.config.js` if your backend runs on different ports.

### Production Build

```sh
npm run build
```

Output goes to `dist/`. In production, serve `dist/` with **nginx** (or any static file server) and proxy `/api` and `/media` to the backend. See the root `nginx.conf` for a pre-built nginx configuration.

## Project Structure

```
src/
├── api/           # API client modules (posts.js, extensions.js, auth.js)
├── components/    # Reusable Vue components
│   ├── charts/    # Legacy chart components
│   ├── dashboards/# Dashboard editor, viewer, chart widgets
│   └── widgets/   # WidgetWrapper, MetricCard
├── composables/   # Vue composables (useAuth)
├── router/        # Vue Router config
└── views/         # Page-level components
```

## Docker Deployment

The front-end is built and served by nginx in the Docker Compose stack at the [repository root](../../). See the root `README.md`.
