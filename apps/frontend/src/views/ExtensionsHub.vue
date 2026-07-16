<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Panel from '@/components/Panel.vue'
import Button from 'primevue/button'
import ProgressSpinner from 'primevue/progressspinner'
import { getExtensions } from '@/api/extensions'

const router = useRouter()
const extensions = ref([])
const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  try {
    const result = await getExtensions()
    extensions.value = Array.isArray(result) ? result : []
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="hub-state">
    <ProgressSpinner style="width: 40px; height: 40px" strokeWidth="4" />
    <span>Loading extensions...</span>
  </div>

  <div v-else-if="error" class="hub-state hub-error">
    <i class="pi pi-exclamation-triangle" style="font-size: 2rem; color: #e57373"></i>
    <span>Failed to load extensions: {{ error }}</span>
  </div>

  <div v-else-if="!extensions.length" class="hub-state">
    <i class="pi pi-box" style="font-size: 2rem; color: #555"></i>
    <span>No extensions available.</span>
  </div>

  <div v-else class="main-container">
    <Panel
      v-for="ext in extensions"
      :key="ext.name"
      :title="ext.name"
      subtitle="Extension"
    >
      <div class="centered">
        <Button label="Configure" icon="pi pi-cog" @click="router.push(`/extensions/${ext.name}`)" />
      </div>
    </Panel>
  </div>
</template>

<style scoped>
.main-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 10px;
}

.hub-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 2rem;
  color: #888;
}

.hub-error {
  color: #e57373;
}
</style>
