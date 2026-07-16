<template>
  <div class="page-container">
    <GridLayout
      v-model:layout="layout"
      :col-num="24"
      :row-height="50"
      :is-draggable="true"
      :is-resizable="true"
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
      >
        <WidgetWrapper :title="item.type === 'card' ? item.cardTitle : item.type === 'extensionslist' ? 'Available Extensions' : item.type === 'dashboardslist' ? 'Custom Dashboards' : ''">
          <MetricCard v-if="item.type === 'card'" :value="item.cardValue" hideTitle />
          <template v-else-if="item.type === 'extensionslist'">
            <div v-if="extensionsLoading" class="ext-state">
              <ProgressSpinner style="width: 32px; height: 32px" strokeWidth="4" />
              <span>Loading extensions...</span>
            </div>
            <div v-else-if="extensionsError" class="ext-state ext-error">
              <i class="pi pi-exclamation-triangle" style="font-size: 1.5rem; color: #e57373"></i>
              <span>Failed to load extensions: {{ extensionsError }}</span>
            </div>
            <div v-else-if="!extensions.length" class="ext-state">
              <i class="pi pi-box" style="font-size: 1.5rem; color: #555"></i>
              <span>No extensions available.</span>
            </div>
            <div v-else class="ext-list">
              <div
                v-for="ext in extensions"
                :key="ext.name"
                class="ext-list-item"
                @click="router.push(`/extensions/${ext.name}`)"
              >
                <div class="ext-name">{{ ext.name }}</div>
                <Button label="Open" icon="pi pi-arrow-up-right" text size="small" @click.stop="router.push(`/extensions/${ext.name}`)" />
              </div>
            </div>
          </template>
          <template v-else-if="item.type === 'dashboardslist'">
            <div v-if="dashboardsLoading" class="ext-state">
              <ProgressSpinner style="width: 32px; height: 32px" strokeWidth="4" />
              <span>Loading dashboards...</span>
            </div>
            <div v-else-if="dashboardsError" class="ext-state ext-error">
              <i class="pi pi-exclamation-triangle" style="font-size: 1.5rem; color: #e57373"></i>
              <span>Failed to load dashboards: {{ dashboardsError }}</span>
            </div>
            <div v-else-if="!dashboards.length" class="ext-state">
              <i class="pi pi-chart-bar" style="font-size: 1.5rem; color: #555"></i>
              <span>No dashboards available.</span>
            </div>
            <div v-else class="ext-list">
              <div
                v-for="d in dashboards"
                :key="d.id"
                class="ext-list-item"
                @click="router.push(`/dashboards/view?id=${d.id}`)"
              >
                <div class="ext-name">{{ d.name || 'Untitled Dashboard' }}</div>
                <Button label="Open" icon="pi pi-arrow-up-right" text size="small" @click.stop="router.push(`/dashboards/view?id=${d.id}`)" />
              </div>
            </div>
          </template>
        </WidgetWrapper>
      </GridItem>
    </GridLayout>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import MetricCard from '../components/widgets/MetricCard.vue'
import WidgetWrapper from '../components/widgets/WidgetWrapper.vue'
import ProgressSpinner from 'primevue/progressspinner'
import Button from 'primevue/button'
import { getAllDevices, getRawIngest, getDashboards } from '../api/posts'
import { getExtensions } from '../api/extensions'

const router = useRouter()
const totalDevices = ref(0)
const totalIngestLogs = ref(0)
const extensions = ref([])
const extensionsLoading = ref(true)
const extensionsError = ref(null)
const dashboards = ref([])
const dashboardsLoading = ref(true)
const dashboardsError = ref(null)

const defaultLayout = [
  { x: 0, y: 0, w: 5, h: 2, i: '0', type: 'card', cardTitle: 'Total Devices', cardValue: 0 },
  { x: 5, y: 0, w: 5, h: 2, i: '1', type: 'card', cardTitle: 'Ingest Logs', cardValue: 0 },
  { x: 0, y: 2, w: 10, h: 4, i: '2', type: 'extensionslist' },
  { x: 0, y: 6, w: 10, h: 4, i: '3', type: 'dashboardslist' },
]

const layout = ref(structuredClone(defaultLayout))

const updateWidgetValue = (widgetId, value) => {
  const widget = layout.value.find(w => w.i === widgetId)
  if (widget) widget.cardValue = value
}

onMounted(async () => {
  try {
    const devices = await getAllDevices()
    const count = Array.isArray(devices) ? devices.length : 0
    totalDevices.value = count
    updateWidgetValue('0', count)
  } catch (e) {
    console.error('Failed to fetch devices count:', e)
  }

  try {
    const result = await getRawIngest({ Limit: 1, Offset: 0, SortByMsgID: 'desc', Status: '' })
    const count = result?.total ?? 0
    totalIngestLogs.value = count
    updateWidgetValue('1', count)
  } catch (e) {
    console.error('Failed to fetch ingest log count:', e)
  }

  try {
    const result = await getExtensions()
    extensions.value = Array.isArray(result) ? result : []
  } catch (e) {
    extensionsError.value = e.message
    console.error('Failed to fetch extensions:', e)
  } finally {
    extensionsLoading.value = false
  }

  try {
    const result = await getDashboards()
    dashboards.value = Array.isArray(result) ? result : []
  } catch (e) {
    dashboardsError.value = e.message
    console.error('Failed to fetch dashboards:', e)
  } finally {
    dashboardsLoading.value = false
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
  transition: box-shadow 0.15s ease;
}

.grid-item:hover {
  box-shadow: 0 0 0 1px #444;
}

.ext-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 0.75rem;
  color: #888;
}

.ext-error {
  color: #e57373;
}

.ext-list {
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  max-height: 100%;
}

.ext-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #2a2a2e;
  cursor: pointer;
  transition: background 0.15s;
}

.ext-list-item:hover {
  background: #202024;
}

.ext-list-item:last-child {
  border-bottom: none;
}

.ext-name {
  font-weight: 600;
  color: #e0e0e0;
}
</style>
