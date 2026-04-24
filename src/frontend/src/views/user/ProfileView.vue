<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">个人信息</h2>
        <p class="text-slate-500">管理您的账户信息和安全设置</p>
      </div>
    </div>

    <div v-if="authStore.user" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Profile Card -->
      <Card class="shadow-sm lg:col-span-2">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-user text-primary-500"></i>
            <span class="font-bold">基本信息</span>
          </div>
        </template>

        <template #content>
          <div class="flex flex-col md:flex-row gap-6">
            <!-- Avatar -->
            <div class="flex flex-col items-center">
              <Avatar 
                :label="authStore.user.username.charAt(0).toUpperCase()"
                size="xlarge"
                shape="circle"
                class="bg-primary-100 text-primary-700 text-2xl mb-3"
              />
              <p class="font-semibold text-slate-800 dark:text-white">{{ authStore.user.username }}</p>
              <Tag 
                :value="userTypeText" 
                :severity="userTypeSeverity"
                class="mt-2"
              />
            </div>

            <!-- Info Grid -->
            <div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-4">
                <div>
                  <label class="block text-sm text-slate-500 mb-1">用户 ID</label>
                  <p class="font-mono text-sm bg-slate-100 dark:bg-slate-700 px-3 py-2 rounded-lg">{{ authStore.user.id }}</p>
                </div>

                <div>
                  <label class="block text-sm text-slate-500 mb-1">用户名</label>
                  <p class="font-semibold text-slate-800 dark:text-white">{{ authStore.user.username }}</p>
                </div>

                <div>
                  <label class="block text-sm text-slate-500 mb-1">邮箱</label>
                  <p class="text-slate-800 dark:text-white">{{ authStore.user.email || '未设置' }}</p>
                </div>
              </div>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm text-slate-500 mb-1">用户类型</label>
                  <p class="font-semibold text-slate-800 dark:text-white">{{ userTypeText }}</p>
                </div>

                <div v-if="authStore.user.role">
                  <label class="block text-sm text-slate-500 mb-1">角色</label>
                  <Tag :value="roleText" severity="info" />
                </div>

                <div>
                  <label class="block text-sm text-slate-500 mb-1">状态</label>
                  <Tag 
                    :value="authStore.user.status === 'active' ? '正常' : '禁用'" 
                    :severity="authStore.user.status === 'active' ? 'success' : 'danger'"
                  />
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Security Card -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-shield text-primary-500"></i>
            <span class="font-bold">安全设置</span>
          </div>
        </template>

        <template #content>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-key text-blue-600"></i>
                </div>
                <div>
                  <p class="font-medium text-slate-800 dark:text-white">修改密码</p>
                  <p class="text-xs text-slate-500">定期更换密码保护账户安全</p>
                </div>
              </div>
              <Button icon="pi pi-chevron-right" text rounded />
            </div>

            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-key text-green-600"></i>
                </div>
                <div>
                  <p class="font-medium text-slate-800 dark:text-white">API Key 管理</p>
                  <p class="text-xs text-slate-500">管理您的 API 访问凭证</p>
                </div>
              </div>
              <Button 
                icon="pi pi-chevron-right" 
                text 
                rounded
                @click="$router.push('/apikeys')"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'

const authStore = useAuthStore()

const userTypeText = computed(() => {
  switch (authStore.user?.user_type) {
    case 'platform_admin': return '平台管理员'
    case 'company': return '企业用户'
    case 'individual': return '个人用户'
    default: return authStore.user?.user_type
  }
})

const userTypeSeverity = computed(() => {
  switch (authStore.user?.user_type) {
    case 'platform_admin': return 'danger'
    case 'company': return 'warning'
    default: return 'info'
  }
})

const roleText = computed(() => {
  switch (authStore.user?.role) {
    case 'admin': return '管理员'
    case 'operator': return '运维'
    case 'finance': return '财务'
    case 'viewer': return '只读'
    default: return authStore.user?.role
  }
})
</script>
