<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>用户管理</h2>
      <p>管理平台内所有用户账户</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <Card class="shadow-none stat-card stat-card-blue">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-blue-50 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-users text-blue-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">总用户数</p>
              <p class="text-lg font-bold text-slate-800 dark:text-white">{{ users.length }}</p>
            </div>
          </div>
        </template>
      </Card>
      <Card class="shadow-none stat-card stat-card-green">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-check text-emerald-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">正常用户</p>
              <p class="text-lg font-bold text-emerald-600">{{ activeUsers }}</p>
            </div>
          </div>
        </template>
      </Card>
      <Card class="shadow-none stat-card stat-card-red">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-red-50 dark:bg-red-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-ban text-red-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">禁用用户</p>
              <p class="text-lg font-bold text-red-500">{{ disabledUsers }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Users Table -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-users"></i>
          <span>用户列表</span>
        </div>
      </template>
      <template #content>
        <DataTable :value="users" stripedRows class="p-datatable-sm" :rows="10" paginator>
          <Column field="id" header="ID" style="width: 60px">
            <template #body="{ data }">
              <code class="text-[10px] bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded">{{ data.id }}</code>
            </template>
          </Column>
          <Column field="username" header="用户名">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar :label="data.username.charAt(0).toUpperCase()" shape="circle" size="small" class="bg-indigo-100 text-indigo-600 text-xs" />
                <span class="font-medium text-sm text-slate-700 dark:text-slate-200">{{ data.username }}</span>
              </div>
            </template>
          </Column>
          <Column field="user_type" header="类型">
            <template #body="{ data }">
              <Tag :value="getUserTypeText(data.user_type)" :severity="getUserTypeSeverity(data.user_type)" class="text-[10px]" />
            </template>
          </Column>
          <Column field="role" header="角色">
            <template #body="{ data }">
              <Tag v-if="data.role" :value="getRoleText(data.role)" severity="info" class="text-[10px]" />
              <span v-else class="text-slate-300 text-xs">-</span>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="data.status === 'active' ? '正常' : '禁用'" :severity="data.status === 'active' ? 'success' : 'danger'" class="text-[10px]" />
            </template>
          </Column>
          <Column header="操作" style="width: 80px">
            <template #body="{ data }">
              <Button
                :icon="data.status === 'active' ? 'pi pi-ban' : 'pi pi-check'"
                :severity="data.status === 'active' ? 'danger' : 'success'"
                text
                rounded
                size="small"
                v-tooltip.top="data.status === 'active' ? '禁用' : '启用'"
                @click="handleToggleStatus(data)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import type { User } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'

const toast = useToast()
const confirm = useConfirm()
const users = ref<User[]>([])

onMounted(() => fetchUsers())

async function fetchUsers() {
  const response = await api.adminListUsers()
  users.value = response.data.data!
}

const activeUsers = computed(() => users.value.filter(u => u.status === 'active').length)
const disabledUsers = computed(() => users.value.filter(u => u.status !== 'active').length)

function getUserTypeText(type: string) {
  const map: Record<string, string> = { platform_admin: '平台管理员', company: '企业用户', individual: '个人用户' }
  return map[type] || type
}

function getUserTypeSeverity(type: string) {
  const map: Record<string, string> = { platform_admin: 'danger', company: 'warn', individual: 'info' }
  return map[type] || 'secondary'
}

function getRoleText(role: string) {
  const map: Record<string, string> = { admin: '管理员', operator: '运维', finance: '财务', viewer: '只读' }
  return map[role] || role
}

function handleToggleStatus(user: User) {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  const action = newStatus === 'disabled' ? '禁用' : '启用'
  confirm.require({
    message: `确定要${action}用户 "${user.username}" 吗？`,
    header: `确认${action}`,
    icon: 'pi pi-exclamation-triangle',
    acceptClass: newStatus === 'disabled' ? 'p-button-danger' : 'p-button-success',
    accept: async () => {
      try {
        await api.adminUpdateUserStatus(user.id, newStatus)
        toast.add({ severity: 'success', summary: '成功', detail: `用户已${action}`, life: 2000 })
        await fetchUsers()
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
      }
    },
  })
}
</script>
