<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 relative overflow-hidden">
    <!-- Background decoration -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-primary-500/10 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-500/10 rounded-full blur-3xl"></div>
    </div>

    <div class="relative z-10 w-full max-w-md px-4">
      <div class="bg-white/95 dark:bg-slate-800/95 backdrop-blur-sm p-8 rounded-2xl shadow-2xl border border-white/20">
        <!-- Logo and Title -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-primary-500 to-primary-600 rounded-2xl shadow-lg mb-4">
            <i class="pi pi-shield text-3xl text-white"></i>
          </div>
          <h1 class="text-3xl font-bold text-slate-800 dark:text-white mb-2">OverstepLab</h1>
          <p class="text-slate-500 dark:text-slate-400">越权漏洞测试靶场</p>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="mb-5">
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">用户名</label>
            <InputText 
              v-model="username" 
              class="w-full" 
              placeholder="请输入用户名" 
              required 
              :class="{ 'p-invalid': error }"
            />
          </div>
          <div class="mb-6">
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">密码</label>
            <Password 
              v-model="password" 
              class="w-full" 
              input-class="w-full"
              placeholder="请输入密码" 
              :feedback="false" 
              required 
              :class="{ 'p-invalid': error }"
            />
          </div>
          <Button 
            label="登录" 
            type="submit" 
            class="w-full" 
            :loading="loading"
            severity="primary"
            size="large"
          />
        </form>

        <div class="mt-6 text-center">
          <router-link to="/register" class="text-primary-500 hover:text-primary-600 text-sm font-medium transition-colors">
            没有账号？立即注册
          </router-link>
        </div>

        <!-- Test Accounts -->
        <div class="mt-8 p-4 bg-slate-50 dark:bg-slate-700/50 rounded-xl border border-slate-200 dark:border-slate-600">
          <p class="font-semibold text-sm text-slate-700 dark:text-slate-300 mb-3 flex items-center gap-2">
            <i class="pi pi-info-circle text-primary-500"></i>
            测试账户
          </p>
          <div class="space-y-2 text-sm">
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800 rounded-lg">
              <span class="text-slate-600 dark:text-slate-400">admin</span>
              <span class="text-slate-400">admin123</span>
              <Tag value="平台管理员" severity="danger" class="text-xs" />
            </div>
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800 rounded-lg">
              <span class="text-slate-600 dark:text-slate-400">acme_admin</span>
              <span class="text-slate-400">pass123</span>
              <Tag value="企业管理员" severity="warning" class="text-xs" />
            </div>
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800 rounded-lg">
              <span class="text-slate-600 dark:text-slate-400">alice</span>
              <span class="text-slate-400">pass123</span>
              <Tag value="个人用户" severity="info" class="text-xs" />
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-center text-slate-500 text-xs mt-6">
        OverstepLab - 越权漏洞测试靶场平台
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Tag from 'primevue/tag'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref(false)

async function handleLogin() {
  loading.value = true
  error.value = false
  try {
    await authStore.login(username.value, password.value)
    toast.add({ severity: 'success', summary: '成功', detail: '登录成功', life: 3000 })
    router.push('/')
  } catch (e: any) {
    error.value = true
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '登录失败', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>
