<template>
  <BlockUI :blocked="loading" fullScreen />
  <ProgressSpinner v-if="loading" class="global-spinner" />

  <div class="page-container">
    <h2 class="page-title">Settings</h2>

    <Tabs value="server-details">
      <TabList>
        <Tab value="server-details">Server Details</Tab>
        <Tab value="user-permissions">User Permissions</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="server-details">
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
        </TabPanel>

        <TabPanel value="user-permissions">
          <div class="permissions-section">
            <p class="permissions-hint">
              Select the pages that <strong>viewer</strong> users can see on the sidebar. Admins always see everything.
            </p>
            <div v-for="page in permissionPages" :key="page.to" class="permission-row">
              <span class="permission-label">{{ page.label }}</span>
              <ToggleSwitch
                :modelValue="viewerPermissions[page.to]"
                :disabled="!isAdmin"
                @update:modelValue="val => setViewerPermission(page.to, val)"
              />
            </div>
            <p v-if="!isAdmin" class="permissions-note">
              <i class="pi pi-lock"></i> Only admins can change these settings.
            </p>
          </div>
        </TabPanel>
      </TabPanels>
    </Tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import ToggleSwitch from 'primevue/toggleswitch'
import { getServerConfig } from '@/api/posts'
import { useAuth } from '@/composables/useAuth'
import { usePermissions, permissionPages } from '@/composables/usePermissions'

const loading = ref(false)
const configEntries = ref([])
const { isAdmin } = useAuth()
const { viewerPermissions, setViewerPermission } = usePermissions()

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
  font-size: var(--font-size-xl);
  font-weight: 600;
}

.config-table :deep(.p-datatable-tbody > tr > td) {
  padding: 0.35rem 0.75rem;
  font-size: var(--font-size-sm);
}

.config-table :deep(.p-datatable-thead > tr > th) {
  padding: 0.4rem 0.75rem;
  font-size: var(--font-size-xs);
}

.config-key {
  color: #a0a0a0;
  font-family: monospace;
  font-size: var(--font-size-sm);
}

.config-value {
  color: #e0e0e0;
  font-family: monospace;
  font-size: var(--font-size-sm);
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
  font-size: var(--font-size-2xl);
  color: #666;
  margin-bottom: 1rem;
}

.placeholder-text {
  color: #666;
  font-size: var(--font-size-md);
  margin: 0;
}

.permissions-section {
  background-color: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 1.5rem;
  max-width: 480px;
}

.permissions-hint {
  color: #a0a0a0;
  font-size: var(--font-size-sm);
  margin: 0 0 1rem 0;
}

.permission-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #212121;
}

.permission-row:last-of-type {
  border-bottom: none;
}

.permission-label {
  color: #e0e0e0;
  font-size: var(--font-size-sm);
}

.permissions-note {
  color: #888;
  font-size: var(--font-size-xs);
  margin: 1rem 0 0 0;
}

.global-spinner {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 9999;
}
</style>
