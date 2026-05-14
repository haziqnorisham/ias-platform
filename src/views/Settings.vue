<template>
  <BlockUI :blocked="loading" fullScreen />
  <ProgressSpinner v-if="loading" class="global-spinner" />

  <div class="page-container">
    <h2 class="page-title">Server Configuration</h2>

    <DataTable v-if="configEntries.length" :value="configEntries" size="small" scrollable scrollHeight="flex" tableStyle="min-width: 40rem" class="config-table">
      <Column field="key" header="Key">
        <template #body="{ data }">
          <span class="config-key" :class="{ 'sensitive-key': data.sensitive }">{{ data.key }}</span>
        </template>
      </Column>
      <Column field="value" header="Value">
        <template #body="{ data }">
          <code class="config-value" :class="{ 'sensitive-value': data.sensitive }">{{ data.value }}</code>
        </template>
      </Column>
    </DataTable>

    <div v-else-if="!loading" class="placeholder-card">
      <i class="pi pi-cog placeholder-icon"></i>
      <p class="placeholder-text">No configuration data available.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { getServerConfig } from '@/api/posts'

const loading = ref(false)
const configEntries = ref([])

onMounted(async () => {
  loading.value = true
  try {
    const config = await getServerConfig()
    configEntries.value = Object.entries(config || {}).map(([key, value]) => ({
      key,
      value,
      sensitive: value === '***'
    }))
  } catch (error) {
    console.error('Failed to fetch server config:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-container {
  padding: 0;
}

.page-title {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.config-table :deep(.p-datatable-tbody > tr > td) {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
}

.config-table :deep(.p-datatable-thead > tr > th) {
  padding: 0.4rem 0.75rem;
  font-size: 0.75rem;
}

.config-key {
  color: #a0a0a0;
  font-family: monospace;
  font-size: 0.82rem;
}

.config-value {
  color: #e0e0e0;
  font-family: monospace;
  font-size: 0.82rem;
}

.sensitive-key {
  color: #d4a844;
}

.sensitive-value {
  color: #888;
  font-style: italic;
}

.placeholder-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 3rem;
  min-height: 200px;
}

.placeholder-icon {
  font-size: 2.5rem;
  color: #666;
  margin-bottom: 1rem;
}

.placeholder-text {
  color: #666;
  font-size: 1rem;
  margin: 0;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}
</style>
