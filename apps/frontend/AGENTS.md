# AGENTS.md

## Architecture

- Vue 3 SPA with Vite, Composition API (`<script setup>`), no TypeScript
- No state management library — all state is local `ref()`/`reactive()` per component
- `@` import alias resolves to `./src` (configured in `vite.config.js`)
- Package manager is **npm**

## Required Backend

The following must be running for the app to function:

| Backend | Port | Proxy prefix |
|---|---|---|
| Main IAS backend | `:9090` | `/api` |

Vite dev server (`:5173`) proxies `/api` to the backend. In production, nginx handles the proxy.

## API Patterns

**Main backend (`src/api/posts.js`):** All endpoints use `POST` with JSON body (RPC-style). Example: `POST /api/get_all_devices`, `POST /api/save_dashboard`.

Both modules share the base `apiFetch()` from `src/api/index.js`.

## Key Commands

```sh
npm install          # install dependencies
npm run dev          # start dev server (hot reload on :5173)
npm run build        # production build (verify correctness)
```

No lint, test, or typecheck commands exist.

## Theming

| Token | Value | Usage |
|---|---|---|
| Primary color | `#48897b` (teal green) | Buttons, chart bars, accent borders, success icons, active tab underline |
| Font | Space Grotesk (weights 300–700) | Loaded from Google Fonts in `index.html` |
| Dark mode | `class="app-dark"` on `<html>` | Activates PrimeVue Aura dark theme preset |
| `--p-button-primary-background` | `#48897b` | Set in `main.css` on `body` |
| `--p-button-primary-border-color` | `#48897b` | Same |
| `--p-button-text-primary-color` | `#48897b` | Text button variant |

**Background palette** (dark greys):

| Element | Color |
|---|---|
| Sidenav | `#18181B` |
| Widget body | `#1a1a1a` |
| Card / widget header | `#202024` |
| Editor canvas | `#0e0e10` / `#0a0a0c` |
| Borders / dividers | `#212121` / `#2a2a2e` / `#3a3a3e` |
| Active chip (green) | `#2d5a4e` with `#4cff88` dot |

Status colors follow PrimeVue severity tokens (`severity="success"`, `"danger"`, `"warn"`, `"secondary"`). Text colors: `#e0e0e0` (primary), `#aaa`/`#888` (muted), `#555`/`#666` (disabled/placeholder).

## Component Conventions

- All `.vue` files use `<script setup>` with Composition API
- PrimeVue 4 (Aura dark theme) for all UI components — theme applied via `class="app-dark"` on `<html>`
- Icons from `primeicons` (`pi pi-*` classes)
- No scoped CSS scoping leaks — use `:deep()` in scoped styles to target PrimeVue internals
- Reusable patterns: `Panel.vue` wraps PrimeVue `Card` with title/subtitle; `ConfirmDialog.vue` for delete confirmations; toast via `useToast()` for user feedback
- **Prefer generalized, reusable components.** Avoid single-use components. When building a new component, make it configurable via props (title, subtitle, icon, variant, etc.) so it can serve multiple use cases. Follow existing reusable patterns (`Panel.vue`, `WidgetWrapper.vue`, `ConfirmDialog.vue`) as references.

## Dashboard Persistence

Dashboards are stored as a single JSON blob (`layout_json`) per dashboard in the backend. The `layout` ref in `DashboardEditor.vue` maps 1:1 to this column — `JSON.stringify(layout.value)` on save, `JSON.parse(data.layout_json)` on load. Widget types and their data shapes are defined in the in-memory `layout` array (see `DashboardEditor.vue` for `addMetric()`, `addChart()`, etc.).

## Unused Dependencies

`@vue-flow/core`, `@vue-flow/background`, `@vue-flow/controls`, `@vue-flow/minimap` are installed but the `src/components/nodes/` directory is empty. Do not remove these — reserved for future flow editor.
