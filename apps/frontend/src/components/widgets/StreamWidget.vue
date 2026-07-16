<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import ProgressSpinner from 'primevue/progressspinner'
import { getCameras, getStreams } from '@/api/media'

const router = useRouter()
const videoRef = ref(null)
const loading = ref(true)

const activeStream = ref(null)
const devices = ref([])

function buildHlsUrl() {
  if (!activeStream.value) return ''
  if (activeStream.value.hls_url) {
    return activeStream.value.hls_url.startsWith('http')
      ? activeStream.value.hls_url
      : `/media${activeStream.value.hls_url}`
  }
  return `/media/hls/${activeStream.value.device_id}/index.m3u8`
}

async function fetchData() {
  loading.value = !activeStream.value
  try {
    const [activeData, deviceData] = await Promise.all([
      getStreams().catch(() => []),
      getCameras().catch(() => [])
    ])

    const running = Array.isArray(activeData)
      ? activeData.find(s => s.status === 'running')
      : null

    const newDeviceId = running ? running.device_id : null
    const currentDeviceId = activeStream.value ? activeStream.value.device_id : null

    if (newDeviceId !== currentDeviceId) {
      activeStream.value = running || null
    }

    devices.value = Array.isArray(deviceData) ? deviceData : []
  } catch {
    activeStream.value = null
  } finally {
    loading.value = false
  }
}

function playVideo() {
  const video = videoRef.value
  if (!video) return
  const url = buildHlsUrl()
  if (!url) return
  video.src = url
  video.load()
  video.play().catch(() => {})
}

watch(activeStream, (val) => {
  if (val) {
    playVideo()
  }
})

let pollInterval = null

onMounted(() => {
  fetchData()
  pollInterval = setInterval(fetchData, 30000)
})

onBeforeUnmount(() => {
  clearInterval(pollInterval)
})
</script>

<template>
  <div class="stream-widget">
    <div v-if="loading && !activeStream" class="widget-state">
      <ProgressSpinner style="width: 32px; height: 32px" strokeWidth="4" />
      <span>Checking for streams...</span>
    </div>

    <template v-else-if="activeStream">
      <video
        ref="videoRef"
        class="stream-video"
        controls
        muted
        playsinline
        autoplay
      />
      <div class="stream-label">{{ activeStream.device_name }}</div>
    </template>

    <div v-else-if="devices.length" class="widget-state">
      <i class="pi pi-play-circle state-icon muted"></i>
      <span>No active streams</span>
      <span class="state-hint">{{ devices.length }} device{{ devices.length !== 1 ? 's' : '' }} configured, none streaming</span>
    </div>

    <div v-else class="widget-state">
      <i class="pi pi-video state-icon"></i>
      <span>No ONVIF streams configured</span>
      <span class="state-link" @click="router.push('/integrations')">
        Configure in Integrations <i class="pi pi-arrow-right"></i>
      </span>
    </div>
  </div>
</template>

<style scoped>
.stream-widget {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0e0e10;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.stream-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

.stream-label {
  position: absolute;
  bottom: 8px;
  left: 12px;
  background: rgba(0, 0, 0, 0.6);
  color: #ccc;
  font-size: 0.7rem;
  padding: 2px 8px;
  border-radius: 4px;
  pointer-events: none;
}

.widget-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  color: #555;
  font-size: 0.9rem;
  text-align: center;
  padding: 1rem;
}

.state-icon {
  font-size: 2.5rem;
  color: #444;
}

.state-icon.muted {
  color: #3a3a3e;
}

.state-hint {
  font-size: 0.75rem;
  color: #444;
}

.state-link {
  color: #64b5f6;
  font-size: 0.8rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.35rem;
  transition: color 0.15s;
}

.state-link:hover {
  color: #90caf9;
}
</style>
