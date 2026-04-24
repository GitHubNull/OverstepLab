<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="page-header">
        <h2>API Key 管理</h2>
        <p>管理您的 API 访问凭证</p>
      </div>
      <Button label="创建 API Key" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
    </div>

    <!-- API Keys Grid -->
    <div v-if="keys.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <Card v-for="key in keys" :key="key.id" class="shadow-none hover:shadow-md">
        <template #content>
          <div class="space-y-3">
            <!-- Header -->
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="w-8 h-8 bg-amber-50 dark:bg-amber-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-key text-amber-500 text-xs"></i>
                </div>
                <span class="font-semibold text-sm text-slate-700 dark:text-white">{{ key.name }}</span>
              </div>
              <Tag :value="key.status === 'active' ? '正常' : '已吊销'" :severity="key.status === 'active' ? 'success' : 'danger'" class="text-[10px]" />
            </div>

            <!-- Details -->
            <div class="space-y-2">
              <div class="info-row !py-1.5">
                <span class="text-[10px] text-slate-400">Key 前缀</span>
                <code class="text-xs bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded font-mono">{{ key.key_prefix }}****</code>
              </div>
              <div class="info-row !py-1.5">
                <span class="text-[10px] text-slate-400">权限</span>
                <span class="text-xs text-slate-600 dark:text-slate-300">{{ key.permissions }}</span>
              </div>
              <div class="info-row !py-1.5">
                <span class="text-[10px] text-slate-400">创建时间</span>
                <span class="text-xs text-slate-400">{{ formatDate(key.created_at) }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-1 border-t border-slate-100 dark:border-slate-700/50">
              <Button label="复制" icon="pi pi-copy" text size="small" class="flex-1 !text-xs" @click="handleCopy(key)" />
              <Button label="删除" icon="pi pi-trash" text size="small" severity="danger" class="flex-1 !text-xs" @click="handleDelete(key)" />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state">
      <div class="empty-state-icon">
        <i class="pi pi-key text-3xl text-slate-300"></i>
      </div>
      <h3 class="text-base font-semibold text-slate-600 dark:text-slate-300 mb-1">暂无 API Key</h3>
      <p class="text-slate-400 text-sm mb-4">创建 API Key 以访问平台 API</p>
      <Button label="创建 API Key" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
    </div>

    <!-- Create API Key Dialog -->
    <Dialog v-model:visible="showCreateDialog" header="创建 API Key" modal :style="{ width: '420px' }">
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">名称</label>
          <InputText v-model="createForm.name" class="w-full" placeholder="例如: my-api-key" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">权限</label>
          <Select v-model="createForm.permissions" :options="permissionOptions" placeholder="选择权限" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showCreateDialog = false" />
        <Button label="创建" icon="pi pi-check" size="small" :loading="creating" @click="handleCreate" />
      </template>
    </Dialog>

    <!-- New Key Display Dialog -->
    <Dialog v-model:visible="showNewKeyDialog" header="API Key 已创建" modal :style="{ width: '480px' }">
      <div class="space-y-3">
        <div class="bg-amber-50 dark:bg-amber-900/10 border border-amber-200 dark:border-amber-800/50 rounded-xl p-4">
          <p class="text-xs text-amber-600 dark:text-amber-400 font-semibold mb-2">请立即复制此 Key，关闭后将无法再次查看！</p>
          <div class="flex items-center gap-2">
            <code class="flex-1 text-xs bg-white dark:bg-slate-800 px-3 py-2 rounded-lg break-all font-mono">{{ newKeyValue }}</code>
            <Button icon="pi pi-copy" text size="small" @click="copyToClipboard(newKeyValue)" />
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="我已保存" size="small" @click="showNewKeyDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import type { APIKey } from '@/types'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { formatDate } from '@/utils/date'

const toast = useToast()
const confirm = useConfirm()
const keys = ref<APIKey[]>([])
const showCreateDialog = ref(false)
const showNewKeyDialog = ref(false)
const creating = ref(false)
const newKeyValue = ref('')
const createForm = ref({ name: '', permissions: 'read' })
const permissionOptions = ['read', 'write', 'read-write']

onMounted(() => fetchKeys())

async function fetchKeys() {
  const response = await api.getApiKeys()
  keys.value = response.data.data!
}

async function handleCreate() {
  if (!createForm.value.name) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请输入名称', life: 2000 })
    return
  }
  creating.value = true
  try {
    const response = await api.createApiKey(createForm.value)
    newKeyValue.value = (response.data.data as any)?.key || 'sk_****'
    showCreateDialog.value = false
    showNewKeyDialog.value = true
    createForm.value = { name: '', permissions: 'read' }
    await fetchKeys()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '创建失败', life: 3000 })
  } finally {
    creating.value = false
  }
}

function handleCopy(key: APIKey) {
  copyToClipboard(key.key_prefix + '****')
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    toast.add({ severity: 'success', summary: '已复制', detail: '已复制到剪贴板', life: 2000 })
  }).catch(() => {
    toast.add({ severity: 'error', summary: '复制失败', life: 2000 })
  })
}

function handleDelete(key: APIKey) {
  confirm.require({
    message: `确定要删除 API Key "${key.name}" 吗？此操作不可恢复。`,
    header: '确认删除',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.deleteApiKey(key.id)
        toast.add({ severity: 'success', summary: '成功', detail: 'API Key 已删除', life: 2000 })
        await fetchKeys()
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '删除失败', life: 3000 })
      }
    },
  })
}
</script>
