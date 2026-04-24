<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">企业管理</h2>
        <p class="text-slate-500">管理平台内所有企业账户</p>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-building text-blue-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">企业总数</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ companies.length }}</p>
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
              <p class="text-sm text-slate-500">正常企业</p>
              <p class="text-xl font-bold text-green-600">{{ activeCompanies }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg flex items-center justify-center">
              <i class="pi pi-wallet text-yellow-600"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">总余额</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">¥{{ totalBalance.toFixed(2) }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Companies Table -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-building text-primary-500"></i>
          <span class="font-bold">企业列表</span>
        </div>
      </template>

      <template #content>
        <DataTable 
          :value="companies" 
          stripedRows
          class="p-datatable-sm"
          :rows="10"
          paginator
        >
          <Column field="id" header="ID">
            <template #body="{ data }">
              <code class="text-xs bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">{{ data.id }}</code>
            </template>
          </Column>

          <Column field="name" header="企业名称">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 bg-primary-100 dark:bg-primary-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-building text-primary-600 text-sm"></i>
                </div>
                <span class="font-medium">{{ data.name }}</span>
              </div>
            </template>
          </Column>

          <Column field="license_no" header="营业执照号">
            <template #body="{ data }">
              <span class="text-sm text-slate-500">{{ data.license_no || '-' }}</span>
            </template>
          </Column>

          <Column field="balance" header="余额">
            <template #body="{ data }">
              <span class="font-semibold text-slate-800 dark:text-white">¥{{ data.balance.toFixed(2) }}</span>
            </template>
          </Column>

          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag 
                :value="data.status === 'active' ? '正常' : '暂停'"
                :severity="data.status === 'active' ? 'success' : 'warning'"
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
