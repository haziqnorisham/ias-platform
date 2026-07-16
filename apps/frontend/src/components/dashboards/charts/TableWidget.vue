<template>
  <div class="table-widget">
    <DataTable :value="tableData" size="small" showGridlines class="inner-table">
      <Column v-for="col in columns" :key="col.field" :field="col.field" :header="col.header"></Column>
    </DataTable>
  </div>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

defineProps({
  title: { type: String, default: 'Table' }
})

const columns = [
  { field: 'name', header: 'Name' },
  { field: 'value', header: 'Value' },
  { field: 'status', header: 'Status' }
]

const tableData = [
  { name: 'Item A', value: '42', status: 'Active' },
  { name: 'Item B', value: '18', status: 'Inactive' },
  { name: 'Item C', value: '76', status: 'Active' }
]
</script>

<style scoped>
.table-widget {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.inner-table :deep(.p-datatable-thead > tr > th) {
  font-size: 0.7rem;
  padding: 0.25rem 0.5rem;
}

.inner-table :deep(.p-datatable-tbody > tr > td) {
  font-size: 0.72rem;
  padding: 0.2rem 0.5rem;
}
</style>
