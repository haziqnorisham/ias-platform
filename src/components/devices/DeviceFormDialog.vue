<script setup>
import { ref, watch, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Select from 'primevue/select';
import ToggleSwitch from 'primevue/toggleswitch';
import { useToast } from 'primevue/usetoast';

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    device: {
        type: Object,
        default: null
    },
    deviceProfiles: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['update:visible', 'save', 'request-delete']);
const toast = useToast();

const dialogMode = computed(() => props.device ? 'edit' : 'create');

const form = ref({
    id: '',
    name: '',
    description: '',
    location_label: '',
    profile_id: null,
    isActive: true
});

watch(() => props.visible, (isVisible) => {
    if (isVisible) {
        if (props.device) {
            form.value = {
                id: String(props.device.Id),
                name: props.device.Name || '',
                description: props.device.Description || '',
                location_label: props.device.LocationLabel || '',
                profile_id: props.device.ProfileID || null,
                isActive: props.device.Status === 'active'
            };
        } else {
            form.value = {
                id: '',
                name: '',
                description: '',
                location_label: '',
                profile_id: null,
                isActive: true
            };
        }
    }
});

function closeDialog() {
    emit('update:visible', false);
}

function handleSave() {
    if (!form.value.profile_id) {
        toast.add({ severity: 'error', summary: 'Validation Error', detail: 'Please select a Device Profile.', life: 3000 });
        return;
    }
    const payload = {
        Id: form.value.id,
        Name: form.value.name,
        Description: form.value.description || null,
        LocationLabel: form.value.location_label || null,
        ProfileID: form.value.profile_id,
        Status: form.value.isActive ? 'active' : 'inactive'
    };
    emit('save', payload);
}

function handleDeleteRequest() {
    emit('request-delete');
}
</script>

<template>
    <Dialog
        :visible="visible"
        @update:visible="emit('update:visible', $event)"
        :header="dialogMode === 'create' ? 'New Device' : 'Edit Device'"
        :modal="true"
        :closable="false"
        :style="{ width: '600px' }"
        class="device-dialog"
    >
        <div class="form-grid">
            <div class="form-field form-field-wide">
                <label for="status">Status</label>
                <ToggleSwitch v-model="form.isActive" inputId="status" />
                <div class="toggle-label">{{ form.isActive ? 'Active' : 'Inactive' }}</div>
            </div>
            <div class="form-field">
                <label for="device_id">Device ID <span class="required">*</span></label>
                <InputText
                    id="device_id"
                    v-model="form.id"
                    placeholder="Enter device ID"
                    :disabled="dialogMode === 'edit'"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="device_name">Name <span class="required">*</span></label>
                <InputText
                    id="device_name"
                    v-model="form.name"
                    placeholder="Enter device name"
                    class="form-input"
                />
            </div>
            <div class="form-field form-field-wide">
                <label for="description">Description</label>
                <Textarea
                    id="description"
                    v-model="form.description"
                    placeholder="Enter description"
                    class="form-input"
                    rows="3"
                    autoResize
                />
            </div>
            <div class="form-field form-field-wide">
                <label for="location_label">Location Label</label>
                <Textarea
                    id="location_label"
                    v-model="form.location_label"
                    placeholder="Enter location"
                    class="form-input"
                    rows="3"
                    autoResize
                />
            </div>
            <div class="form-field">
                <label for="profile">Device Profile <span class="required">*</span></label>
                <Select
                    id="profile"
                    v-model="form.profile_id"
                    :options="deviceProfiles"
                    optionLabel="label"
                    optionValue="value"
                    placeholder="Select a device profile"
                    class="form-input"
                />
            </div>
        </div>
        <template #footer>
            <Button v-if="dialogMode === 'edit'" label="Delete" icon="pi pi-trash" @click="handleDeleteRequest" severity="danger" class="p-button-text" />
            <div class="footer-right">
                <Button label="Cancel" icon="pi pi-times" @click="closeDialog" class="p-button-text" />
                <Button :label="dialogMode === 'create' ? 'Save Device' : 'Update Device'" icon="pi pi-save" @click="handleSave" />
            </div>
        </template>
    </Dialog>
</template>

<style scoped>
.device-dialog :deep(.p-dialog-header) {
    border-bottom: 1px solid #212121;
    padding: 1.25rem 1.5rem;
}

.device-dialog :deep(.p-dialog-content) {
    padding: 1.5rem;
}

.device-dialog :deep(.p-dialog-footer) {
    display: flex;
    align-items: center;
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

.form-field-wide {
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

.toggle-label {
    font-size: 0.8rem;
    color: #a0a0a0;
    margin-top: 0.25rem;
}

.footer-right {
    display: flex;
    gap: 0.5rem;
    margin-left: auto;
}
</style>
