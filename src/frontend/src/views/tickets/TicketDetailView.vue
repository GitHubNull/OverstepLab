<template>
  <div v-if="ticket" class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" text rounded size="small" class="text-slate-400" @click="$router.push('/tickets')" />
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-lg font-bold text-slate-800 dark:text-white">{{ ticket.title }}</h2>
            <Tag :value="getStatusText(ticket.status)" :severity="getStatusSeverity(ticket.status)" class="text-[10px]" />
          </div>
          <p class="text-xs text-slate-400 mt-0.5">创建于 {{ formatDate(ticket.created_at) }}</p>
        </div>
      </div>
      <Button
        v-if="ticket.status !== 'closed'"
        label="关闭工单"
        icon="pi pi-check"
        severity="success"
        size="small"
        outlined
        @click="handleClose"
      />
    </div>

    <!-- Original Content -->
    <Card class="shadow-none">
      <template #content>
        <div class="flex gap-3">
          <Avatar
            :label="(ticket as any).user?.username?.charAt(0).toUpperCase() || 'U'"
            shape="circle"
            class="bg-indigo-100 text-indigo-600 text-xs flex-shrink-0"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1.5">
              <span class="font-semibold text-sm text-slate-700 dark:text-white">{{ (ticket as any).user?.username || '用户' }}</span>
              <span class="text-[10px] text-slate-400">{{ formatDate(ticket.created_at) }}</span>
            </div>
            <p class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed">{{ ticket.content }}</p>
          </div>
        </div>
      </template>
    </Card>

    <!-- Replies -->
    <div v-if="replies.length > 0" class="space-y-3">
      <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider">回复 ({{ replies.length }})</h3>

      <Card v-for="reply in replies" :key="reply.id" class="shadow-none">
        <template #content>
          <div class="flex gap-3">
            <Avatar
              :label="reply.user?.username?.charAt(0).toUpperCase() || 'U'"
              shape="circle"
              :class="reply.user?.user_type === 'platform_admin' ? 'bg-red-100 text-red-600 text-xs' : 'bg-indigo-100 text-indigo-600 text-xs'"
            />
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1.5">
                <span class="font-semibold text-sm text-slate-700 dark:text-white">{{ reply.user?.username || '用户' }}</span>
                <Tag v-if="reply.user?.user_type === 'platform_admin'" value="管理员" severity="danger" class="text-[10px]" />
                <span class="text-[10px] text-slate-400">{{ formatDate(reply.created_at) }}</span>
              </div>
              <p class="text-sm text-slate-600 dark:text-slate-300 leading-relaxed">{{ reply.content }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Reply Form -->
    <Card v-if="ticket.status !== 'closed'" class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-reply"></i>
          <span>回复工单</span>
        </div>
      </template>
      <template #content>
        <div class="space-y-3">
          <Textarea v-model="replyContent" rows="4" placeholder="请输入回复内容..." class="w-full" />
          <div class="flex justify-end">
            <Button label="提交回复" icon="pi pi-send" size="small" :loading="submitting" @click="submitReply" />
          </div>
        </div>
      </template>
    </Card>
  </div>

  <!-- Loading -->
  <div v-else class="flex items-center justify-center py-20">
    <ProgressSpinner />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import * as api from '@/api'
import type { Ticket, TicketReply } from '@/types'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Textarea from 'primevue/textarea'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { formatDate } from '@/utils/date'

const route = useRoute()
const toast = useToast()
const confirm = useConfirm()

const ticket = ref<Ticket | null>(null)
const replies = ref<TicketReply[]>([])
const replyContent = ref('')
const submitting = ref(false)

onMounted(async () => {
  const response = await api.getTicketDetail(Number(route.params.id))
  ticket.value = response.data.data!.ticket
  replies.value = response.data.data!.replies
})

function getStatusText(status: string) {
  const map: Record<string, string> = { open: '待处理', replied: '已回复', closed: '已关闭' }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = { open: 'warn', replied: 'info', closed: 'secondary' }
  return map[status] || 'secondary'
}

async function submitReply() {
  if (!replyContent.value.trim()) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请输入回复内容' })
    return
  }
  submitting.value = true
  try {
    await api.replyTicket(Number(route.params.id), replyContent.value)
    toast.add({ severity: 'success', summary: '成功', detail: '回复已提交' })
    replyContent.value = ''
    const response = await api.getTicketDetail(Number(route.params.id))
    ticket.value = response.data.data!.ticket
    replies.value = response.data.data!.replies
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '提交失败' })
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  confirm.require({
    message: '确定要关闭此工单吗？',
    header: '确认关闭',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await api.closeTicket(Number(route.params.id))
        toast.add({ severity: 'success', summary: '成功', detail: '工单已关闭' })
        const response = await api.getTicketDetail(Number(route.params.id))
        ticket.value = response.data.data!.ticket
        replies.value = response.data.data!.replies
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '关闭失败' })
      }
    },
  })
}
</script>
