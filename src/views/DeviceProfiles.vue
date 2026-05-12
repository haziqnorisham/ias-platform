<script setup>
import { ref, onMounted } from 'vue'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Textarea from 'primevue/textarea';

import { getDeviceProfiles, createDeviceProfile, updateDeviceProfile, deleteDeviceProfile } from '@/api/posts'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from "primevue/usetoast";
const toast = useToast();
const showError = (summary, detail) => {
    toast.add({ severity: 'error', summary: summary, detail: detail, life: 3000 });
};
const loading = ref(false)

const profiles = ref([])

/* hc_device_profiles table structure:
    - profile_id SERIAL PRIMARY KEY
    - profile_name VARCHAR(100) NOT NULL
    - manufacturer VARCHAR(100) NOT NULL DEFAULT ''
    - communications_protocol VARCHAR(50) NOT NULL DEFAULT ''
    - decoder TEXT DEFAULT ''
*/

// ── New Profile Dialog ──────────────────────────────────────────
const dialogVisible = ref(false)
const dialogMode = ref('create') // 'create' | 'edit'
const editingProfileId = ref(null)
const protocols = ref(['MQTT', 'HTTP', 'TCP', 'UDP', 'LoRaWAN', 'BLE'])

const newProfile = ref({
    profile_name: '',
    manufacturer: '',
    model_number: '',
    communications_protocol: '',
    decoder: ''
})

function openNewProfileDialog() {
    dialogMode.value = 'create'
    editingProfileId.value = null
    newProfile.value = {
        profile_name: '',
        manufacturer: '',
        model_number: '',
        communications_protocol: '',
        decoder: ''
    }
    dialogVisible.value = true
}

function openEditProfileDialog(rowData) {
    dialogMode.value = 'edit'
    editingProfileId.value = rowData.ProfileID
    newProfile.value = {
        profile_name: rowData.ProfileName || '',
        manufacturer: rowData.Manufacturer || '',
        model_number: rowData.ModelNumber || '',
        communications_protocol: rowData.CommunicationsProtocol || '',
        decoder: rowData.Decoder || ''
    }
    dialogVisible.value = true
}

function closeDialog() {
    dialogVisible.value = false
}

// ── Delete Confirmation ────────────────────────────────────────
const confirmDeleteVisible = ref(false)

function openDeleteConfirmation() {
    confirmDeleteVisible.value = true
}

async function handleDeleteProfile() {
    loading.value = false
    confirmDeleteVisible.value = false
    try {
        const payload = { profile_id: editingProfileId.value }
        const result = await deleteDeviceProfile(payload)
        console.log('Profile deleted:', result)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile deleted successfully.', life: 3000 })
        closeDialog()
        await fetchDeviceProfiles()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to delete device profile.')
    } finally {
        loading.value = false
    }
}

async function handleSaveProfile() {
    loading.value = false
    try {
        const payload = {
            ProfileName: newProfile.value.profile_name,
            Manufacturer: newProfile.value.manufacturer,
            ModelNumber: newProfile.value.model_number,
            CommunicationsProtocol: newProfile.value.communications_protocol,
            Decoder: newProfile.value.decoder
        }

        if (dialogMode.value === 'create') {
            const result = await createDeviceProfile(payload)
            console.log('Profile created:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile created successfully.', life: 3000 })
        } else {
            payload.ProfileID = editingProfileId.value
            const result = await updateDeviceProfile(payload)
            console.log('Profile updated:', result)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Device profile updated successfully.', life: 3000 })
        }

        closeDialog()
        await fetchDeviceProfiles()
    } catch (error) {
        console.error(error)
        showError('API Call Failed', 'Failed to save device profile.')
    } finally {
        loading.value = false
    }
}
// ─────────────────────────────────────────────────────────────────

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

    <!-- New Profile Dialog -->
    <Dialog 
        v-model:visible="dialogVisible" 
        :header="dialogMode === 'create' ? 'New Device Profile' : 'Edit Device Profile'" 
        :modal="true" 
        :closable="false"
        :style="{ width: '720px' }"
        class="profile-dialog"
    >
        <div class="form-grid">
            <div class="form-field">
                <label for="profile_name">Profile Name <span class="required">*</span></label>
                <InputText 
                    id="profile_name" 
                    v-model="newProfile.profile_name" 
                    placeholder="Enter profile name"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="manufacturer">Manufacturer</label>
                <InputText 
                    id="manufacturer" 
                    v-model="newProfile.manufacturer" 
                    placeholder="Enter manufacturer name"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="model_number">Model Number</label>
                <InputText 
                    id="model_number" 
                    v-model="newProfile.model_number" 
                    placeholder="Enter model number"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="protocol">Communications Protocol</label>
                <Select 
                    id="protocol"
                    v-model="newProfile.communications_protocol" 
                    :options="protocols"
                    placeholder="Select a protocol"
                    class="form-input"
                />
            </div>
            <div class="form-field full-width">
                <label for="decoder">Decoder</label>
                <Textarea 
                    id="decoder"
                    v-model="newProfile.decoder" 
                    placeholder="Enter decoder logic or script..."
                    :autoResize="true"
                    rows="5"
                    class="form-input"
                />
            </div>
        </div>
        <template #footer>
            <Button v-if="dialogMode === 'edit'" label="Delete" icon="pi pi-trash" @click="openDeleteConfirmation" severity="danger" class="p-button-text" />
            <div class="footer-right">
                <Button label="Cancel" icon="pi pi-times" @click="closeDialog" class="p-button-text" />
                <Button :label="dialogMode === 'create' ? 'Save Profile' : 'Update Profile'" icon="pi pi-save" @click="handleSaveProfile" />
            </div>
        </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog
        v-model:visible="confirmDeleteVisible"
        header="Delete Device Profile"
        :modal="true"
        :closable="false"
        :style="{ width: '420px' }"
    >
        <p class="confirm-text">Are you sure you want to delete this device profile? This action cannot be undone.</p>
        <template #footer>
            <Button label="Cancel" icon="pi pi-times" @click="confirmDeleteVisible = false" class="p-button-text" />
            <Button label="Delete" icon="pi pi-trash" @click="handleDeleteProfile" severity="danger" />
        </template>
    </Dialog>
</template>

<style scoped>
.toolbar {
    display: flex;
    margin-bottom: 1rem;
}

.device_card {
    background-image: linear-gradient(to right, #2A2A2E 0%, #18181B 30%);
    border-radius: 12px;
    display: flex;
    flex-direction: row;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    margin-bottom: 30px;
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
    padding: 20px;
    gap: 1rem;
    margin-left: 40px;
    font-weight: bold;
}

.device_info {
    display: flex;
    flex-direction: column;
    white-space: nowrap;
    height: 100%;
    padding: 20px;
    padding-right: 5px;
}

.italic {
    font-style: italic;
    font-size: 0.8rem;
}

.bold {
    margin-top: 0.5rem;
    font-weight: bold;
    font-size: 1.5rem;
    text-decoration: underline;
}

.chip_style {
    font-size: 0.8rem;
    margin-top: 1rem;
    background-color: darkgreen;
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}

/* ── New Profile Dialog ──────────────────────────────────── */
.profile-dialog :deep(.p-dialog-header) {
    border-bottom: 1px solid #212121;
    padding: 1.25rem 1.5rem;
}

.profile-dialog :deep(.p-dialog-content) {
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

.form-field.full-width {
    grid-column: 1 / -1;
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

/* ── Clickable Profile ID ──────────────────────────────────── */
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

/* ── Dialog Footer ──────────────────────────────────── */
.profile-dialog :deep(.p-dialog-footer) {
    display: flex;
    align-items: center;
}

.footer-right {
    display: flex;
    gap: 0.5rem;
    margin-left: auto;
}

.confirm-text {
    color: #ccc;
    margin: 0;
    line-height: 1.5;
}
</style>