<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import ProgressSpinner from 'primevue/progressspinner'
import { loadExtensionComponent } from '@/api/extensions'

const route = useRoute()
const loading = ref(true)
const error = ref(null)
const tagName = ref(null)

const extensionName = computed(() => route.params.name)

onMounted(async () => {
  try {
    await loadExtensionComponent(extensionName.value)
    tagName.value = `extension-${extensionName.value}`
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="extension-detail">
    <div v-if="loading" class="detail-state">
      <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="4" />
      <span>Loading extension...</span>
    </div>

    <div v-else-if="error" class="detail-state detail-error">
      <i class="pi pi-exclamation-triangle" style="font-size: 2rem; color: #e57373"></i>
      <span>{{ error }}</span>
    </div>

    <component v-else :is="tagName" />
  </div>
</template>

<style scoped>
.extension-detail {
  min-height: 200px;
}

.detail-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 2rem;
  color: #888;
}

.detail-error {
  color: #e57373;
}
</style>
