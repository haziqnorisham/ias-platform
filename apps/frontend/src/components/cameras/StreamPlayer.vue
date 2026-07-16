<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  hlsUrl: {
    type: String,
    default: ''
  },
  streamStatus: {
    type: String,
    default: null
  }
})

const videoRef = ref(null)

function buildFullUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `/media${path}`
}

watch(() => props.hlsUrl, (newUrl) => {
  const video = videoRef.value
  if (!video) return

  if (!newUrl) {
    video.removeAttribute('src')
    video.load()
    return
  }

  const url = buildFullUrl(newUrl)
  if (video.src !== url) {
    video.src = url
    video.load()
    video.play().catch(() => {})
  }
}, { immediate: true })
</script>

<template>
  <div class="player-container">
    <div v-if="!hlsUrl" class="player-placeholder">
      <i class="pi pi-video"></i>
      <span>Stream not started</span>
    </div>
    <video
      ref="videoRef"
      class="video-element"
      :class="{ hidden: !hlsUrl }"
      controls
      muted
      playsinline
      autoplay
    />
  </div>
</template>

<style scoped>
.player-container {
  background: #0e0e10;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 360px;
  position: relative;
}

.video-element {
  width: 100%;
  max-height: 70vh;
  display: block;
}

.video-element.hidden {
  display: none;
}

.player-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  color: #555;
  font-size: 1rem;
}

.player-placeholder i {
  font-size: 3rem;
}
</style>
