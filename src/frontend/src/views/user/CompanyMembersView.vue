<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">企业成员</h2>
        <p class="text-slate-500">管理企业内的用户和权限</p>
      </div>
      <Button 
        label="添加成员" 
        icon="pi pi-plus" 
        severity="primary"
      />
    </div>

    <!-- Members Table -->
    <Card class="shadow-sm">
      <template #content>
        <DataTable 
          :value="members" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
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

          <Column field="role" header="角色">
            <template #body="{ data }">
              <Tag 
                :value="getRoleText(data.role)"
                :severity="getRoleSeverity(data.role)"
                class="text-xs"
              />
            </template>
          </Column>

          <Column field="email" header="邮箱">
            <template #body="{ data }">
              <span class="text-sm text-slate-600">{{ data.email || '-' }}</span>
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

          <Column header="操作" style="width: 100px">
            <template #body="{ data }">
              <Button 
                icon="pi pi-pencil" 
                text 
                rounded
                size="small"
                v-tooltip.top="'编辑'"
              />
              <Button 
                icon="pi pi-trash" 
                text 
                rounded
                size="small"
                severity="danger"
                v-tooltip.top="'删除'"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import type { User } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'

const members = ref<User[]>([])

onMounted(async () => {
  const response = await api.getMembers()
  members.value = response.data.data!
})

function getRoleText(role: string) {
  const map: Record<string, string> = {
    admin: '管理员',
    operator: '运维',
    finance: '财务',
    viewer: '只读',
  }
  return map[role] || role
}

function getRoleSeverity(role: string) {
  const map: Record<string, string> = {
    admin: 'danger',
    operator: 'warning',
    finance: 'info',
    viewer: 'secondary',
  }
  return map[role] || 'secondary'
}
</script>
