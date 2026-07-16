<template>
      <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />
  <div class="main-container">
    <Panel title="UI Version" subtitle="ALPHA v0.01" />
    <Panel title="Backend Version" subtitle="ALPHA v0.01" />
    <Panel title="Schema Version" subtitle="ALPHA v0.01" />
    <Panel title="Tenant ID" subtitle="DEV-ENV" />
    <Panel title="Instance UUID" subtitle="93cc7f26-14e8-441f-bacd-4ab63d5a9adb" />

  </div>

  <!-- PrimeVue Dialog Modal -->
  <Dialog v-model:visible="modalVisible" header="Test Result" :modal="true" :closable="true" :style="{ width: '450px' }">
    <div class="modal-content">
      <i class="pi pi-check-circle" style="font-size: 3rem; color: #4CAF50; margin-bottom: 1rem;"></i>
      <p style="font-size: 1.2rem; margin: 0;">Test Successful</p>
    </div>
    <template #footer>
      <Button label="Close" icon="pi pi-times" @click="closeModal" class="p-button-text"/>
    </template>
  </Dialog>
</template>

<script setup>
import { ref } from 'vue';
import Panel from "../components/Panel.vue";
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'

// Page blocking state
const loading = ref(false)

// Modal visibility state
const modalVisible = ref(false);

// Function to show the modal
async function showTestModal() {
  loading.value = true

    try {
        await new Promise(r => setTimeout(r, 2000));
    } finally {
        loading.value = false
        modalVisible.value = true;
    }
  
};

// Function to close the modal
const closeModal = () => {
  modalVisible.value = false;
};
</script>

<style scoped>
.main-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    gap: 10px;
}

.modal-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem 0;
}
.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}
</style>