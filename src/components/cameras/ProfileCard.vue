<script setup>
defineProps({
  profile: {
    type: Object,
    required: true
  },
  selected: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['select'])
</script>

<template>
  <div class="profile-card" :class="{ 'profile-selected': selected }">
    <div class="profile-header">
      <span class="profile-name">{{ profile.name || 'Unnamed' }}</span>
      <span v-if="selected" class="profile-badge">Selected</span>
    </div>
    <div class="profile-body">
      <div class="profile-row">
        <span class="profile-label">Token</span>
        <code class="profile-token">{{ profile.token }}</code>
      </div>
      <div class="profile-row">
        <span class="profile-label">Encoding</span>
        <span class="profile-value">{{ profile.encoding || '—' }}</span>
      </div>
      <div class="profile-row">
        <span class="profile-label">Resolution</span>
        <span class="profile-value">{{ profile.width }} x {{ profile.height }}</span>
      </div>
    </div>
    <div class="profile-footer">
      <button
        v-if="!selected"
        class="select-btn"
        @click="emit('select', profile.token)"
      >
        Select Profile
      </button>
    </div>
  </div>
</template>

<style scoped>
.profile-card {
  background: #202024;
  border: 1px solid #3a3a3e;
  border-radius: 10px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-width: 180px;
  transition: border-color 0.15s;
}

.profile-card:hover {
  border-color: #555;
}

.profile-selected {
  border-color: #48897b;
  box-shadow: 0 0 12px rgba(72, 137, 123, 0.3);
}

.profile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-name {
  font-weight: 600;
  font-size: 1rem;
}

.profile-badge {
  font-size: 0.65rem;
  background: #2d5a4e;
  color: #4cff88;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.profile-body {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.profile-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-label {
  font-size: 0.75rem;
  color: #888;
}

.profile-value {
  font-size: 0.8rem;
  color: #e0e0e0;
}

.profile-token {
  font-size: 0.75rem;
  color: #64b5f6;
  background: #18181B;
  padding: 1px 6px;
  border-radius: 4px;
}

.profile-footer {
  margin-top: auto;
  padding-top: 0.5rem;
  border-top: 1px solid #2a2a2e;
}

.select-btn {
  width: 100%;
  padding: 6px 12px;
  background: #2a2a2e;
  border: 1px solid #3a3a3e;
  border-radius: 6px;
  color: #e0e0e0;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.15s;
}

.select-btn:hover {
  background: #3a3a3e;
  border-color: #48897b;
  color: #fff;
}
</style>
