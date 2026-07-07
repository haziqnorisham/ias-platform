import { ref, computed } from 'vue'
import { login as apiLogin, logout as apiLogout } from '../api/auth'

const user = ref(null)
const token = ref(null)
const isAuthenticated = ref(false)
const loading = ref(false)
const error = ref(null)

const STORAGE_KEY = 'ias_auth'

function initFromStorage() {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (!saved) return
  try {
    const data = JSON.parse(saved)
    token.value = data.token
    user.value = data.user
    isAuthenticated.value = true
  } catch {
    localStorage.removeItem(STORAGE_KEY)
  }
}

initFromStorage()

export function useAuth() {
  function checkAuth() {
    initFromStorage()
  }

  async function login(credentials, rememberMe = false) {
    loading.value = true
    error.value = null
    try {
      const data = await apiLogin(credentials.email, credentials.password)
      user.value = data.user
      token.value = data.token
      isAuthenticated.value = true

      if (rememberMe) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify({
          user: data.user,
          token: data.token
        }))
      }
      return data
    } catch (err) {
      error.value = err.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    isAuthenticated.value = false
    error.value = null
    localStorage.removeItem(STORAGE_KEY)
    apiLogout()
  }

  function clearError() {
    error.value = null
  }

  const currentUser = computed(() => user.value)

  return {
    user,
    token,
    isAuthenticated,
    loading,
    error,
    currentUser,
    login,
    logout,
    checkAuth,
    clearError
  }
}
