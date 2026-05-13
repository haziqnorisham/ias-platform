<template>
  <div class="page-container">
    <h2 class="page-title">Settings</h2>
    <div class="placeholder-card">
      <i class="pi pi-cog placeholder-icon"></i>
      <p class="placeholder-text">Settings will be available here.</p>
    </div>
  </div>
</template>

<script setup>
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
</style>
