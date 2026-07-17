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
    <div class="glow-layer" aria-hidden="true">
      <div class="glow-blob glow-blob--primary"></div>
      <div class="glow-blob glow-blob--secondary"></div>
    </div>

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
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  overflow-x: hidden;
  font-family: 'Space Grotesk', sans-serif;

  /* 1. Dotted grid background */
  background-color: #0e0e10;
  background-image: radial-gradient(#28282c 1px, transparent 1px);
  background-size: 24px 24px;
}

/* 2. Atmospheric glow blobs layer */
.glow-layer {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
}

.glow-blob {
  position: absolute;
  width: 24rem;
  height: 24rem;
  border-radius: 50%;
  filter: blur(120px);
  opacity: 0.18;
}

.glow-blob--primary {
  top: 25%;
  left: 25%;
  background-color: #48897b;
}

.glow-blob--secondary {
  bottom: 25%;
  right: 25%;
  background-color: #2d5a4e;
}

/* 3. Login card — renders above the blobs */
.login-card {
  position: relative;
  z-index: 1;
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
