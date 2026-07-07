// SWAP: Replace mock implementations with real API calls.
// import { apiFetch } from './index'

export async function login(email, password) {
  // Real API integration:
  // return apiFetch('/api/login', {
  //   method: 'POST',
  //   body: JSON.stringify({ email, password })
  // })

  await new Promise(r => setTimeout(r, 800 + Math.random() * 400))

  if (!email || !email.trim()) {
    throw new Error('Email or Username is required')
  }

  if (!password) {
    throw new Error('Password is required')
  }

  if (password.length < 8) {
    throw new Error('Password must be at least 8 characters')
  }

  return {
    user: {
      id: 1,
      name: email.includes('@') ? email.split('@')[0] : email,
      email: email.includes('@') ? email : `${email}@demo.local`,
      role: 'admin'
    },
    token: 'mock-jwt-' + Date.now() + '-' + Math.random().toString(36).slice(2)
  }
}

export async function logout() {
  // SWAP: Call backend logout endpoint if needed
  // return apiFetch('/api/logout', { method: 'POST' })
}
