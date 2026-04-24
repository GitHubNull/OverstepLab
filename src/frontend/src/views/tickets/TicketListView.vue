<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">工单系统</h2>
        <p class="text-slate-500">提交和管理技术支持请求</p>
      </div>
      <Button 
        label="新建工单" 
        icon="pi pi-plus" 
        severity="primary"
      />
    </div>

    <!-- Tickets List -->
    <Card class="shadow-sm">
      <template #content>
        <DataTable 
          :value="tickets" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
          <Column field="title" header="标题">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i 
                  class="pi"
                  :class="getStatusIcon(data.status)"
                ></i>
                <span class="font-medium">{{ data.title }}</span>
              </div>
            </template>
          </Column>

          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag 
                :value="getStatusText(data.status)"
                :severity="getStatusSeverity(data.status)"
                class="text-xs"
              />
            </template>
          </Column>

          <Column header="创建时间">
            <template #body="{ data }">
              <span class="text-sm text-slate-500">{{ formatDate(data.created_at) }}</span>
            </template>
          </Column>

          <Column header="操作" style="width: 100px">
            <template #body="{ data }">
              <Button 
                label="查看" 
                icon="pi pi-eye" 
                text 
                size="small"
                @click="$router.push(`/tickets/${data.id}`)"
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
import type { Ticket } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import { formatDate } from '@/utils/date'

const tickets = ref<Ticket[]>([])

onMounted(async () => {
  const response = await api.getTickets()
  tickets.value = response.data.data!
})

function getStatusText(status: string) {
  const map: Record<string, string> = {
    open: '待处理',
    replied: '已回复',
    closed: '已关闭',
  }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = {
    open: 'warning',
    replied: 'info',
    closed: 'secondary',
  }
  return map[status] || 'secondary'
}

function getStatusIcon(status: string) {
  const map: Record<string, string> = {
    open: 'pi-envelope text-yellow-500',
    replied: 'pi-envelope-open text-blue-500',
    closed: 'pi-check-circle text-green-500',
  }
  return map[status] || 'pi-envelope'
}
</script>
