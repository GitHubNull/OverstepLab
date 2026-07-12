import apiClient from './client'
import type { ApiResponse, VPSInstance, User, Order, Ticket, APIKey } from '@/types'

// ==================== Encoded VPS Endpoints ====================

export const getEncodedVps = (encodedId: string, encodingType: string) =>
  apiClient.get<ApiResponse<VPSInstance>>('/encoded/vps', {
    params: { v: encodedId },
    headers: { 'X-Encoding-Type': encodingType },
  })

export const startEncodedVps = (encodedId: string, encodingType: string) =>
  apiClient.post<ApiResponse>('/encoded/vps/start', { v: encodedId }, {
    headers: { 'X-Encoding-Type': encodingType },
  })

export const stopEncodedVps = (encodedId: string, encodingType: string) =>
  apiClient.post<ApiResponse>('/encoded/vps/stop', { v: encodedId }, {
    headers: { 'X-Encoding-Type': encodingType },
  })

export const reinstallEncodedVps = (encodedId: string, osImage: string, encodingType: string) =>
  apiClient.post<ApiResponse>('/encoded/vps/reinstall', { v: encodedId, os_image: osImage }, {
    headers: { 'X-Encoding-Type': encodingType },
  })

// ==================== Encoded User Endpoints ====================

export const getEncodedUser = (encodedId: string, encodingType: string) =>
  apiClient.get<ApiResponse<User>>('/encoded/users', {
    params: { v: encodedId },
    headers: { 'X-Encoding-Type': encodingType },
  })

// ==================== Encoded Order Endpoints ====================

export const getEncodedOrder = (encodedId: string, encodingType: string) =>
  apiClient.get<ApiResponse<Order>>('/encoded/orders', {
    params: { v: encodedId },
    headers: { 'X-Encoding-Type': encodingType },
  })

// ==================== Encoded Ticket Endpoints ====================

export const getEncodedTicket = (encodedId: string, encodingType: string) =>
  apiClient.get<ApiResponse<Ticket>>('/encoded/tickets', {
    params: { v: encodedId },
    headers: { 'X-Encoding-Type': encodingType },
  })

// ==================== Encoded API Key Endpoints ====================

export const deleteEncodedApiKey = (encodedId: string, encodingType: string) =>
  apiClient.delete<ApiResponse>('/encoded/apikeys', {
    data: { v: encodedId },
    headers: { 'X-Encoding-Type': encodingType },
  })

// ==================== Encoded Company Endpoints ====================

export const addEncodedMember = (data: { username: string; password: string; email: string; role: string }, encodingType: string) =>
  apiClient.post<ApiResponse>('/encoded/company/members', data, {
    headers: { 'X-Encoding-Type': encodingType },
  })

export const changeEncodedRole = (encodedId: string, role: string, encodingType: string) =>
  apiClient.put<ApiResponse>('/encoded/company/members/role', { v: encodedId, role }, {
    headers: { 'X-Encoding-Type': encodingType },
  })
