<template>
  <div class="text-widget">
    <p class="text-widget-content">{{ text }}</p>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: 'Text' },
  text: { type: String, default: 'Double-click to edit this text.' }
})
</script>

<style scoped>
.text-widget {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 12px;
  box-sizing: border-box;
}

.text-widget-content {
  color: #ccc;
  font-size: var(--font-size-md);
  font-weight: 500;
  letter-spacing: 0.03em;
  margin: 0;
  text-align: center;
  user-select: none;
}
</style>
