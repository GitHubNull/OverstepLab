import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type EncodingType = 'none' | 'base64' | 'base32' | 'caesar' | 'custom-base64' | 'multi' | 'aes' | 'hmac' | 'sm4' | 'hash-sign'

interface EncodingChallenge {
  id: string
  encodingType: EncodingType
}

export const useEncodingChallengeStore = defineStore('encodingChallenge', () => {
  const activeEncodingType = ref<EncodingType>('none')
  const activeChallengeId = ref<string | null>(null)

  // Load from localStorage on init
  const savedType = localStorage.getItem('activeEncodingType') as EncodingType | null
  const savedId = localStorage.getItem('activeChallengeId')
  if (savedType && savedType !== 'none') {
    activeEncodingType.value = savedType
    activeChallengeId.value = savedId
  }

  function activate(challenge: EncodingChallenge) {
    activeEncodingType.value = challenge.encodingType
    activeChallengeId.value = challenge.id
    localStorage.setItem('activeEncodingType', challenge.encodingType)
    localStorage.setItem('activeChallengeId', challenge.id)
  }

  function deactivate() {
    activeEncodingType.value = 'none'
    activeChallengeId.value = null
    localStorage.removeItem('activeEncodingType')
    localStorage.removeItem('activeChallengeId')
  }

  const isActive = computed(() => activeEncodingType.value !== 'none')

  function getActiveEncodingType(): EncodingType {
    return activeEncodingType.value
  }

  function getActiveEncodingChallenge(): EncodingChallenge | null {
    if (!activeChallengeId.value) return null
    return {
      id: activeChallengeId.value,
      encodingType: activeEncodingType.value,
    }
  }

  return {
    activeEncodingType,
    activeChallengeId,
    activate,
    deactivate,
    isActive,
    getActiveEncodingType,
    getActiveEncodingChallenge,
  }
})
