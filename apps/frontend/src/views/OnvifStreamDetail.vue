<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Button from 'primevue/button'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import Chip from 'primevue/chip'
import { useToast } from 'primevue/usetoast'

import {
  getCamera, getCameraProfiles, setStreamProfile,
  startStream, stopStream, updateCamera, deleteCamera
} from '@/api/media'

import StreamFormDialog from '@/components/cameras/StreamFormDialog.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ProfileCard from '@/components/cameras/ProfileCard.vue'
import StreamPlayer from '@/components/cameras/StreamPlayer.vue'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const loading = ref(false)

const showError = (summary, detail) => {
  toast.add({ severity: 'error', summary, detail, life: 3000 })
}

const stream = ref(null)
const profiles = ref([])
const streamState = ref(null)
const profilesLoading = ref(false)

async function fetchStream() {
  try {
    stream.value = await getCamera(route.params.id)
  } catch (error) {
    console.error(error)
    showError('Error', 'Failed to load stream details.')
  }
}

async function fetchProfiles() {
  profilesLoading.value = true
  try {
    const result = await getCameraProfiles(route.params.id)
    profiles.value = result.profiles || []
  } catch (error) {
    console.error(error)
    showError('ONVIF Error', 'Failed to fetch stream profiles. Device may be offline.')
  } finally {
    profilesLoading.value = false
  }
}

onMounted(() => {
  fetchStream()
  fetchProfiles()
})

watch(() => route.params.id, () => {
  fetchStream()
  fetchProfiles()
})

async function handleSelectProfile(token) {
  try {
    await setStreamProfile(route.params.id, token)
    stream.value.stream_profile_token = token
    toast.add({ severity: 'success', summary: 'Profile Set', detail: `Stream profile set to ${token}.`, life: 3000 })
  } catch (error) {
    console.error(error)
    showError('Error', 'Failed to set stream profile.')
  }
}

async function handleStartStream() {
  loading.value = true
  try {
    const result = await startStream(route.params.id)
    streamState.value = result
    toast.add({ severity: 'success', summary: 'Stream Started', detail: 'HLS stream is starting...', life: 3000 })
  } catch (error) {
    console.error(error)
    showError('Stream Error', 'Failed to start stream. Ensure a profile is selected and device is reachable.')
  } finally {
    loading.value = false
  }
}

async function handleStopStream() {
  loading.value = true
  try {
    await stopStream(route.params.id)
    streamState.value = null
    toast.add({ severity: 'success', summary: 'Stopped', detail: 'Stream stopped.', life: 3000 })
  } catch (error) {
    console.error(error)
    showError('Stream Error', 'Failed to stop stream.')
  } finally {
    loading.value = false
  }
}

const formVisible = ref(false)
const confirmVisible = ref(false)

function openEditDialog() {
  formVisible.value = true
}

async function handleSaveStream(payload) {
  loading.value = true
  try {
    await updateCamera(route.params.id, payload)
    toast.add({ severity: 'success', summary: 'Updated', detail: 'Stream updated.', life: 3000 })
    formVisible.value = false
    await fetchStream()
  } catch (error) {
    console.error(error)
    showError('Error', 'Failed to update stream.')
  } finally {
    loading.value = false
  }
}

function handleDeleteRequest() {
  confirmVisible.value = true
}

async function handleDeleteConfirm() {
  loading.value = true
  try {
    await deleteCamera(route.params.id)
    toast.add({ severity: 'success', summary: 'Deleted', detail: 'Stream deleted.', life: 3000 })
    router.push('/integrations/onvif-streams')
  } catch (error) {
    console.error(error)
    showError('Error', 'Failed to delete stream.')
  } finally {
    loading.value = false
  }
}

function hlsUrl() {
  if (!streamState.value) return ''
  return streamState.value.hls_url || `/hls/${route.params.id}/index.m3u8`
}
</script>

<template>
  <BlockUI :blocked="loading" fullScreen />
  <ProgressSpinner v-if="loading" class="global-spinner" />

  <div v-if="!stream" class="empty-state">
    <ProgressSpinner v-if="loading" />
    <p v-else>Stream not found.</p>
  </div>

  <template v-else>
    <div class="detail-header">
      <Button
        icon="pi pi-arrow-left"
        label="Back"
        severity="secondary"
        text
        @click="router.push('/integrations/onvif-streams')"
      />
      <div class="header-info">
        <h2 class="stream-name">{{ stream.name }}</h2>
        <div class="stream-meta">
          <span>{{ stream.ip }}:{{ stream.port || 80 }}</span>
          <Chip
            v-if="streamState"
            label="Live"
            icon="pi pi-circle-fill"
            class="chip-live"
          />
        </div>
      </div>
      <div class="header-actions">
        <Button icon="pi pi-pencil" label="Edit" severity="secondary" text @click="openEditDialog" />
      </div>
    </div>

    <div v-if="stream.notes" class="stream-notes">
      {{ stream.notes }}
    </div>

    <div class="section">
      <div class="section-header">
        <h3>ONVIF Profiles</h3>
        <Button
          icon="pi pi-refresh"
          label="Refresh"
          size="small"
          severity="secondary"
          text
          :loading="profilesLoading"
          @click="fetchProfiles"
        />
      </div>
      <div v-if="!profiles.length && !profilesLoading" class="profiles-empty">
        <i class="pi pi-info-circle"></i>
        <span>No profiles loaded. Click Refresh to fetch profiles from the device.</span>
      </div>
      <div v-else class="profiles-grid">
        <ProfileCard
          v-for="p in profiles"
          :key="p.token"
          :profile="p"
          :selected="stream.stream_profile_token === p.token"
          @select="handleSelectProfile"
        />
      </div>
    </div>

    <div class="section">
      <div class="stream-controls">
        <Button
          v-if="!streamState"
          icon="pi pi-play"
          label="Start Stream"
          severity="success"
          :disabled="!stream.stream_profile_token"
          @click="handleStartStream"
        />
        <Button
          v-else
          icon="pi pi-stop-circle"
          label="Stop Stream"
          severity="danger"
          outlined
          @click="handleStopStream"
        />
        <Chip
          v-if="streamState"
          :label="streamState.status"
          class="stream-status-chip"
        />
      </div>

      <div v-if="streamState" class="stream-info">
        <div class="info-row">
          <span class="info-label">Profile:</span>
          <code>{{ streamState.profile_token }}</code>
        </div>
        <div class="info-row">
          <span class="info-label">RTSP:</span>
          <code class="rtsp-url">{{ streamState.rtsp_url }}</code>
        </div>
      </div>

      <StreamPlayer
        :hls-url="hlsUrl()"
        :stream-status="streamState?.status"
      />
    </div>

    <StreamFormDialog
      v-model:visible="formVisible"
      :stream="stream"
      @save="handleSaveStream"
      @request-delete="handleDeleteRequest"
    />

    <ConfirmDialog
      v-model:visible="confirmVisible"
      title="Delete Stream"
      message="Are you sure you want to delete this stream? This will also stop any active stream."
      confirm-label="Delete"
      severity="danger"
      @confirm="handleDeleteConfirm"
    />
  </template>
</template>

<style scoped>
.detail-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.header-info {
  flex: 1;
}

.stream-name {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.stream-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #888;
  font-size: 0.85rem;
  margin-top: 0.25rem;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.stream-notes {
  background: #202024;
  border-left: 3px solid #48897b;
  padding: 0.75rem 1rem;
  border-radius: 0 8px 8px 0;
  color: #aaa;
  font-size: 0.85rem;
  margin-bottom: 1.5rem;
}

.section {
  margin-bottom: 1.5rem;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.section-header h3 {
  margin: 0;
  font-size: 1rem;
  color: #ccc;
}

.profiles-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.profiles-empty {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #666;
  font-size: 0.85rem;
  padding: 1.5rem;
  background: #18181B;
  border-radius: 8px;
}

.stream-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.stream-info {
  background: #18181B;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8rem;
}

.info-label {
  color: #888;
  min-width: 4rem;
}

.info-row code {
  color: #e0e0e0;
  font-size: 0.8rem;
}

.rtsp-url {
  color: #64b5f6;
  background: #0e0e10;
  padding: 1px 6px;
  border-radius: 4px;
  max-width: 45ch;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}

.chip-live {
  font-size: 0.7rem;
  background-color: #2d5a4e;
}

.chip-live :deep(.p-chip-icon) {
  color: #4cff88;
}

.stream-status-chip {
  font-size: 0.8rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem;
  color: #555;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}
</style>
