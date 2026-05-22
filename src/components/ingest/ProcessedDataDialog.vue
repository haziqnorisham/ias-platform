<script setup>
import { ref, watch } from 'vue';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Tag from 'primevue/tag';
import ProgressSpinner from 'primevue/progressspinner';

import { getProcessedData } from '@/api/posts'

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    record: {
        type: Object,
        default: null
    }
});

const emit = defineEmits(['update:visible']);

const processedRecords = ref([])
const loading = ref(false)

watch(() => props.visible, async (isVisible) => {
    if (isVisible && props.record) {
        loading.value = true
        processedRecords.value = []
        try {
            const result = await getProcessedData({ raw_message_id: props.record.MessageID })
            if (result && result.records) {
                processedRecords.value = result.records
            }
        } catch (error) {
            console.error('Failed to fetch processed data:', error)
        } finally {
            loading.value = false
        }
    }
});

function closeDialog() {
    emit('update:visible', false);
}
</script>

<template>
    <Dialog
        :visible="visible"
        @update:visible="emit('update:visible', $event)"
        :header="'Processed Data — Msg ID: ' + (record?.MessageID ?? '')"
        :modal="true"
        :closable="false"
        :style="{ width: '95vw', maxWidth: '1400px' }"
        class="processed-data-dialog"
    >
        <ProgressSpinner v-if="loading" class="dialog-spinner" />

        <p v-else-if="!processedRecords.length" class="empty-state">No processed data available.</p>

        <DataTable v-else :value="processedRecords" size="small" scrollable scrollHeight="500px" tableStyle="min-width: 40rem" class="tight-table">
            <Column field="ID" header="ID"></Column>
            <Column field="DeviceID" header="Device ID"></Column>
            <Column field="ProfileID" header="Profile ID"></Column>
            <Column field="Success" header="Status">
                <template #body="{ data }">
                    <Tag
                        :value="data.Success ? 'Success' : 'Failed'"
                        :severity="data.Success ? 'success' : 'danger'"
                    />
                </template>
            </Column>
            <Column field="ErrorMessage" header="Error">
                <template #body="{ data }">
                    <span :class="{ 'error-text': !data.Success }">{{ data.ErrorMessage || '—' }}</span>
                </template>
            </Column>
            <Column field="ProcessedPayload" header="Processed Payload">
                <template #body="{ data }">
                    <span class="payload-cell">{{ data.ProcessedPayload || '—' }}</span>
                </template>
            </Column>
            <Column field="ProcessedAt" header="Processed At"></Column>
        </DataTable>

        <template #footer>
            <Button label="Close" icon="pi pi-times" @click="closeDialog" class="p-button-text" />
        </template>
    </Dialog>
</template>

<style scoped>
.processed-data-dialog :deep(.p-dialog-header) {
    border-bottom: 1px solid #212121;
    padding: 1.25rem 1.5rem;
}

.processed-data-dialog :deep(.p-dialog-content) {
    padding: 1rem 1.5rem;
}

.tight-table :deep(.p-datatable-thead > tr > th) {
    font-size: 0.75rem;
    padding: 0.4rem 0.6rem;
}

.tight-table :deep(.p-datatable-tbody > tr > td) {
    font-size: 0.78rem;
    padding: 0.3rem 0.6rem;
}

.processed-data-dialog :deep(.p-dialog-footer) {
    display: flex;
    align-items: center;
    justify-content: flex-end;
}

.dialog-spinner {
    display: flex;
    justify-content: center;
    padding: 2rem 0;
}

.empty-state {
    display: flex;
    justify-content: center;
    padding: 2rem 0;
    color: #888;
}

.error-text {
    color: #f44336;
}

.payload-cell {
    display: block;
    max-width: 250px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 0.85rem;
}
</style>
