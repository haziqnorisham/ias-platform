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
        <WidgetWrapper :title="item.type === 'card' ? item.cardTitle : item.type === 'barchart' ? item.chartTitle : item.type === 'table' ? item.tableTitle : item.type === 'text' ? item.textTitle : item.type === 'stream' ? 'Live Stream' : ''">
          <MetricCard v-if="item.type === 'card'" :title="item.cardTitle" :value="item.cardValue" />
          <BarChartWidget v-else-if="item.type === 'barchart'" :title="item.chartTitle" />
          <TableWidget v-else-if="item.type === 'table'" :title="item.tableTitle" />
          <TextWidget v-else-if="item.type === 'text'" :title="item.textTitle" :text="item.textContent" />
          <StreamWidget v-else-if="item.type === 'stream'" />
        </WidgetWrapper>
      </GridItem>
    </GridLayout>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import WidgetWrapper from '@/components/widgets/WidgetWrapper.vue'
import MetricCard from '@/components/widgets/MetricCard.vue'
import BarChartWidget from '@/components/dashboards/charts/BarChartWidget.vue'
import TableWidget from '@/components/dashboards/charts/TableWidget.vue'
import TextWidget from '@/components/dashboards/charts/TextWidget.vue'
import StreamWidget from '@/components/widgets/StreamWidget.vue'
import { getDashboard } from '@/api/posts'

const route = useRoute()
const loading = ref(false)
const layout = ref([])

onMounted(async () => {
  const idParam = route.query.id
  if (!idParam) return

  loading.value = true
  try {
    const data = await getDashboard({ id: parseInt(idParam) })
    if (data && data.layout_json) {
      const parsed = JSON.parse(data.layout_json)
      layout.value = Array.isArray(parsed) ? parsed : []
    }
  } catch (e) {
    console.error('Failed to load dashboard:', e)
  } finally {
    loading.value = false
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
