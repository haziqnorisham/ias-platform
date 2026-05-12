<script setup>
import { ref, watch, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Select from 'primevue/select';

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    profile: {
        type: Object,
        default: null
    },
    protocols: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits(['update:visible', 'save', 'request-delete']);

const dialogMode = computed(() => props.profile ? 'edit' : 'create');

const form = ref({
    profile_name: '',
    manufacturer: '',
    model_number: '',
    communications_protocol: '',
    decoder: ''
});

watch(() => props.visible, (isVisible) => {
    if (isVisible) {
        if (props.profile) {
            form.value = {
                profile_name: props.profile.ProfileName || '',
                manufacturer: props.profile.Manufacturer || '',
                model_number: props.profile.ModelNumber || '',
                communications_protocol: props.profile.CommunicationsProtocol || '',
                decoder: props.profile.Decoder || ''
            };
        } else {
            form.value = {
                profile_name: '',
                manufacturer: '',
                model_number: '',
                communications_protocol: '',
                decoder: ''
            };
        }
    }
});

function closeDialog() {
    emit('update:visible', false);
}

function handleSave() {
    const payload = {
        ProfileName: form.value.profile_name,
        Manufacturer: form.value.manufacturer,
        ModelNumber: form.value.model_number,
        CommunicationsProtocol: form.value.communications_protocol,
        Decoder: form.value.decoder
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
                    v-model="form.profile_name"
                    placeholder="Enter profile name"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="manufacturer">Manufacturer</label>
                <InputText
                    id="manufacturer"
                    v-model="form.manufacturer"
                    placeholder="Enter manufacturer name"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="model_number">Model Number</label>
                <InputText
                    id="model_number"
                    v-model="form.model_number"
                    placeholder="Enter model number"
                    class="form-input"
                />
            </div>
            <div class="form-field">
                <label for="protocol">Communications Protocol</label>
                <Select
                    id="protocol"
                    v-model="form.communications_protocol"
                    :options="protocols"
                    placeholder="Select a protocol"
                    class="form-input"
                />
            </div>
            <div class="form-field full-width">
                <label for="decoder">Decoder</label>
                <Textarea
                    id="decoder"
                    v-model="form.decoder"
                    placeholder="Enter decoder logic or script..."
                    :autoResize="true"
                    rows="5"
                    class="form-input"
                />
            </div>
        </div>
        <template #footer>
            <Button v-if="dialogMode === 'edit'" label="Delete" icon="pi pi-trash" @click="handleDeleteRequest" severity="danger" class="p-button-text" />
            <div class="footer-right">
                <Button label="Cancel" icon="pi pi-times" @click="closeDialog" class="p-button-text" />
                <Button :label="dialogMode === 'create' ? 'Save Profile' : 'Update Profile'" icon="pi pi-save" @click="handleSave" />
            </div>
        </template>
    </Dialog>
</template>

<style scoped>
.profile-dialog :deep(.p-dialog-header) {
    border-bottom: 1px solid #212121;
    padding: 1.25rem 1.5rem;
}

.profile-dialog :deep(.p-dialog-content) {
    padding: 1.5rem;
}

.profile-dialog :deep(.p-dialog-footer) {
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

.footer-right {
    display: flex;
    gap: 0.5rem;
    margin-left: auto;
}
</style>
