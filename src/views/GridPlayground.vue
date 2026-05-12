<script setup>
import { ref } from 'vue'
import { GridLayout, GridItem } from 'vue3-grid-layout'

// 3 blank black widgets — draggable and resizable
const layout = ref([
  { x: 0, y: 0, w: 4, h: 3, i: '0' },
  { x: 4, y: 0, w: 4, h: 3, i: '1' },
  { x: 8, y: 0, w: 4, h: 3, i: '2' },
])

const colNum = 12
const rowHeight = 150
</script>

<template>
  <div class="page-container">
    <h2 class="page-title">Grid Playground</h2>

    <GridLayout
      v-model:layout="layout"
      :col-num="colNum"
      :row-height="rowHeight"
      :is-draggable="true"
      :is-resizable="true"
      :margin="[16, 16]"
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
        class="widget"
      >
        <div class="widget-content">
          <span class="widget-label">Widget {{ parseInt(item.i) + 1 }}</span>
        </div>
      </GridItem>
    </GridLayout>
  </div>
</template>

<style scoped>
.page-container {
  padding: 0;
}

.page-title {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.grid-container {
  background-color: transparent;
  min-height: 400px;
}

.widget {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  transition: box-shadow 0.15s ease;
  cursor: grab;
}

.widget:active {
  cursor: grabbing;
}

.widget:hover {
  box-shadow: 0 0 0 1px #444;
}

.widget-content {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
}

.widget-label {
  color: #666;
  font-size: 1rem;
  font-weight: 500;
  user-select: none;
  pointer-events: none;
}
</style>
