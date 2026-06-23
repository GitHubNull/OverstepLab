<template>
  <div v-if="ticket" class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-3">
        <Button icon="pi pi-arrow-left" text rounded size="small" class="text-[var(--text-tertiary)]" @click="$router.push('/tickets')" />
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-lg font-bold text-[var(--text-primary)]">{{ ticket.title }}</h1>
            <Tag :value="getStatusText(ticket.status)" :severity="getStatusSeverity(ticket.status)" class="text-[10px]" />
          </div>
          <p class="text-xs text-[var(--text-tertiary)] mt-0.5">创建于 {{ formatDate(ticket.created_at) }}</p>
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
    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl p-5">
      <div class="flex gap-3">
        <Avatar
          :label="(ticket as any).user?.username?.charAt(0).toUpperCase() || 'U'"
          shape="circle"
          class="bg-[var(--primary-subtle)] text-[var(--primary)] text-xs flex-shrink-0"
        />
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 mb-1.5">
            <span class="font-semibold text-sm text-[var(--text-primary)]">{{ (ticket as any).user?.username || '用户' }}</span>
            <span class="text-[10px] text-[var(--text-tertiary)]">{{ formatDate(ticket.created_at) }}</span>
          </div>
          <p class="text-sm text-[var(--text-secondary)] leading-relaxed">{{ ticket.content }}</p>
        </div>
      </div>
    </div>

    <!-- Replies -->
    <div v-if="replies.length > 0" class="space-y-3">
      <h3 class="text-xs font-semibold text-[var(--text-tertiary)] uppercase tracking-wider">回复 ({{ replies.length }})</h3>

      <div
        v-for="reply in replies"
        :key="reply.id"
        class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl p-5 relative"
        :class="reply.user?.user_type === 'platform_admin' ? 'border-l-[3px]' : ''"
        :style="reply.user?.user_type === 'platform_admin' ? { borderLeftColor: 'var(--danger)' } : {}"
      >
        <div class="flex gap-3">
          <Avatar
            :label="reply.user?.username?.charAt(0).toUpperCase() || 'U'"
            shape="circle"
            :class="reply.user?.user_type === 'platform_admin' ? 'bg-[var(--danger-subtle)] text-[var(--danger)] text-xs' : 'bg-[var(--primary-subtle)] text-[var(--primary)] text-xs'"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1.5">
              <span class="font-semibold text-sm text-[var(--text-primary)]">{{ reply.user?.username || '用户' }}</span>
              <Tag v-if="reply.user?.user_type === 'platform_admin'" value="管理员" severity="danger" class="text-[10px]" />
              <span class="text-[10px] text-[var(--text-tertiary)]">{{ formatDate(reply.created_at) }}</span>
            </div>
            <p class="text-sm text-[var(--text-secondary)] leading-relaxed">{{ reply.content }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Reply Form -->
    <div v-if="ticket.status !== 'closed'" class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div class="px-5 py-4 border-b border-[var(--border-default)]">
        <div class="section-title">
          <i class="pi pi-reply"></i>
          <span>回复工单</span>
        </div>
      </div>
      <div class="p-4 space-y-3">
        <Textarea v-model="replyContent" rows="4" placeholder="请输入回复内容..." class="w-full" />
        <div class="flex justify-end">
          <Button label="提交回复" icon="pi pi-send" size="small" :loading="submitting" @click="submitReply" />
        </div>
      </div>
    </div>
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
  try {
    const response = await api.getTicketDetail(Number(route.params.id))
    ticket.value = response.data.data!.ticket
    replies.value = response.data.data!.replies
  } catch (e: any) {
    toast.add({
      severity: 'error',
      summary: '加载失败',
      detail: e.response?.data?.message || '无法获取工单详情',
      life: 3000,
    })
  }
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
