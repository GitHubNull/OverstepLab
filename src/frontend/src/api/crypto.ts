import apiClient from './client'
import type { ApiResponse } from '@/types'

export const cryptoEncode = (value: string, encodingType: string) =>
  apiClient.post<ApiResponse<{ encoding_type: string; original: string; encoded: string }>>('/crypto/encode', { value, encoding_type: encodingType })

export const cryptoDecode = (value: string, encodingType: string) =>
  apiClient.post<ApiResponse<{ encoding_type: string; encoded: string; decoded: string }>>('/crypto/decode', { value, encoding_type: encodingType })

export const getCryptoKeys = () =>
  apiClient.get<ApiResponse<{ aes_key_base64: string; hmac_key_base64: string; sm4_key_base64: string; rsa_public_key: string; sm2_public_key: string }>>('/crypto/keys')
