<script setup>
import { ref, onMounted } from 'vue'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

import Button from 'primevue/button';
import Tag from 'primevue/tag';
import Dialog from 'primevue/dialog';

import { getRawIngest } from '@/api/posts'
import ProcessedDataDialog from '@/components/ingest/ProcessedDataDialog.vue'
import BlockUI from 'primevue/blockui'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from "primevue/usetoast";
const toast = useToast();
const showError = (summary, detail) => {
    toast.add({ severity: 'error', summary: summary, detail: detail, life: 3000 });
};
const showInfo = (summary, detail) => {
    toast.add({ severity: 'info', summary: summary, detail: detail, life: 3000 });
};
const showWarn = (summary, detail) => {
    toast.add({ severity: 'warn', summary: summary, detail: detail, life: 3000 });
};
const loading = ref(false)

const dialogVisible = ref(false)
const selectedPayload = ref('')

const processedDialogVisible = ref(false)
const selectedRecord = ref(null)

function onMsgIdClick(event, data) {
    event.stopPropagation()
    selectedRecord.value = data
    processedDialogVisible.value = true
}

function openPayloadDialog(payload, event) {
    event.stopPropagation()
    selectedPayload.value = payload
    dialogVisible.value = true
}

const records = ref([])
const total = ref(0)
const first = ref(0)
const rows = ref(10)
const selectedRecords = ref([])

function onPage(event) {
    first.value = event.first
    rows.value = event.rows
    pagination.value.limit = event.rows
    pagination.value.offset = event.first
    fetchIngestLogs()
}

const pagination = ref({
    limit: 10,
    offset: 0,
    sortByMsgID: 'desc',
    status: ''
})

async function fetchIngestLogs() {
    loading.value = false

    try {
        const result = await getRawIngest({
            Limit: pagination.value.limit,
            Offset: pagination.value.offset,
            SortByMsgID: pagination.value.sortByMsgID,
            Status: pagination.value.status
        })
        if (result && result.records) {
            records.value = result.records
            total.value = result.total || 0
        }
    } catch (error) { 
        console.log(error)
        showError('API Call Failed', 'Failed to fetch ingest logs.')
    } finally {
        loading.value = false
    }
}

function refreshLogs() {
    first.value = 0
    pagination.value.offset = 0
    selectedRecords.value = []
    fetchIngestLogs()
}

function onRefreshClick() {
    refreshLogs()
    showInfo('Refreshed', 'Ingest logs have been refreshed.')
}

// Dummy reprocess function — backend API not yet implemented
function reprocessLogs() {
    if (!selectedRecords.value.length) {
        showWarn('No Selection', 'Please select at least one record to re-process.')
        return
    }

    const msgIDs = selectedRecords.value.map(r => r.MessageID)
    console.log('[Dummy] Reprocessing MessageIDs:', msgIDs)
    showInfo(
        'Reprocessing Triggered',
        `Selected ${msgIDs.length} record(s): [${msgIDs.join(', ')}]. (Backend API not yet implemented.)`
    )
    selectedRecords.value = []
}

onMounted(() => {
    refreshLogs()
})
</script>

<template>
    <BlockUI :blocked="loading" fullScreen />
    <ProgressSpinner v-if="loading" class="global-spinner" />

    <div class="toolbar">
        <div class="toolbar-left">
            <Button label="Refresh" icon="pi pi-refresh" @click="onRefreshClick" />
            <Button label="Re-Process" icon="pi pi-replay" severity="warn" :disabled="!selectedRecords.length" @click="reprocessLogs" />
        </div>
    </div>

    <DataTable scrollHeight="flex" :value="records" stripedRows paginator :rows="rows" :rowsPerPageOptions="[5, 10, 20, 50]" :totalRecords="total" :first="first" lazy @page="onPage" tableStyle="min-width: 50rem" dataKey="MessageID" v-model:selection="selectedRecords" selectionMode="multiple">
        <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>
        <Column field="MessageID" header="Msg ID">
            <template #body="{ data }">
                <span class="msgid-cell" @click="onMsgIdClick($event, data)">{{ data.MessageID }}</span>
            </template>
        </Column>
        <Column field="Topic" header="Topic">
            <template #body="{ data }">
                <span class="topic-cell">{{ data.Topic }}</span>
            </template>
        </Column>
        <Column field="Payload" header="Payload">
            <template #body="{ data }">
                <span class="payload-cell" @click="openPayloadDialog(data.Payload, $event)">{{ data.Payload }}</span>
            </template>
        </Column>
        <Column field="DeviceID" header="Device ID"></Column>
        <Column field="IngestMethod" header="Ingest Method"></Column>
        <Column field="Status" header="Status">
            <template #body="{ data }">
                <Tag
                    :value="data.Status"
                    :severity="data.Status === 'processed' ? 'success' : 'danger'"
                />
            </template>
        </Column>
        <Column field="ReceivedAt" header="Received At"></Column>
    </DataTable>
   
    <Dialog v-model:visible="dialogVisible" header="Payload Details" :modal="true" :closable="false" :style="{ width: '600px' }">
        <div class="modal-content">
            <div class="modal-details">
                <pre>{{ selectedPayload }}</pre>
            </div>
        </div>
        <template #footer>
            <Button label="Close" icon="pi pi-times" @click="dialogVisible = false" class="p-button-text"/>
        </template>
    </Dialog>

    <ProcessedDataDialog
        v-model:visible="processedDialogVisible"
        :record="selectedRecord"
    />

</template>

<style scoped>


.toolbar {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
}

.toolbar-left {
    display: flex;
    gap: 0.5rem;
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

.topic-cell {
    display: block;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.msgid-cell {
    cursor: pointer;
    transition: color 0.15s, text-decoration 0.15s;
}

.msgid-cell:hover {
    color: #90caf9;
    text-decoration: underline;
}

.payload-cell {
    display: block;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    cursor: pointer;
    padding: 2px 6px;
    border-radius: 4px;
    background-color: #212121;
    transition: background-color 0.2s ease;
    font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', 'Courier New', monospace;
    font-size: 0.85rem;
}

.payload-cell:hover {
    background-color: rgba(72, 137, 123, 0.15);
}

.modal-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem 0;
}

.modal-details {
    margin-top: 1rem;
    padding: 1rem;
    border-radius: 4px;
    max-width: 100%;
    background-color: #212121;
    overflow-x: auto;
    width: 100%;
    box-sizing: border-box;
}

.modal-details pre {
    margin: 0;
    font-size: 0.85rem;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.global-spinner {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
}
</style>