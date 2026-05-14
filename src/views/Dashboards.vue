<template>
  <div>
    <div class="toolbar">
      <Button label="New Dashboard" icon="pi pi-plus" @click="router.push('/dashboards/edit')" />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Button from 'primevue/button';

const router = useRouter()
</script>

<style scoped>
.toolbar {
  display: flex;
  margin-bottom: 1rem;
}
</style>
