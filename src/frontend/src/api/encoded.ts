import { apiClient } from './client'

export function encodeValue(type: string, value: string) {
  return apiClient.get(`/api/v1/encoded/${type}/encode`, { params: { value } })
}

export function decodeValue(type: string, value: string) {
  return apiClient.get(`/api/v1/encoded/${type}/decode`, { params: { value } })
}

export function getEncodedKeys() {
  return apiClient.get('/api/v1/encoded/keys')
}
