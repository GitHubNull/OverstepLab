import apiClient from './client'
import type { ApiResponse, VPSInstance } from '@/types'

export const getVpsList = () =>
  apiClient.get<ApiResponse<VPSInstance[]>>('/vps')

export const getVpsDetail = (id: number) =>
  apiClient.get<ApiResponse<VPSInstance>>(`/vps/${id}`)

export const createVps = (data: { name: string; cpu: number; memory: number; disk: number; bandwidth: number; os_image: string }) =>
  apiClient.post<ApiResponse>('/vps', data)

export const startVps = (id: number) =>
  apiClient.post<ApiResponse>(`/vps/${id}/start`)

export const stopVps = (id: number) =>
  apiClient.post<ApiResponse>(`/vps/${id}/stop`)

export const restartVps = (id: number) =>
  apiClient.post<ApiResponse>(`/vps/${id}/restart`)

export const reinstallVps = (id: number, osImage: string) =>
  apiClient.post<ApiResponse>(`/vps/${id}/reinstall`, { os_image: osImage })

export const deleteVps = (id: number) =>
  apiClient.delete<ApiResponse>(`/vps/${id}`)

export const getConsole = (id: number) =>
  apiClient.get<ApiResponse>(`/vps/${id}/console`)
