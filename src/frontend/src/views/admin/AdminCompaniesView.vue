<template>
  <div>
    <PageHeader title="企业管理" description="管理平台内所有企业账户" />

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-5">
      <StatCard color="#4f46e5" icon="pi pi-building" :value="companies.length" label="企业总数" />
      <StatCard color="#10b981" icon="pi pi-check" :value="activeCompanies" label="正常企业" />
      <StatCard color="#f59e0b" icon="pi pi-wallet" :value="'¥' + totalBalance.toFixed(2)" label="总余额" />
    </div>

    <!-- Companies Table -->
    <Card>
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
              <span class="mono text-[10px] px-1.5 py-0.5 rounded bg-[var(--bg-base)] text-[var(--text-secondary)]">{{ data.id }}</span>
            </template>
          </Column>
          <Column field="name" header="企业名称">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 bg-[var(--primary-subtle)] rounded-lg flex items-center justify-center">
                  <i class="pi pi-building text-[var(--primary)] text-[10px]"></i>
                </div>
                <span class="font-medium text-sm text-[var(--text-primary)]">{{ data.name }}</span>
              </div>
            </template>
          </Column>
          <Column field="license_no" header="营业执照号">
            <template #body="{ data }">
              <span class="text-xs text-[var(--text-secondary)]">{{ data.license_no || '-' }}</span>
            </template>
          </Column>
          <Column field="balance" header="余额">
            <template #body="{ data }">
              <span class="font-semibold text-sm mono text-[var(--text-primary)]">¥{{ data.balance.toFixed(2) }}</span>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="data.status === 'active' ? '正常' : '暂停'" :severity="data.status === 'active' ? 'success' : 'warn'" class="text-[10px]" />
            </template>
          </Column>
          <Column header="创建时间">
            <template #body="{ data }">
              <span class="text-xs text-[var(--text-tertiary)]">{{ formatDate(data.created_at) }}</span>
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
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
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
