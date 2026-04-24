<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">账单管理</h2>
        <p class="text-slate-500">查看您的消费记录和账户余额</p>
      </div>
      <div class="flex gap-2">
        <Button 
          label="充值" 
          icon="pi pi-plus" 
          severity="primary"
        />
        <Button 
          label="导出" 
          icon="pi pi-download" 
          severity="secondary"
          outlined
        />
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-wallet text-xl text-green-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">账户余额</p>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">¥{{ balance.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-red-100 dark:bg-red-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-arrow-down text-xl text-red-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">总支出</p>
              <p class="text-2xl font-bold text-red-600">¥{{ totalExpense.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-receipt text-xl text-blue-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">账单数量</p>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ bills.length }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Bills Table -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-list text-primary-500"></i>
          <span class="font-bold">账单明细</span>
        </div>
      </template>

      <template #content>
        <DataTable 
          :value="bills" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
          <Column field="id" header="账单号">
            <template #body="{ data }">
              <code class="text-xs bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">#{{ data.id }}</code>
            </template>
          </Column>

          <Column field="type" header="类型">
            <template #body="{ data }">
              <Tag 
                :value="data.type === 'expense' ? '消费' : '充值'"
                :severity="data.type === 'expense' ? 'danger' : 'success'"
                class="text-xs"
              />
            </template>
          </Column>

          <Column field="amount" header="金额">
            <template #body="{ data }">
              <span 
                class="font-semibold"
                :class="data.amount > 0 ? 'text-green-600' : 'text-red-600'"
              >
                {{ data.amount > 0 ? '+' : '' }}¥{{ data.amount.toFixed(2) }}
              </span>
            </template>
          </Column>

          <Column field="description" header="描述">
            <template #body="{ data }">
              <span class="text-sm text-slate-600">{{ data.description }}</span>
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
import type { Bill } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import { formatDate } from '@/utils/date'

const bills = ref<Bill[]>([])

onMounted(async () => {
  const response = await api.getBills()
  bills.value = response.data.data!
})

const balance = computed(() => {
  return bills.value.reduce((sum, b) => sum + b.amount, 0)
})

const totalExpense = computed(() => {
  return Math.abs(bills.value.filter(b => b.type === 'expense').reduce((sum, b) => sum + b.amount, 0))
})
</script>
