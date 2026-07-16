<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  value: {
    type: [String, Number],
    default: '0'
  },
  icon: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  hideTitle: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  error: {
    type: Boolean,
    default: false
  }
});
</script>

<template>
  <div class="metric-card">
    <div v-if="!hideTitle && title" class="metric-card__title">{{ title }}</div>
    <div class="metric-card__value" :class="{ 'metric-card__value--loading': loading, 'metric-card__value--error': error }">
      {{ value }}
      <i v-if="error" class="pi pi-exclamation-triangle error-icon"></i>
    </div>
    <div v-if="subtitle" class="metric-card__subtitle">{{ subtitle }}</div>
  </div>
</template>

<style scoped>
.metric-card {
  --mc-value-color: #6d6d6d;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  padding: 12px;
  box-sizing: border-box;
}

.metric-card__title {
  font-size: var(--font-size-md);
  font-weight: 600;
  color: #aaa;
  margin-bottom: 8px;
}

.metric-card__value {
  font-size: var(--font-size-2xl);
  font-weight: bold;
  color: var(--mc-value-color);
  line-height: 1.2;
  transition: opacity 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.metric-card__value--loading {
  opacity: 0.5;
}

.metric-card__value--error .error-icon {
  font-size: var(--font-size-sm);
  color: #f44336;
  flex-shrink: 0;
}

.metric-card__subtitle {
  font-size: var(--font-size-xs);
  color: #777;
  margin-top: 4px;
}
</style>
