import apiClient from './client'
import type { ApiResponse, Order, Bill, Ticket, TicketReply, APIKey, AuditLog, User, Company, Challenge } from '@/types'

// Company
export const getMembers = () =>
  apiClient.get<ApiResponse<User[]>>('/company/members')

export const addMember = (data: { username: string; password: string; email: string; role: string }) =>
  apiClient.post<ApiResponse>('/company/members', data)

export const updateMember = (id: number, data: { email?: string; phone?: string; status?: string }) =>
  apiClient.put<ApiResponse>(`/company/members/${id}`, data)

export const deleteMember = (id: number) =>
  apiClient.delete<ApiResponse>(`/company/members/${id}`)

export const changeRole = (id: number, role: string) =>
  apiClient.put<ApiResponse>(`/company/members/${id}/role`, { role })

// Orders
export const getOrders = () =>
  apiClient.get<ApiResponse<Order[]>>('/orders')

export const getOrderDetail = (id: number) =>
  apiClient.get<ApiResponse<Order>>(`/orders/${id}`)

// Bills
export const getBills = () =>
  apiClient.get<ApiResponse<Bill[]>>('/bills')

export const recharge = (amount: number) =>
  apiClient.post<ApiResponse>('/bills/recharge', { amount })

export const exportBills = () =>
  apiClient.get('/bills/export', { responseType: 'blob' })

// Tickets
export const getTickets = () =>
  apiClient.get<ApiResponse<Ticket[]>>('/tickets')

export const createTicket = (data: { title: string; content: string }) =>
  apiClient.post<ApiResponse>('/tickets', data)

export const getTicketDetail = (id: number) =>
  apiClient.get<ApiResponse<{ ticket: Ticket; replies: TicketReply[] }>>(`/tickets/${id}`)

export const replyTicket = (id: number, content: string) =>
  apiClient.post<ApiResponse>(`/tickets/${id}/reply`, { content })

export const closeTicket = (id: number) =>
  apiClient.put<ApiResponse>(`/tickets/${id}/close`)

// API Keys
export const getApiKeys = () =>
  apiClient.get<ApiResponse<APIKey[]>>('/apikeys')

export const createApiKey = (data: { name: string; permissions: string }) =>
  apiClient.post<ApiResponse>('/apikeys', data)

export const deleteApiKey = (id: number) =>
  apiClient.delete<ApiResponse>(`/apikeys/${id}`)

// Audit Logs
export const getAuditLogs = () =>
  apiClient.get<ApiResponse<AuditLog[]>>('/audit-logs')

// Admin
export const adminListUsers = () =>
  apiClient.get<ApiResponse<User[]>>('/admin/users')

export const adminUpdateUserStatus = (id: number, status: string) =>
  apiClient.put<ApiResponse>(`/admin/users/${id}/status`, { status })

export const adminListCompanies = () =>
  apiClient.get<ApiResponse<Company[]>>('/admin/companies')

export const adminListVps = () =>
  apiClient.get<ApiResponse>('/admin/vps')

export const adminResetDb = () =>
  apiClient.post<ApiResponse>('/admin/reset')

// Challenges
export const getChallenges = () =>
  apiClient.get<ApiResponse<Challenge[]>>('/challenges')

export const getChallengeDetail = (id: string) =>
  apiClient.get<ApiResponse>(`/challenges/${id}`)

export const getHint = (id: string, level: number) =>
  apiClient.get<ApiResponse<{ hint: string; level: number }>>(`/challenges/${id}/hints/${level}`)

export const markChallengeComplete = (id: string) =>
  apiClient.post<ApiResponse>(`/challenges/${id}/complete`)

// Security Mode
export const getSecurityMode = () =>
  apiClient.get<ApiResponse<{ mode: string }>>('/security-mode')

export const setSecurityMode = (mode: string) =>
  apiClient.put<ApiResponse>('/security-mode', { mode })
