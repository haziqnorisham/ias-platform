<template>
  <div class="editor-widget">
    <div class="editor-widget-header">
      <div class="header-leading">
        <span
          class="header-title"
          :class="{ 'header-title--editing': editing }"
          @dblclick="startEdit"
        >
          {{ widgetTitle }}
        </span>
        <Button
          icon="pi pi-pencil"
          severity="secondary"
          text
          size="small"
          class="edit-btn"
          @click="startEdit"
        />
      </div>
      <div class="header-trailing">
        <Button
          icon="pi pi-times"
          severity="danger"
          text
          rounded
          size="small"
          @click="emit('delete', widget.i)"
        />
        <span class="drag-handle" title="Drag to reposition">⠿</span>
      </div>
    </div>

    <div class="editor-widget-body">
      <MetricCard v-if="widget.type === 'card'" :title="widget.cardTitle" :value="widget.cardValue" />
      <BarChartWidget v-else-if="widget.type === 'barchart'" :title="widget.chartTitle" />
      <LineChartWidget v-else-if="widget.type === 'linechart'" :title="widget.lineChartTitle" />
      <TableWidget v-else-if="widget.type === 'table'" :title="widget.tableTitle" />
      <TextWidget v-else-if="widget.type === 'text'" :title="widget.textTitle" :text="widget.textContent" />
    </div>

    <Dialog
      v-model:visible="editing"
      header="Edit Widget"
      :modal="true"
      :closable="false"
      :style="{ width: '600px' }"
      class="edit-dialog"
    >
      <div class="edit-section">
        <div class="edit-section-label">Display</div>
        <div class="edit-field">
          <label for="widgetTitle">Title</label>
          <InputText id="widgetTitle" v-model="draftTitle" placeholder="Widget title" class="form-input" @keydown.enter="saveTitle" />
        </div>
        <div v-if="widget.type === 'card' && !widgetQueryConfig" class="edit-field">
          <label for="widgetValue">Value</label>
          <InputText id="widgetValue" v-model="draftValue" placeholder="Widget value" class="form-input" @keydown.enter="saveTitle" />
        </div>
        <div v-if="widget.type === 'text'" class="edit-field">
          <label for="widgetValue">Content</label>
          <InputText id="widgetValue" v-model="draftValue" placeholder="Text content" class="form-input" @keydown.enter="saveTitle" />
        </div>
      </div>

      <template v-if="widget.type === 'card' || widget.type === 'barchart' || widget.type === 'linechart'">
        <hr class="edit-divider" />
        <DataSourceConfig
          v-model="widgetQueryConfig"
          :widget-type="widget.type"
          :devices="deviceOptions"
        />
      </template>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" @click="editing = false" class="p-button-text" />
        <Button label="Save" icon="pi pi-check" @click="saveTitle" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import MetricCard from '../widgets/MetricCard.vue'
import BarChartWidget from './charts/BarChartWidget.vue'
import LineChartWidget from './charts/LineChartWidget.vue'
import TableWidget from './charts/TableWidget.vue'
import TextWidget from './charts/TextWidget.vue'
import DataSourceConfig from './DataSourceConfig.vue'

const props = defineProps({
  widget: {
    type: Object,
    required: true
  },
  devices: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['delete', 'update'])

const editing = ref(false)
const draftTitle = ref('')
const draftValue = ref('')
const widgetQueryConfig = ref(null)
const deviceOptions = ref([])

const widgetTitle = computed(() => {
  switch (props.widget.type) {
    case 'barchart': return props.widget.chartTitle
    case 'linechart': return props.widget.lineChartTitle
    case 'table': return props.widget.tableTitle
    case 'text': return props.widget.textTitle
    default: return props.widget.cardTitle || '—'
  }
})

watch(() => props.widget, (w) => {
  switch (w.type) {
    case 'barchart':
      draftTitle.value = w.chartTitle || ''
      draftValue.value = ''
      break
    case 'linechart':
      draftTitle.value = w.lineChartTitle || ''
      draftValue.value = ''
      break
    case 'table':
      draftTitle.value = w.tableTitle || ''
      draftValue.value = ''
      break
    case 'text':
      draftTitle.value = w.textTitle || ''
      draftValue.value = w.textContent || ''
      break
    default:
      draftTitle.value = w.cardTitle || ''
      draftValue.value = w.cardValue || ''
  }
}, { immediate: true })

function startEdit() {
  const w = props.widget
  switch (w.type) {
    case 'barchart':
      draftTitle.value = w.chartTitle || ''
      draftValue.value = ''
      break
    case 'linechart':
      draftTitle.value = w.lineChartTitle || ''
      draftValue.value = ''
      break
    case 'table':
      draftTitle.value = w.tableTitle || ''
      draftValue.value = ''
      break
    case 'text':
      draftTitle.value = w.textTitle || ''
      draftValue.value = w.textContent || ''
      break
    default:
      draftTitle.value = w.cardTitle || ''
      draftValue.value = w.cardValue || ''
  }

  if (w.type === 'card' || w.type === 'barchart' || w.type === 'linechart') {
    deviceOptions.value = props.devices
    widgetQueryConfig.value = w.config?.query ?? null
  }

  editing.value = true
}

function saveTitle() {
  const newTitle = draftTitle.value.trim()
  const newValue = draftValue.value.trim()
  const updatePayload = { i: props.widget.i }

  switch (props.widget.type) {
    case 'barchart':
      if (newTitle) updatePayload.chartTitle = newTitle
      break
    case 'linechart':
      if (newTitle) updatePayload.lineChartTitle = newTitle
      break
    case 'table':
      if (newTitle) updatePayload.tableTitle = newTitle
      break
    case 'text':
      updatePayload.textTitle = newTitle || props.widget.textTitle
      updatePayload.textContent = newValue || props.widget.textContent
      break
    default:
      updatePayload.cardTitle = newTitle || props.widget.cardTitle
      updatePayload.cardValue = newValue || props.widget.cardValue
      break
  }

  if (widgetQueryConfig.value && widgetQueryConfig.value.deviceID) {
    updatePayload.config = { query: { ...widgetQueryConfig.value } }
  }

  if (Object.keys(updatePayload).length > 1) {
    emit('update', updatePayload)
  }
  editing.value = false
}
</script>

<style scoped>
.editor-widget {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background-color: #1a1a1a;
  border-radius: 8px;
}

.editor-widget-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  flex-shrink: 0;
  background-color: #202024;
  border-bottom: 1px solid #2a2a2e;
  min-height: 28px;
}

.header-leading {
  display: flex;
  align-items: center;
  gap: 4px;
  overflow: hidden;
  min-width: 0;
}

.header-title {
  font-size: var(--font-size-xs);
  font-weight: 500;
  color: #aaa;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: default;
  user-select: none;
}

.header-title--editing {
  color: #90caf9;
}

.edit-btn {
  opacity: 0;
  transition: opacity 0.15s ease;
}

.editor-widget-header:hover .edit-btn {
  opacity: 1;
}

.header-trailing {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.drag-handle {
  color: #555;
  font-size: var(--font-size-md);
  letter-spacing: 2px;
  user-select: none;
  line-height: 1;
  cursor: grab;
  padding: 2px 4px;
  border-radius: 4px;
  transition: color 0.15s ease;
}

.drag-handle:hover {
  color: #888;
}

.drag-handle:active {
  cursor: grabbing;
}

.editor-widget-body {
  flex: 1;
  overflow: hidden;
  min-height: 0;
  position: relative;
}

.edit-dialog :deep(.p-dialog-header) {
  border-bottom: 1px solid #212121;
  padding: 1.25rem 1.5rem;
}

.edit-dialog :deep(.p-dialog-content) {
  padding: 1.5rem;
}

.edit-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.edit-section-label {
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: #777;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.edit-field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.edit-field label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: #a0a0a0;
}

.form-input {
  width: 100%;
}

.edit-divider {
  border: none;
  border-top: 1px solid #212121;
  margin: 1.5rem 0;
}
</style>
