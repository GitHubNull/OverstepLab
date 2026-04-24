<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>系统设置</h2>
      <p>管理系统配置和数据</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Database Management -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-database"></i>
            <span>数据库管理</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div class="bg-red-50 dark:bg-red-900/10 border border-red-100 dark:border-red-800/30 rounded-xl p-4">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-9 h-9 bg-red-100 dark:bg-red-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-exclamation-triangle text-red-500 text-sm"></i>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-700 dark:text-white">重置数据库</p>
                  <p class="text-[10px] text-slate-400">清除所有数据并恢复初始状态</p>
                </div>
              </div>
              <p class="text-xs text-slate-500 dark:text-slate-400 mb-4 leading-relaxed">
                此操作将删除所有用户数据、VPS 实例、订单记录等，并将数据库恢复到初始种子状态。<strong class="text-red-500">此操作不可撤销！</strong>
              </p>
              <Button
                label="重置数据库"
                icon="pi pi-refresh"
                severity="danger"
                size="small"
                @click="handleReset"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- Security Mode -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-shield"></i>
            <span>安全模式</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div
              class="rounded-xl p-4 border-2 transition-colors"
              :class="authStore.securityMode === 'vulnerable'
                ? 'bg-red-50 dark:bg-red-900/10 border-red-200 dark:border-red-800/50'
                : 'bg-emerald-50 dark:bg-emerald-900/10 border-emerald-200 dark:border-emerald-800/50'"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-11 h-11 rounded-xl flex items-center justify-center"
                  :class="authStore.securityMode === 'vulnerable' ? 'bg-red-100 dark:bg-red-900/30' : 'bg-emerald-100 dark:bg-emerald-900/30'"
                >
                  <i
                    class="pi text-xl"
                    :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-500' : 'pi-lock text-emerald-500'"
                  ></i>
                </div>
                <div>
                  <p class="font-semibold text-sm text-slate-700 dark:text-white">
                    当前模式: {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
                  </p>
                  <p class="text-[10px] text-slate-400 mt-0.5">
                    {{ authStore.securityMode === 'vulnerable'
                      ? '所有越权漏洞均可触发'
                      : '所有漏洞已被修复，权限校验严格' }}
                  </p>
                </div>
              </div>
            </div>

            <div class="info-row">
              <span class="text-xs text-slate-500 dark:text-slate-400">切换安全模式</span>
              <ToggleSwitch
                v-model="isSecureMode"
                @change="toggleMode"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- System Info -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-info-circle"></i>
            <span>系统信息</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row">
              <span class="text-xs text-slate-400">版本</span>
              <span class="text-xs font-medium text-slate-700 dark:text-slate-200">v1.0.0</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-slate-400">技术栈</span>
              <span class="text-xs font-medium text-slate-700 dark:text-slate-200">Go + Vue3 + PrimeVue</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-slate-400">数据库</span>
              <span class="text-xs font-medium text-slate-700 dark:text-slate-200">SQLite3</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-slate-400">框架</span>
              <span class="text-xs font-medium text-slate-700 dark:text-slate-200">Gin + GORM</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Quick Links -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-directions"></i>
            <span>快速导航</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row cursor-pointer group" @click="$router.push('/admin/users')">
              <div class="flex items-center gap-2">
                <i class="pi pi-users text-indigo-400 text-xs"></i>
                <span class="text-xs font-medium text-slate-600 dark:text-slate-300 group-hover:text-indigo-500 transition-colors">用户管理</span>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/admin/companies')">
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-indigo-400 text-xs"></i>
                <span class="text-xs font-medium text-slate-600 dark:text-slate-300 group-hover:text-indigo-500 transition-colors">企业管理</span>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/challenges')">
              <div class="flex items-center gap-2">
                <i class="pi pi-flag text-indigo-400 text-xs"></i>
                <span class="text-xs font-medium text-slate-600 dark:text-slate-300 group-hover:text-indigo-500 transition-colors">漏洞挑战</span>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/audit')">
              <div class="flex items-center gap-2">
                <i class="pi pi-history text-indigo-400 text-xs"></i>
                <span class="text-xs font-medium text-slate-600 dark:text-slate-300 group-hover:text-indigo-500 transition-colors">审计日志</span>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-[10px]"></i>
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
import { useConfirm } from 'primevue/useconfirm'
import * as api from '@/api'
import Card from 'primevue/card'
import Button from 'primevue/button'
import ToggleSwitch from 'primevue/toggleswitch'

const authStore = useAuthStore()
const toast = useToast()
const confirm = useConfirm()
const isSecureMode = ref(authStore.securityMode === 'secure')

async function toggleMode() {
  await authStore.toggleSecurityMode()
  isSecureMode.value = authStore.securityMode === 'secure'
}

function handleReset() {
  confirm.require({
    message: '确定要重置数据库吗？所有数据将被清除并恢复初始状态，此操作不可撤销！',
    header: '确认重置数据库',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.adminResetDb()
        toast.add({ severity: 'success', summary: '成功', detail: '数据库已重置', life: 3000 })
        window.location.reload()
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: '重置失败', life: 3000 })
      }
    },
  })
}
</script>
