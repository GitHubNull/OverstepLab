<template>
  <div class="space-y-6">
    <!-- Welcome Banner -->
    <div class="bg-gradient-to-r from-primary-600 to-primary-700 rounded-2xl p-6 text-white shadow-lg">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold mb-2">欢迎回来，{{ authStore.user?.username }}！</h1>
          <p class="text-primary-100">
            {{ authStore.securityMode === 'vulnerable' 
              ? '当前处于漏洞模式，所有越权漏洞均可触发' 
              : '当前处于安全模式，所有漏洞已被修复' }}
          </p>
        </div>
        <div class="hidden md:block">
          <div 
            class="w-16 h-16 rounded-2xl flex items-center justify-center"
            :class="authStore.securityMode === 'vulnerable' ? 'bg-red-500/20' : 'bg-green-500/20'"
          >
            <i 
              class="pi text-3xl"
              :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-300' : 'pi-lock text-green-300'"
            ></i>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <Card class="shadow-sm hover:shadow-md transition-shadow">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-server text-xl text-blue-600 dark:text-blue-400"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">VPS 总数</p>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.vpsList.length }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm hover:shadow-md transition-shadow">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-green-100 dark:bg-green-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-check-circle text-xl text-green-600 dark:text-green-400"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">运行中</p>
              <p class="text-2xl font-bold text-green-600">{{ runningCount }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm hover:shadow-md transition-shadow">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-red-100 dark:bg-red-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-times-circle text-xl text-red-600 dark:text-red-400"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">已停止</p>
              <p class="text-2xl font-bold text-red-600">{{ stoppedCount }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm hover:shadow-md transition-shadow cursor-pointer">
        <template #content>
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-purple-100 dark:bg-purple-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-flag text-xl text-purple-600 dark:text-purple-400"></i>
            </div>
            <div>
              <p class="text-sm text-slate-500">挑战完成</p>
              <p class="text-2xl font-bold text-purple-600">{{ completedChallenges }}/13</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- VPS List -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <i class="pi pi-server text-primary-500"></i>
            <span class="font-bold">我的 VPS</span>
          </div>
          <Button 
            label="查看全部" 
            icon="pi pi-arrow-right" 
            text 
            size="small"
            @click="$router.push('/vps')"
          />
        </div>
      </template>
      <template #content>
        <DataTable 
          :value="vpsStore.vpsList" 
          stripedRows 
          :loading="vpsStore.loading"
          class="p-datatable-sm"
          :rows="5"
          paginator
          :rowsPerPageOptions="[5, 10, 20]"
        >
          <Column field="name" header="名称">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <div 
                  class="w-2 h-2 rounded-full"
                  :class="data.status === 'running' ? 'bg-green-500' : 'bg-red-500'"
                ></div>
                <span class="font-medium">{{ data.name }}</span>
              </div>
            </template>
          </Column>
          <Column field="ip_address" header="IP 地址">
            <template #body="{ data }">
              <code class="bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded text-sm">{{ data.ip_address }}</code>
            </template>
          </Column>
          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag 
                :severity="data.status === 'running' ? 'success' : 'danger'" 
                :value="data.status === 'running' ? '运行中' : '已停止'"
                class="text-xs"
              />
            </template>
          </Column>
          <Column field="os_image" header="系统">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i class="pi pi-desktop text-slate-400"></i>
                <span class="text-sm">{{ data.os_image }}</span>
              </div>
            </template>
          </Column>
          <Column header="操作" style="width: 100px">
            <template #body="{ data }">
              <Button 
                icon="pi pi-eye" 
                text 
                rounded 
                size="small"
                @click="$router.push(`/vps/${data.id}`)"
                title="查看详情"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card class="shadow-sm hover:shadow-md transition-shadow cursor-pointer">
        <template #content>
          <div class="flex items-center gap-4 p-2">
            <div class="w-12 h-12 bg-orange-100 dark:bg-orange-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-flag text-xl text-orange-600"></i>
            </div>
            <div>
              <p class="font-semibold text-slate-800 dark:text-white">漏洞挑战</p>
              <p class="text-sm text-slate-500">开始您的越权测试练习</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm hover:shadow-md transition-shadow cursor-pointer">
        <template #content>
          <div class="flex items-center gap-4 p-2">
            <div class="w-12 h-12 bg-cyan-100 dark:bg-cyan-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-book text-xl text-cyan-600"></i>
            </div>
            <div>
              <p class="font-semibold text-slate-800 dark:text-white">学习文档</p>
              <p class="text-sm text-slate-500">查看漏洞利用教程</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-sm hover:shadow-md transition-shadow cursor-pointer">
        <template #content>
          <div class="flex items-center gap-4 p-2">
            <div class="w-12 h-12 bg-pink-100 dark:bg-pink-900/30 rounded-xl flex items-center justify-center">
              <i class="pi pi-cog text-xl text-pink-600"></i>
            </div>
            <div>
              <p class="font-semibold text-slate-800 dark:text-white">系统设置</p>
              <p class="text-sm text-slate-500">管理您的账户信息</p>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useVpsStore } from '@/stores/vps'
import { useAuthStore } from '@/stores/auth'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'

const vpsStore = useVpsStore()
const authStore = useAuthStore()

onMounted(() => {
  vpsStore.fetchList()
})

const runningCount = computed(() => vpsStore.vpsList.filter(v => v.status === 'running').length)
const stoppedCount = computed(() => vpsStore.vpsList.filter(v => v.status === 'stopped').length)
const completedChallenges = computed(() => 0) // TODO: implement challenge tracking
</script>
