<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  widgetKey: { type: String, required: true },
  config: { type: Object, default: () => ({}) },
})

const host = ref(null)
let instance = null

function mount() {
  const factory = window.__ias_getWidget?.(props.widgetKey)
  if (!factory || !host.value) return
  host.value.innerHTML = ''
  instance = factory.create(host.value, props.config)
  if (instance?.mount) instance.mount()
}

onMounted(() => mount())

onUnmounted(() => {
  if (instance?.unmount) instance.unmount()
  instance = null
})
</script>

<template>
  <div ref="host" class="ext-widget-slot"></div>
</template>

<style scoped>
.ext-widget-slot {
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding: 4px 8px;
}
</style>
