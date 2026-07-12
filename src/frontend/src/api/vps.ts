import apiClient from './client'
import type { ApiResponse, VPSInstance } from '@/types'

export const getVpsList = () =>
  apiClient.get<ApiResponse<VPSInstance[]>>('/vps')

export const getVpsDetail = (id: number) =>
  apiClient.get<ApiResponse<VPSInstance>>('/vps/detail', { params: { vpsId: id } })

export const createVps = (data: { name: string; cpu: number; memory: number; disk: number; bandwidth: number; os_image: string }) =>
  apiClient.post<ApiResponse>('/vps', data)

export const startVps = (id: number) =>
  apiClient.post<ApiResponse>('/vps/start', { vpsId: id })

export const stopVps = (id: number) =>
  apiClient.post<ApiResponse>('/vps/stop', { vpsId: id })

export const restartVps = (id: number) =>
  apiClient.post<ApiResponse>('/vps/restart', { vpsId: id })

export const reinstallVps = (id: number, osImage: string) =>
  apiClient.post<ApiResponse>('/vps/reinstall', { vpsId: id, os_image: osImage })

export const deleteVps = (id: number) =>
  apiClient.delete<ApiResponse>('/vps', { data: { vpsId: id } })

export const getConsole = (id: number) =>
  apiClient.get<ApiResponse<{ view_url: string; expires_at: number }>>('/vps/console', { params: { vpsId: id } })
