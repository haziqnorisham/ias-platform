<script setup>
import { ref, computed } from 'vue';
import Chip from 'primevue/chip';

const props = defineProps({
    device: {
        type: Object,
        required: true
    }
});

const emit = defineEmits(['edit']);

const fallbackActive = ref(false);
const imageSrc = computed(() => {
    if (fallbackActive.value || !props.device.ProfileID) return '/MCS.png';
    return `/device_profile_${props.device.ProfileID}.png`;
});

function onImageError() {
    fallbackActive.value = true;
}
</script>

<template>
    <div class="device_card">
        <div class="device_info">
            <div class="italic">Device ID: {{ device.Id }}</div>
            <div class="bold device-name-link" @click="emit('edit', device)">{{ device.Name }}</div>
            <Chip :label="device.Status || 'Unknown'" icon="pi pi-circle-fill" class="chip_style" :class="device.Status === 'active' ? 'chip_active' : 'chip_inactive'" />
        </div>
        <img :src="imageSrc" style="height: 10rem;" @error="onImageError">
        <div class="device_details">
            <div class="detail-row"><span class="detail-label">Name:</span> <span class="detail-value">{{ device.Name || '—' }}</span></div>
            <div class="detail-row"><span class="detail-label">Description:</span> <span class="detail-value">{{ device.Description || '—' }}</span></div>
            <div class="detail-row"><span class="detail-label">Profile ID:</span> <span class="detail-value">{{ device.ProfileID || '—' }}</span></div>
            <div class="detail-row"><span class="detail-label">Location:</span> <span class="detail-value">{{ device.LocationLabel || '—' }}</span></div>
        </div>
    </div>
</template>

<style scoped>
.device_card {
    background-image: linear-gradient(to right, #2A2A2E 0%, #18181B 30%);
    border-radius: 12px;
    display: flex;
    flex-direction: row;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
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

.detail-row {
    font-size: 0.8rem;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 20ch;
}

.detail-label {
    font-style: italic;
    color: #888;
}

.detail-value {
    font-style: normal;
    color: #e0e0e0;
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
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 11ch;
}

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
</style>
