import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile, logout as apiLogout } from '@/api/auth'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const storedUser = localStorage.getItem('user')
  const user = ref<User | null>(storedUser ? JSON.parse(storedUser) : null)
  const securityMode = ref<string>('vulnerable')

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.user_type === 'platform_admin')
  const isCompanyAdmin = computed(() => user.value?.role === 'admin')
  const isViewer = computed(() => user.value?.role === 'viewer')
  const isOperator = computed(() => user.value?.role === 'operator')
  const isFinance = computed(() => user.value?.role === 'finance')
  const isIndividual = computed(() => user.value?.user_type === 'individual')

  async function login(username: string, password: string) {
    const response = await apiLogin(username, password)
    const loginData = response.data.data!
    token.value = loginData.token
    user.value = loginData.user
    localStorage.setItem('token', loginData.token)
    localStorage.setItem('user', JSON.stringify(loginData.user))
  }

  async function register(formData: { username: string; password: string; email: string; user_type: string }) {
    await apiRegister(formData)
  }

  async function fetchProfile() {
    const response = await getProfile()
    user.value = response.data.data!
    localStorage.setItem('user', JSON.stringify(response.data.data))
  }

  async function logout() {
    try {
      await apiLogout()
    } catch {}
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function toggleSecurityMode() {
    const newMode = securityMode.value === 'vulnerable' ? 'secure' : 'vulnerable'
    const { getSecurityMode: getMode, setSecurityMode: setMode } = await import('@/api')
    await setMode(newMode)
    securityMode.value = newMode
  }

  function init() {
    const storedUser = localStorage.getItem('user')
    const storedToken = localStorage.getItem('token')
    if (storedUser && storedToken) {
      try {
        user.value = JSON.parse(storedUser)
        token.value = storedToken
      } catch {
        logout()
      }
    }
  }

  return {
    token, user, securityMode, isAuthenticated, isAdmin, isCompanyAdmin,
    isViewer, isOperator, isFinance, isIndividual,
    login, register, fetchProfile, logout, toggleSecurityMode, init,
  }
})
