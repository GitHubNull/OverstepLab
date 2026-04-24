<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900 relative overflow-hidden">
    <!-- Background decoration -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-96 h-96 bg-indigo-500/8 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-purple-500/8 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-indigo-500/5 rounded-full blur-3xl"></div>
    </div>

    <div class="relative z-10 w-full max-w-[420px] px-4">
      <div class="bg-white dark:bg-slate-800 p-8 rounded-2xl shadow-2xl shadow-black/20 border border-white/10">
        <!-- Logo and Title -->
        <div class="text-center mb-7">
          <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-2xl shadow-lg shadow-indigo-200/50 mb-4">
            <i class="pi pi-shield text-2xl text-white"></i>
          </div>
          <h1 class="text-2xl font-bold text-slate-800 dark:text-white mb-1">OverstepLab</h1>
          <p class="text-sm text-slate-400">越权漏洞测试靶场</p>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="mb-4">
            <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">用户名</label>
            <span class="p-input-icon-left w-full block">
              <i class="pi pi-user text-slate-400 z-10" />
              <InputText
                v-model="username"
                class="w-full pl-10"
                placeholder="请输入用户名"
                required
                :class="{ 'p-invalid': error }"
              />
            </span>
          </div>
          <div class="mb-6">
            <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">密码</label>
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

        <div class="mt-5 text-center">
          <router-link to="/register" class="text-indigo-500 hover:text-indigo-600 text-sm font-medium transition-colors">
            没有账号？立即注册
          </router-link>
        </div>

        <!-- Test Accounts -->
        <div class="mt-6 p-4 bg-slate-50 dark:bg-slate-700/40 rounded-xl border border-slate-100 dark:border-slate-600/50">
          <p class="font-semibold text-xs text-slate-500 dark:text-slate-400 mb-3 flex items-center gap-1.5 uppercase tracking-wider">
            <i class="pi pi-info-circle text-indigo-400 text-sm"></i>
            测试账户
          </p>
          <div class="space-y-1.5 text-sm">
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800/80 rounded-lg cursor-pointer hover:bg-indigo-50 dark:hover:bg-indigo-900/10 transition-colors" @click="fillLogin('admin', 'admin123')">
              <span class="font-medium text-slate-700 dark:text-slate-300">admin</span>
              <Tag value="平台管理员" severity="danger" class="text-[10px]" />
            </div>
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800/80 rounded-lg cursor-pointer hover:bg-indigo-50 dark:hover:bg-indigo-900/10 transition-colors" @click="fillLogin('acme_admin', 'pass123')">
              <span class="font-medium text-slate-700 dark:text-slate-300">acme_admin</span>
              <Tag value="企业管理员" severity="warn" class="text-[10px]" />
            </div>
            <div class="flex items-center justify-between p-2 bg-white dark:bg-slate-800/80 rounded-lg cursor-pointer hover:bg-indigo-50 dark:hover:bg-indigo-900/10 transition-colors" @click="fillLogin('alice', 'pass123')">
              <span class="font-medium text-slate-700 dark:text-slate-300">alice</span>
              <Tag value="个人用户" severity="info" class="text-[10px]" />
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-center text-slate-500/60 text-xs mt-5">
        OverstepLab - 仅供安全研究与学习使用
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

function fillLogin(user: string, pass: string) {
  username.value = user
  password.value = pass
}

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
