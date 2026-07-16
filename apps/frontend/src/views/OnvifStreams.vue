<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

import Button from 'primevue/button'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Chip from 'primevue/chip'
import { useToast } from 'primevue/usetoast'

import {
  getCameras, createCamera, updateCamera, deleteCamera,
  getStreams, stopStream
} from '@/api/media'

import StreamCard from '@/components/cameras/StreamCard.vue'
import StreamFormDialog from '@/components/cameras/StreamFormDialog.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const toast = useToast()
const router = useRouter()
const loading = ref(false)

const showError = (summary, detail) => {
  toast.add({ severity: 'error', summary, detail, life: 3000 })
}

const streams = ref([])
const activeStreams = ref([])

const streamMap = computed(() => {
  const map = {}
  activeStreams.value.forEach(s => { map[s.device_id] = s })
  return map
})

async function fetchStreams() {
  try {
    const result = await getCameras()
    streams.value = Array.isArray(result) ? result : []
  } catch (error) {
    console.error(error)
    showError('API Call Failed', 'Failed to fetch streams.')
  }
}

async function fetchActiveStreams() {
  try {
    const result = await getStreams()
    activeStreams.value = Array.isArray(result) ? result : []
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  fetchStreams()
  fetchActiveStreams()
})

const formVisible = ref(false)
const editingStream = ref(null)

function openAddDialog() {
  editingStream.value = null
  formVisible.value = true
}

function openEditDialog(stream) {
  editingStream.value = stream
  formVisible.value = true
}

async function handleSave(payload) {
  loading.value = true
  try {
    if (editingStream.value) {
      await updateCamera(editingStream.value.id, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Stream updated.', life: 3000 })
    } else {
      await createCamera(payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Stream created.', life: 3000 })
    }
    formVisible.value = false
    await fetchStreams()
  } catch (error) {
    console.error(error)
    showError('API Call Failed', 'Failed to save stream.')
  } finally {
    loading.value = false
  }
}

const confirmVisible = ref(false)

function handleDeleteRequest() {
  confirmVisible.value = true
}

async function handleDeleteConfirm() {
  loading.value = true
  try {
    await deleteCamera(editingStream.value.id)
    toast.add({ severity: 'success', summary: 'Success', detail: 'Stream deleted.', life: 3000 })
    formVisible.value = false
    confirmVisible.value = false
    await fetchStreams()
  } catch (error) {
    console.error(error)
    showError('API Call Failed', 'Failed to delete stream.')
  } finally {
    loading.value = false
  }
}

function viewStream(stream) {
  router.push(`/integrations/onvif-streams/${stream.id}`)
}

async function handleStopStream(deviceId) {
  try {
    await stopStream(deviceId)
    toast.add({ severity: 'success', summary: 'Stopped', detail: 'Stream stopped.', life: 3000 })
    await fetchActiveStreams()
  } catch (error) {
    console.error(error)
    showError('Stream Error', 'Failed to stop stream.')
  }
}

function streamStatusSeverity(status) {
  switch (status) {
    case 'running': return 'success'
    case 'starting': return 'warn'
    case 'error': return 'danger'
    default: return 'secondary'
  }
}

const activeTab = ref(0)
</script>

<template>
  <BlockUI :blocked="loading" fullScreen />
  <ProgressSpinner v-if="loading" class="global-spinner" />

  <div class="toolbar">
    <Button label="Add Stream" icon="pi pi-plus" @click="openAddDialog" />
  </div>

  <div class="tabs">
    <button
      class="tab-btn"
      :class="{ active: activeTab === 0 }"
      @click="activeTab = 0"
    >
      <i class="pi pi-video"></i> Streams
      <span class="tab-count">{{ streams.length }}</span>
    </button>
    <button
      class="tab-btn"
      :class="{ active: activeTab === 1 }"
      @click="activeTab = 1"
    >
      <i class="pi pi-play-circle"></i> Active
      <span class="tab-count">{{ activeStreams.length }}</span>
    </button>
  </div>

  <div v-if="activeTab === 0">
    <div v-if="!streams.length" class="empty-state">
      <i class="pi pi-video empty-icon"></i>
      <p>No streams configured. Add your first ONVIF stream to get started.</p>
    </div>
    <div v-else class="streams-grid">
      <StreamCard
        v-for="stream in streams"
        :key="stream.id"
        :stream="stream"
        :stream-status="streamMap[stream.id]"
        @edit="openEditDialog"
        @view="viewStream"
      />
    </div>
  </div>

  <div v-if="activeTab === 1">
    <div v-if="!activeStreams.length" class="empty-state">
      <i class="pi pi-play-circle empty-icon"></i>
      <p>No active streams. Start a stream from the stream detail page.</p>
    </div>
    <DataTable v-else :value="activeStreams" stripedRows size="small" tableStyle="min-width: 40rem">
      <Column field="device_name" header="Stream" />
      <Column field="profile_token" header="Profile" />
      <Column header="Status">
        <template #body="{ data }">
          <Chip :label="data.status" :severity="streamStatusSeverity(data.status)" />
        </template>
      </Column>
      <Column field="rtsp_url" header="RTSP URL">
        <template #body="{ data }">
          <code class="rtsp-url">{{ data.rtsp_url }}</code>
        </template>
      </Column>
      <Column header="Actions">
        <template #body="{ data }">
          <Button
            icon="pi pi-stop-circle"
            label="Stop"
            size="small"
            severity="danger"
            outlined
            @click="handleStopStream(data.device_id)"
          />
        </template>
      </Column>
    </DataTable>
  </div>

  <StreamFormDialog
    v-model:visible="formVisible"
    :stream="editingStream"
    @save="handleSave"
    @request-delete="handleDeleteRequest"
  />

  <ConfirmDialog
    v-model:visible="confirmVisible"
    title="Delete Stream"
    message="Are you sure you want to delete this stream? This action cannot be undone."
    confirm-label="Delete"
    severity="danger"
    @confirm="handleDeleteConfirm"
  />
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.tabs {
  display: flex;
  gap: 0;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #2a2a2e;
}

.tab-btn {
  background: none;
  border: none;
  color: #888;
  padding: 0.6rem 1.25rem;
  font-size: 0.9rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-bottom: 2px solid transparent;
  transition: all 0.15s;
  font-family: inherit;
}

.tab-btn:hover {
  color: #ccc;
}

.tab-btn.active {
  color: #48897b;
  border-bottom-color: #48897b;
}

.tab-count {
  background: #2a2a2e;
  color: #888;
  font-size: 0.7rem;
  padding: 1px 6px;
  border-radius: 8px;
}

.tab-btn.active .tab-count {
  background: #2d5a4e;
  color: #4cff88;
}

.streams-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
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

.rtsp-url {
  font-size: 0.75rem;
  color: #64b5f6;
  background: #18181B;
  padding: 1px 6px;
  border-radius: 4px;
  max-width: 30ch;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}
</style>
