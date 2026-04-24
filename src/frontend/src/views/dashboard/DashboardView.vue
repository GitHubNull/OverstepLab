<template>
  <div class="space-y-5">
    <!-- Welcome Banner -->
    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl p-5 relative overflow-hidden">
      <div class="absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full bg-[var(--primary)]"></div>
      <div class="flex items-center justify-between pl-3">
        <div>
          <h1 class="text-lg font-bold text-[var(--text-primary)]">欢迎回来，{{ authStore.user?.username }}</h1>
          <p class="text-sm text-[var(--text-secondary)] mt-1">
            {{ authStore.securityMode === 'vulnerable'
              ? '当前处于漏洞模式，所有越权漏洞均可触发'
              : '当前处于安全模式，所有漏洞已被修复' }}
          </p>
        </div>
        <div class="hidden md:flex">
          <div
            class="w-12 h-12 rounded-xl flex items-center justify-center"
            :class="authStore.securityMode === 'vulnerable' ? 'bg-[var(--danger-subtle)]' : 'bg-[var(--success-subtle)]'"
          >
            <i
              class="pi text-xl"
              :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-[var(--danger)]' : 'pi-lock text-[var(--success)]'"
            ></i>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard color="#3b82f6" icon="pi pi-server" :value="vpsStore.vpsList.length" label="VPS 总数" />
      <StatCard color="#10b981" icon="pi pi-check-circle" :value="runningCount" label="运行中" />
      <StatCard color="#f43f5e" icon="pi pi-times-circle" :value="stoppedCount" label="已停止" />
      <StatCard color="#8b5cf6" icon="pi pi-flag" :value="`${completedChallenges}/13`" label="挑战完成" />
    </div>

    <!-- VPS List & Quick Actions -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- VPS List -->
      <div class="lg:col-span-2 bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-5 py-4 border-b border-[var(--border-default)]">
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
        <div class="p-0">
          <DataTable
            :value="vpsStore.vpsList"
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
                    :class="data.status === 'running' ? 'bg-[var(--success)]' : 'bg-[var(--danger)]'"
                    :style="data.status === 'running' ? 'box-shadow: 0 0 0 2px var(--success-subtle)' : ''"
                  ></div>
                  <span class="font-medium text-[var(--text-primary)]">{{ data.name }}</span>
                </div>
              </template>
            </Column>
            <Column field="ip_address" header="IP 地址">
              <template #body="{ data }">
                <code class="text-xs bg-[var(--bg-base)] px-1.5 py-0.5 rounded text-[var(--text-secondary)] mono">{{ data.ip_address }}</code>
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
                  class="text-[var(--text-tertiary)]"
                  @click="$router.push(`/vps/${data.id}`)"
                />
              </template>
            </Column>
          </DataTable>

          <!-- Empty State -->
          <div v-if="vpsStore.vpsList.length === 0 && !vpsStore.loading" class="empty-state">
            <div class="empty-state-icon">
              <i class="pi pi-server text-2xl text-[var(--text-tertiary)]"></i>
            </div>
            <p class="text-[var(--text-secondary)] text-sm">暂无 VPS 实例</p>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="space-y-2">
        <div
          v-for="item in quickActions"
          :key="item.label"
          class="flex items-center gap-3 px-4 py-3 rounded-xl cursor-pointer transition-all duration-200 border border-transparent hover:border-[var(--border-default)] hover:bg-[var(--bg-surface-hover)]"
          @click="$router.push(item.to)"
        >
          <div
            class="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0"
            :style="{ backgroundColor: `${item.color}14` }"
          >
            <i :class="[item.icon, 'text-sm']" :style="{ color: item.color }"></i>
          </div>
          <div>
            <p class="font-semibold text-sm text-[var(--text-primary)]">{{ item.label }}</p>
            <p class="text-[11px] text-[var(--text-secondary)]">{{ item.desc }}</p>
          </div>
          <i class="pi pi-chevron-right text-[10px] text-[var(--text-tertiary)] ml-auto"></i>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useVpsStore } from '@/stores/vps'
import { useAuthStore } from '@/stores/auth'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import StatCard from '@/components/StatCard.vue'

const vpsStore = useVpsStore()
const authStore = useAuthStore()

onMounted(() => {
  vpsStore.fetchList()
})

const runningCount = computed(() => vpsStore.vpsList.filter(v => v.status === 'running').length)
const stoppedCount = computed(() => vpsStore.vpsList.filter(v => v.status === 'stopped').length)
const completedChallenges = computed(() => 0) // TODO: implement challenge tracking

const quickActions = [
  { label: '漏洞挑战', desc: '开始您的越权测试练习', icon: 'pi pi-flag', color: '#f97316', to: '/challenges' },
  { label: '个人资料', desc: '管理您的账户信息', icon: 'pi pi-user', color: '#06b6d4', to: '/profile' },
  { label: 'API Key', desc: '管理 API 访问凭证', icon: 'pi pi-key', color: '#eab308', to: '/apikeys' },
  { label: '工单系统', desc: '提交技术支持请求', icon: 'pi pi-comments', color: '#ec4899', to: '/tickets' },
]
</script>
