<script setup>
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    title: {
        type: String,
        default: 'Confirm'
    },
    message: {
        type: String,
        default: 'Are you sure?'
    },
    confirmLabel: {
        type: String,
        default: 'Confirm'
    },
    severity: {
        type: String,
        default: 'danger'
    }
});

const emit = defineEmits(['update:visible', 'confirm']);

function onCancel() {
    emit('update:visible', false);
}

function onConfirm() {
    emit('confirm');
    emit('update:visible', false);
}
</script>

<template>
    <Dialog
        :visible="visible"
        @update:visible="emit('update:visible', $event)"
        :header="title"
        :modal="true"
        :closable="false"
        :style="{ width: '420px' }"
    >
        <p class="confirm-text">{{ message }}</p>
        <template #footer>
            <Button label="Cancel" icon="pi pi-times" @click="onCancel" class="p-button-text" />
            <Button :label="confirmLabel" icon="pi pi-trash" @click="onConfirm" :severity="severity" />
        </template>
    </Dialog>
</template>

<style scoped>
.confirm-text {
    color: #ccc;
    margin: 0;
    line-height: 1.5;
}
</style>
