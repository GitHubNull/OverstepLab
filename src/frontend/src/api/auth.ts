import apiClient from './client'
import type { ApiResponse, User } from '@/types'

export const login = (username: string, password: string) =>
  apiClient.post<ApiResponse<{ token: string; refresh_token: string; user: User }>>('/auth/login', { username, password })

export const register = (data: { username: string; password: string; email: string; user_type: string; company_name?: string }) =>
  apiClient.post<ApiResponse>('/auth/register', data)

export const refreshToken = (refreshToken: string) =>
  apiClient.post<ApiResponse<{ token: string; refresh_token: string }>>('/auth/refresh', { refresh_token: refreshToken })

export const logout = () =>
  apiClient.post<ApiResponse>('/logout')

export const getProfile = () =>
  apiClient.get<ApiResponse<User>>('/user/profile')

export const updateProfile = (data: { email?: string; phone?: string }) =>
  apiClient.put<ApiResponse>('/user/profile', data)

export const changePassword = (data: { old_password: string; new_password: string }) =>
  apiClient.put<ApiResponse>('/user/password', data)

export const getUserById = (id: number) =>
  apiClient.get<ApiResponse<User>>(`/users/${id}`)
