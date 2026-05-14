<template>
  <div class="editor-widget">
    <div class="editor-widget-header">
      <div class="header-leading">
        <span
          class="header-title"
          :class="{ 'header-title--editing': editing }"
          @dblclick="startEdit"
        >
          {{ widget.cardTitle || '—' }}
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
      <MetricCard :title="widget.cardTitle" :value="widget.cardValue" />
    </div>

    <Dialog
      v-model:visible="editing"
      header="Edit Widget"
      :modal="true"
      :closable="false"
      :style="{ width: '420px' }"
      class="edit-dialog"
    >
      <div class="edit-form">
        <div class="edit-field">
          <label for="widgetTitle">Title</label>
          <InputText id="widgetTitle" v-model="draftTitle" placeholder="Widget title" class="form-input" @keydown.enter="saveTitle" />
        </div>
        <div class="edit-field">
          <label for="widgetValue">Value</label>
          <InputText id="widgetValue" v-model="draftValue" placeholder="Widget value" class="form-input" @keydown.enter="saveTitle" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" icon="pi pi-times" @click="editing = false" class="p-button-text" />
        <Button label="Save" icon="pi pi-check" @click="saveTitle" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import MetricCard from '../widgets/MetricCard.vue'

const props = defineProps({
  widget: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['delete', 'update'])

const editing = ref(false)
const draftTitle = ref('')
const draftValue = ref('')

watch(() => props.widget, (w) => {
  draftTitle.value = w.cardTitle || ''
  draftValue.value = w.cardValue || ''
}, { immediate: true })

function startEdit() {
  draftTitle.value = props.widget.cardTitle || ''
  draftValue.value = props.widget.cardValue || ''
  editing.value = true
}

function saveTitle() {
  const newTitle = draftTitle.value.trim()
  const newValue = draftValue.value.trim()
  if (newTitle || newValue) {
    emit('update', {
      i: props.widget.i,
      cardTitle: newTitle || props.widget.cardTitle,
      cardValue: newValue || props.widget.cardValue
    })
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
  font-size: 12px;
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
  font-size: 14px;
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
}

.edit-dialog :deep(.p-dialog-header) {
  border-bottom: 1px solid #212121;
  padding: 1.25rem 1.5rem;
}

.edit-dialog :deep(.p-dialog-content) {
  padding: 1.5rem;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.edit-field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.edit-field label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #a0a0a0;
}

.form-input {
  width: 100%;
}
</style>
