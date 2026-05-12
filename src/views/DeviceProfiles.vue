<script setup>
import { ref, onMounted } from 'vue'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

import Button from 'primevue/button';
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from "primevue/usetoast";

import { getDeviceProfiles, createDeviceProfile, updateDeviceProfile, deleteDeviceProfile } from '@/api/posts'

import DeviceProfileFormDialog from '@/components/devices/DeviceProfileFormDialog.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const toast = useToast();
const showError = (summary, detail) => {
    toast.add({ severity: 'error', summary: summary, detail: detail, life: 3000 });
};
const loading = ref(false)

const profiles = ref([])

const protocols = ref(['MQTT', 'HTTP', 'TCP', 'UDP', 'LoRaWAN', 'BLE'])

const formVisible = ref(false)
const editingProfile = ref(null)

function openNewProfileDialog() {
    editingProfile.value = null
    formVisible.value = true
}

function openEditProfileDialog(rowData) {
    editingProfile.value = rowData
    formVisible.value = true
}

const confirmVisible = ref(false)

function handleDeleteRequest() {
    confirmVisible.value = true
}

async function handleDeleteConfirm() {
    loading.value = false
    confirmVisible.value = false
    try {
        const payload = { profile_id: editingProfile.value.ProfileID }
        const result = await deleteDeviceProfile(payload)
        console.log('Profile deleted:', result)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile deleted successfully.', life: 3000 })
        formVisible.value = false
        await fetchDeviceProfiles()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to delete device profile.')
    } finally {
        loading.value = false
    }
}

async function handleSaveProfile(payload) {
    loading.value = false
    try {
        if (editingProfile.value) {
            payload.ProfileID = editingProfile.value.ProfileID
            const result = await updateDeviceProfile(payload)
            console.log('Profile updated:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile updated successfully.', life: 3000 })
        } else {
            const result = await createDeviceProfile(payload)
            console.log('Profile created:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile created successfully.', life: 3000 })
        }

        formVisible.value = false
        await fetchDeviceProfiles()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to save device profile.')
    } finally {
        loading.value = false
    }
}

async function fetchDeviceProfiles() {
    loading.value = false

    try {
        const result = await getDeviceProfiles()
        if (result) {
            profiles.value = result
        }
    } catch (error) {
        showError('API Call Failed', 'Failed to fetch device profiles.')
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchDeviceProfiles()
})
</script>

<template>
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />

    <div class="toolbar">
        <Button label="New Profile" icon="pi pi-plus" @click="openNewProfileDialog" />
    </div>

    <DataTable :value="profiles" stripedRows tableStyle="min-width: 50rem">
        <Column field="ProfileID" header="Profile ID">
            <template #body="{ data }">
                <span class="profile-id-link" @click.stop="openEditProfileDialog(data)">{{ data.ProfileID }}</span>
            </template>
        </Column>
        <Column field="ProfileName" header="Profile Name"></Column>
        <Column field="Manufacturer" header="Manufacturer"></Column>
        <Column field="ModelNumber" header="Model Number"></Column>
        <Column field="CommunicationsProtocol" header="Communications Protocol"></Column>
        <Column field="Decoder" header="Decoder"></Column>
    </DataTable>

    <DeviceProfileFormDialog
        v-model:visible="formVisible"
        :profile="editingProfile"
        :protocols="protocols"
        @save="handleSaveProfile"
        @request-delete="handleDeleteRequest"
    />

    <ConfirmDialog
        v-model:visible="confirmVisible"
        title="Delete Device Profile"
        message="Are you sure you want to delete this device profile? This action cannot be undone."
        confirm-label="Delete"
        severity="danger"
        @confirm="handleDeleteConfirm"
    />
</template>

<style scoped>
.toolbar {
    display: flex;
    margin-bottom: 1rem;
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}

.profile-id-link {
    color: #64b5f6;
    cursor: pointer;
    font-weight: 600;
    transition: color 0.15s, text-decoration 0.15s;
}

.profile-id-link:hover {
    color: #90caf9;
    text-decoration: underline;
}
</style>
