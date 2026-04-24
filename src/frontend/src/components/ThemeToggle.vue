<template>
  <div class="relative">
    <Button
      icon="pi pi-palette"
      text
      rounded
      size="small"
      class="text-[var(--text-secondary)]"
      @click="toggleMenu"
      aria-label="切换主题"
    />
    <Menu
      ref="menu"
      :model="items"
      popup
      :pt="{
        root: { class: 'min-w-[140px]' }
      }"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from 'primevue/button'
import Menu from 'primevue/menu'
import { useTheme, type ThemeMode } from '@/composables/useTheme'

const { theme, setTheme } = useTheme()
const menu = ref<InstanceType<typeof Menu>>()

const items = computed(() => [
  {
    label: '浅色',
    icon: theme.value === 'light' ? 'pi pi-check' : 'pi pi-sun',
    command: () => setTheme('light'),
    class: theme.value === 'light' ? 'font-semibold' : '',
  },
  {
    label: '深色',
    icon: theme.value === 'dark' ? 'pi pi-check' : 'pi pi-moon',
    command: () => setTheme('dark'),
    class: theme.value === 'dark' ? 'font-semibold' : '',
  },
  {
    label: '跟随系统',
    icon: theme.value === 'system' ? 'pi pi-check' : 'pi pi-desktop',
    command: () => setTheme('system'),
    class: theme.value === 'system' ? 'font-semibold' : '',
  },
])

function toggleMenu(event: Event) {
  menu.value?.toggle(event)
}
</script>
