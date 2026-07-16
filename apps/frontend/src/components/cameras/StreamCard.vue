<script setup>
import Chip from 'primevue/chip'

defineProps({
  stream: {
    type: Object,
    required: true
  },
  streamStatus: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['edit', 'view'])
</script>

<template>
  <div class="stream-card">
    <div class="stream-info">
      <div class="italic">Stream ID: {{ stream.id }}</div>
      <div class="bold stream-name-link" @click="emit('view', stream)">{{ stream.name }}</div>
      <div class="stream-ip">{{ stream.ip }}{{ stream.port && stream.port !== 80 ? ':' + stream.port : '' }}</div>
      <Chip
        v-if="streamStatus && streamStatus.status === 'running'"
        label="Live"
        icon="pi pi-circle-fill"
        class="chip-live"
      />
      <Chip
        v-else-if="streamStatus"
        :label="streamStatus.status"
        icon="pi pi-circle-fill"
        class="chip-stream"
        :class="'chip-' + streamStatus.status"
      />
      <Chip
        v-else
        label="Offline"
        icon="pi pi-circle-fill"
        class="chip-inactive"
      />
    </div>
    <div class="stream-details">
      <div class="detail-row"><span class="detail-label">Profile:</span> <span class="detail-value">{{ stream.stream_profile_token || '—' }}</span></div>
      <div class="detail-row"><span class="detail-label">IP:</span> <span class="detail-value">{{ stream.ip }}:{{ stream.port || 80 }}</span></div>
      <div class="detail-row"><span class="detail-label">User:</span> <span class="detail-value">{{ stream.username || '—' }}</span></div>
      <div class="detail-row"><span class="detail-label">Notes:</span> <span class="detail-value">{{ stream.notes || '—' }}</span></div>
    </div>
    <div class="stream-actions">
      <button class="action-btn view-btn" @click="emit('view', stream)">
        <i class="pi pi-play-circle"></i>
      </button>
      <button class="action-btn edit-btn" @click="emit('edit', stream)">
        <i class="pi pi-pencil"></i>
      </button>
    </div>
  </div>
</template>

<style scoped>
.stream-card {
  background-image: linear-gradient(to right, #2A2A2E 0%, #18181B 30%);
  border-radius: 12px;
  display: flex;
  flex-direction: row;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  width: 480px;
  position: relative;
}

.stream-info {
  display: flex;
  flex-direction: column;
  white-space: nowrap;
  padding: 12px;
  padding-right: 5px;
  min-width: 0;
  overflow: hidden;
}

.stream-details {
  display: flex;
  flex-direction: column;
  padding: 12px;
  gap: 0.5rem;
  margin-left: 16px;
  font-weight: bold;
  flex: 1;
}

.detail-row {
  font-size: 0.8rem;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 18ch;
}

.detail-label {
  font-style: italic;
  color: #888;
}

.detail-value {
  font-style: normal;
  color: #e0e0e0;
}

.italic {
  font-style: italic;
  font-size: 0.8rem;
}

.bold {
  margin-top: 0.25rem;
  font-weight: bold;
  font-size: 1.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 11ch;
}

.stream-name-link {
  cursor: pointer;
  transition: color 0.15s;
  text-decoration: underline;
}

.stream-name-link:hover {
  color: #90caf9;
}

.stream-ip {
  font-size: 0.75rem;
  color: #888;
  margin-top: 0.25rem;
}

.stream-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  padding-left: 0;
  justify-content: center;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid #3a3a3e;
  background: #202024;
  color: #e0e0e0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  font-size: 1rem;
}

.action-btn:hover {
  background: #3a3a3e;
  border-color: #555;
}

.view-btn:hover {
  color: #64b5f6;
}

.edit-btn:hover {
  color: #a5d6a7;
}

.chip-live {
  font-size: 0.75rem;
  margin-top: 0.5rem;
  background-color: #2d5a4e;
  box-shadow: 0 0 8px rgba(72, 137, 123, 0.45);
}

.chip-live :deep(.p-chip-icon) {
  color: #4cff88;
}

.chip-stream {
  font-size: 0.75rem;
  margin-top: 0.5rem;
}

.chip-running {
  background-color: #2d5a4e;
}

.chip-running :deep(.p-chip-icon) {
  color: #4cff88;
}

.chip-starting {
  background-color: #5a4e2d;
}

.chip-starting :deep(.p-chip-icon) {
  color: #ffc107;
}

.chip-error {
  background-color: #5a2d2d;
}

.chip-error :deep(.p-chip-icon) {
  color: #f44336;
}

.chip-stopped {
  background-color: #3a3a3e;
}

.chip-stopped :deep(.p-chip-icon) {
  color: #888;
}

.chip-inactive {
  font-size: 0.75rem;
  margin-top: 0.5rem;
  background-color: #3a3a3e;
}

.chip-inactive :deep(.p-chip-icon) {
  color: #888;
}
</style>
