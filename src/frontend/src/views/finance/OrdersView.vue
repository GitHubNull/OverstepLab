<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>订单管理</h2>
      <p>查看您的购买和续费记录</p>
    </div>

    <!-- Orders Table -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-receipt"></i>
          <span>订单列表</span>
        </div>
      </template>
      <template #content>
        <DataTable :value="orders" stripedRows class="p-datatable-sm" :rows="10" paginator>
          <Column field="order_no" header="订单号">
            <template #body="{ data }">
              <code class="text-xs bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded text-slate-600 dark:text-slate-300">{{ data.order_no }}</code>
            </template>
          </Column>
          <Column field="type" header="类型">
            <template #body="{ data }">
              <Tag :value="getTypeText(data.type)" :severity="getTypeSeverity(data.type)" class="text-[10px]" />
            </template>
          </Column>
          <Column field="amount" header="金额">
            <template #body="{ data }">
              <span class="font-semibold text-sm text-slate-700 dark:text-slate-200">¥{{ data.amount.toFixed(2) }}</span>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="getStatusText(data.status)" :severity="getStatusSeverity(data.status)" class="text-[10px]" />
            </template>
          </Column>
          <Column header="创建时间">
            <template #body="{ data }">
              <span class="text-xs text-slate-400">{{ formatDate(data.created_at) }}</span>
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
import type { Order } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import { formatDate } from '@/utils/date'

const orders = ref<Order[]>([])

onMounted(async () => {
  const response = await api.getOrders()
  orders.value = response.data.data!
})

function getTypeText(type: string) {
  const map: Record<string, string> = { purchase: '购买', renew: '续费', upgrade: '升级' }
  return map[type] || type
}

function getTypeSeverity(type: string) {
  const map: Record<string, string> = { purchase: 'primary', renew: 'info', upgrade: 'warn' }
  return map[type] || 'secondary'
}

function getStatusText(status: string) {
  const map: Record<string, string> = { pending: '待支付', paid: '已支付', cancelled: '已取消' }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = { pending: 'warn', paid: 'success', cancelled: 'danger' }
  return map[status] || 'secondary'
}
</script>
