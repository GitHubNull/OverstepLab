import { ref, computed, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

const STORAGE_KEY = 'oversteplab-theme'

const theme = ref<ThemeMode>('system')
const systemDark = ref(false)

const isDark = computed(() => {
  if (theme.value === 'system') {
    return systemDark.value
  }
  return theme.value === 'dark'
})

function updateDOM() {
  const root = document.documentElement
  if (isDark.value) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

function init() {
  // Load saved preference
  const saved = localStorage.getItem(STORAGE_KEY) as ThemeMode | null
  if (saved && ['light', 'dark', 'system'].includes(saved)) {
    theme.value = saved
  }

  // Listen to system preference
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  systemDark.value = mediaQuery.matches

  mediaQuery.addEventListener('change', (e) => {
    systemDark.value = e.matches
    updateDOM()
  })

  updateDOM()
}

function setTheme(mode: ThemeMode) {
  theme.value = mode
  localStorage.setItem(STORAGE_KEY, mode)
  updateDOM()
}

watch(isDark, updateDOM, { immediate: false })

export function useTheme() {
  return {
    theme,
    isDark,
    init,
    setTheme,
  }
}
