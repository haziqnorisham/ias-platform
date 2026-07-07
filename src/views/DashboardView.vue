<template>
  <div class="page-container">
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />

    <div v-if="!layout.length" class="empty-state">
      <i class="pi pi-chart-bar empty-icon"></i>
      <p>This dashboard has no widgets.</p>
    </div>

    <GridLayout
      v-else
      v-model:layout="layout"
      :col-num="12"
      :row-height="80"
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
        <WidgetWrapper :title="item.type === 'card' ? item.cardTitle : item.type === 'barchart' ? item.chartTitle : item.type === 'linechart' ? item.lineChartTitle : item.type === 'table' ? item.tableTitle : item.type === 'text' ? item.textTitle : item.type === 'stream' ? 'Live Stream' : ''">
          <MetricCard
            v-if="item.type === 'card'"
            :title="item.cardTitle"
            :value="widgetData[item.i]?.value ?? item.cardValue"
            :loading="widgetData[item.i]?.loading ?? false"
            :error="widgetData[item.i]?.error ?? false"
          />
          <BarChartWidget v-else-if="item.type === 'barchart'" :title="item.chartTitle" />
          <LineChartWidget v-else-if="item.type === 'linechart'" :title="item.lineChartTitle" />
          <TableWidget v-else-if="item.type === 'table'" :title="item.tableTitle" />
          <TextWidget v-else-if="item.type === 'text'" :title="item.textTitle" :text="item.textContent" />
          <StreamWidget v-else-if="item.type === 'stream'" />
        </WidgetWrapper>
      </GridItem>
    </GridLayout>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import WidgetWrapper from '@/components/widgets/WidgetWrapper.vue'
import MetricCard from '@/components/widgets/MetricCard.vue'
import BarChartWidget from '@/components/dashboards/charts/BarChartWidget.vue'
import LineChartWidget from '@/components/dashboards/charts/LineChartWidget.vue'
import TableWidget from '@/components/dashboards/charts/TableWidget.vue'
import TextWidget from '@/components/dashboards/charts/TextWidget.vue'
import StreamWidget from '@/components/widgets/StreamWidget.vue'
import { getDashboard, getDashboardMetric } from '@/api/posts'

const route = useRoute()
const loading = ref(false)
const layout = ref([])
const widgetData = ref({})

let pollInterval = null

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
  return { metrics }
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

async function fetchWidgetData() {
  const payload = buildMetricsPayload()
  if (!payload.metrics.length) {
    console.log('[DashboardView] fetchWidgetData — no metrics to fetch, layout has no configured deviceIDs')
    return
  }

  for (const item of layout.value) {
    if (item.type === 'card' && item.config?.query?.deviceID) {
      widgetData.value = {
        ...widgetData.value,
        [item.i]: { ...widgetData.value[item.i], loading: true, error: false }
      }
    }
  }

  try {
    console.log('[DashboardView] fetchWidgetData — calling getDashboardMetric with:', JSON.stringify(payload))
    const result = await getDashboardMetric(payload)
    console.log('[DashboardView] fetchWidgetData — API response:', JSON.stringify(result))
    const metrics = result?.metrics || []

    for (const item of layout.value) {
      if (item.type === 'card' && item.config?.query?.deviceID) {
        const q = item.config.query
        const liveValue = lookupResult(metrics, q.deviceID, q.column_name)
        const fallback = widgetData.value[item.i]?.value ?? item.cardValue
        console.log('[DashboardView] fetchWidgetData — widget', item.i, 'liveValue:', liveValue, 'fallback:', fallback)
        widgetData.value = {
          ...widgetData.value,
          [item.i]: {
            value: liveValue !== undefined ? formatValue(liveValue, q.data_type) : fallback,
            loading: false,
            error: liveValue === undefined
          }
        }
      }
    }
    console.log('[DashboardView] fetchWidgetData — final widgetData:', JSON.stringify(widgetData.value))
  } catch (err) {
    console.error('[DashboardView] fetchWidgetData — API error:', err)
    for (const item of layout.value) {
      if (item.type === 'card' && item.config?.query?.deviceID) {
        widgetData.value = {
          ...widgetData.value,
          [item.i]: { ...widgetData.value[item.i], loading: false, error: true }
        }
      }
    }
  }
}

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
  pollInterval = setInterval(fetchWidgetData, 5000)
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
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
  font-size: 3rem;
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
</style>
