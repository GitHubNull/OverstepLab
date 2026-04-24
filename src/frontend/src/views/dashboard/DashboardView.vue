<template>
  <div class="space-y-5">
    <!-- Welcome Banner -->
    <div class="bg-gradient-to-r from-indigo-600 to-indigo-700 rounded-2xl p-5 text-white shadow-lg shadow-indigo-200/30">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-bold mb-1">欢迎回来，{{ authStore.user?.username }}</h1>
          <p class="text-indigo-200 text-sm">
            {{ authStore.securityMode === 'vulnerable'
              ? '当前处于漏洞模式，所有越权漏洞均可触发'
              : '当前处于安全模式，所有漏洞已被修复' }}
          </p>
        </div>
        <div class="hidden md:flex">
          <div
            class="w-14 h-14 rounded-2xl flex items-center justify-center"
            :class="authStore.securityMode === 'vulnerable' ? 'bg-red-500/20' : 'bg-emerald-500/20'"
          >
            <i
              class="pi text-2xl"
              :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-300' : 'pi-lock text-emerald-300'"
            ></i>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <Card class="shadow-none stat-card stat-card-blue">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-blue-50 dark:bg-blue-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-server text-blue-500"></i>
            </div>
            <div class="min-w-0">
              <p class="text-[11px] text-slate-400 font-medium">VPS 总数</p>
              <p class="text-xl font-bold text-slate-800 dark:text-white">{{ vpsStore.vpsList.length }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-none stat-card stat-card-green">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-check-circle text-emerald-500"></i>
            </div>
            <div class="min-w-0">
              <p class="text-[11px] text-slate-400 font-medium">运行中</p>
              <p class="text-xl font-bold text-emerald-600">{{ runningCount }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-none stat-card stat-card-red">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-red-50 dark:bg-red-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-times-circle text-red-500"></i>
            </div>
            <div class="min-w-0">
              <p class="text-[11px] text-slate-400 font-medium">已停止</p>
              <p class="text-xl font-bold text-red-500">{{ stoppedCount }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="shadow-none stat-card stat-card-purple cursor-pointer" @click="$router.push('/challenges')">
        <template #content>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-purple-50 dark:bg-purple-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
              <i class="pi pi-flag text-purple-500"></i>
            </div>
            <div class="min-w-0">
              <p class="text-[11px] text-slate-400 font-medium">挑战完成</p>
              <p class="text-xl font-bold text-purple-600">{{ completedChallenges }}/13</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- VPS List & Quick Actions -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- VPS List -->
      <Card class="shadow-none lg:col-span-2">
        <template #title>
          <div class="flex items-center justify-between">
            <div class="section-title">
              <i class="pi pi-server"></i>
              <span>我的 VPS</span>
            </div>
            <Button
              label="查看全部"
              icon="pi pi-arrow-right"
              iconPos="right"
              text
              size="small"
              class="text-xs"
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
                    class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                    :class="data.status === 'running' ? 'bg-emerald-500' : 'bg-red-400'"
                  ></div>
                  <span class="font-medium text-slate-700 dark:text-slate-300">{{ data.name }}</span>
                </div>
              </template>
            </Column>
            <Column field="ip_address" header="IP 地址">
              <template #body="{ data }">
                <code class="text-xs bg-slate-100 dark:bg-slate-700 px-1.5 py-0.5 rounded text-slate-600 dark:text-slate-300">{{ data.ip_address }}</code>
              </template>
            </Column>
            <Column field="status" header="状态">
              <template #body="{ data }">
                <Tag
                  :severity="data.status === 'running' ? 'success' : 'danger'"
                  :value="data.status === 'running' ? '运行中' : '已停止'"
                  class="text-[10px]"
                />
              </template>
            </Column>
            <Column header="操作" style="width: 60px">
              <template #body="{ data }">
                <Button
                  icon="pi pi-eye"
                  text
                  rounded
                  size="small"
                  class="text-slate-400"
                  @click="$router.push(`/vps/${data.id}`)"
                />
              </template>
            </Column>
          </DataTable>

          <!-- Empty State -->
          <div v-if="vpsStore.vpsList.length === 0 && !vpsStore.loading" class="empty-state">
            <div class="empty-state-icon">
              <i class="pi pi-server text-2xl text-slate-400"></i>
            </div>
            <p class="text-slate-500 text-sm">暂无 VPS 实例</p>
          </div>
        </template>
      </Card>

      <!-- Quick Actions -->
      <div class="space-y-3">
        <Card class="shadow-none cursor-pointer hover:border-indigo-200 dark:hover:border-indigo-800" @click="$router.push('/challenges')">
          <template #content>
            <div class="flex items-center gap-3 p-1">
              <div class="w-10 h-10 bg-orange-50 dark:bg-orange-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
                <i class="pi pi-flag text-orange-500"></i>
              </div>
              <div>
                <p class="font-semibold text-sm text-slate-700 dark:text-slate-200">漏洞挑战</p>
                <p class="text-[11px] text-slate-400">开始您的越权测试练习</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="shadow-none cursor-pointer hover:border-indigo-200 dark:hover:border-indigo-800" @click="$router.push('/profile')">
          <template #content>
            <div class="flex items-center gap-3 p-1">
              <div class="w-10 h-10 bg-cyan-50 dark:bg-cyan-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
                <i class="pi pi-user text-cyan-500"></i>
              </div>
              <div>
                <p class="font-semibold text-sm text-slate-700 dark:text-slate-200">个人资料</p>
                <p class="text-[11px] text-slate-400">管理您的账户信息</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="shadow-none cursor-pointer hover:border-indigo-200 dark:hover:border-indigo-800" @click="$router.push('/apikeys')">
          <template #content>
            <div class="flex items-center gap-3 p-1">
              <div class="w-10 h-10 bg-yellow-50 dark:bg-yellow-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
                <i class="pi pi-key text-yellow-500"></i>
              </div>
              <div>
                <p class="font-semibold text-sm text-slate-700 dark:text-slate-200">API Key</p>
                <p class="text-[11px] text-slate-400">管理 API 访问凭证</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="shadow-none cursor-pointer hover:border-indigo-200 dark:hover:border-indigo-800" @click="$router.push('/tickets')">
          <template #content>
            <div class="flex items-center gap-3 p-1">
              <div class="w-10 h-10 bg-pink-50 dark:bg-pink-900/20 rounded-xl flex items-center justify-center flex-shrink-0">
                <i class="pi pi-comments text-pink-500"></i>
              </div>
              <div>
                <p class="font-semibold text-sm text-slate-700 dark:text-slate-200">工单系统</p>
                <p class="text-[11px] text-slate-400">提交技术支持请求</p>
              </div>
            </div>
          </template>
        </Card>
      </div>
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
