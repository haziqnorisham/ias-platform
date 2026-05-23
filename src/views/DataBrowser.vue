<template>
  <BlockUI :blocked="loading" fullScreen />
  <ProgressSpinner v-if="loading" class="global-spinner" />

  <div class="main-container">
    <h1 class="page-title">Data Browser</h1>

    <div class="selection-panel">
      <label class="field-label">Select Device</label>
      <Select
        v-model="selectedDevice"
        :options="deviceOptions"
        optionLabel="label"
        optionValue="value"
        placeholder="Choose a device..."
        class="device-select"
      />
    </div>

    <div v-if="selectedDevice" class="device-info-panel">
      <div class="device-info-body">
        <div class="device-info-left">
          <h2 class="section-title">Device Information</h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">Device ID</span>
              <span class="info-value">{{ selectedDevice.Id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Name</span>
              <span class="info-value">{{ selectedDevice.Name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Description</span>
              <span class="info-value">{{ selectedDevice.Description || '—' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Profile ID</span>
              <span class="info-value">{{ selectedDevice.ProfileID || '—' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Status</span>
              <Tag
                :value="selectedDevice.Status"
                :severity="selectedDevice.Status === 'active' ? 'success' : 'danger'"
              />
            </div>
            <div class="info-item">
              <span class="info-label">Location</span>
              <span class="info-value">{{ selectedDevice.LocationLabel || '—' }}</span>
            </div>
          </div>
        </div>
        <div class="device-info-right">
          <img :src="imageSrc" class="device-image" @error="onImageError" />
        </div>
      </div>
    </div>

    <div v-if="selectedDevice" class="data-section">
      <h2 class="section-title">Device Data</h2>
      <ProgressSpinner v-if="dataLoading" class="table-spinner" />

      <div v-else-if="!deviceData.length" class="empty-table">
        No data available for this device.
      </div>

      <DataTable v-else :value="deviceData" size="small" scrollable scrollHeight="500px"
        paginator :rows="20" :rowsPerPageOptions="[10, 20, 50]"
        tableStyle="min-width: 50rem" class="tight-table">
        <Column field="message_id" header="Msg ID" />
        <Column field="topic" header="Topic">
          <template #body="{ data }">
            <span class="topic-cell">{{ data.topic }}</span>
          </template>
        </Column>
        <Column field="ingest_method" header="Ingest Method" />
        <Column field="status" header="Status">
          <template #body="{ data }">
            <Tag
              :value="data.status"
              :severity="data.status === 'processed' ? 'success' : 'danger'"
            />
          </template>
        </Column>
        <Column field="success" header="Success">
          <template #body="{ data }">
            <Tag
              :value="data.success ? 'Yes' : 'No'"
              :severity="data.success ? 'success' : 'danger'"
            />
          </template>
        </Column>
        <Column field="processed_payload" header="Processed Payload">
          <template #body="{ data }">
            <span class="payload-cell">{{ data.processed_payload || '—' }}</span>
          </template>
        </Column>
        <Column field="received_at" header="Received At" />
      </DataTable>
    </div>

    <div v-else-if="devices.length" class="empty-selection">
      Select a device from the dropdown above to browse its data.
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useToast } from 'primevue/usetoast'
import { getAllDevices, getDeviceSuccessfulIngest } from '@/api/posts'

const toast = useToast()
const loading = ref(false)
const devices = ref([])
const selectedDevice = ref(null)

const deviceData = ref([])
const dataLoading = ref(false)
const fallbackActive = ref(false)

const imageSrc = computed(() => {
  if (fallbackActive.value || !selectedDevice.value?.ProfileID) return '/MCS.png'
  return `/api/image/device_profile_${selectedDevice.value.ProfileID}.png`
})

function onImageError() {
  fallbackActive.value = true
}

const deviceOptions = computed(() =>
  devices.value.map(d => ({
    label: `${d.Name} (${d.Id})`,
    value: d
  }))
)

async function fetchDevices() {
  loading.value = true
  try {
    const result = await getAllDevices()
    if (result) {
      devices.value = result
    }
  } catch (error) {
    console.error(error)
    toast.add({ severity: 'error', summary: 'API Call Failed', detail: 'Failed to fetch devices.', life: 3000 })
  } finally {
    loading.value = false
  }
}

async function fetchDeviceData() {
  if (!selectedDevice.value) {
    deviceData.value = []
    return
  }
  dataLoading.value = true
  deviceData.value = []
  try {
    const result = await getDeviceSuccessfulIngest({ device_id: selectedDevice.value.Id })
    if (result) {
      deviceData.value = result
    }
  } catch (error) {
    console.error(error)
    toast.add({ severity: 'error', summary: 'API Call Failed', detail: 'Failed to fetch device data.', life: 3000 })
  } finally {
    dataLoading.value = false
  }
}

watch(selectedDevice, () => {
  fetchDeviceData()
})

onMounted(() => {
  fetchDevices()
})
</script>

<style scoped>
.main-container {
  display: flex;
  flex-direction: column;
  padding: 1rem 0;
  gap: 1.5rem;
}

.page-title {
  font-family: "Space Grotesk", sans-serif;
  font-size: 1.5rem;
  font-weight: 700;
  color: #e0e0e0;
  margin: 0;
}

.selection-panel {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field-label {
  font-family: "Space Grotesk", sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  color: #a0a0a0;
}

.device-select {
  width: 100%;
  max-width: 500px;
}

.section-title {
  font-family: "Space Grotesk", sans-serif;
  font-size: 1.1rem;
  font-weight: 600;
  color: #e0e0e0;
  margin: 0 0 0.75rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #2a2a2e;
}

.device-info-panel {
  background: #1a1a1a;
  border: 1px solid #2a2a2e;
  border-radius: 8px;
  padding: 1.25rem;
}

.device-info-body {
  display: flex;
  gap: 1.5rem;
  align-items: flex-start;
}

.device-info-left {
  flex: 1;
  min-width: 0;
}

.device-info-right {
  flex-shrink: 0;
}

.device-image {
  height: 10rem;
  border-radius: 4px;
  object-fit: contain;
}

.info-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-value {
  font-size: 0.9rem;
  color: #e0e0e0;
}

.data-section {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.tight-table :deep(.p-datatable-thead > tr > th) {
  font-size: 0.75rem;
  padding: 0.4rem 0.6rem;
}

.tight-table :deep(.p-datatable-tbody > tr > td) {
  font-size: 0.78rem;
  padding: 0.3rem 0.6rem;
}

.empty-table,
.empty-selection {
  color: #666;
  font-size: 0.9rem;
  text-align: center;
  padding: 2rem 0;
}

.topic-cell {
  display: block;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.payload-cell {
  display: block;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
  font-size: 0.85rem;
}

.table-spinner {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}
</style>
