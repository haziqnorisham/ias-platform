<script setup>
import { ref, watch, computed } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Select from 'primevue/select';
import { useToast } from 'primevue/usetoast';

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
const toast = useToast();

const dialogMode = computed(() => props.profile ? 'edit' : 'create');

const fileInput = ref(null);
const imageBase64 = ref(null);
const serverImageFailed = ref(false);
const previewSrc = computed(() => imageBase64.value ? `data:image/png;base64,${imageBase64.value}` : null);
const serverImageUrl = computed(() => {
    if (dialogMode.value === 'edit' && props.profile?.profile_id) {
        return `/api/image/device_profile_${props.profile.profile_id}.png`;
    }
    return null;
});

const form = ref({
    profile_name: '',
    manufacturer: '',
    model_number: '',
    communications_protocol: '',
    decoder: ''
});

const dataMappings = ref([]);
const dataTypeOptions = ['number', 'text'];

function addMapping() {
    dataMappings.value.push({ name: '', type: 'number' });
}

function removeMapping(index) {
    dataMappings.value.splice(index, 1);
}

watch(() => props.visible, (isVisible) => {
    if (isVisible) {
        imageBase64.value = null;
        serverImageFailed.value = false;
        if (fileInput.value) fileInput.value.value = '';
        if (props.profile) {
            form.value = {
                profile_name: props.profile.profile_name || '',
                manufacturer: props.profile.manufacturer || '',
                model_number: props.profile.model_number || '',
                communications_protocol: props.profile.communications_protocol || '',
                decoder: props.profile.decoder || ''
            };
            if (props.profile.image_base64) {
                imageBase64.value = props.profile.image_base64;
            }
            dataMappings.value = props.profile.data_mappings
                ? props.profile.data_mappings.map(m => ({ ...m }))
                : [];
        } else {
            form.value = {
                profile_name: '',
                manufacturer: '',
                model_number: '',
                communications_protocol: '',
                decoder: ''
            };
            dataMappings.value = [];
        }
    }
});

function triggerFileInput() {
    fileInput.value?.click();
}

function handleImageSelect(event) {
    const file = event.target.files[0];
    if (!file) return;

    const ext = file.name.split('.').pop().toLowerCase();
    if (ext !== 'png') {
        toast.add({ severity: 'error', summary: 'Invalid File Type', detail: 'Only .png files are allowed.', life: 3000 });
        event.target.value = '';
        return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
        const dataUrl = e.target.result;
        const img = new Image();
        img.onload = () => {
            if (img.naturalWidth !== 216 || img.naturalHeight !== 329) {
                toast.add({ severity: 'error', summary: 'Invalid Dimensions', detail: `Image must be 216×329px. Selected image is ${img.naturalWidth}×${img.naturalHeight}px.`, life: 3000 });
                imageBase64.value = null;
                event.target.value = '';
                return;
            }
            imageBase64.value = dataUrl.replace('data:image/png;base64,', '');
        };
        img.onerror = () => {
            toast.add({ severity: 'error', summary: 'Invalid Image', detail: 'The selected file is not a valid image.', life: 3000 });
            imageBase64.value = null;
            event.target.value = '';
        };
        img.src = dataUrl;
    };
    reader.readAsDataURL(file);
}

function handleRemoveImage() {
    imageBase64.value = null;
    serverImageFailed.value = true;
    if (fileInput.value) fileInput.value.value = '';
}

function closeDialog() {
    emit('update:visible', false);
}

function handleSave() {
    const payload = {
        profile_name: form.value.profile_name,
        manufacturer: form.value.manufacturer,
        model_number: form.value.model_number,
        communications_protocol: form.value.communications_protocol,
        decoder: form.value.decoder,
        data_mappings: dataMappings.value
    };
    if (imageBase64.value) {
        payload.image_base64 = imageBase64.value;
    }
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
            <div class="form-field full-width image-upload-row">
                <label>Profile Image</label>
                <div class="image-upload-area">
                    <input
                        type="file"
                        ref="fileInput"
                        accept=".png"
                        style="display: none"
                        @change="handleImageSelect"
                    />
                    <div class="image-preview-wrapper" :class="{ 'has-image': previewSrc || (serverImageUrl && !serverImageFailed) }">
                        <img v-if="previewSrc" :src="previewSrc" alt="Profile Image" class="image-preview" />
                        <img v-else-if="serverImageUrl && !serverImageFailed" :src="serverImageUrl" alt="Existing Profile Image" class="image-preview" @error="serverImageFailed = true" />
                        <div v-else class="image-placeholder">
                            <i class="pi pi-image" style="font-size: 2rem; color: #555"></i>
                            <span>216 × 329px</span>
                            <span class="placeholder-hint">PNG only</span>
                        </div>
                    </div>
                    <div class="image-upload-actions">
                        <Button
                            label="Upload Image"
                            icon="pi pi-upload"
                            size="small"
                            @click="triggerFileInput"
                        />
                        <Button
                            v-if="previewSrc || (serverImageUrl && !serverImageFailed)"
                            label="Remove"
                            icon="pi pi-times"
                            size="small"
                            severity="secondary"
                            class="p-button-text"
                            @click="handleRemoveImage"
                        />
                    </div>
                </div>
            </div>
            <div class="form-field full-width">
                <label for="decoder">Decoder</label>
                <Textarea
                    id="decoder"
                    v-model="form.decoder"
                    placeholder="Enter decoder logic or script..."
                    :autoResize="true"
                    rows="5"
                    class="form-input decoder-textarea"
                />
            </div>
            <div class="form-field full-width">
                <label>Data Mappings</label>
                <div class="mappings-container">
                    <div v-if="!dataMappings.length" class="mappings-empty">
                        No columns defined. Click "Add Column" to define a data mapping.
                    </div>
                    <div v-for="(mapping, index) in dataMappings" :key="index" class="mapping-row">
                        <InputText
                            v-model="mapping.name"
                            placeholder="Column name"
                            class="mapping-name-input"
                        />
                        <Select
                            v-model="mapping.type"
                            :options="dataTypeOptions"
                            placeholder="Type"
                            class="mapping-type-select"
                        />
                        <Button
                            icon="pi pi-trash"
                            size="small"
                            severity="secondary"
                            class="p-button-text mapping-remove-btn"
                            @click="removeMapping(index)"
                        />
                    </div>
                    <Button
                        label="Add Column"
                        icon="pi pi-plus"
                        size="small"
                        class="p-button-text mapping-add-btn"
                        @click="addMapping"
                    />
                </div>
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

.image-upload-area {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
}

.image-preview-wrapper {
    width: 108px;
    height: 165px;
    border: 2px dashed #3a3a3e;
    border-radius: 6px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    background: #121214;
    transition: border-color 0.2s;
}

.image-preview-wrapper.has-image {
    border-style: solid;
    border-color: #2a2a2e;
}

.image-preview {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.image-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    color: #666;
    font-size: 0.75rem;
    font-family: "Space Grotesk", sans-serif;
    text-align: center;
}

.image-placeholder .placeholder-hint {
    color: #555;
    font-size: 0.7rem;
}

.image-upload-actions {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.decoder-textarea :deep(textarea) {
    resize: vertical;
    font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 0.82rem;
    line-height: 1.55;
    background: #0a0a0c;
    color: #d4d4d4;
    border-color: #2a2a2e;
    tab-size: 4;
}

.mappings-container {
    background: #121214;
    border: 1px solid #2a2a2e;
    border-radius: 6px;
    padding: 0.75rem;
}

.mappings-empty {
    color: #666;
    font-size: 0.8rem;
    padding: 0.5rem 0;
    text-align: center;
}

.mapping-row {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    margin-bottom: 0.5rem;
}

.mapping-row:last-of-type {
    margin-bottom: 0;
}

.mapping-name-input {
    flex: 1;
}

.mapping-type-select {
    width: 130px;
    flex-shrink: 0;
}

.mapping-remove-btn {
    flex-shrink: 0;
}

.mapping-add-btn {
    margin-top: 0.25rem;
}
</style>
