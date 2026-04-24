<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>审计日志</h2>
      <p>查看系统操作记录和安全审计信息</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <Card class="shadow-none stat-card stat-card-blue">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-blue-50 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-list text-blue-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">总记录数</p>
              <p class="text-lg font-bold text-slate-800 dark:text-white">{{ logs.length }}</p>
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
              <p class="text-[10px] text-slate-400 font-medium">今日操作</p>
              <p class="text-lg font-bold text-emerald-600">{{ todayLogs }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Logs Table -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-history"></i>
          <span>操作记录</span>
        </div>
      </template>
      <template #content>
        <DataTable :value="logs" stripedRows class="p-datatable-sm" :rows="15" paginator>
          <Column field="user_id" header="用户">
            <template #body="{ data }">
              <div class="flex items-center gap-1.5">
                <Avatar :label="String(data.user_id)" shape="circle" size="small" class="bg-indigo-100 text-indigo-600 text-[10px]" />
                <span class="text-xs text-slate-600 dark:text-slate-300">用户 #{{ data.user_id }}</span>
              </div>
            </template>
          </Column>
          <Column field="action" header="操作">
            <template #body="{ data }">
              <Tag :value="data.action" severity="info" class="text-[10px] font-mono" />
            </template>
          </Column>
          <Column field="resource_type" header="资源类型">
            <template #body="{ data }">
              <span class="text-xs text-slate-500">{{ data.resource_type }}</span>
            </template>
          </Column>
          <Column field="resource_id" header="资源 ID">
            <template #body="{ data }">
              <code class="text-[10px] bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded">{{ data.resource_id }}</code>
            </template>
          </Column>
          <Column field="ip_address" header="IP 地址">
            <template #body="{ data }">
              <code class="text-[10px] text-slate-400">{{ data.ip_address }}</code>
            </template>
          </Column>
          <Column header="时间">
            <template #body="{ data }">
              <span class="text-[10px] text-slate-400">{{ formatDate(data.created_at) }}</span>
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
