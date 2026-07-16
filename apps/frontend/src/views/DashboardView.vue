<template>
  <div class="page-container">
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />

    <div v-if="!layout.length" class="empty-state">
      <i class="pi pi-chart-bar empty-icon"></i>
      <p>This dashboard has no widgets.</p>
    </div>

    <template v-else>
    <div class="controls-panel">
      <div class="controls-panel__label">Dashboard Controls</div>
      <div class="time-range-controls">
        <Select
          v-model="selectedTimeRange"
          :options="timeRangeOptions"
          optionLabel="label"
          optionValue="value"
          size="small"
          class="time-range-select"
        />
        <DatePicker
          v-if="selectedTimeRange === 'custom'"
          v-model="customStartDate"
          showTime
          hourFormat="24"
          size="small"
          class="time-range-datepicker"
        />
        <span class="toolbar-separator"></span>
        <Select
          v-model="refreshInterval"
          :options="refreshIntervalOptions"
          optionLabel="label"
          optionValue="value"
          size="small"
          class="refresh-interval-select"
        />
        <span v-if="refreshInterval !== 'off'" class="countdown-group">
          <svg class="countdown-ring" width="20" height="20" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="8" fill="none" stroke="#2a2a2e" stroke-width="2" />
            <circle cx="10" cy="10" r="8" fill="none" stroke="#48897b" stroke-width="2"
              :stroke-dasharray="50.27"
              :stroke-dashoffset="50.27 * (1 - countdown / intervalsMap[refreshInterval])"
              stroke-linecap="round"
              transform="rotate(-90 10 10)" />
          </svg>
          <span class="countdown-text">{{ countdown }}s</span>
        </span>
      </div>
    </div>

    <GridLayout
      v-model:layout="layout"
      :col-num="24"
      :row-height="50"
      :is-draggable="false"
      :is-resizable="false"
      :margin="[12, 12]"
      :use-css-transforms="true"
      class="grid-container"
    >
      <GridItem
        v-for="item in layout"
        :key="item.i"
        :x="item.x"
        :y="item.y"
        :w="item.w"
        :h="item.h"
        :i="item.i"
        class="grid-item"
        :static="true"
      >
        <WidgetWrapper :title="item.type === 'card' ? item.cardTitle : item.type === 'barchart' ? item.chartTitle : item.type === 'linechart' ? item.lineChartTitle : item.type === 'table' ? item.tableTitle : item.type === 'text' ? item.textTitle : ''">
          <MetricCard
            v-if="item.type === 'card'"
            :title="item.cardTitle"
            :value="widgetData[item.i]?.value ?? item.cardValue"
            :loading="widgetData[item.i]?.loading ?? false"
            :error="widgetData[item.i]?.error ?? false"
          />
          <BarChartWidget v-else-if="item.type === 'barchart'" :title="item.chartTitle" :dataPoints="widgetData[item.i]?.dataPoints" :loading="widgetData[item.i]?.loading ?? false" :error="widgetData[item.i]?.error ?? false" />
          <LineChartWidget v-else-if="item.type === 'linechart'" :title="item.lineChartTitle" :dataPoints="widgetData[item.i]?.dataPoints" :loading="widgetData[item.i]?.loading ?? false" :error="widgetData[item.i]?.error ?? false" />
          <TableWidget v-else-if="item.type === 'table'" :title="item.tableTitle" />
          <TextWidget v-else-if="item.type === 'text'" :title="item.textTitle" :text="item.textContent" />
        </WidgetWrapper>
      </GridItem>
    </GridLayout>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import WidgetWrapper from '@/components/widgets/WidgetWrapper.vue'
import MetricCard from '@/components/widgets/MetricCard.vue'
import BarChartWidget from '@/components/dashboards/charts/BarChartWidget.vue'
import LineChartWidget from '@/components/dashboards/charts/LineChartWidget.vue'
import TableWidget from '@/components/dashboards/charts/TableWidget.vue'
import TextWidget from '@/components/dashboards/charts/TextWidget.vue'
import { getDashboard, getDashboardMetric } from '@/api/posts'

const route = useRoute()
const loading = ref(false)
const layout = ref([])
const widgetData = ref({})
const selectedTimeRange = ref('24h')
const customStartDate = ref(null)
const refreshInterval = ref('5s')
const countdown = ref(5)

const intervalsMap = {
  '5s': 5,
  '10s': 10,
  '30s': 30,
  '60s': 60,
}

const timeRangeOptions = [
  { label: 'Last 5 min', value: '5m' },
  { label: 'Last 30 min', value: '30m' },
  { label: 'Last 1 hour', value: '1h' },
  { label: 'Last 12 hours', value: '12h' },
  { label: 'Last 24 hours', value: '24h' },
  { label: 'Last 7 days', value: '7d' },
  { label: 'Last 30 days', value: '30d' },
  { label: 'Last 1 year', value: '365d' },
  { label: 'Custom...', value: 'custom' },
]

const refreshIntervalOptions = [
  { label: '5s', value: '5s' },
  { label: '10s', value: '10s' },
  { label: '30s', value: '30s' },
  { label: '60s', value: '60s' },
  { label: 'Off', value: 'off' },
]

let pollInterval = null
let countdownInterval = null

function startPolling() {
  stopPolling()
  if (refreshInterval.value === 'off') return

  const seconds = intervalsMap[refreshInterval.value]
  if (!seconds) return

  pollInterval = setInterval(fetchWidgetData, seconds * 1000)

  countdown.value = seconds
  countdownInterval = setInterval(() => {
    countdown.value = countdown.value > 0 ? countdown.value - 1 : seconds
  }, 1000)
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
  countdown.value = 0
}

function computeStartTime() {
  if (selectedTimeRange.value === 'custom') {
    return customStartDate.value ? customStartDate.value.toISOString() : null
  }

  const durations = {
    '5m': 5 * 60 * 1000,
    '30m': 30 * 60 * 1000,
    '1h': 60 * 60 * 1000,
    '12h': 12 * 60 * 60 * 1000,
    '24h': 24 * 60 * 60 * 1000,
    '7d': 7 * 24 * 60 * 60 * 1000,
    '30d': 30 * 24 * 60 * 60 * 1000,
    '365d': 365 * 24 * 60 * 60 * 1000,
  }

  const ms = durations[selectedTimeRange.value]
  return ms ? new Date(Date.now() - ms).toISOString() : null
}

function buildMetricsPayload() {
  const metrics = []
  for (const item of layout.value) {
    console.log('[DashboardView] buildMetricsPayload — item:', JSON.stringify({ i: item.i, type: item.type, config: item.config }))
    if (item.config?.query?.deviceID) {
      if (item.type === 'card') {
        metrics.push({
          type: 'card',
          deviceID: item.config.query.deviceID,
          column_name: item.config.query.column_name || '',
          data_type: item.config.query.data_type || 'number'
        })
      } else if (item.type === 'barchart') {
        metrics.push({
          type: 'barchart',
          deviceID: item.config.query.deviceID,
          x_axis: item.config.query.x_axis || '',
          y_axis: item.config.query.y_axis || ''
        })
      } else if (item.type === 'linechart') {
        metrics.push({
          type: 'linechart',
          deviceID: item.config.query.deviceID,
          x_axis: item.config.query.x_axis || '',
          y_axis: item.config.query.y_axis || ''
        })
      }
    } else {
      console.log('[DashboardView] buildMetricsPayload — SKIPPED (no deviceID):', item.i, item.type)
    }
  }
  console.log('[DashboardView] buildMetricsPayload — payload:', JSON.stringify({ metrics }))
  return { metrics, start: computeStartTime() }
}

function formatValue(value, dataType) {
  if (dataType === 'number' && typeof value === 'number') {
    return value % 1 === 0 ? String(value) : value.toFixed(1)
  }
  return String(value ?? '\u2014')
}

function lookupResult(results, deviceID, columnName) {
  console.log('[DashboardView] lookupResult — looking for deviceID:', deviceID, 'column_name:', columnName, 'in results:', JSON.stringify(results))
  if (!results) return undefined
  for (const m of results) {
    if (m.deviceID === deviceID && m.column_name === columnName) {
      console.log('[DashboardView] lookupResult — FOUND:', m.value)
      return m.value
    }
  }
  console.log('[DashboardView] lookupResult — NOT FOUND')
  return undefined
}

function lookupChartResult(results, deviceID) {
  if (!results) return null
  for (const m of results) {
    if (m.deviceID === deviceID && m.data_points) {
      return m.data_points
    }
  }
  return null
}

async function fetchWidgetData() {
  const payload = buildMetricsPayload()
  if (!payload.metrics.length) {
    return
  }

  for (const item of layout.value) {
    if (item.config?.query?.deviceID) {
      widgetData.value = {
        ...widgetData.value,
        [item.i]: { ...widgetData.value[item.i], loading: true, error: false }
      }
    }
  }

  try {
    const result = await getDashboardMetric(payload)
    const metrics = result?.metrics || []

    for (const item of layout.value) {
      const q = item.config?.query
      if (!q?.deviceID) continue

      if (item.type === 'card') {
        const liveValue = lookupResult(metrics, q.deviceID, q.column_name)
        const fallback = widgetData.value[item.i]?.value ?? item.cardValue
        widgetData.value = {
          ...widgetData.value,
          [item.i]: {
            value: liveValue !== undefined ? formatValue(liveValue, q.data_type) : fallback,
            loading: false,
            error: liveValue === undefined
          }
        }
      } else if (item.type === 'barchart' || item.type === 'linechart') {
        const dataPoints = lookupChartResult(metrics, q.deviceID)
        widgetData.value = {
          ...widgetData.value,
          [item.i]: {
            ...widgetData.value[item.i],
            dataPoints: dataPoints,
            loading: false,
            error: !dataPoints
          }
        }
      }
    }
  } catch (err) {
    console.error('[DashboardView] fetchWidgetData — API error:', err)
    for (const item of layout.value) {
      if (item.config?.query?.deviceID) {
        widgetData.value = {
          ...widgetData.value,
          [item.i]: { ...widgetData.value[item.i], loading: false, error: true }
        }
      }
    }
  }
  countdown.value = intervalsMap[refreshInterval.value] || 0
}

watch([selectedTimeRange, customStartDate], () => {
  fetchWidgetData()
})

watch(refreshInterval, () => {
  startPolling()
})

onMounted(async () => {
  const idParam = route.query.id
  console.log('[DashboardView] onMounted — route.query.id:', idParam)
  if (!idParam) return

  loading.value = true
  try {
    const data = await getDashboard({ id: parseInt(idParam) })
    console.log('[DashboardView] onMounted — getDashboard raw response:', JSON.stringify(data))
    if (data && data.layout_json) {
      const parsed = JSON.parse(data.layout_json)
      console.log('[DashboardView] onMounted — parsed layout:', JSON.stringify(parsed))
      layout.value = Array.isArray(parsed) ? parsed : []
      console.log('[DashboardView] onMounted — layout.value count:', layout.value.length)
    } else {
      console.log('[DashboardView] onMounted — no layout_json in response')
    }
  } catch (e) {
    console.error('Failed to load dashboard:', e)
  } finally {
    loading.value = false
  }

  fetchWidgetData()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.page-container {
  padding: 8px;
  min-height: 100%;
}

.grid-container {
  background-color: transparent;
  min-height: 600px;
}

.grid-item {
  border-radius: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: #555;
  text-align: center;
}

.empty-icon {
  font-size: var(--font-size-3xl);
  margin-bottom: 0.75rem;
  color: #444;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}

.time-range-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.controls-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: #1a1a1a;
  border: 1px solid #2a2a2e;
  border-radius: 8px;
}

.controls-panel__label {
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: #777;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  white-space: nowrap;
}

.time-range-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-range-select {
  min-width: 160px;
}

.time-range-datepicker {
  min-width: 220px;
}

.toolbar-separator {
  width: 1px;
  height: 24px;
  background: #2a2a2e;
  margin: 0 4px;
}

.refresh-interval-select {
  min-width: 80px;
}

.countdown-group {
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 56px;
  justify-content: flex-end;
}

.countdown-ring {
  flex-shrink: 0;
}

.countdown-text {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: #888;
  font-variant-numeric: tabular-nums;
}
</style>
