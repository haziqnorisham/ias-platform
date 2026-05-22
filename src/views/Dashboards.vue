<template>
  <div>
    <div class="toolbar">
      <Button label="New Dashboard" icon="pi pi-plus" @click="router.push('/dashboards/edit')" />
    </div>

    <div v-if="!loading && !dashboards.length" class="empty-state">
      <i class="pi pi-chart-bar empty-icon"></i>
      <p>No dashboards yet. Create your first dashboard to get started.</p>
    </div>

    <div v-else class="dashboards-grid">
      <div v-for="d in dashboards" :key="d.id" class="dashboard-card">
        <div class="card-header">
          <i class="pi pi-palette card-icon"></i>
          <span class="card-name">{{ d.name || 'Untitled Dashboard' }}</span>
        </div>
        <div class="card-meta">
          <span class="meta-label">Updated</span>
          <span class="meta-value">{{ d.updated_at || '—' }}</span>
        </div>
        <div class="card-actions">
          <Button label="View" icon="pi pi-eye" size="small" text @click="router.push(`/dashboards/view?id=${d.id}`)" />
          <Button label="Edit" icon="pi pi-pencil" size="small" @click="router.push(`/dashboards/edit?id=${d.id}`)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Button from 'primevue/button'
import { getDashboards } from '@/api/posts'

const router = useRouter()
const loading = ref(false)
const dashboards = ref([])

onMounted(async () => {
  loading.value = true
  try {
    const result = await getDashboards()
    dashboards.value = Array.isArray(result) ? result : []
  } catch (e) {
    console.error('Failed to fetch dashboards:', e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.toolbar {
  display: flex;
  margin-bottom: 1rem;
}

.dashboards-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.dashboard-card {
  background: #202024;
  border: 1px solid #2a2a2e;
  border-radius: 12px;
  padding: 1rem;
  width: 280px;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  transition: border-color 0.15s;
}

.dashboard-card:hover {
  border-color: #3a3a3e;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.card-icon {
  font-size: 1.1rem;
  color: #888;
}

.card-name {
  font-weight: 600;
  font-size: 1rem;
  color: #e0e0e0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.meta-label {
  font-size: 0.65rem;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.meta-value {
  font-size: 0.75rem;
  color: #888;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid #2a2a2e;
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
</style>
