import { ref, computed } from 'vue'
import { login as apiLogin, logout as apiLogout, checkSession } from '../api/auth'

const user = ref(null)
const isAuthenticated = ref(false)
const loading = ref(false)
const error = ref(null)

let initPromise = null

async function doInit() {
  try {
    const data = await checkSession()
    user.value = data
    isAuthenticated.value = true
  } catch {
    user.value = null
    isAuthenticated.value = false
  }
}

export function useAuth() {
  function checkAuth() {
    if (!initPromise) {
      initPromise = doInit()
    }
    return initPromise
  }

  function waitForInit() {
    return initPromise || Promise.resolve()
  }

  async function login(credentials) {
    loading.value = true
    error.value = null
    try {
      const data = await apiLogin(credentials.username, credentials.password)
      user.value = data
      isAuthenticated.value = true
      return data
    } catch (err) {
      error.value = err.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await apiLogout()
    } finally {
      user.value = null
      isAuthenticated.value = false
      error.value = null
    }
  }

  function clearError() {
    error.value = null
  }

  const currentUser = computed(() => user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  return {
    user,
    isAuthenticated,
    loading,
    error,
    currentUser,
    isAdmin,
    login,
    logout,
    checkAuth,
    waitForInit,
    clearError
  }
}

initPromise = doInit()
