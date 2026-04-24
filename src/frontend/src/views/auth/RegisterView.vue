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
          <h1 class="text-3xl font-bold text-slate-800 dark:text-white mb-2">注册账号</h1>
          <p class="text-slate-500 dark:text-slate-400">创建您的 OverstepLab 账户</p>
        </div>

        <Tabs value="0" class="mb-6">
          <TabList class="grid grid-cols-2">
            <Tab value="0" class="flex items-center justify-center gap-2">
              <i class="pi pi-user"></i>
              个人注册
            </Tab>
            <Tab value="1" class="flex items-center justify-center gap-2">
              <i class="pi pi-building"></i>
              企业注册
            </Tab>
          </TabList>
          <TabPanels>
            <TabPanel value="0">
              <form @submit.prevent="handleRegister">
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">用户名</label>
                    <InputText v-model="form.username" class="w-full" placeholder="请输入用户名" required />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">密码</label>
                    <Password v-model="form.password" class="w-full" placeholder="请输入密码" :feedback="true" required />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">邮箱</label>
                    <InputText v-model="form.email" type="email" class="w-full" placeholder="请输入邮箱" required />
                  </div>
                </div>
                <Button 
                  label="注册" 
                  type="submit" 
                  class="w-full mt-6" 
                  :loading="loading"
                  severity="primary"
                  size="large"
                />
              </form>
            </TabPanel>
            <TabPanel value="1">
              <form @submit.prevent="handleRegister">
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">企业名称</label>
                    <InputText v-model="form.company_name" class="w-full" placeholder="请输入企业名称" required />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">用户名</label>
                    <InputText v-model="form.username" class="w-full" placeholder="请输入管理员用户名" required />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">密码</label>
                    <Password v-model="form.password" class="w-full" placeholder="请输入密码" :feedback="true" required />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">邮箱</label>
                    <InputText v-model="form.email" type="email" class="w-full" placeholder="请输入邮箱" required />
                  </div>
                </div>
                <Button 
                  label="注册企业账号" 
                  type="submit" 
                  class="w-full mt-6" 
                  :loading="loading"
                  severity="primary"
                  size="large"
                />
              </form>
            </TabPanel>
          </TabPanels>
        </Tabs>

        <div class="mt-4 text-center">
          <router-link to="/login" class="text-primary-500 hover:text-primary-600 text-sm font-medium transition-colors">
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
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'

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
