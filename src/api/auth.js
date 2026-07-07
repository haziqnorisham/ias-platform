import { apiFetch } from './index'

export async function login(username, password) {
  return apiFetch('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password })
  })
}

export async function logout() {
  return apiFetch('/api/auth/logout', { method: 'POST' })
}

export async function checkSession() {
  return apiFetch('/api/auth/session', { method: 'POST' })
}
