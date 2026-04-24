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

const { init } = useTheme()

onMounted(() => {
  init()
})
</script>
