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
  }
});
</script>

<template>
  <div class="metric-card">
    <div v-if="!hideTitle && title" class="metric-card__title">{{ title }}</div>
    <div class="metric-card__value">{{ value }}</div>
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
  font-size: 1rem;
  font-weight: 600;
  color: #aaa;
  margin-bottom: 8px;
}

.metric-card__value {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--mc-value-color);
  line-height: 1.2;
}

.metric-card__subtitle {
  font-size: 0.75rem;
  color: #777;
  margin-top: 4px;
}
</style>
