<template>
  <div ref="chartContainer" class="chart-container">
    <div v-if="loading" class="chart-overlay">
      <ProgressSpinner style="width: 28px; height: 28px" strokeWidth="4" />
    </div>
    <div v-else-if="error" class="chart-overlay chart-overlay--error">
      <i class="pi pi-exclamation-triangle" style="font-size: 1.5rem; color: #e57373"></i>
    </div>
    <div v-else-if="hasData" class="chart-actions">
      <Button
        :icon="zoomActive ? 'pi pi-times' : 'pi pi-search-plus'"
        text
        size="small"
        :title="zoomActive ? 'Cancel zoom' : 'Drag to zoom'"
        @click.stop="toggleZoom"
      />
      <Button
        icon="pi pi-undo"
        text
        size="small"
        title="Reset zoom"
        @click.stop="resetZoom"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import ProgressSpinner from 'primevue/progressspinner'
import Button from 'primevue/button'

const props = defineProps({
  title: { type: String, default: 'Bar Chart' },
  dataPoints: { type: Array, default: null },
  loading: { type: Boolean, default: false },
  error: { type: Boolean, default: false }
})

const chartContainer = ref(null)
const zoomActive = ref(false)
let chart = null
let resizeObserver = null

const hasData = computed(() => props.dataPoints && props.dataPoints.length > 0)

const PRIMARY = '#48897b'
const PRIMARY_LIGHT = '#5fa89a'

function toggleZoom() {
  if (!chart) return
  zoomActive.value = !zoomActive.value
  chart.dispatchAction({
    type: 'takeGlobalCursor',
    key: 'dataZoomSelect',
    dataZoomSelectActive: zoomActive.value,
  })
}

function resetZoom() {
  if (!chart) return
  zoomActive.value = false
  chart.dispatchAction({
    type: 'takeGlobalCursor',
    key: 'dataZoomSelect',
    dataZoomSelectActive: false,
  })
  chart.dispatchAction({ type: 'restore' })
}

function buildOption(dataPoints) {
  const hasData = dataPoints && dataPoints.length > 0
  const seriesData = hasData
    ? dataPoints.map(p => [p.x, typeof p.y === 'number' ? p.y : parseFloat(p.y)])
    : [120, 200, 150, 80, 70, 110, 130]

  const xAxis = hasData ? {
    type: 'time',
    axisLabel: {
      color: '#6d6d6d',
      fontFamily: '"Space Grotesk", sans-serif',
      fontSize: 9,
    },
    axisLine: { lineStyle: { color: '#2a2a2e' } },
    axisTick: { show: false },
    splitLine: { show: false },
  } : {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    axisLabel: {
      color: '#6d6d6d',
      fontFamily: '"Space Grotesk", sans-serif',
      fontSize: 9,
    },
    axisLine: { lineStyle: { color: '#2a2a2e' } },
    axisTick: { show: false },
  }

  return {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: '#202024',
      borderColor: '#2a2a2e',
      borderWidth: 1,
      textStyle: {
        color: '#ffffff',
        fontFamily: '"Space Grotesk", sans-serif',
        fontSize: 10,
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '8%',
      containLabel: true,
    },
    xAxis,
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#6d6d6d',
        fontFamily: '"Space Grotesk", sans-serif',
        fontSize: 9,
      },
      splitLine: {
        lineStyle: { color: 'rgba(255,255,255,0.06)', type: 'dashed' },
      },
    },
    series: [
      {
        name: 'Series',
        type: 'bar',
        data: seriesData,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: PRIMARY_LIGHT },
            { offset: 1, color: PRIMARY },
          ]),
          borderRadius: [6, 6, 0, 0],
        },
        barWidth: '55%',
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#7dc4b5' },
              { offset: 1, color: PRIMARY_LIGHT },
            ]),
          },
        },
      },
    ],
    dataZoom: [
      {
        type: 'inside',
        xAxisIndex: [0],
        zoomOnMouseWheel: true,
        moveOnMouseMove: true,
      },
    ],
  }
}

onMounted(() => {
  if (!chartContainer.value) return

  chart = echarts.init(chartContainer.value, 'dark')
  chart.setOption(buildOption(props.dataPoints))

  resizeObserver = new ResizeObserver(() => {
    chart?.resize()
  })
  resizeObserver.observe(chartContainer.value)

  window.addEventListener('resize', handleResize)
})

watch(() => props.dataPoints, (newVal) => {
  if (chart) {
    chart.setOption(buildOption(newVal))
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  if (chart) {
    chart.dispose()
    chart = null
  }
})

function handleResize() {
  chart?.resize()
}
</script>

<style scoped>
.chart-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.chart-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(26, 26, 26, 0.85);
  z-index: 2;
  pointer-events: none;
}

.chart-overlay--error {
  background: rgba(26, 26, 26, 0.7);
}

.chart-actions {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 2px;
  z-index: 1;
}
</style>
