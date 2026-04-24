<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">系统设置</h2>
        <p class="text-slate-500">管理系统配置和数据</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Database Management -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-database text-primary-500"></i>
            <span class="font-bold">数据库管理</span>
          </div>
        </template>

        <template #content>
          <div class="space-y-4">
            <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-10 h-10 bg-red-100 dark:bg-red-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-exclamation-triangle text-red-600"></i>
                </div>
                <div>
                  <p class="font-semibold text-slate-800 dark:text-white">重置数据库</p>
                  <p class="text-sm text-slate-500">清除所有数据并恢复初始状态</p>
                </div>
              </div>
              <p class="text-sm text-slate-600 dark:text-slate-400 mb-4">
                此操作将删除所有用户数据、VPS 实例、订单记录等，并将数据库恢复到初始种子状态。此操作不可撤销！
              </p>
              <Button 
                label="重置数据库" 
                icon="pi pi-refresh" 
                severity="danger"
                @click="handleReset"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- Security Mode -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-shield text-primary-500"></i>
            <span class="font-bold">安全模式</span>
          </div>
        </template>

        <template #content>
          <div class="space-y-4">
            <div 
              class="rounded-xl p-4 border-2"
              :class="authStore.securityMode === 'vulnerable' 
                ? 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800' 
                : 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800'"
            >
              <div class="flex items-center gap-3">
                <div 
                  class="w-12 h-12 rounded-xl flex items-center justify-center"
                  :class="authStore.securityMode === 'vulnerable' ? 'bg-red-100' : 'bg-green-100'"
                >
                  <i 
                    class="pi text-2xl"
                    :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-600' : 'pi-lock text-green-600'"
                  ></i>
                </div>
                <div>
                  <p class="font-semibold text-slate-800 dark:text-white">
                    当前模式: {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
                  </p>
                  <p class="text-sm text-slate-500">
                    {{ authStore.securityMode === 'vulnerable' 
                      ? '所有越权漏洞均可触发' 
                      : '所有漏洞已被修复，权限校验严格' }}
                  </p>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-600 dark:text-slate-400">切换安全模式</span>
              <ToggleSwitch 
                v-model="isSecureMode"
                @change="toggleMode"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- System Info -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-info-circle text-primary-500"></i>
            <span class="font-bold">系统信息</span>
          </div>
        </template>

        <template #content>
          <div class="space-y-3">
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <span class="text-sm text-slate-500">版本</span>
              <span class="text-sm font-medium text-slate-800 dark:text-white">v1.0.0</span>
            </div>
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <span class="text-sm text-slate-500">技术栈</span>
              <span class="text-sm font-medium text-slate-800 dark:text-white">Go + Vue3 + PrimeVue</span>
            </div>
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <span class="text-sm text-slate-500">数据库</span>
              <span class="text-sm font-medium text-slate-800 dark:text-white">SQLite3</span>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import * as api from '@/api'
import Card from 'primevue/card'
import Button from 'primevue/button'
import ToggleSwitch from 'primevue/toggleswitch'

const authStore = useAuthStore()
const toast = useToast()

const isSecureMode = ref(authStore.securityMode === 'secure')

async function toggleMode() {
  await authStore.toggleSecurityMode()
  isSecureMode.value = authStore.securityMode === 'secure'
}

async function handleReset() {
  try {
    await api.adminResetDb()
    toast.add({ severity: 'success', summary: '成功', detail: '数据库已重置', life: 3000 })
    window.location.reload()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: '重置失败', life: 3000 })
  }
}
</script>
