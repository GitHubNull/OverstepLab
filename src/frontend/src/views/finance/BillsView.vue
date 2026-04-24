<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="page-header">
        <h2>账单管理</h2>
        <p>查看您的消费记录和账户余额</p>
      </div>
      <div class="flex gap-2">
        <Button label="充值" icon="pi pi-plus" size="small" @click="showRechargeDialog = true" />
        <Button label="导出" icon="pi pi-download" size="small" severity="secondary" outlined @click="handleExport" />
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
      <Card class="shadow-none stat-card stat-card-green">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-wallet text-emerald-500"></i>
            </div>
            <div>
              <p class="text-[11px] text-slate-400 font-medium">账户余额</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">¥{{ balance.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-none stat-card stat-card-red">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-red-50 dark:bg-red-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-arrow-down text-red-500"></i>
            </div>
            <div>
              <p class="text-[11px] text-slate-400 font-medium">总支出</p>
              <p class="text-xl font-bold text-red-500">¥{{ totalExpense.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-none stat-card stat-card-blue">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-50 dark:bg-blue-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-receipt text-blue-500"></i>
            </div>
            <div>
              <p class="text-[11px] text-slate-400 font-medium">账单数量</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ bills.length }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Bills Table -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-list"></i>
          <span>账单明细</span>
        </div>
      </template>
      <template #content>
        <DataTable :value="bills" stripedRows class="p-datatable-sm" :rows="10" paginator>
          <Column field="id" header="账单号">
            <template #body="{ data }">
              <code class="text-xs bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded text-slate-500">#{{ data.id }}</code>
            </template>
          </Column>
          <Column field="type" header="类型">
            <template #body="{ data }">
              <Tag :value="data.type === 'expense' ? '消费' : '充值'" :severity="data.type === 'expense' ? 'danger' : 'success'" class="text-[10px]" />
            </template>
          </Column>
          <Column field="amount" header="金额">
            <template #body="{ data }">
              <span class="font-semibold text-sm" :class="data.amount > 0 ? 'text-emerald-600' : 'text-red-500'">
                {{ data.amount > 0 ? '+' : '' }}¥{{ data.amount.toFixed(2) }}
              </span>
            </template>
          </Column>
          <Column field="description" header="描述">
            <template #body="{ data }">
              <span class="text-sm text-slate-500">{{ data.description }}</span>
            </template>
          </Column>
          <Column header="时间">
            <template #body="{ data }">
              <span class="text-xs text-slate-400">{{ formatDate(data.created_at) }}</span>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Recharge Dialog -->
    <Dialog v-model:visible="showRechargeDialog" header="余额充值" modal :style="{ width: '400px' }">
      <div class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-2 uppercase tracking-wider">充值金额</label>
          <div class="grid grid-cols-3 gap-2 mb-3">
            <Button v-for="amount in [50, 100, 200, 500, 1000, 2000]" :key="amount"
              :label="'¥' + amount"
              :severity="rechargeAmount === amount ? 'primary' : 'secondary'"
              :outlined="rechargeAmount !== amount"
              size="small"
              @click="rechargeAmount = amount; rechargeAmountInput = String(amount)"
              class="!text-xs"
            />
          </div>
          <InputText v-model="rechargeAmountInput" class="w-full" placeholder="或输入自定义金额" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showRechargeDialog = false" />
        <Button label="确认充值" icon="pi pi-check" size="small" :loading="recharging" @click="handleRecharge" />
      </template>
    </Dialog>
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
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'
import { formatDate } from '@/utils/date'

const toast = useToast()
const bills = ref<Bill[]>([])
const showRechargeDialog = ref(false)
const rechargeAmount = ref(100)
const rechargeAmountInput = ref('100')
const recharging = ref(false)

onMounted(async () => {
  const response = await api.getBills()
  bills.value = response.data.data!
})

const balance = computed(() => bills.value.reduce((sum, b) => sum + b.amount, 0))
const totalExpense = computed(() => Math.abs(bills.value.filter(b => b.type === 'expense').reduce((sum, b) => sum + b.amount, 0)))

async function handleRecharge() {
  const amount = Number(rechargeAmountInput.value) || rechargeAmount.value
  if (!amount || amount <= 0) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请输入有效金额', life: 2000 })
    return
  }
  recharging.value = true
  try {
    await api.recharge(amount)
    toast.add({ severity: 'success', summary: '成功', detail: `已充值 ¥${amount}`, life: 2000 })
    showRechargeDialog.value = false
    rechargeAmount.value = 100
    rechargeAmountInput.value = '100'
    const response = await api.getBills()
    bills.value = response.data.data!
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '充值失败', life: 3000 })
  } finally {
    recharging.value = false
  }
}

async function handleExport() {
  try {
    const response = await api.exportBills()
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'bills.csv')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    toast.add({ severity: 'success', summary: '成功', detail: '账单已导出', life: 2000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: '导出失败', life: 3000 })
  }
}
</script>
