import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import { getActiveEncodingType } from '@/stores/encodingChallenge'
import {
  base64Encode,
  base32Encode,
  caesarEncode,
  customBase64Encode,
  multiEncode,
} from '@/utils/crypto'
import * as cryptoApi from './crypto'
import CryptoJS from 'crypto-js'

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

// ---- Encoding helpers ----

// Paths excluded from encoding (auth endpoints, crypto utility endpoints)
const ENCODING_SKIP_PREFIXES = ['/auth/', '/crypto/']

function shouldSkipEncoding(url: string): boolean {
  return ENCODING_SKIP_PREFIXES.some(prefix => url.startsWith(prefix))
}

function encodeValue(value: string, type: string): string {
  switch (type) {
    case 'base64':
      return base64Encode(value)
    case 'base32':
      return base32Encode(value)
    case 'caesar':
      return caesarEncode(value, 3)
    case 'custom_base64':
      return customBase64Encode(value)
    case 'multi':
      return multiEncode(value)
    default:
      return value
  }
}

async function getEncodedValueAsync(value: string, type: string): Promise<string> {
  if (type === 'aes' || type === 'signed' || type === 'sm4' || type === 'hash') {
    const response = await cryptoApi.cryptoEncode(value, type)
    return response.data.data!.encoded
  }
  return encodeValue(value, type)
}

// Recursively encode all string and number values in an arbitrary object/array
async function encodeObjectValues(obj: unknown, encType: string): Promise<unknown> {
  if (obj === null || obj === undefined) {
    return obj
  }
  if (typeof obj === 'string') {
    return getEncodedValueAsync(obj, encType)
  }
  if (typeof obj === 'number') {
    const encoded = await getEncodedValueAsync(String(obj), encType)
    // Try to preserve numeric type when encoding only affects digits
    const parsed = Number(encoded)
    if (!isNaN(parsed) && String(parsed) === encoded) {
      return parsed
    }
    return encoded
  }
  if (Array.isArray(obj)) {
    return Promise.all(obj.map(v => encodeObjectValues(v, encType)))
  }
  if (typeof obj === 'object') {
    const result: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(obj as Record<string, unknown>)) {
      // Skip axios internal config keys
      if (key === 'signal' || key === 'validateStatus' || key.startsWith('_')) {
        result[key] = value
      } else {
        result[key] = await encodeObjectValues(value, encType)
      }
    }
    return result
  }
  return obj
}

// ---- Hash Sign helper (for E-09 challenge) ----
// Real-world style: JSON body as string → MD5+salt → X-Hash-Sign header
// GET: params sorted as key=value&key2=value2 → MD5+salt → X-Hash-Sign header
const HASH_SALT = 'OverstepLabHashSalt2024'

// ---- HMAC Sign helper (for E-07 challenge) ----
// Real-world style: same as hash but using HMAC-SHA256 with a secret key
const HMAC_KEY = 'OverstepLabHMACSecretKey!@#2024'

function serializeForSign(data: unknown, isJsonBody: boolean): string {
  if (data === null || data === undefined) return ''
  if (isJsonBody) {
    // POST/PUT: JSON body as raw string (same as axios will send)
    return JSON.stringify(data)
  }
  // GET: params sorted as key=value&key2=value2
  if (typeof data !== 'object') return String(data)
  const entries: string[] = []
  for (const [key, value] of Object.entries(data as Record<string, unknown>)) {
    if (key === 'signal' || key === 'validateStatus' || key.startsWith('_')) continue
    const valStr = typeof value === 'object' ? JSON.stringify(value) : String(value)
    entries.push(`${key}=${valStr}`)
  }
  return entries.sort().join('&')
}

function computeHashSign(data: unknown, isJsonBody: boolean): string {
  const payload = serializeForSign(data, isJsonBody) + HASH_SALT
  return CryptoJS.MD5(payload).toString()
}

function computeHmacSign(data: unknown, isJsonBody: boolean): string {
  const payload = serializeForSign(data, isJsonBody)
  return CryptoJS.HmacSHA256(payload, HMAC_KEY).toString()
}

// ---- Request interceptor: universal encoding for all requests ----
apiClient.interceptors.request.use(async (config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  // Apply universal encoding when a challenge is active
  const encodingType = getActiveEncodingType()
  if (!encodingType) {
    return config
  }

  const url = config.url || ''
  if (shouldSkipEncoding(url)) {
    return config
  }

  config.headers['X-Encoding-Type'] = encodingType
  const method = config.method?.toUpperCase() || 'GET'

  // Special handling for hash type (E-09): real-world header signature
  if (encodingType === 'hash') {
    const isGet = method === 'GET' || method === 'HEAD'
    const dataToSign = isGet ? config.params : config.data
    if (dataToSign !== null && dataToSign !== undefined) {
      config.headers['X-Hash-Sign'] = computeHashSign(dataToSign, !isGet)
    }
    return config
  }

  // Special handling for signed type (E-07): real-world HMAC header signature
  if (encodingType === 'signed') {
    const isGet = method === 'GET' || method === 'HEAD'
    const dataToSign = isGet ? config.params : config.data
    if (dataToSign !== null && dataToSign !== undefined) {
      config.headers['X-HMAC-Sign'] = computeHmacSign(dataToSign, !isGet)
    }
    return config
  }

  if (method === 'GET' || method === 'HEAD') {
    // Encode all query param values
    if (config.params && typeof config.params === 'object') {
      config.params = await encodeObjectValues(config.params, encodingType)
    }
  } else {
    // Encode all string values in request body
    if (config.data !== undefined && config.data !== null) {
      config.data = await encodeObjectValues(config.data, encodingType)
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

export default apiClient
