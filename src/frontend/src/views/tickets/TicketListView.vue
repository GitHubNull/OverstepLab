<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="page-header">
        <h2>工单系统</h2>
        <p>提交和管理技术支持请求</p>
      </div>
      <Button label="新建工单" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
    </div>

    <!-- Tickets List -->
    <Card class="shadow-none">
      <template #content>
        <DataTable :value="tickets" stripedRows class="p-datatable-sm" :rows="10" paginator>
          <Column field="title" header="标题">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i class="pi text-xs" :class="getStatusIcon(data.status)"></i>
                <span class="font-medium text-sm text-slate-700 dark:text-slate-200">{{ data.title }}</span>
              </div>
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
          <Column header="操作" style="width: 80px">
            <template #body="{ data }">
              <Button label="查看" icon="pi pi-eye" text size="small" class="!text-xs" @click="$router.push(`/tickets/${data.id}`)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Create Ticket Dialog -->
    <Dialog v-model:visible="showCreateDialog" header="新建工单" modal :style="{ width: '500px' }">
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">标题</label>
          <InputText v-model="createForm.title" class="w-full" placeholder="请简述您的问题" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">内容</label>
          <Textarea v-model="createForm.content" rows="5" class="w-full" placeholder="请详细描述您遇到的问题..." />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showCreateDialog = false" />
        <Button label="提交工单" icon="pi pi-send" size="small" :loading="creating" @click="handleCreate" />
      </template>
    </Dialog>
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
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import { useToast } from 'primevue/usetoast'
import { formatDate } from '@/utils/date'

const toast = useToast()
const tickets = ref<Ticket[]>([])
const showCreateDialog = ref(false)
const creating = ref(false)
const createForm = ref({ title: '', content: '' })

onMounted(() => fetchTickets())

async function fetchTickets() {
  const response = await api.getTickets()
  tickets.value = response.data.data!
}

async function handleCreate() {
  if (!createForm.value.title || !createForm.value.content) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请填写标题和内容', life: 2000 })
    return
  }
  creating.value = true
  try {
    await api.createTicket(createForm.value)
    toast.add({ severity: 'success', summary: '成功', detail: '工单已创建', life: 2000 })
    showCreateDialog.value = false
    createForm.value = { title: '', content: '' }
    await fetchTickets()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '创建失败', life: 3000 })
  } finally {
    creating.value = false
  }
}

function getStatusText(status: string) {
  const map: Record<string, string> = { open: '待处理', replied: '已回复', closed: '已关闭' }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = { open: 'warn', replied: 'info', closed: 'secondary' }
  return map[status] || 'secondary'
}

function getStatusIcon(status: string) {
  const map: Record<string, string> = {
    open: 'pi-envelope text-amber-500',
    replied: 'pi-envelope-open text-blue-500',
    closed: 'pi-check-circle text-slate-400',
  }
  return map[status] || 'pi-envelope'
}
</script>
