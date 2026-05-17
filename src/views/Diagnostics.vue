<template>
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />
  <div class="main-container">
    <Panel title="Backend Link Test" subtitle="Send a HTTP request to backend /api/test endpoint.">
      <div class="centered">
        <Button label="Run Test" icon="pi pi-play" class="btn" @click="runBackendLinkTest"/>
      </div>
    </Panel>
  </div>

  <!-- PrimeVue Dialog Modal -->
  <Dialog v-model:visible="modalVisible" :header="modalHeader" :modal="true" :closable="false" :style="{ width: '500px' }">
    <div class="modal-content">
      <i :class="modalIconClass" :style="{ fontSize: '3rem', color: modalIconColor, marginBottom: '1rem' }"></i>
      <p style="font-size: 1.2rem; margin: 0 0 1rem 0;">{{ modalMessage }}</p>
      <div v-if="modalDetails" class="modal-details">
        <pre>{{ modalDetails }}</pre>
      </div>
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

import { testApi } from '@/api/posts'

// Page blocking state
const loading = ref(false)

// Modal visibility state
const modalVisible = ref(false);
const modalHeader = ref('Test Result');
const modalMessage = ref('');
const modalIconClass = ref('pi pi-check-circle');
const modalIconColor = ref('#48897b');
const modalDetails = ref(null);

// Generic function to show the modal with custom content
async function showTestModal(testName, success = true, message = 'Test Successful', details = null) {
  if (success) {
    modalHeader.value = `${testName} - Success`;
    modalMessage.value = message;
    modalIconClass.value = 'pi pi-check-circle';
    modalIconColor.value = '#48897b';
  } else {
    modalHeader.value = `${testName} - Failed`;
    modalMessage.value = message;
    modalIconClass.value = 'pi pi-times-circle';
    modalIconColor.value = '#f44336';
  }
  modalDetails.value = details;
  modalVisible.value = true;
}

// Function to run the backend link test
async function runBackendLinkTest() {
  loading.value = true
  try {
    const response = await testApi()
    
    if (response.status == "OK") {
      console.log("Backend Link Test Successful")
      // Pass the response data into the modal
      showTestModal(
        'Backend Link Test', 
        true, 
        'Backend connection successful!', 
        JSON.stringify(response, null, 2)
      )
    } else {
      console.log("Backend Link Test Failed")
      showTestModal(
        'Backend Link Test', 
        false, 
        `Backend test failed: ${response.message || 'Unknown error'}`, 
        JSON.stringify(response, null, 2)
      )
    }
  } catch (error) {
    console.error("Backend Link Test Error:", error)
    showTestModal(
      'Backend Link Test', 
      false, 
      `Error: ${error.message || 'Failed to connect to backend'}`, 
      error.stack || JSON.stringify(error, null, 2)
    )
  } finally {
    loading.value = false
  }
}

// Function to close the modal
const closeModal = () => {
  modalVisible.value = false;
  modalDetails.value = null; // Clear details when closing
};
</script>

<style scoped>
.main-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 10px;
}

.modal-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem 0;
}

.modal-details {
    margin-top: 1rem;
    padding: 1rem;
    border-radius: 4px;
    max-width: 100%;
    background-color: #212121;
    overflow-x: auto;
}

.modal-details pre {
    margin: 0;
    font-size: 0.85rem;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}
</style>