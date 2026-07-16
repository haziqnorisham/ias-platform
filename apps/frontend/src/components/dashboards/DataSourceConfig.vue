<script setup>
import { ref, watch } from 'vue'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import ToggleSwitch from 'primevue/toggleswitch'

const props = defineProps({
  modelValue: {
    type: Object,
    default: null
  },
  widgetType: {
    type: String,
    required: true
  },
  devices: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const enabled = ref(false)
const deviceId = ref('')
const columnName = ref('')
const dataType = ref('number')
const yAxis = ref('')

const dataTypeOptions = [
  { label: 'Number', value: 'number' },
  { label: 'String', value: 'string' },
  { label: 'Boolean', value: 'boolean' }
]

function loadFromValue(val) {
  if (val && val.deviceID) {
    enabled.value = true
    deviceId.value = val.deviceID || ''
    if (props.widgetType === 'barchart' || props.widgetType === 'linechart') {
      yAxis.value = val.y_axis || ''
    } else {
      columnName.value = val.column_name || ''
      dataType.value = val.data_type || 'number'
    }
  } else {
    enabled.value = false
    deviceId.value = ''
    columnName.value = ''
    dataType.value = 'number'
    yAxis.value = ''
  }
}

watch(() => props.modelValue, loadFromValue, { immediate: true })

function emitUpdate() {
  if (!enabled.value || !deviceId.value) {
    emit('update:modelValue', null)
    return
  }

  if (props.widgetType === 'barchart' || props.widgetType === 'linechart') {
    emit('update:modelValue', {
      deviceID: deviceId.value,
      y_axis: yAxis.value.trim()
    })
  } else {
    emit('update:modelValue', {
      deviceID: deviceId.value,
      column_name: columnName.value.trim(),
      data_type: dataType.value
    })
  }
}

function onEnabledChange() {
  if (!enabled.value) {
    emit('update:modelValue', null)
    return
  }
  emitUpdate()
}
</script>

<template>
  <div class="ds-section">
    <div class="ds-section-label">Data Source</div>

    <div class="toggle-row">
      <ToggleSwitch v-model="enabled" @change="onEnabledChange" />
      <span class="toggle-label">Dynamic Data</span>
    </div>

    <div v-if="enabled" class="query-grid">
      <div class="edit-field query-field-device">
        <label>Device</label>
        <Select
          v-model="deviceId"
          :options="devices"
          optionLabel="label"
          optionValue="value"
          placeholder="Select device"
          size="small"
          class="form-input"
          @change="emitUpdate"
        />
      </div>

      <template v-if="widgetType === 'card'">
        <div class="edit-field query-field-column">
          <label>Column name</label>
          <InputText v-model="columnName" placeholder="object.temperature" size="small" class="form-input" @blur="emitUpdate" @keydown.enter="$event.target.blur()" />
        </div>
        <div class="edit-field query-field-type">
          <label>Data type</label>
          <Select
            v-model="dataType"
            :options="dataTypeOptions"
            optionLabel="label"
            optionValue="value"
            size="small"
            class="form-input"
            @change="emitUpdate"
          />
        </div>
      </template>

      <template v-if="widgetType === 'barchart' || widgetType === 'linechart'">
        <div class="edit-field query-field-column">
          <label>Column Name</label>
          <InputText v-model="yAxis" placeholder="object.temperature" size="small" class="form-input" @blur="emitUpdate" @keydown.enter="$event.target.blur()" />
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.ds-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.ds-section-label {
  font-size: var(--font-size-xs);
  font-weight: 700;
  color: #777;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.toggle-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toggle-label {
  font-size: var(--font-size-sm);
  color: #ccc;
}

.query-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
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

.query-field-device {
  grid-column: 1 / -1;
}
</style>
