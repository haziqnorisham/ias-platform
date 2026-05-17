<template>
  <div class="editor-container">
    <div class="editor-header">
      <div class="header-left">
        <i class="pi pi-palette header-icon"></i>
        <InputText placeholder="Untitled Dashboard" size="small" class="dashboard-name" />
      </div>
      <div class="header-right">
        <Button label="Cancel" icon="pi pi-times" severity="secondary" text size="small" />
        <Button label="Save Dashboard" icon="pi pi-save" size="small" />
      </div>
    </div>

    <div class="editor-toolbar">
      <div class="toolbar-section">
        <ButtonGroup>
          <Button icon="pi pi-chart-bar" label="Bar Chart" size="small" severity="secondary" @click="addChart" />
          <Button icon="pi pi-hashtag" label="Metric" size="small" severity="secondary" @click="addMetric" />
          <Button icon="pi pi-table" label="Table" size="small" severity="secondary" @click="addTable" />
          <Button icon="pi pi-pen-to-square" label="Text" size="small" severity="secondary" @click="addText" />
        </ButtonGroup>
      </div>
    </div>

    <div class="canvas-area">
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
          <EditorWidget :widget="item" :devices="deviceOptions" @delete="deleteWidget" @update="updateWidget" />
        </GridItem>
      </GridLayout>

      <div v-if="!layout.length" class="canvas-placeholder">
        <i class="pi pi-plus-circle placeholder-icon"></i>
        <p class="placeholder-text">Drag widgets here or use the toolbar to add content</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GridLayout, GridItem } from 'vue3-grid-layout'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import EditorWidget from './EditorWidget.vue'
import { getAllDevices } from '@/api/posts'

const deviceOptions = ref([])

onMounted(() => {
  getAllDevices().then(result => {
    deviceOptions.value = (result || []).map(d => ({
      label: `${d.Name || d.Id} (${d.Id})`,
      value: String(d.Id)
    }))
  }).catch(() => { deviceOptions.value = [] })
})

const layout = ref([
  { x: 0, y: 0, w: 4, h: 2, i: '0', type: 'card', cardTitle: 'Starter Metric', cardValue: '—' }
])

let nextId = 1

function addMetric() {
  const maxY = layout.value.reduce((max, item) => Math.max(max, item.y + item.h), 0)
  layout.value.push({
    x: 0,
    y: maxY,
    w: 4,
    h: 2,
    i: String(nextId++),
    type: 'card',
    cardTitle: 'New Metric',
    cardValue: '—'
  })
}

function addChart() {
  const maxY = layout.value.reduce((max, item) => Math.max(max, item.y + item.h), 0)
  layout.value.push({
    x: 0,
    y: maxY,
    w: 6,
    h: 4,
    i: String(nextId++),
    type: 'barchart',
    chartTitle: 'New Bar Chart'
  })
}

function addTable() {
  const maxY = layout.value.reduce((max, item) => Math.max(max, item.y + item.h), 0)
  layout.value.push({
    x: 0,
    y: maxY,
    w: 6,
    h: 4,
    i: String(nextId++),
    type: 'table',
    tableTitle: 'New Table'
  })
}

function addText() {
  const maxY = layout.value.reduce((max, item) => Math.max(max, item.y + item.h), 0)
  layout.value.push({
    x: 0,
    y: maxY,
    w: 4,
    h: 2,
    i: String(nextId++),
    type: 'text',
    textTitle: 'New Text',
    textContent: 'Double-click to edit this text.'
  })
}

function deleteWidget(id) {
  const idx = layout.value.findIndex(item => item.i === id)
  if (idx !== -1) layout.value.splice(idx, 1)
}

function updateWidget(updated) {
  const item = layout.value.find(item => item.i === updated.i)
  if (item) {
    if (updated.cardTitle !== undefined) item.cardTitle = updated.cardTitle
    if (updated.cardValue !== undefined) item.cardValue = updated.cardValue
    if (updated.chartTitle !== undefined) item.chartTitle = updated.chartTitle
    if (updated.tableTitle !== undefined) item.tableTitle = updated.tableTitle
    if (updated.textTitle !== undefined) item.textTitle = updated.textTitle
    if (updated.textContent !== undefined) item.textContent = updated.textContent
    if (updated.config !== undefined) item.config = updated.config
  }
}
</script>

<style scoped>
.editor-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  border-radius: 10px;
  border: 1px solid #212121;
  background-color: #0e0e10;
  overflow: hidden;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  border-bottom: 1px solid #212121;
  background-color: #141416;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dashboard-name {
  width: 260px;
  font-weight: 600;
  font-size: 0.9rem;
}

.header-icon {
  font-size: 1.1rem;
  color: #888;
}

.header-right {
  display: flex;
  gap: 0.5rem;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  padding: 0.4rem 1rem;
  border-bottom: 1px solid #212121;
  background-color: #141416;
  gap: 0.75rem;
  flex-shrink: 0;
}

.toolbar-section {
  display: flex;
  align-items: center;
}

.canvas-area {
  flex: 1;
  position: relative;
  overflow: auto;
  background-color: #0a0a0c;
  background-image:
    radial-gradient(circle, #1a1a1e 1px, transparent 1px);
  background-size: 40px 40px;
}

.grid-container {
  background-color: transparent;
  min-height: 100%;
}

.grid-item {
  border-radius: 8px;
  transition: box-shadow 0.15s ease;
}

.grid-item:hover {
  box-shadow: 0 0 0 1px #555;
}

.canvas-placeholder {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.placeholder-icon {
  font-size: 3rem;
  color: #2a2a2e;
  margin-bottom: 0.75rem;
}

.placeholder-text {
  color: #444;
  font-size: 0.9rem;
  margin: 0;
}
</style>
