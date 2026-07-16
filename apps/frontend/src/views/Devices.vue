<script setup>
import { ref, onMounted } from 'vue'

import Button from 'primevue/button';
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from "primevue/usetoast";

import { getAllDevices, createDevice, updateDevice, deleteDevice, getDeviceProfiles } from '@/api/posts'

import DeviceCard from '@/components/devices/DeviceCard.vue'
import DeviceFormDialog from '@/components/devices/DeviceFormDialog.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const toast = useToast();
const showError = (summary, detail) => {
    toast.add({ severity: 'error', summary: summary, detail: detail, life: 3000 });
};
const loading = ref(false)

const devices = ref([])

async function getAllDevicesHandler() {
    loading.value = false

    try {
        const result = await getAllDevices()
        result && result.forEach(function(obj){
            devices.value.push(obj)
        })
    } catch (error) {
        console.log(error)
        showError('API Call Failed', 'Failed to fetch devices details.')
    } finally {
        loading.value = false
    }
}

const deviceProfiles = ref([])

async function fetchDeviceProfiles() {
    try {
        const result = await getDeviceProfiles()
        if (result) {
            deviceProfiles.value = result.map(p => ({
                label: `${p.profile_name} (${p.profile_id})`,
                value: parseInt(p.profile_id, 10)
            }))
        }
    } catch (error) {
        console.error(error)
    }
}

onMounted(() => {
    getAllDevicesHandler()
    fetchDeviceProfiles()
})

// ── Device Form Dialog ──────────────────────────────────────
const formVisible = ref(false)
const editingDevice = ref(null)

function openAddDeviceDialog() {
    editingDevice.value = null
    formVisible.value = true
}

function openEditDeviceDialog(device) {
    editingDevice.value = device
    formVisible.value = true
}

async function handleSaveDevice(payload) {
    loading.value = false
    try {
        if (editingDevice.value) {
            const result = await updateDevice(payload)
            console.log('Device updated:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device updated successfully.', life: 3000 })
        } else {
            const result = await createDevice(payload)
            console.log('Device created:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device created successfully.', life: 3000 })
        }

        formVisible.value = false
        devices.value = []
        await getAllDevicesHandler()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to save device.')
    } finally {
        loading.value = false
    }
}

// ── Delete Device ───────────────────────────────────────────
const confirmVisible = ref(false)
const deviceToDelete = ref(null)

function handleDeleteRequest() {
    deviceToDelete.value = editingDevice.value
    confirmVisible.value = true
}

async function handleDeleteConfirm() {
    loading.value = false
    try {
        const payload = { device_id: String(deviceToDelete.value.Id) }
        const result = await deleteDevice(payload)
        console.log('Device deleted:', result)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Device deleted successfully.', life: 3000 })

        formVisible.value = false
        devices.value = []
        await getAllDevicesHandler()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to delete device.')
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />

    <div class="toolbar">
        <Button label="Add Device" icon="pi pi-plus" @click="openAddDeviceDialog" />
    </div>

    <div v-if="!devices.length" class="device_card empty_card">
        <div class="device_info">
            <div class="bold">No Device Found</div>
            <div class="italic">There are no devices to display.</div>
        </div>
    </div>

    <div v-else class="devices-grid">
        <DeviceCard
            v-for="device in devices"
            :key="device.Id"
            :device="device"
            @edit="openEditDeviceDialog"
        />
    </div>

    <DeviceFormDialog
        v-model:visible="formVisible"
        :device="editingDevice"
        :device-profiles="deviceProfiles"
        @save="handleSaveDevice"
        @request-delete="handleDeleteRequest"
    />

    <ConfirmDialog
        v-model:visible="confirmVisible"
        title="Delete Device"
        message="Are you sure you want to delete this device? This action cannot be undone."
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

.devices-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
}

.device_card {
    background-image: linear-gradient(to right, #2A2A2E 0%, #18181B 30%);
    border-radius: 12px;
    display: flex;
    flex-direction: row;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    width: 480px;
}

.empty_card {
    background-image: linear-gradient(to right, #18181B 0%, #18181B 30%);
    align-items: center;
    justify-content: center;
    min-height: 12rem;
}

.device_info {
    display: flex;
    flex-direction: column;
    white-space: nowrap;
    height: 100%;
    padding: 12px;
    padding-right: 5px;
    min-width: 0;
    overflow: hidden;
}

.italic {
    font-style: italic;
    font-size: 0.8rem;
}

.bold {
    margin-top: 0.25rem;
    font-weight: bold;
    font-size: 1.25rem;
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}
</style>
