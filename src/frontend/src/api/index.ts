import apiClient, { rawClient } from './client'
import type { ApiResponse, Order, Bill, Ticket, TicketReply, APIKey, AuditLog, User, Company, Challenge, ChallengeDetail, Announcement } from '@/types'

// Company
export const getMembers = () =>
  apiClient.get<ApiResponse<User[]>>('/company/members')

export const addMember = (data: { username: string; password: string; email: string; role: string }) =>
  apiClient.post<ApiResponse>('/company/members', data)

export const updateMember = (id: number, data: { email?: string; phone?: string; status?: string }) =>
  apiClient.put<ApiResponse>('/company/members', { id, ...data })

export const deleteMember = (id: number) =>
  apiClient.delete<ApiResponse>('/company/members', { data: { id } })

export const changeRole = (id: number, role: string) =>
  apiClient.put<ApiResponse>('/company/members/role', { id, role })

// Orders
export const getOrders = () =>
  apiClient.get<ApiResponse<Order[]>>('/orders')

export const getOrderDetail = (id: number) =>
  apiClient.get<ApiResponse<Order>>('/orders/detail', { params: { orderId: id } })

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
  apiClient.get<ApiResponse<{ ticket: Ticket; replies: TicketReply[] }>>('/tickets/detail', { params: { ticketId: id } })

export const replyTicket = (id: number, content: string) =>
  apiClient.post<ApiResponse>('/tickets/reply', { ticketId: id, content })

export const closeTicket = (id: number) =>
  apiClient.put<ApiResponse>('/tickets/close', { ticketId: id })

// API Keys
export const getApiKeys = () =>
  apiClient.get<ApiResponse<APIKey[]>>('/apikeys')

export const createApiKey = (data: { name: string; permissions: string }) =>
  apiClient.post<ApiResponse>('/apikeys', data)

export const deleteApiKey = (id: number) =>
  apiClient.delete<ApiResponse>('/apikeys', { data: { id } })

// Audit Logs
export const getAuditLogs = () =>
  apiClient.get<ApiResponse<AuditLog[]>>('/audit-logs')

// Admin
export const adminListUsers = () =>
  apiClient.get<ApiResponse<User[]>>('/admin/users')

export const adminUpdateUserStatus = (id: number, status: string) =>
  apiClient.put<ApiResponse>('/admin/users/status', { id, status })

export const adminResetUserPassword = (id: number, password: string) =>
  apiClient.put<ApiResponse>('/admin/users/password', { id, password })

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
  apiClient.get<ApiResponse<ChallengeDetail>>('/challenges/detail', { params: { id } })

export const getHint = (id: string, level: number) =>
  apiClient.get<ApiResponse<{ hint: string; level: number }>>('/challenges/hints', { params: { id, level } })

export const markChallengeComplete = (id: string) =>
  apiClient.post<ApiResponse>('/challenges/complete', { id })

// Security Mode
export const getSecurityMode = () =>
  apiClient.get<ApiResponse<{ mode: string }>>('/security-mode')

export const setSecurityMode = (mode: string) =>
  apiClient.put<ApiResponse>('/security-mode', { mode })

// Announcements
export const getAnnouncements = () =>
  apiClient.get<ApiResponse<Announcement[]>>('/announcements')

export const adminCreateAnnouncement = (data: { title: string; content: string; is_pinned: boolean }) =>
  apiClient.post<ApiResponse<Announcement>>('/admin/announcements', data)

export const adminUpdateAnnouncement = (id: number, data: { title: string; content: string; is_pinned: boolean }) =>
  apiClient.put<ApiResponse>('/admin/announcements', { id, ...data })

export const adminDeleteAnnouncement = (id: number) =>
  apiClient.delete<ApiResponse>('/admin/announcements', { data: { id } })

export { rawClient }

// System Config
export const getSystemConfig = () =>
  apiClient.get<ApiResponse<Record<string, string>>>('/admin/config')

export const updateSystemConfig = (data: { key: string; value: string }) =>
  apiClient.put<ApiResponse>('/admin/config', data)

// Encoding Challenge State (admin managed, persistent)
export const getEncodingChallengeState = () =>
  apiClient.get<ApiResponse<{ active: boolean; challenge_id: string | null; encoding_type: string; challenge_name: string | null }>>('/encoding-challenge-state')

export const setEncodingChallengeState = (data: { challenge_id: string; encoding_type: string; challenge_name: string; active: boolean }) =>
  apiClient.put<ApiResponse>('/admin/encoding-challenge-state', data)
