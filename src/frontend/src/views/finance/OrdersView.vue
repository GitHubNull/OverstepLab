<template>
  <div class="space-y-5">
    <PageHeader title="订单管理" description="查看您的购买和续费记录" />

    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div class="px-5 py-4 border-b border-[var(--border-default)]">
        <div class="section-title">
          <i class="pi pi-receipt"></i>
          <span>订单列表</span>
        </div>
      </div>
      <div class="p-0">
        <DataTable :value="orders" class="p-datatable-sm" :rows="10" paginator :loading="loading">
          <template #empty>
            <div class="text-center py-10">
              <div class="w-14 h-14 rounded-2xl bg-[var(--bg-surface-hover)] flex items-center justify-center mx-auto mb-3">
                <i class="pi pi-inbox text-2xl text-[var(--text-tertiary)]"></i>
              </div>
              <p class="text-[var(--text-secondary)] text-sm font-medium mb-1">暂无订单记录</p>
              <p class="text-[var(--text-tertiary)] text-xs">购买 VPS 实例后会自动生成订单</p>
              <router-link to="/vps" class="inline-flex items-center gap-1.5 mt-3 text-xs font-medium text-[var(--primary)] hover:underline">
                <i class="pi pi-plus text-[10px]"></i>
                去创建 VPS
              </router-link>
            </div>
          </template>
          <Column field="order_no" header="订单号">
            <template #body="{ data }">
              <code class="text-xs bg-[var(--bg-base)] px-1.5 py-0.5 rounded text-[var(--text-secondary)] mono">{{ data.order_no }}</code>
            </template>
          </Column>
          <Column field="type" header="类型">
            <template #body="{ data }">
              <Tag :value="getTypeText(data.type)" :severity="getTypeSeverity(data.type)" class="text-[10px]" />
            </template>
          </Column>
          <Column field="amount" header="金额">
            <template #body="{ data }">
              <span class="font-semibold text-sm text-[var(--text-primary)] mono">¥{{ data.amount.toFixed(2) }}</span>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="getStatusText(data.status)" :severity="getStatusSeverity(data.status)" class="text-[10px]" />
            </template>
          </Column>
          <Column header="创建时间">
            <template #body="{ data }">
              <span class="text-xs text-[var(--text-tertiary)]">{{ formatDate(data.created_at) }}</span>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import type { Order } from '@/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import { formatDate } from '@/utils/date'
import PageHeader from '@/components/PageHeader.vue'

const orders = ref<Order[]>([])
const loading = ref(false)

onMounted(async () => {
  try {
    const response = await api.getOrders()
    orders.value = response.data.data!
  } catch {
    // Error handled silently - orders will remain empty
  }
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
