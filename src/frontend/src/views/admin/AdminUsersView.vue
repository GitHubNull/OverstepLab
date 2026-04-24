<template>
  <div>
    <PageHeader title="用户管理" description="管理平台内所有用户账户" />

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-5">
      <StatCard color="#4f46e5" icon="pi pi-users" :value="users.length" label="总用户数" />
      <StatCard color="#10b981" icon="pi pi-check" :value="activeUsers" label="正常用户" />
      <StatCard color="#f43f5e" icon="pi pi-ban" :value="disabledUsers" label="禁用用户" />
    </div>

    <!-- Users Table -->
    <Card>
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
              <span class="mono text-[10px] px-1.5 py-0.5 rounded bg-[var(--bg-base)] text-[var(--text-secondary)]">{{ data.id }}</span>
            </template>
          </Column>
          <Column field="username" header="用户名">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar :label="data.username.charAt(0).toUpperCase()" shape="circle" size="small" class="bg-[var(--primary-subtle)] text-[var(--primary)] text-xs" />
                <span class="font-medium text-sm text-[var(--text-primary)]">{{ data.username }}</span>
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
              <span v-else class="text-[var(--text-tertiary)] text-xs">-</span>
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
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
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
