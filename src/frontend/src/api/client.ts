import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'

const apiClient: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const rawClient: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let isRefreshing = false
let refreshPromise: Promise<string> | null = null

async function doRefresh(): Promise<string> {
  const refreshTokenValue = localStorage.getItem('refresh_token')
  if (!refreshTokenValue) {
    throw new Error('No refresh token')
  }

  const response = await rawClient.post('/auth/refresh', {
    refresh_token: refreshTokenValue,
  })

  const { token, refresh_token } = response.data.data
  localStorage.setItem('token', token)
  localStorage.setItem('refresh_token', refresh_token)
  return token
}

function refreshToken(): Promise<string> {
  if (!isRefreshing) {
    isRefreshing = true
    refreshPromise = doRefresh()
      .finally(() => {
        isRefreshing = false
        refreshPromise = null
      })
  }
  if (!refreshPromise) {
    refreshPromise = Promise.reject(new Error('Refresh failed'))
  }
  return refreshPromise
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const newToken = await refreshToken()
        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return apiClient(originalRequest)
      } catch (refreshError) {
        localStorage.removeItem('token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user')
        const publicPaths = ['/login', '/register']
        if (!publicPaths.includes(window.location.pathname)) {
          window.location.href = '/login'
        }
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient
