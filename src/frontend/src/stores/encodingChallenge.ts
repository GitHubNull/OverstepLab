import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type EncodingChallengeId = 'E-01' | 'E-02' | 'E-03' | 'E-04' | 'E-05' | 'E-06' | 'E-07' | 'E-08' | 'E-09' | null

interface EncodingChallenge {
  id: EncodingChallengeId
  name: string
  encodingType: string
  description: string
  globalEncoding: boolean
  targetPage?: string
  targetPageLabel?: string
}

export const encodingChallenges: EncodingChallenge[] = [
  {
    id: 'E-01',
    name: 'Base64 编码 VPS ID 越权',
    encodingType: 'base64',
    description: 'VPS 详情查询参数使用 Base64 编码传输',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
  {
    id: 'E-02',
    name: 'Base32 编码用户 ID 越权',
    encodingType: 'base32',
    description: '用户资料查询参数使用 Base32 编码传输',
    globalEncoding: true,
    targetPage: '/profile',
    targetPageLabel: '个人资料',
  },
  {
    id: 'E-03',
    name: '凯撒密码订单 ID 越权',
    encodingType: 'caesar',
    description: '订单详情查询参数使用凯撒密码编码传输',
    globalEncoding: true,
    targetPage: '/orders',
    targetPageLabel: '订单管理',
  },
  {
    id: 'E-04',
    name: '自定义 Base64 编码表绕过',
    encodingType: 'custom_base64',
    description: 'VPS 详情查询参数使用自定义 Base64 编码表传输',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
  {
    id: 'E-05',
    name: '多层嵌套编码工单越权',
    encodingType: 'multi',
    description: '工单详情查询参数使用 Base64→Base32 双层编码传输',
    globalEncoding: true,
    targetPage: '/tickets',
    targetPageLabel: '工单系统',
  },
  {
    id: 'E-06',
    name: 'AES 加密 VPS ID 操作',
    encodingType: 'aes',
    description: 'VPS 启动操作参数使用 AES-256-GCM 加密传输',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
  {
    id: 'E-07',
    name: 'HMAC 签名验证绕过',
    encodingType: 'signed',
    description: 'VPS 停止操作参数附带 HMAC-SHA256 签名保护',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
  {
    id: 'E-08',
    name: 'SM4 国密加密跨企业操作',
    encodingType: 'sm4',
    description: 'VPS 重装操作参数使用 SM4-CBC 加密传输',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
  {
    id: 'E-09',
    name: '简单哈希签名绕过',
    encodingType: 'hash',
    description: '所有请求参数附带 MD5 哈希签名，后端校验参数完整性',
    globalEncoding: true,
    targetPage: '/vps',
    targetPageLabel: 'VPS 管理',
  },
]

// Module-level state for interceptor access (outside Pinia reactivity)
// Initialize from localStorage if available
const STORAGE_KEY = 'active_encoding_challenge'
let _activeChallengeId: EncodingChallengeId = (() => {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored && stored.match(/^E-0[1-9]$/)) {
      return stored as EncodingChallengeId
    }
  }
  return null
})()

export function getActiveEncodingChallenge(): EncodingChallengeId {
  return _activeChallengeId
}

export function getActiveEncodingType(): string | null {
  if (!_activeChallengeId) return null
  const ch = encodingChallenges.find(c => c.id === _activeChallengeId)
  return ch?.encodingType || null
}

export function setActiveEncodingChallenge(id: EncodingChallengeId) {
  _activeChallengeId = id
  if (typeof window !== 'undefined') {
    if (id) {
      localStorage.setItem(STORAGE_KEY, id)
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }
}

export const useEncodingChallengeStore = defineStore('encodingChallenge', () => {
  const activeChallengeId = ref<EncodingChallengeId>(_activeChallengeId)

  const activeChallenge = computed(() => {
    if (!activeChallengeId.value) return null
    return encodingChallenges.find(ch => ch.id === activeChallengeId.value) || null
  })

  const isActive = computed(() => activeChallengeId.value !== null)
  const encodingType = computed(() => activeChallenge.value?.encodingType || null)
  const isGlobal = computed(() => activeChallenge.value?.globalEncoding ?? false)

  function activateChallenge(id: EncodingChallengeId) {
    _activeChallengeId = id
    activeChallengeId.value = id
    if (typeof window !== 'undefined') {
      if (id) {
        localStorage.setItem(STORAGE_KEY, id)
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
    }
  }

  function deactivateChallenge() {
    _activeChallengeId = null
    activeChallengeId.value = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function toggleChallenge(id: EncodingChallengeId) {
    if (_activeChallengeId === id) {
      _activeChallengeId = null
      activeChallengeId.value = null
      if (typeof window !== 'undefined') {
        localStorage.removeItem(STORAGE_KEY)
      }
    } else if (id) {
      _activeChallengeId = id
      activeChallengeId.value = id
      if (typeof window !== 'undefined') {
        localStorage.setItem(STORAGE_KEY, id)
      }
    }
  }

  return {
    activeChallengeId,
    activeChallenge,
    isActive,
    isGlobal,
    encodingType,
    activateChallenge,
    deactivateChallenge,
    toggleChallenge,
  }
})
