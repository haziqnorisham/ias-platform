<script setup>
    import { ref, onMounted, onUnmounted } from 'vue';
    import * as echarts from 'echarts';
    import BaseCard from "./BaseCard.vue";

    defineProps({
      title: { type: String, default: 'ECharts Bar Chart' },
      subtitle: { type: String, default: '' },
      width: { type: String, default: '50vw' },
      height: { type: String, default: 'auto' },
    });

    // Theme colors matching the app
    const THEME = {
      primary: '#48897b',
      primaryLight: '#5fa89a',
      bgDark: '#18181B',
      cardBg: '#202024',
      border: '#2a2a2e',
      textMuted: '#6d6d6d',
      textLight: '#a0a0a0',
      textWhite: '#ffffff',
    };

    const chartContainer = ref(null);
    let chart = null;

    onMounted(() => {
      if (!chartContainer.value) return;

      chart = echarts.init(chartContainer.value, 'dark');

      const option = {
        backgroundColor: 'transparent',
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow',
            shadowStyle: {
              color: 'rgba(72, 137, 123, 0.08)',
            },
          },
          backgroundColor: THEME.cardBg,
          borderColor: THEME.border,
          borderWidth: 1,
          textStyle: {
            color: THEME.textWhite,
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
            color: THEME.textMuted,
            fontFamily: '"Space Grotesk", sans-serif',
            fontSize: 11,
          },
          axisLine: {
            lineStyle: {
              color: THEME.border,
            },
          },
          axisTick: {
            show: false,
          },
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            color: THEME.textMuted,
            fontFamily: '"Space Grotesk", sans-serif',
            fontSize: 11,
          },
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              color: 'rgba(255, 255, 255, 0.06)',
              type: 'dashed',
            },
          },
        },
        series: [
          {
            name: 'Sales',
            type: 'bar',
            data: [120, 200, 150, 80, 70, 110, 130],
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: THEME.primaryLight },
                { offset: 1, color: THEME.primary },
              ]),
              borderRadius: [6, 6, 0, 0],
            },
            barWidth: '55%',
            emphasis: {
              itemStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#7dc4b5' },
                  { offset: 1, color: THEME.primaryLight },
                ]),
              },
            },
          },
        ],
      };

      chart.setOption(option);

      window.addEventListener('resize', handleResize);
    });

    onUnmounted(() => {
      window.removeEventListener('resize', handleResize);
      if (chart) {
        chart.dispose();
        chart = null;
      }
    });

    function handleResize() {
      chart?.resize();
    }
</script>

<template>
    <BaseCard :title="title" :width="width" :height="height" :subtitle="subtitle">
        <div ref="chartContainer" class="chart-container"></div>
    </BaseCard>
</template>

<style scoped>
.chart-container {
    width: 100%;
    height: 100%;
    min-height: 300px;
}
</style>