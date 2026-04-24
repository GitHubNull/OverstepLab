<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-indigo-950 to-slate-900 relative overflow-hidden">
    <!-- Background decoration -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-96 h-96 bg-indigo-500/8 rounded-full blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-purple-500/8 rounded-full blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-indigo-500/5 rounded-full blur-3xl"></div>
    </div>

    <div class="relative z-10 w-full max-w-[440px] px-4">
      <div class="bg-white dark:bg-slate-800 p-8 rounded-2xl shadow-2xl shadow-black/20 border border-white/10">
        <!-- Logo and Title -->
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-2xl shadow-lg shadow-indigo-200/50 mb-4">
            <i class="pi pi-shield text-2xl text-white"></i>
          </div>
          <h1 class="text-2xl font-bold text-slate-800 dark:text-white mb-1">注册账号</h1>
          <p class="text-sm text-slate-400">创建您的 OverstepLab 账户</p>
        </div>

        <!-- User Type Selector -->
        <div class="grid grid-cols-2 gap-2 mb-6">
          <div
            class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl cursor-pointer transition-all duration-200 border-2"
            :class="form.user_type === 'individual'
              ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400'
              : 'border-slate-200 dark:border-slate-600 bg-slate-50 dark:bg-slate-700/40 text-slate-500 hover:border-slate-300'"
            @click="form.user_type = 'individual'"
          >
            <i class="pi pi-user text-sm"></i>
            <span class="text-sm font-semibold">个人注册</span>
          </div>
          <div
            class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl cursor-pointer transition-all duration-200 border-2"
            :class="form.user_type === 'company'
              ? 'border-indigo-500 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400'
              : 'border-slate-200 dark:border-slate-600 bg-slate-50 dark:bg-slate-700/40 text-slate-500 hover:border-slate-300'"
            @click="form.user_type = 'company'"
          >
            <i class="pi pi-building text-sm"></i>
            <span class="text-sm font-semibold">企业注册</span>
          </div>
        </div>

        <form @submit.prevent="handleRegister">
          <div class="space-y-4">
            <!-- Company name (company registration) -->
            <div v-if="form.user_type === 'company'">
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">企业名称</label>
              <InputText v-model="form.company_name" class="w-full" placeholder="请输入企业名称" required />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">用户名</label>
              <InputText v-model="form.username" class="w-full" :placeholder="form.user_type === 'company' ? '请输入管理员用户名' : '请输入用户名'" required />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">密码</label>
              <Password v-model="form.password" class="w-full" placeholder="请输入密码" :feedback="true" required />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-500 dark:text-slate-400 mb-1.5 uppercase tracking-wider">邮箱</label>
              <InputText v-model="form.email" type="email" class="w-full" placeholder="请输入邮箱" required />
            </div>
          </div>

          <Button
            :label="form.user_type === 'company' ? '注册企业账号' : '注册'"
            type="submit"
            class="w-full mt-6"
            :loading="loading"
            severity="primary"
            size="large"
          />
        </form>

        <div class="mt-5 text-center">
          <router-link to="/login" class="text-indigo-500 hover:text-indigo-600 text-sm font-medium transition-colors">
            已有账号？立即登录
          </router-link>
        </div>
      </div>
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

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const loading = ref(false)
const form = ref({
  username: '',
  password: '',
  email: '',
  company_name: '',
  user_type: 'individual'
})

async function handleRegister() {
  loading.value = true
  try {
    await authStore.register({
      username: form.value.username,
      password: form.value.password,
      email: form.value.email,
      user_type: form.value.user_type,
      ...(form.value.user_type === 'company' && { company_name: form.value.company_name }),
    })
    toast.add({ severity: 'success', summary: '成功', detail: '注册成功，请登录', life: 3000 })
    router.push('/login')
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '注册失败', life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>
