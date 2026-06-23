<template>
  <router-view />
  <Toast />
  <ConfirmDialog />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import { useTheme } from '@/composables/useTheme'
import { useAuthStore } from '@/stores/auth'

const { init } = useTheme()
const authStore = useAuthStore()

onMounted(async () => {
  init()
  await authStore.init()
})
</script>
