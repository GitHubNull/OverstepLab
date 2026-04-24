<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">用户管理</h2>
        <p class="text-slate-500">管理平台内所有用户账户</p>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-users text-blue-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">总用户数</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ users.length }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-check text-green-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">正常用户</p>
              <p class="text-xl font-bold text-green-600">{{ activeUsers }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-red-100 dark:bg-red-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-ban text-red-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">禁用用户</p>
              <p class="text-xl font-bold text-red-600">{{ disabledUsers }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Users Table -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-users text-primary-500"></i>
          <span class="font-bold">用户列表</span>
        </div>
      </template>

      <template #content>
        <DataTable 
          :value="users" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
          <Column field="id" header="ID">
            <template #body="{ data }">
              <code class="text-xs bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">{{ data.id }}</code>
            </template>
          </Column>

          <Column field="username" header="用户名">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar 
                  :label="data.username.charAt(0).toUpperCase()"
                  shape="circle"
                  size="small"
                  class="bg-primary-100 text-primary-700"
                />
                <span class="font-medium">{{ data.username }}</span>
              </div>
            </template>
          </Column>

          <Column field="user_type" header="类型">
            <template #body="{ data }">
              <Tag 
                :value="getUserTypeText(data.user_type)"
                :severity="getUserTypeSeverity(data.user_type)"
                class="text-xs"
              />
            </template>
          </Column>

          <Column field="role" header="角色">
            <template #body="{ data }">
              <Tag 
                v-if="data.role"
                :value="getRoleText(data.role)"
                severity="info"
                class="text-xs"
              />
              <span v-else class="text-slate-400">-</span>
            </template>
          </Column>

          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag 
                :value="data.status === 'active' ? '正常' : '禁用'"
                :severity="data.status === 'active' ? 'success' : 'danger'"
                class="text-xs"
              />
            </template>
          </Column>

          <Column header="操作" style="width: 120px">
            <template #body="{ data }">
              <Button 
                :icon="data.status === 'active' ? 'pi pi-ban' : 'pi pi-check'"
                :severity="data.status === 'active' ? 'danger' : 'success'"
                text 
                rounded
                size="small"
                v-tooltip.top="data.status === 'active' ? '禁用' : '启用'"
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

const users = ref<User[]>([])

onMounted(async () => {
  const response = await api.adminListUsers()
  users.value = response.data.data!
})

const activeUsers = computed(() => users.value.filter(u => u.status === 'active').length)
const disabledUsers = computed(() => users.value.filter(u => u.status !== 'active').length)

function getUserTypeText(type: string) {
  const map: Record<string, string> = {
    platform_admin: '平台管理员',
    company: '企业用户',
    individual: '个人用户',
  }
  return map[type] || type
}

function getUserTypeSeverity(type: string) {
  const map: Record<string, string> = {
    platform_admin: 'danger',
    company: 'warning',
    individual: 'info',
  }
  return map[type] || 'secondary'
}

function getRoleText(role: string) {
  const map: Record<string, string> = {
    admin: '管理员',
    operator: '运维',
    finance: '财务',
    viewer: '只读',
  }
  return map[role] || role
}
</script>
