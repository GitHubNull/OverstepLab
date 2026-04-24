<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">订单管理</h2>
        <p class="text-slate-500">查看您的购买和续费记录</p>
      </div>
    </div>

    <!-- Orders Table -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-shopping-cart text-primary-500"></i>
          <span class="font-bold">订单列表</span>
        </div>
      </template>

      <template #content>
        <DataTable 
          :value="orders" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
          <Column field="order_no" header="订单号">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i class="pi pi-file text-slate-400"></i>
                <code class="text-xs bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">{{ data.order_no }}</code>
              </div>
            </template>
          </Column>

          <Column field="type" header="类型">
            <template #body="{ data }">
              <Tag 
                :value="getTypeText(data.type)"
                :severity="getTypeSeverity(data.type)"
                class="text-xs"
              />
            </template>
          </Column>

          <Column field="amount" header="金额">
            <template #body="{ data }">
              <span class="font-semibold text-slate-800 dark:text-white">
                ¥{{ data.amount.toFixed(2) }}
              </span>
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
  const map: Record<string, string> = {
    purchase: '购买',
    renew: '续费',
    upgrade: '升级',
  }
  return map[type] || type
}

function getTypeSeverity(type: string) {
  const map: Record<string, string> = {
    purchase: 'primary',
    renew: 'info',
    upgrade: 'warning',
  }
  return map[type] || 'secondary'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
  }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'danger',
  }
  return map[status] || 'secondary'
}
</script>
