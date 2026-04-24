<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>个人信息</h2>
      <p>管理您的账户信息和安全设置</p>
    </div>

    <div v-if="authStore.user" class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Profile Card -->
      <Card class="shadow-none lg:col-span-2">
        <template #title>
          <div class="section-title">
            <i class="pi pi-user"></i>
            <span>基本信息</span>
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
                class="bg-indigo-100 text-indigo-600 text-2xl mb-2"
              />
              <p class="font-semibold text-sm text-slate-700 dark:text-white">{{ authStore.user.username }}</p>
              <Tag :value="userTypeText" :severity="userTypeSeverity" class="text-[10px] mt-1" />
            </div>

            <!-- Info Grid -->
            <div class="flex-1 space-y-3">
              <!-- Editable fields -->
              <div v-if="!editing" class="space-y-2">
                <div class="info-row">
                  <span class="text-xs text-slate-400">用户 ID</span>
                  <code class="text-xs bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded">{{ authStore.user.id }}</code>
                </div>
                <div class="info-row">
                  <span class="text-xs text-slate-400">用户名</span>
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ authStore.user.username }}</span>
                </div>
                <div class="info-row">
                  <span class="text-xs text-slate-400">邮箱</span>
                  <span class="text-sm text-slate-700 dark:text-slate-200">{{ authStore.user.email || '未设置' }}</span>
                </div>
                <div class="info-row">
                  <span class="text-xs text-slate-400">手机</span>
                  <span class="text-sm text-slate-700 dark:text-slate-200">{{ authStore.user.phone || '未设置' }}</span>
                </div>
                <div class="info-row">
                  <span class="text-xs text-slate-400">用户类型</span>
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-200">{{ userTypeText }}</span>
                </div>
                <div v-if="authStore.user.role" class="info-row">
                  <span class="text-xs text-slate-400">角色</span>
                  <Tag :value="roleText" severity="info" class="text-[10px]" />
                </div>
                <div class="info-row">
                  <span class="text-xs text-slate-400">状态</span>
                  <Tag :value="authStore.user.status === 'active' ? '正常' : '禁用'" :severity="authStore.user.status === 'active' ? 'success' : 'danger'" class="text-[10px]" />
                </div>
              </div>

              <!-- Edit mode -->
              <div v-else class="space-y-3">
                <div>
                  <label class="block text-xs font-semibold text-slate-500 mb-1 uppercase tracking-wider">邮箱</label>
                  <InputText v-model="editForm.email" class="w-full" placeholder="请输入邮箱" />
                </div>
                <div>
                  <label class="block text-xs font-semibold text-slate-500 mb-1 uppercase tracking-wider">手机号</label>
                  <InputText v-model="editForm.phone" class="w-full" placeholder="请输入手机号" />
                </div>
                <div class="flex gap-2 pt-1">
                  <Button label="保存" icon="pi pi-check" size="small" :loading="saving" @click="handleSaveProfile" />
                  <Button label="取消" text size="small" @click="editing = false" />
                </div>
              </div>

              <div v-if="!editing" class="pt-2">
                <Button label="编辑资料" icon="pi pi-pencil" text size="small" @click="startEdit" />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Security Card -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-shield"></i>
            <span>安全设置</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row cursor-pointer group" @click="showPasswordDialog = true">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-blue-50 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-key text-blue-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-sm font-medium text-slate-700 dark:text-white group-hover:text-indigo-500 transition-colors">修改密码</p>
                  <p class="text-[10px] text-slate-400">定期更换密码保护账户安全</p>
                </div>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-xs"></i>
            </div>

            <div class="info-row cursor-pointer group" @click="$router.push('/apikeys')">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-key text-emerald-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-sm font-medium text-slate-700 dark:text-white group-hover:text-indigo-500 transition-colors">API Key 管理</p>
                  <p class="text-[10px] text-slate-400">管理您的 API 访问凭证</p>
                </div>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-xs"></i>
            </div>

            <div class="info-row cursor-pointer group" @click="$router.push('/audit')">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-orange-50 dark:bg-orange-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-history text-orange-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-sm font-medium text-slate-700 dark:text-white group-hover:text-indigo-500 transition-colors">操作日志</p>
                  <p class="text-[10px] text-slate-400">查看您的操作记录</p>
                </div>
              </div>
              <i class="pi pi-chevron-right text-slate-300 text-xs"></i>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Change Password Dialog -->
    <Dialog v-model:visible="showPasswordDialog" header="修改密码" modal :style="{ width: '400px' }">
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">当前密码</label>
          <Password v-model="passwordForm.old_password" class="w-full" placeholder="请输入当前密码" :feedback="false" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">新密码</label>
          <Password v-model="passwordForm.new_password" class="w-full" placeholder="请输入新密码" :feedback="true" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">确认密码</label>
          <Password v-model="passwordForm.confirm_password" class="w-full" placeholder="再次输入新密码" :feedback="false" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showPasswordDialog = false" />
        <Button label="确认修改" icon="pi pi-check" size="small" :loading="changingPassword" @click="handleChangePassword" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { updateProfile as apiUpdateProfile, changePassword as apiChangePassword } from '@/api/auth'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { useToast } from 'primevue/usetoast'

const authStore = useAuthStore()
const toast = useToast()

const editing = ref(false)
const saving = ref(false)
const showPasswordDialog = ref(false)
const changingPassword = ref(false)

const editForm = ref({ email: '', phone: '' })
const passwordForm = ref({ old_password: '', new_password: '', confirm_password: '' })

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
    case 'company': return 'warn'
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

function startEdit() {
  editForm.value.email = authStore.user?.email || ''
  editForm.value.phone = authStore.user?.phone || ''
  editing.value = true
}

async function handleSaveProfile() {
  saving.value = true
  try {
    await apiUpdateProfile(editForm.value)
    await authStore.fetchProfile()
    toast.add({ severity: 'success', summary: '成功', detail: '资料已更新', life: 2000 })
    editing.value = false
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '更新失败', life: 3000 })
  } finally {
    saving.value = false
  }
}

async function handleChangePassword() {
  if (!passwordForm.value.old_password || !passwordForm.value.new_password) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请填写完整信息', life: 2000 })
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    toast.add({ severity: 'warn', summary: '提示', detail: '两次密码不一致', life: 2000 })
    return
  }
  changingPassword.value = true
  try {
    await apiChangePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password,
    })
    toast.add({ severity: 'success', summary: '成功', detail: '密码已修改', life: 2000 })
    showPasswordDialog.value = false
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '修改失败', life: 3000 })
  } finally {
    changingPassword.value = false
  }
}
</script>
