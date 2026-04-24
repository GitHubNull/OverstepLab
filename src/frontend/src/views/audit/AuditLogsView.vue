<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">审计日志</h2>
        <p class="text-slate-500">查看系统操作记录和安全审计信息</p>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-list text-blue-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">总记录数</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ logs.length }}</p>
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
              <p class="text-sm text-slate-500">今日操作</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ todayLogs }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Logs Table -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-history text-primary-500"></i>
          <span class="font-bold">操作记录</span>
        </div>
      </template>

      <template #content>
        <DataTable 
          :value="logs" 
          stripedRows
          class="p-datatable-sm"
          :rows="15"
          paginator
        >
          <Column field="user_id" header="用户">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar 
                  :label="String(data.user_id)"
                  shape="circle"
                  size="small"
                  class="bg-primary-100 text-primary-700 text-xs"
                />
                <span class="text-sm">用户 #{{ data.user_id }}</span>
              </div>
            </template>
          </Column>

          <Column field="action" header="操作">
            <template #body="{ data }">
              <Tag 
                :value="data.action"
                severity="info"
                class="text-xs font-mono"
              />
            </template>
          </Column>

          <Column field="resource_type" header="资源类型">
            <template #body="{ data }">
              <span class="text-sm text-slate-600">{{ data.resource_type }}</span>
            </template>
          </Column>

          <Column field="resource_id" header="资源 ID">
            <template #body="{ data }">
              <code class="text-xs bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">{{ data.resource_id }}</code>
            </template>
          </Column>

          <Column field="ip_address" header="IP 地址">
            <template #body="{ data }">
              <code class="text-xs text-slate-500">{{ data.ip_address }}</code>
            </template>
          </Column>

          <Column header="时间">
            <template #body="{ data }">
              <span class="text-sm text-slate-500">{{ formatDate(data.created_at) }}</span>
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
import type { AuditLog } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import { formatDate } from '@/utils/date'

const logs = ref<AuditLog[]>([])

onMounted(async () => {
  const response = await api.getAuditLogs()
  logs.value = response.data.data!
})

const todayLogs = computed(() => {
  const today = new Date().toDateString()
  return logs.value.filter(log => new Date(log.created_at).toDateString() === today).length
})
</script>
