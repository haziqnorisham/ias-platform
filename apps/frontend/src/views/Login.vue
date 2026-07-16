<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const { login, loading } = useAuth()

const username = ref('')
const password = ref('')

async function handleLogin() {
  if (!username.value || !username.value.trim()) {
    toast.add({ severity: 'error', summary: 'Validation Error', detail: 'Username is required', life: 3000 })
    return
  }

  if (!password.value) {
    toast.add({ severity: 'error', summary: 'Validation Error', detail: 'Password is required', life: 3000 })
    return
  }

  try {
    await login({
      username: username.value.trim(),
      password: password.value
    })

    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Login Failed', detail: err.message, life: 3000 })
  }
}
</script>

<template>
  <div class="login-wrapper">
    <Card class="login-card">
      <template #content>
        <div class="login-header">
          <img src="/bitmap.png" alt="IAS Logo" class="login-logo" />
          <h1 class="login-title">Sign In</h1>
          <p class="login-subtitle">Enter your credentials to access IAS Health Center</p>
        </div>

        <div class="login-form">
          <div class="field">
            <label for="username">Username</label>
            <InputText
              id="username"
              v-model="username"
              placeholder="Enter your username"
              class="form-input"
              @keyup.enter="handleLogin"
            />
          </div>

          <div class="field">
            <label for="password">Password</label>
            <Password
              id="password"
              v-model="password"
              placeholder="Enter your password"
              class="form-input"
              toggleMask
              :feedback="false"
              @keyup.enter="handleLogin"
            />
          </div>

          <Button
            label="Sign In"
            icon="pi pi-sign-in"
            class="login-button"
            :loading="loading"
            @click="handleLogin"
          />

          <div class="forgot-password">
            <a href="#">Forgot Password?</a>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.login-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #0e0e10;
  font-family: 'Space Grotesk', sans-serif;
}

.login-card {
  width: 420px;
  max-width: 90vw;
  border: 1px solid #212121;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-logo {
  width: 180px;
  margin-bottom: 1.5rem;
}

.login-title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  margin: 0 0 0.5rem;
  color: #e0e0e0;
}

.login-subtitle {
  font-size: var(--font-size-sm);
  color: #888;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
}

.field label {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: #a0a0a0;
  margin-bottom: 0.5rem;
}

.form-input {
  width: 100%;
}

:deep(.p-password) {
  width: 100%;
}

:deep(.p-password-input) {
  width: 100%;
}

.login-button {
  width: 100%;
  justify-content: center;
}

.forgot-password {
  text-align: center;
}

.forgot-password a {
  color: #888;
  font-size: 0.85rem;
  text-decoration: none;
  transition: color 0.2s;
}

.forgot-password a:hover {
  color: #48897b;
}
</style>
