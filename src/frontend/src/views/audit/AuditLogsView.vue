<template>
  <div class="space-y-5">
    <PageHeader title="审计日志" description="查看系统操作记录和安全审计信息" />

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard color="#3b82f6" icon="pi pi-list" :value="logs.length" label="总记录数" />
      <StatCard color="#10b981" icon="pi pi-check" :value="todayLogs" label="今日操作" />
    </div>

    <!-- Logs Table -->
    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div class="px-5 py-4 border-b border-[var(--border-default)]">
        <div class="section-title">
          <i class="pi pi-history"></i>
          <span>操作记录</span>
        </div>
      </div>
      <div class="p-0">
        <DataTable :value="logs" class="p-datatable-sm" :rows="15" paginator>
          <Column field="user_id" header="用户">
            <template #body="{ data }">
              <div class="flex items-center gap-1.5">
                <Avatar :label="String(data.user_id)" shape="circle" size="small" class="bg-[var(--primary-subtle)] text-[var(--primary)] text-[10px]" />
                <span class="text-xs text-[var(--text-secondary)]">用户 #{{ data.user_id }}</span>
              </div>
            </template>
          </Column>
          <Column field="action" header="操作">
            <template #body="{ data }">
              <code class="text-[10px] px-1.5 py-0.5 rounded bg-[var(--bg-base)] text-[var(--text-secondary)] mono font-semibold">{{ data.action }}</code>
            </template>
          </Column>
          <Column field="resource_type" header="资源类型">
            <template #body="{ data }">
              <span class="text-xs text-[var(--text-secondary)]">{{ data.resource_type }}</span>
            </template>
          </Column>
          <Column field="resource_id" header="资源 ID">
            <template #body="{ data }">
              <code class="text-[10px] bg-[var(--bg-base)] px-1.5 py-0.5 rounded mono text-[var(--text-tertiary)]">{{ data.resource_id }}</code>
            </template>
          </Column>
          <Column field="ip_address" header="IP 地址">
            <template #body="{ data }">
              <code class="text-[10px] text-[var(--text-tertiary)] mono">{{ data.ip_address }}</code>
            </template>
          </Column>
          <Column header="时间">
            <template #body="{ data }">
              <span class="text-[10px] text-[var(--text-tertiary)]">{{ formatDate(data.created_at) }}</span>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import type { AuditLog } from '@/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Avatar from 'primevue/avatar'
import { formatDate } from '@/utils/date'
import StatCard from '@/components/StatCard.vue'
import PageHeader from '@/components/PageHeader.vue'

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
