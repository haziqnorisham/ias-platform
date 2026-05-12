<script setup>
import { ref, onMounted } from 'vue'

import Chip from 'primevue/chip';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';

import { getAllDevices, createDevice, updateDevice, deleteDevice, getDeviceProfiles } from '@/api/posts'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from "primevue/usetoast";
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

/* hc_devices table structure:
    - id INT NOT NULL (user-provided, not auto-assigned)
    - name VARCHAR(100) NOT NULL
    - description TEXT DEFAULT ''
    - profile_id INT REFERENCES hc_device_profiles(profile_id)
    - status VARCHAR(50) NOT NULL DEFAULT 'inactive'
    - location_label VARCHAR(100) DEFAULT ''
*/

// ── Add / Edit Device Dialog ───────────────────────────────────
const dialogVisible = ref(false)
const dialogMode = ref('create') // 'create' | 'edit'
const editingDeviceId = ref(null)
const deviceProfiles = ref([])

const newDevice = ref({
    id: '',
    name: '',
    description: '',
    location_label: '',
    profile_id: null
})

function openAddDeviceDialog() {
    dialogMode.value = 'create'
    editingDeviceId.value = null
    newDevice.value = {
        id: '',
        name: '',
        description: '',
        location_label: '',
        profile_id: null
    }
    dialogVisible.value = true
}

function openEditDeviceDialog(device) {
    dialogMode.value = 'edit'
    editingDeviceId.value = device.Id
    newDevice.value = {
        id: String(device.Id),
        name: device.Name || '',
        description: device.Description || '',
        location_label: device.LocationLabel || '',
        profile_id: device.ProfileID || null
    }
    dialogVisible.value = true
}

function closeDialog() {
    dialogVisible.value = false
}

async function fetchDeviceProfiles() {
    try {
        const result = await getDeviceProfiles()
        if (result) {
            deviceProfiles.value = result.map(p => ({
                label: `${p.ProfileName} (${p.ProfileID})`,
                value: p.ProfileID
            }))
        }
    } catch (error) {
        console.error(error)
    }
}

async function handleSaveDevice() {
    loading.value = false
    try {
        const payload = {
            Id: parseInt(newDevice.value.id),
            Name: newDevice.value.name,
            Description: newDevice.value.description || null,
            LocationLabel: newDevice.value.location_label || null,
            ProfileID: newDevice.value.profile_id
        }

        if (dialogMode.value === 'create') {
            const result = await createDevice(payload)
            console.log('Device created:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device created successfully.', life: 3000 })
        } else {
            const result = await updateDevice(payload)
            console.log('Device updated:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device updated successfully.', life: 3000 })
        }

        closeDialog()
        devices.value = []
        await getAllDevicesHandler()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to save device.')
    } finally {
        loading.value = false
    }
}

// ── Delete Device ──────────────────────────────────────────────
const confirmDeleteVisible = ref(false)

function openDeleteConfirmation() {
    confirmDeleteVisible.value = true
}

async function handleDeleteDevice() {
    loading.value = false
    confirmDeleteVisible.value = false
    try {
        const payload = { device_id: String(editingDeviceId.value) }
        const result = await deleteDevice(payload)
        console.log('Device deleted:', result)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Device deleted successfully.', life: 3000 })
        closeDialog()
        devices.value = []
        await getAllDevicesHandler()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to delete device.')
    } finally {
        loading.value = false
    }
}
// ─────────────────────────────────────────────────────────────────

onMounted(() => {
  getAllDevicesHandler()
  fetchDeviceProfiles()
})
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
        <div v-for="device in devices">
            <div class="device_card">
                <div class="device_info">
                <div class="italic">Device ID: {{ device.Id }}</div>
                <div class="bold device-name-link" @click="openEditDeviceDialog(device)">{{ device.Name }}</div>
                <Chip :label="device.Status || 'Unknown'" icon="pi pi-circle-fill" class="chip_style" :class="device.Status === 'active' ? 'chip_active' : 'chip_inactive'" />
            </div>
            <img src="/MCS.png" style="height: 10rem;">
            <div class="device_details">
                <div class="italic">Description: {{ device.Description || '—' }}</div>
                <div class="italic">Profile ID: {{ device.ProfileID || '—' }}</div>
                <div class="italic">Location: {{ device.LocationLabel || '—' }}</div>
            </div>
            </div>
        </div>
    </div>

    <!-- Add / Edit Device Dialog -->
    <Dialog 
        v-model:visible="dialogVisible" 
        :header="dialogMode === 'create' ? 'New Device' : 'Edit Device'" 
        :modal="true" 
        :closable="false"
        :style="{ width: '600px' }"
        class="device-dialog"
    >
        <div class="form-grid">
            <div class="form-field">
                <label for="device_id">Device ID <span class="required">*</span></label>
                <InputText 
                    id="device_id" 
                    v-model="newDevice.id" 
                    placeholder="Enter device ID (numeric)"
                    :disabled="dialogMode === 'edit'"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="device_name">Name <span class="required">*</span></label>
                <InputText 
                    id="device_name" 
                    v-model="newDevice.name" 
                    placeholder="Enter device name"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="description">Description</label>
                <InputText 
                    id="description" 
                    v-model="newDevice.description" 
                    placeholder="Enter description"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="location_label">Location Label</label>
                <InputText 
                    id="location_label" 
                    v-model="newDevice.location_label" 
                    placeholder="Enter location"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="profile">Device Profile <span class="required">*</span></label>
                <Select 
                    id="profile"
                    v-model="newDevice.profile_id" 
                    :options="deviceProfiles"
                    optionLabel="label"
                    optionValue="value"
                    placeholder="Select a device profile"
                    class="form-input"
                />
            </div>
        </div>
        <template #footer>
            <Button v-if="dialogMode === 'edit'" label="Delete" icon="pi pi-trash" @click="openDeleteConfirmation" severity="danger" class="p-button-text" />
            <div class="footer-right">
                <Button label="Cancel" icon="pi pi-times" @click="closeDialog" class="p-button-text" />
                <Button :label="dialogMode === 'create' ? 'Save Device' : 'Update Device'" icon="pi pi-save" @click="handleSaveDevice" />
            </div>
        </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog
        v-model:visible="confirmDeleteVisible"
        header="Delete Device"
        :modal="true"
        :closable="false"
        :style="{ width: '420px' }"
    >
        <p class="confirm-text">Are you sure you want to delete this device? This action cannot be undone.</p>
        <template #footer>
            <Button label="Cancel" icon="pi pi-times" @click="confirmDeleteVisible = false" class="p-button-text" />
            <Button label="Delete" icon="pi pi-trash" @click="handleDeleteDevice" severity="danger" />
        </template>
    </Dialog>
    
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

.devices-grid > div {
    flex: 0 0 480px;
}

.device_card {
    background-image: linear-gradient(to right, #2A2A2E 0%, #18181B 30%);
    border-radius: 12px;
    display: flex;
    flex-direction: row;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.empty_card {
    background-image: linear-gradient(to right, #18181B 0%, #18181B 30%);
    align-items: center;
    justify-content: center;
    min-height: 12rem;
}

.device_details {
    display: flex;
    flex-direction: column;
    white-space: nowrap;
    height: 100%;
    padding: 12px;
    gap: 0.65rem;
    margin-left: 20px;
    font-weight: bold;
}

.device_info {
    display: flex;
    flex-direction: column;
    white-space: nowrap;
    height: 100%;
    padding: 12px;
    padding-right: 5px;
}

.italic {
    font-style: italic;
    font-size: 0.8rem;
}

.bold {
    margin-top: 0.25rem;
    font-weight: bold;
    font-size: 1.25rem;
    text-decoration: underline;
}

/* ── Clickable Device Name ──────────────────────────────────── */
.device-name-link {
    cursor: pointer;
    transition: color 0.15s, text-decoration 0.15s;
}

.device-name-link:hover {
    color: #90caf9;
    text-decoration: underline;
}

.chip_style {
    font-size: 0.8rem;
    margin-top: 1rem;
}

.chip_active {
    background-color: #2d5a4e;
    box-shadow: 0 0 8px rgba(72, 137, 123, 0.45);
}

.chip_active :deep(.p-chip-icon) {
    color: #4cff88;
}

.chip_inactive {
    background-color: #3a3a3e;
    box-shadow: 0 0 8px rgba(72, 137, 123, 0.15);
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}

/* ── Add Device Dialog ──────────────────────────────────── */
.device-dialog :deep(.p-dialog-header) {
    border-bottom: 1px solid #212121;
    padding: 1.25rem 1.5rem;
}

.device-dialog :deep(.p-dialog-content) {
    padding: 1.5rem;
}

.form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
}

.form-field {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
}

.form-field label {
    font-family: "Space Grotesk", sans-serif;
    font-size: 0.85rem;
    font-weight: 600;
    color: #a0a0a0;
}

.form-field .required {
    color: #f44336;
}

.form-input {
    width: 100%;
}

.footer-right {
    display: flex;
    gap: 0.5rem;
    margin-left: auto;
}

/* ── Dialog Footer ──────────────────────────────────── */
.device-dialog :deep(.p-dialog-footer) {
    display: flex;
    align-items: center;
}

.confirm-text {
    color: #ccc;
    margin: 0;
    line-height: 1.5;
}
</style>