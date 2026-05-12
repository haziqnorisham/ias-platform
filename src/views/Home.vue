<template>
  <div class="page-container">
    <GridLayout
      v-model:layout="layout"
      :col-num="12"
      :row-height="80"
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
        <WidgetWrapper>
          <Card v-if="item.type === 'card'" :title="item.cardTitle" :value="item.cardValue" />
          <BarChart
            v-else-if="item.type === 'barchart'"
            title="Example Bar Chart"
            subtitle="This is a sample subtitle for the bar chart."
          />
          <eSampleChart
            v-else-if="item.type === 'esamplechart'"
            title="Example ECharts Bar Chart"
            subtitle="This is a sample subtitle for the ECharts bar chart."
          />
          <div v-else-if="item.type === 'statictext'" class="static-text-widget">
            {{ item.text }}
          </div>
        </WidgetWrapper>
      </GridItem>
    </GridLayout>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import Card from '../components/Card.vue'
import BarChart from '../components/charts/BarChart.vue'
import eSampleChart from '../components/charts/eSampleChart.vue'
import WidgetWrapper from '../components/widgets/WidgetWrapper.vue'
import { getAllDevices, getRawIngest } from '../api/posts'

const totalDevices = ref(0)
const totalIngestLogs = ref(0)

const defaultLayout = [
  { x: 0, y: 0, w: 3, h: 2, i: '0', type: 'card', cardTitle: 'Total Devices', cardValue: 0 },
  { x: 3, y: 0, w: 3, h: 2, i: '1', type: 'card', cardTitle: 'Ingest Logs', cardValue: 0 },
  { x: 6, y: 0, w: 3, h: 2, i: '2', type: 'card', cardTitle: 'Active Gateways', cardValue: '6' },
  { x: 9, y: 0, w: 3, h: 2, i: '3', type: 'card', cardTitle: 'Active Alarms', cardValue: '3' },
  { x: 0, y: 2, w: 5, h: 4, i: '4', type: 'barchart' },
  { x: 5, y: 2, w: 7, h: 4, i: '5', type: 'esamplechart' },
  { x: 0, y: 6, w: 12, h: 1, i: '6', type: 'statictext', text: 'IAS - IOT DEV' },
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
})
</script>

<style scoped>
.page-container {
  padding: 8px;
  height: calc(100vh - 70px);
  overflow: auto;
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

.static-text-widget {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  font-size: 1.5rem;
  font-weight: 600;
  color: #ccc;
  letter-spacing: 0.05em;
  user-select: none;
}
</style>
