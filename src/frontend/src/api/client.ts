import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import {
  base64Encode,
  base32Encode,
  caesarEncode,
  customBase64Encode,
} from '@/utils/crypto'

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

// ---- Global encoding state (managed by backend, fetched on boot) ----
let globalEncodingType: string = 'none'
let globalEncodingChallengeId: string | null = null

// Paths excluded from encoding
const ENCODING_SKIP_PREFIXES = ['/auth/', '/crypto/', '/encoding-challenge-state']

function shouldSkipEncoding(url: string): boolean {
  return ENCODING_SKIP_PREFIXES.some(prefix => url.includes(prefix))
}

// Fetch encoding state from backend
async function fetchEncodingState(): Promise<void> {
  try {
    const res = await rawClient.get('/encoding-challenge-state')
    if (res.data?.data?.active) {
      globalEncodingType = res.data.data.encoding_type || 'none'
      globalEncodingChallengeId = res.data.data.challenge_id || null
    } else {
      globalEncodingType = 'none'
      globalEncodingChallengeId = null
    }
  } catch {
    globalEncodingType = 'none'
    globalEncodingChallengeId = null
  }
}

// Encode a single value
function encodeValue(value: string, type: string): string {
  switch (type) {
    case 'base64':
      return base64Encode(value)
    case 'base32':
      return base32Encode(value)
    case 'caesar':
      return caesarEncode(value, 3)
    case 'custom-base64':
      return customBase64Encode(value)
    default:
      return value
  }
}

// Recursively encode all string values in an object
function encodeObjectValues(obj: any, type: string): any {
  if (obj === null || obj === undefined) return obj
  if (typeof obj === 'string') {
    // Don't encode empty strings or JWT tokens
    if (obj === '' || obj.startsWith('Bearer ')) return obj
    return encodeValue(obj, type)
  }
  if (Array.isArray(obj)) {
    return obj.map(v => encodeObjectValues(v, type))
  }
  if (typeof obj === 'object') {
    const result: any = {}
    for (const [key, value] of Object.entries(obj)) {
      result[key] = encodeObjectValues(value, type)
    }
    return result
  }
  return obj
}

// ---- Request interceptor: token injection + encoding ----
apiClient.interceptors.request.use(async (config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  // Apply encoding if active
  if (globalEncodingType !== 'none' && globalEncodingType !== '') {
    const url = config.url || ''
    if (shouldSkipEncoding(url)) {
      return config
    }

    config.headers['X-Encoding-Type'] = globalEncodingType
    const method = config.method?.toUpperCase() || 'GET'

    if (method === 'GET' || method === 'HEAD') {
      if (config.params && typeof config.params === 'object') {
        config.params = encodeObjectValues(config.params, globalEncodingType)
      }
    } else {
      if (config.data !== undefined && config.data !== null) {
        config.data = encodeObjectValues(config.data, globalEncodingType)
      }
    }
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

// Export function to refresh encoding state
export function refreshEncodingState(): Promise<void> {
  return fetchEncodingState()
}

// Export current encoding state for components
export function getGlobalEncodingType(): string {
  return globalEncodingType
}

export function getGlobalEncodingChallengeId(): string | null {
  return globalEncodingChallengeId
}

export function isEncodingActive(): boolean {
  return globalEncodingType !== 'none' && globalEncodingType !== ''
}

// Initialize encoding state on load
fetchEncodingState()

export default apiClient
