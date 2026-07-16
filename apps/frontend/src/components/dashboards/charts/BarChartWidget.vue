<template>
  <div ref="chartContainer" class="chart-container"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'

defineProps({
  title: { type: String, default: 'Bar Chart' }
})

const chartContainer = ref(null)
let chart = null
let resizeObserver = null

const PRIMARY = '#48897b'
const PRIMARY_LIGHT = '#5fa89a'

onMounted(() => {
  if (!chartContainer.value) return

  chart = echarts.init(chartContainer.value, 'dark')

  const option = {
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
        fontSize: 12,
      },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '8%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
      axisLabel: {
        color: '#6d6d6d',
        fontFamily: '"Space Grotesk", sans-serif',
        fontSize: 10,
      },
      axisLine: { lineStyle: { color: '#2a2a2e' } },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: '#6d6d6d',
        fontFamily: '"Space Grotesk", sans-serif',
        fontSize: 10,
      },
      splitLine: {
        lineStyle: { color: 'rgba(255,255,255,0.06)', type: 'dashed' },
      },
    },
    series: [
      {
        name: 'Series',
        type: 'bar',
        data: [120, 200, 150, 80, 70, 110, 130],
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
  }

  chart.setOption(option)

  resizeObserver = new ResizeObserver(() => {
    chart?.resize()
  })
  resizeObserver.observe(chartContainer.value)

  window.addEventListener('resize', handleResize)
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
</style>
