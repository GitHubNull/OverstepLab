<template>
  <div class="dark min-h-screen flex items-center justify-center relative overflow-hidden bg-[#0a0a0a]">
    <!-- Animated grid background -->
    <div class="absolute inset-0 grid-bg opacity-60"></div>
    <!-- Subtle radial glow -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(79,70,229,0.08)_0%,transparent_70%)]"></div>

    <div class="relative z-10 w-full max-w-[420px] px-4">
      <!-- Glass card -->
      <div class="rounded-2xl p-8 shadow-2xl shadow-black/40 bg-[rgba(10,10,10,0.88)] backdrop-blur-xl border border-white/[0.08]">
        <!-- Logo -->
        <div class="text-center mb-7">
          <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-2xl shadow-lg shadow-indigo-500/20 mb-4">
            <i class="pi pi-shield text-2xl text-white"></i>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">OverstepLab</h1>
          <p class="text-sm text-neutral-400">越权漏洞测试靶场</p>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="mb-4">
            <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">用户名</label>
            <span class="p-input-icon-left w-full block">
              <i class="pi pi-user text-neutral-500 z-10" />
              <InputText
                v-model="username"
                class="w-full pl-10 bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
                placeholder="请输入用户名"
                required
                :class="{ 'p-invalid': error }"
              />
            </span>
          </div>
          <div class="mb-6">
            <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">密码</label>
            <Password
              v-model="password"
              class="w-full"
              input-class="w-full bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
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
          <router-link to="/register" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium transition-colors">
            没有账号？立即注册
          </router-link>
        </div>

        <!-- Test Accounts -->
        <div class="mt-6 p-3 bg-white/8 rounded-xl border border-white/10">
          <p class="font-semibold text-[10px] text-neutral-400 mb-2 flex items-center gap-1.5 uppercase tracking-wider">
            <i class="pi pi-info-circle text-indigo-400 text-xs"></i>
            测试账户
          </p>
          <div class="space-y-1 text-sm">
            <div
              v-for="account in testAccounts"
              :key="account.user"
              class="flex items-center justify-between px-2.5 py-1.5 bg-white/8 rounded-lg cursor-pointer hover:bg-indigo-500/15 transition-colors border border-transparent hover:border-indigo-500/20 group relative"
              title="点击自动填充"
              @click="fillLogin(account.user, account.pass)"
            >
              <div class="flex items-center gap-1.5 min-w-0">
                <span class="font-medium text-neutral-200 text-xs">{{ account.user }}</span>
                <span class="text-neutral-500 text-[10px]">/</span>
                <span class="text-neutral-400 text-[11px] font-mono">{{ account.pass }}</span>
              </div>
              <Tag :value="account.role" :severity="account.severity" class="text-[10px] shrink-0" />
            </div>
          </div>
        </div>
      </div>

      <p class="text-center text-neutral-600 text-xs mt-5">
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

const testAccounts = [
  { user: 'admin', pass: 'admin123', role: '平台管理员', severity: 'danger' as const },
  { user: 'acme_admin', pass: 'pass123', role: '企业管理员', severity: 'warn' as const },
  { user: 'alice', pass: 'pass123', role: '个人用户', severity: 'info' as const },
]

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
