<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">API Key 管理</h2>
        <p class="text-slate-500">管理您的 API 访问凭证</p>
      </div>
      <Button 
        label="创建 API Key" 
        icon="pi pi-plus" 
        severity="primary"
      />
    </div>

    <!-- API Keys Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Card 
        v-for="key in keys" 
        :key="key.id"
        class="shadow-sm hover:shadow-md transition-shadow"
      >
        <template #title>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <div class="w-10 h-10 bg-yellow-100 dark:bg-yellow-900/30 rounded-xl flex items-center justify-center">
                <i class="pi pi-key text-lg text-yellow-600"></i>
              </div>
              <span class="font-semibold text-slate-800 dark:text-white">{{ key.name }}</span>
            </div>
            <Tag 
              :value="key.status === 'active' ? '正常' : '已吊销'"
              :severity="key.status === 'active' ? 'success' : 'danger'"
              class="text-xs"
            />
          </div>
        </template>

        <template #content>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-500">Key 前缀</span>
              <code class="text-sm bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded">{{ key.key_prefix }}****</code>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-500">权限</span>
              <span class="text-sm text-slate-700 dark:text-slate-300">{{ key.permissions }}</span>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-500">创建时间</span>
              <span class="text-sm text-slate-500">{{ formatDate(key.created_at) }}</span>
            </div>

            <div class="flex gap-2 pt-2">
              <Button 
                label="复制" 
                icon="pi pi-copy" 
                text 
                size="small"
                class="flex-1"
              />
              <Button 
                label="删除" 
                icon="pi pi-trash" 
                text 
                size="small"
                severity="danger"
                class="flex-1"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-if="keys.length === 0" class="text-center py-16">
      <div class="w-20 h-20 bg-slate-100 dark:bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4">
        <i class="pi pi-key text-4xl text-slate-400"></i>
      </div>
      <h3 class="text-lg font-semibold text-slate-700 dark:text-slate-300 mb-2">暂无 API Key</h3>
      <p class="text-slate-500 mb-4">创建 API Key 以访问平台 API</p>
      <Button label="创建 API Key" icon="pi pi-plus" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import type { APIKey } from '@/types'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import { formatDate } from '@/utils/date'

const keys = ref<APIKey[]>([])

onMounted(async () => {
  const response = await api.getApiKeys()
  keys.value = response.data.data!
})
</script>
