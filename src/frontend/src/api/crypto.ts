import apiClient from './client'

export function cryptoEncode(type: string, value: string) {
  return apiClient.get(`/api/v1/crypto/${type}/encode`, { params: { value } })
}

export function cryptoDecode(type: string, value: string) {
  return apiClient.get(`/api/v1/crypto/${type}/decode`, { params: { value } })
}

export function getCryptoKeys() {
  return apiClient.get('/api/v1/crypto/keys')
}
