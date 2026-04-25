<template>
  <div class="dark min-h-screen flex items-center justify-center relative overflow-hidden bg-[#0a0a0a]">
    <!-- Animated grid background -->
    <div class="absolute inset-0 grid-bg opacity-60"></div>
    <!-- Subtle radial glow -->
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,rgba(79,70,229,0.08)_0%,transparent_70%)]"></div>

    <div class="relative z-10 w-full max-w-[440px] px-4">
      <div class="rounded-2xl p-8 shadow-2xl shadow-black/40 bg-[rgba(10,10,10,0.88)] backdrop-blur-xl border border-white/[0.08]">
        <!-- Logo -->
        <div class="text-center mb-6">
          <div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-2xl shadow-lg shadow-indigo-500/20 mb-4">
            <i class="pi pi-shield text-2xl text-white"></i>
          </div>
          <h1 class="text-2xl font-bold text-white mb-1">注册账号</h1>
          <p class="text-sm text-neutral-400">创建您的 OverstepLab 账户</p>
        </div>

        <!-- User Type Selector -->
        <div class="grid grid-cols-2 gap-2 mb-6">
          <div
            class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl cursor-pointer transition-all duration-200 border"
            :class="form.user_type === 'individual'
              ? 'border-indigo-500/50 bg-indigo-500/10 text-indigo-400'
              : 'border-white/10 bg-white/10 text-neutral-400 hover:border-white/20 hover:text-neutral-300'"
            @click="form.user_type = 'individual'"
          >
            <i class="pi pi-user text-sm"></i>
            <span class="text-sm font-semibold">个人注册</span>
          </div>
          <div
            class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl cursor-pointer transition-all duration-200 border"
            :class="form.user_type === 'company'
              ? 'border-indigo-500/50 bg-indigo-500/10 text-indigo-400'
              : 'border-white/10 bg-white/10 text-neutral-400 hover:border-white/20 hover:text-neutral-300'"
            @click="form.user_type = 'company'"
          >
            <i class="pi pi-building text-sm"></i>
            <span class="text-sm font-semibold">企业注册</span>
          </div>
        </div>

        <form @submit.prevent="handleRegister">
          <div class="space-y-4">
            <div v-if="form.user_type === 'company'">
              <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">企业名称</label>
              <InputText
                v-model="form.company_name"
                class="w-full bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
                placeholder="请输入企业名称"
                required
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">用户名</label>
              <InputText
                v-model="form.username"
                class="w-full bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
                :placeholder="form.user_type === 'company' ? '请输入管理员用户名' : '请输入用户名'"
                required
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">密码</label>
              <Password
                v-model="form.password"
                class="w-full"
                input-class="w-full bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
                placeholder="请输入密码"
                :feedback="true"
                required
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-neutral-400 mb-1.5 uppercase tracking-wider">邮箱</label>
              <InputText
                v-model="form.email"
                type="email"
                class="w-full bg-white/10 border-white/15 text-white placeholder:text-neutral-400 focus:border-indigo-500 focus:shadow-[0_0_0_3px_rgba(99,102,241,0.15)]"
                placeholder="请输入邮箱"
                required
              />
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
          <router-link to="/login" class="text-indigo-400 hover:text-indigo-300 text-sm font-medium transition-colors">
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
