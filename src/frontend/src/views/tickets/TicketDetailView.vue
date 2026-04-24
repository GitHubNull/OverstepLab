<template>
  <div v-if="ticket" class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <Button 
          icon="pi pi-arrow-left" 
          text 
          rounded
          @click="$router.push('/tickets')"
        />
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-2xl font-bold text-slate-800 dark:text-white">{{ ticket.title }}</h2>
            <Tag 
              :value="getStatusText(ticket.status)"
              :severity="getStatusSeverity(ticket.status)"
            />
          </div>
          <p class="text-sm text-slate-500">创建于 {{ formatDate(ticket.created_at) }}</p>
        </div>
      </div>

      <div class="flex gap-2">
        <Button 
          v-if="ticket.status !== 'closed'"
          label="关闭工单" 
          icon="pi pi-check" 
          severity="success"
        />
      </div>
    </div>

    <!-- Original Content -->
    <Card class="shadow-sm">
      <template #content>
        <div class="flex gap-4">
          <Avatar 
            :label="(ticket as any).user?.username?.charAt(0).toUpperCase() || 'U'"
            shape="circle"
            class="bg-primary-100 text-primary-700"
          />
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <span class="font-semibold text-slate-800 dark:text-white">{{ (ticket as any).user?.username || '用户' }}</span>
              <span class="text-xs text-slate-500">{{ formatDate(ticket.created_at) }}</span>
            </div>
            <p class="text-slate-700 dark:text-slate-300">{{ ticket.content }}</p>
          </div>
        </div>
      </template>
    </Card>

    <!-- Replies -->
    <div v-if="replies.length > 0" class="space-y-4">
      <h3 class="font-semibold text-slate-700 dark:text-slate-300">回复 ({{ replies.length }})</h3>
      
      <Card 
        v-for="reply in replies" 
        :key="reply.id"
        class="shadow-sm"
      >
        <template #content>
          <div class="flex gap-4">
            <Avatar 
              :label="reply.user?.username?.charAt(0).toUpperCase() || 'U'"
              shape="circle"
              :class="reply.user?.user_type === 'platform_admin' ? 'bg-red-100 text-red-700' : 'bg-primary-100 text-primary-700'"
            />
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <span class="font-semibold text-slate-800 dark:text-white">{{ reply.user?.username || '用户' }}</span>
                <Tag 
                  v-if="reply.user?.user_type === 'platform_admin'"
                  value="管理员"
                  severity="danger"
                  class="text-xs"
                />
                <span class="text-xs text-slate-500">{{ formatDate(reply.created_at) }}</span>
              </div>
              <p class="text-slate-700 dark:text-slate-300">{{ reply.content }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Reply Form -->
    <Card v-if="ticket.status !== 'closed'" class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-reply text-primary-500"></i>
          <span class="font-bold">回复工单</span>
        </div>
      </template>

      <template #content>
        <div class="space-y-4">
          <Textarea 
            v-model="replyContent"
            rows="4"
            placeholder="请输入回复内容..."
            class="w-full"
          />
          <div class="flex justify-end">
            <Button 
              label="提交回复" 
              icon="pi pi-send"
              severity="primary"
              :loading="submitting"
              @click="submitReply"
            />
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
import { formatDate } from '@/utils/date'

const route = useRoute()
const toast = useToast()

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
  const map: Record<string, string> = {
    open: '待处理',
    replied: '已回复',
    closed: '已关闭',
  }
  return map[status] || status
}

function getStatusSeverity(status: string) {
  const map: Record<string, string> = {
    open: 'warning',
    replied: 'info',
    closed: 'secondary',
  }
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
    // Refresh ticket data
    const response = await api.getTicketDetail(Number(route.params.id))
    ticket.value = response.data.data!.ticket
    replies.value = response.data.data!.replies
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '提交失败' })
  } finally {
    submitting.value = false
  }
}
</script>
