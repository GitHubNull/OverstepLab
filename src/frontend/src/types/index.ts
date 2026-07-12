export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface User {
  id: number
  username: string
  email: string
  phone: string
  avatar: string
  user_type: string
  company_id: number | null
  role: string
  status: string
  created_at: string
  updated_at: string
  company?: Company
}

export interface Company {
  id: number
  name: string
  license_no: string
  contact_name: string
  contact_phone: string
  balance: number
  status: string
  created_at: string
  updated_at: string
}

export interface VPSInstance {
  id: number
  name: string
  owner_id: number
  company_id: number | null
  cpu: number
  memory: number
  disk: number
  bandwidth: number
  ip_address: string
  os_image: string
  status: string
  expire_at: string
  created_at: string
  updated_at: string
  owner?: User
  company?: Company
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  company_id: number | null
  vps_id: number | null
  type: string
  amount: number
  status: string
  created_at: string
  updated_at: string
}

export interface Bill {
  id: number
  company_id: number | null
  user_id: number
  type: string
  amount: number
  balance_after: number
  description: string
  created_at: string
}

export interface Ticket {
  id: number
  title: string
  content: string
  user_id: number
  company_id: number | null
  status: string
  created_at: string
  updated_at: string
  user?: User
}

export interface TicketReply {
  id: number
  ticket_id: number
  user_id: number
  content: string
  created_at: string
  user?: User
}

export interface APIKey {
  id: number
  user_id: number
  name: string
  key_prefix: string
  permissions: string
  status: string
  last_used_at: string | null
  expire_at: string | null
  created_at: string
}

export interface AuditLog {
  id: number
  user_id: number
  company_id: number | null
  action: string
  resource_type: string
  resource_id: number
  detail: string
  ip_address: string
  created_at: string
  user?: User
}

export interface Challenge {
  id: string
  title: string
  category: string
  difficulty: number
  description: string
  completed: boolean
  endpoint?: string
  method?: string
  hints?: string[]
  writeup?: string
}

export interface ChallengeDetail extends Challenge {
  endpoint: string
  method: string
  hints: string[]
  writeup: string
}

export interface Announcement {
  id: number
  title: string
  content: string
  user_id: number
  is_pinned: boolean
  created_at: string
  updated_at: string
  user?: User
}

export interface SystemConfig {
  key: string
  value: string
}

export interface APIKeyCreateResponse extends APIKey {
  key_value: string
}
