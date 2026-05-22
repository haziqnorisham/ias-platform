<script setup>
import { ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'

const props = defineProps({
  visible: Boolean,
  stream: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'save', 'request-delete'])

const form = ref({
  name: '',
  ip: '',
  port: 80,
  username: '',
  password: '',
  notes: ''
})

const isEditing = ref(false)

watch(() => props.visible, (val) => {
  if (val) {
    if (props.stream) {
      isEditing.value = true
      form.value = {
        name: props.stream.name || '',
        ip: props.stream.ip || '',
        port: props.stream.port || 80,
        username: props.stream.username || '',
        password: '',
        notes: props.stream.notes || ''
      }
    } else {
      isEditing.value = false
      form.value = { name: '', ip: '', port: 80, username: '', password: '', notes: '' }
    }
  }
})

function handleSave() {
  const payload = { ...form.value }
  if (isEditing.value && !payload.password) {
    delete payload.password
  }
  emit('save', payload)
}

function close() {
  emit('update:visible', false)
}
</script>

<template>
  <Dialog
    :visible="visible"
    :header="isEditing ? 'Edit Stream' : 'Add Stream'"
    :modal="true"
    :style="{ width: '480px' }"
    @update:visible="close"
  >
    <div class="form-grid">
      <div class="field">
        <label>Name *</label>
        <InputText v-model="form.name" placeholder="e.g. Front Gate" class="w-full" />
      </div>
      <div class="field">
        <label>IP Address *</label>
        <InputText v-model="form.ip" placeholder="e.g. 192.168.1.100" class="w-full" />
      </div>
      <div class="field">
        <label>Port</label>
        <InputText v-model.number="form.port" placeholder="80" type="number" class="w-full" />
      </div>
      <div class="field">
        <label>Username *</label>
        <InputText v-model="form.username" placeholder="ONVIF username" class="w-full" />
      </div>
      <div class="field">
        <label>{{ isEditing ? 'Password (leave blank to keep)' : 'Password *' }}</label>
        <InputText v-model="form.password" type="password" placeholder="ONVIF password" class="w-full" />
      </div>
      <div class="field full-width">
        <label>Notes</label>
        <Textarea v-model="form.notes" placeholder="Optional notes..." class="w-full" rows="3" />
      </div>
    </div>
    <template #footer>
      <div class="footer-row">
        <Button
          v-if="isEditing"
          label="Delete"
          icon="pi pi-trash"
          severity="danger"
          outlined
          @click="emit('request-delete')"
        />
        <div class="footer-spacer" />
        <Button label="Cancel" icon="pi pi-times" severity="secondary" @click="close" />
        <Button label="Save" icon="pi pi-check" @click="handleSave" />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.form-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.field {
  flex: 1 1 calc(50% - 0.5rem);
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.field label {
  font-size: 0.8rem;
  color: #aaa;
}

.field.full-width {
  flex: 1 1 100%;
}

.w-full {
  width: 100%;
}

.footer-row {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.5rem;
}

.footer-spacer {
  flex: 1;
}
</style>
