<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="page-header">
      <h2>企业管理</h2>
      <p>管理平台内所有企业账户</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <Card class="shadow-none stat-card stat-card-blue">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-blue-50 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-building text-blue-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">企业总数</p>
              <p class="text-lg font-bold text-slate-800 dark:text-white">{{ companies.length }}</p>
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
              <p class="text-[10px] text-slate-400 font-medium">正常企业</p>
              <p class="text-lg font-bold text-emerald-600">{{ activeCompanies }}</p>
            </div>
          </div>
        </template>
      </Card>
      <Card class="shadow-none stat-card stat-card-orange">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 bg-amber-50 dark:bg-amber-900/20 rounded-lg flex items-center justify-center">
              <i class="pi pi-wallet text-amber-500 text-sm"></i>
            </div>
            <div>
              <p class="text-[10px] text-slate-400 font-medium">总余额</p>
              <p class="text-lg font-bold text-slate-800 dark:text-white">¥{{ totalBalance.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Companies Table -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-building"></i>
          <span>企业列表</span>
        </div>
      </template>
      <template #content>
        <DataTable :value="companies" stripedRows class="p-datatable-sm" :rows="10" paginator>
          <Column field="id" header="ID" style="width: 60px">
            <template #body="{ data }">
              <code class="text-[10px] bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded">{{ data.id }}</code>
            </template>
          </Column>
          <Column field="name" header="企业名称">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 bg-indigo-50 dark:bg-indigo-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-building text-indigo-500 text-[10px]"></i>
                </div>
                <span class="font-medium text-sm text-slate-700 dark:text-slate-200">{{ data.name }}</span>
              </div>
            </template>
          </Column>
          <Column field="license_no" header="营业执照号">
            <template #body="{ data }">
              <span class="text-xs text-slate-500">{{ data.license_no || '-' }}</span>
            </template>
          </Column>
          <Column field="balance" header="余额">
            <template #body="{ data }">
              <span class="font-semibold text-sm text-slate-700 dark:text-slate-200">¥{{ data.balance.toFixed(2) }}</span>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="data.status === 'active' ? '正常' : '暂停'" :severity="data.status === 'active' ? 'success' : 'warn'" class="text-[10px]" />
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
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import type { Company } from '@/types'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import { formatDate } from '@/utils/date'

const companies = ref<Company[]>([])

onMounted(async () => {
  const response = await api.adminListCompanies()
  companies.value = response.data.data!
})

const activeCompanies = computed(() => companies.value.filter(c => c.status === 'active').length)
const totalBalance = computed(() => companies.value.reduce((sum, c) => sum + c.balance, 0))
</script>
