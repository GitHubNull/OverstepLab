<template>
  <div class="space-y-5">
    <PageHeader title="个人信息" description="管理您的账户信息和安全设置" />

    <div v-if="authStore.user" class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Profile Card -->
      <div class="lg:col-span-2 bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-5 py-4 border-b border-[var(--border-default)]">
          <div class="section-title">
            <i class="pi pi-user"></i>
            <span>基本信息</span>
          </div>
        </div>
        <div class="p-5">
          <div class="flex flex-col md:flex-row gap-6">
            <!-- Avatar -->
            <div class="flex flex-col items-center">
              <Avatar
                :label="authStore.user.username.charAt(0).toUpperCase()"
                size="xlarge"
                shape="circle"
                class="bg-[var(--primary-subtle)] text-[var(--primary)] text-2xl mb-2 w-20 h-20"
              />
              <p class="font-semibold text-sm text-[var(--text-primary)]">{{ authStore.user.username }}</p>
              <Tag :value="userTypeText" :severity="userTypeSeverity" class="text-[10px] mt-1" />
            </div>

            <!-- Info Grid -->
            <div class="flex-1 space-y-3">
              <div v-if="!editing" class="space-y-2">
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">用户 ID</span>
                  <code class="text-xs bg-[var(--bg-base)] px-1.5 py-0.5 rounded mono text-[var(--text-secondary)]">{{ authStore.user.id }}</code>
                </div>
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">用户名</span>
                  <span class="text-sm font-medium text-[var(--text-primary)]">{{ authStore.user.username }}</span>
                </div>
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">邮箱</span>
                  <span class="text-sm text-[var(--text-primary)]">{{ authStore.user.email || '未设置' }}</span>
                </div>
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">手机</span>
                  <span class="text-sm text-[var(--text-primary)]">{{ authStore.user.phone || '未设置' }}</span>
                </div>
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">用户类型</span>
                  <span class="text-sm font-medium text-[var(--text-primary)]">{{ userTypeText }}</span>
                </div>
                <div v-if="authStore.user.role" class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">角色</span>
                  <Tag :value="roleText" severity="info" class="text-[10px]" />
                </div>
                <div class="info-row">
                  <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">状态</span>
                  <Tag :value="authStore.user.status === 'active' ? '正常' : '禁用'" :severity="authStore.user.status === 'active' ? 'success' : 'danger'" class="text-[10px]" />
                </div>
              </div>

              <div v-else class="space-y-3">
                <div>
                  <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1 uppercase tracking-wider">邮箱</label>
                  <InputText v-model="editForm.email" class="w-full" placeholder="请输入邮箱" />
                </div>
                <div>
                  <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1 uppercase tracking-wider">手机号</label>
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
        </div>
      </div>

      <!-- Security Card -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-5 py-4 border-b border-[var(--border-default)]">
          <div class="section-title">
            <i class="pi pi-shield"></i>
            <span>安全设置</span>
          </div>
        </div>
        <div class="p-2">
          <div class="space-y-1">
            <div class="flex items-center gap-3 px-3 py-3 rounded-lg cursor-pointer hover:bg-[var(--bg-surface-hover)] transition-colors" @click="showPasswordDialog = true">
              <div class="w-9 h-9 bg-[var(--info-subtle)] rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="pi pi-key text-[var(--info)] text-xs"></i>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-[var(--text-primary)]">修改密码</p>
                <p class="text-[11px] text-[var(--text-secondary)]">定期更换密码保护账户安全</p>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-xs"></i>
            </div>

            <div class="flex items-center gap-3 px-3 py-3 rounded-lg cursor-pointer hover:bg-[var(--bg-surface-hover)] transition-colors" @click="$router.push('/apikeys')">
              <div class="w-9 h-9 bg-[var(--success-subtle)] rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="pi pi-key text-[var(--success)] text-xs"></i>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-[var(--text-primary)]">API Key 管理</p>
                <p class="text-[11px] text-[var(--text-secondary)]">管理您的 API 访问凭证</p>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-xs"></i>
            </div>

            <div class="flex items-center gap-3 px-3 py-3 rounded-lg cursor-pointer hover:bg-[var(--bg-surface-hover)] transition-colors" @click="$router.push('/audit')">
              <div class="w-9 h-9 bg-[var(--warning-subtle)] rounded-lg flex items-center justify-center flex-shrink-0">
                <i class="pi pi-history text-[var(--warning)] text-xs"></i>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-[var(--text-primary)]">操作日志</p>
                <p class="text-[11px] text-[var(--text-secondary)]">查看您的操作记录</p>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-xs"></i>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Change Password Dialog -->
    <Dialog v-model:visible="showPasswordDialog" header="修改密码" modal :style="{ width: '400px' }">
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">当前密码</label>
          <Password v-model="passwordForm.old_password" class="w-full" placeholder="请输入当前密码" :feedback="false" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">新密码</label>
          <Password v-model="passwordForm.new_password" class="w-full" placeholder="请输入新密码" :feedback="true" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">确认密码</label>
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
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { useToast } from 'primevue/usetoast'
import PageHeader from '@/components/PageHeader.vue'

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
