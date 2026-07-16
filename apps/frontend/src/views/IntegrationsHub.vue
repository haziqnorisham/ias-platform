<script setup>
import { useRouter } from 'vue-router'
import Panel from '@/components/Panel.vue'
import Button from 'primevue/button'

const router = useRouter()

const integrations = [
  {
    id: 'onvif',
    title: 'ONVIF Streams',
    subtitle: 'Manage and stream ONVIF device feeds',
    route: '/integrations/onvif-streams'
  }
]
</script>

<template>
  <div class="main-container">
    <Panel :title="item.title" :subtitle="item.subtitle" v-for="item in integrations" :key="item.id">
      <div class="centered">
        <Button label="Configure" icon="pi pi-cog" @click="router.push(item.route)"/>
      </div>
    </Panel>
  </div>
</template>

<style scoped>
.main-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 10px;
}
</style>
