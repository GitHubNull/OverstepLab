<template>
  <div v-if="vpsStore.currentVps" class="space-y-5">
    <!-- Header with actions -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-3">
        <Button icon="pi pi-arrow-left" text rounded size="small" class="text-[var(--text-tertiary)]" @click="$router.push('/vps')" />
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-lg font-bold text-[var(--text-primary)]">{{ vpsStore.currentVps.name }}</h1>
            <Tag
              :severity="vpsStore.currentVps.status === 'running' ? 'success' : 'danger'"
              :value="vpsStore.currentVps.status === 'running' ? '运行中' : '已停止'"
              class="text-[10px]"
            />
          </div>
          <code class="text-xs text-[var(--text-tertiary)] mono">{{ vpsStore.currentVps.ip_address }}</code>
        </div>
      </div>

      <div class="flex gap-2 flex-wrap">
        <Button v-if="vpsStore.currentVps.status === 'stopped'" label="启动" icon="pi pi-play" severity="success" size="small" :loading="starting" @click="handleStart" />
        <Button v-if="vpsStore.currentVps.status === 'running'" label="停止" icon="pi pi-stop" severity="warn" size="small" :loading="stopping" @click="handleStop" />
        <Button label="重启" icon="pi pi-refresh" size="small" outlined :loading="restarting" @click="handleRestart" />
        <Button label="重装系统" icon="pi pi-download" severity="info" size="small" outlined @click="showReinstallDialog = true" />
        <Button label="删除" icon="pi pi-trash" severity="danger" size="small" outlined @click="deleteVps" />
      </div>
    </div>

    <!-- Info Cards -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Configuration -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-5 py-4 border-b border-[var(--border-default)]">
          <div class="section-title">
            <i class="pi pi-cog"></i>
            <span>配置信息</span>
          </div>
        </div>
        <div class="p-4">
          <div class="grid grid-cols-2 gap-3">
            <div class="spec-box">
              <p class="text-[10px] text-[var(--text-tertiary)] font-medium uppercase tracking-wider mb-1">CPU</p>
              <p class="text-2xl font-bold text-[var(--text-primary)] mono">{{ vpsStore.currentVps.cpu }} <span class="text-xs font-normal text-[var(--text-tertiary)]">核</span></p>
            </div>
            <div class="spec-box">
              <p class="text-[10px] text-[var(--text-tertiary)] font-medium uppercase tracking-wider mb-1">内存</p>
              <p class="text-2xl font-bold text-[var(--text-primary)] mono">{{ vpsStore.currentVps.memory }} <span class="text-xs font-normal text-[var(--text-tertiary)]">MB</span></p>
            </div>
            <div class="spec-box">
              <p class="text-[10px] text-[var(--text-tertiary)] font-medium uppercase tracking-wider mb-1">磁盘</p>
              <p class="text-2xl font-bold text-[var(--text-primary)] mono">{{ vpsStore.currentVps.disk }} <span class="text-xs font-normal text-[var(--text-tertiary)]">GB</span></p>
            </div>
            <div class="spec-box">
              <p class="text-[10px] text-[var(--text-tertiary)] font-medium uppercase tracking-wider mb-1">带宽</p>
              <p class="text-2xl font-bold text-[var(--text-primary)] mono">{{ vpsStore.currentVps.bandwidth }} <span class="text-xs font-normal text-[var(--text-tertiary)]">Mbps</span></p>
            </div>
          </div>
        </div>
      </div>

      <!-- Network and System -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-5 py-4 border-b border-[var(--border-default)]">
          <div class="section-title">
            <i class="pi pi-globe"></i>
            <span>网络与系统</span>
          </div>
        </div>
        <div class="p-4 space-y-2">
          <div class="info-row">
            <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">操作系统</span>
            <span class="text-sm font-semibold text-[var(--text-primary)]">{{ vpsStore.currentVps.os_image }}</span>
          </div>
          <div class="info-row">
            <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">IP 地址</span>
            <code class="text-sm font-semibold text-[var(--text-primary)] mono">{{ vpsStore.currentVps.ip_address }}</code>
          </div>
          <div class="info-row">
            <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">到期时间</span>
            <span class="text-sm font-semibold text-[var(--text-primary)]">{{ formatDate(vpsStore.currentVps.expire_at) }}</span>
          </div>
          <div class="info-row">
            <span class="text-[11px] text-[var(--text-tertiary)] uppercase tracking-wider">创建时间</span>
            <span class="text-sm font-semibold text-[var(--text-primary)]">{{ formatDate(vpsStore.currentVps.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Console Access -->
    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div class="px-5 py-4 border-b border-[var(--border-default)]">
        <div class="section-title">
          <i class="pi pi-terminal"></i>
          <span>控制台访问</span>
        </div>
      </div>
      <div class="p-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 bg-[var(--bg-base)] rounded-xl flex items-center justify-center">
            <i class="pi pi-desktop text-[var(--text-tertiary)]"></i>
          </div>
          <div>
            <p class="font-semibold text-sm text-[var(--text-primary)]">Web 控制台</p>
            <p class="text-[11px] text-[var(--text-secondary)]">通过浏览器访问 VPS 终端</p>
          </div>
        </div>
        <Button label="打开控制台" icon="pi pi-external-link" size="small" outlined @click="handleConsole" />
      </div>
    </div>

    <!-- Reinstall Dialog -->
    <Dialog v-model:visible="showReinstallDialog" header="重装系统" modal :style="{ width: '400px' }">
      <div class="space-y-3">
        <p class="text-sm text-[var(--text-secondary)]">选择要安装的操作系统镜像：</p>
        <Select v-model="reinstallOs" :options="osOptions" placeholder="选择操作系统" class="w-full" />
        <p class="text-xs text-[var(--danger)]">注意：重装系统将清除所有数据，此操作不可恢复！</p>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showReinstallDialog = false" />
        <Button label="确认重装" icon="pi pi-download" severity="danger" size="small" :loading="reinstalling" @click="handleReinstall" />
      </template>
    </Dialog>
  </div>

  <!-- Loading State -->
  <div v-else class="flex items-center justify-center py-20">
    <ProgressSpinner />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVpsStore } from '@/stores/vps'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import ProgressSpinner from 'primevue/progressspinner'
import { formatDate } from '@/utils/date'

const route = useRoute()
const router = useRouter()
const vpsStore = useVpsStore()
const confirm = useConfirm()
const toast = useToast()
const showReinstallDialog = ref(false)
const reinstallOs = ref('Ubuntu 22.04')
const reinstalling = ref(false)
const starting = ref(false)
const stopping = ref(false)
const restarting = ref(false)
const osOptions = ['Ubuntu 22.04', 'Ubuntu 20.04', 'CentOS 8', 'Debian 11', 'Debian 12']

onMounted(() => {
  vpsStore.fetchDetail(Number(route.params.id))
})

async function handleStart() {
  starting.value = true
  try {
    await vpsStore.start(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已启动', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  } finally {
    starting.value = false
  }
}

async function handleStop() {
  stopping.value = true
  try {
    await vpsStore.stop(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已停止', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  } finally {
    stopping.value = false
  }
}

async function handleRestart() {
  restarting.value = true
  try {
    await vpsStore.restart(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已重启', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  } finally {
    restarting.value = false
  }
}

function handleConsole() {
  const id = route.params.id
  window.open(`/api/v1/vps/${id}/console`, '_blank')
}

async function handleReinstall() {
  reinstalling.value = true
  try {
    await vpsStore.reinstall(vpsStore.currentVps!.id, reinstallOs.value)
    toast.add({ severity: 'success', summary: '成功', detail: '系统重装已开始', life: 2000 })
    showReinstallDialog.value = false
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '重装失败', life: 3000 })
  } finally {
    reinstalling.value = false
  }
}

function deleteVps() {
  confirm.require({
    message: '确定要删除此 VPS 吗？此操作不可恢复。',
    header: '确认删除',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      await vpsStore.remove(Number(route.params.id))
      toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已删除' })
      router.push('/vps')
    },
  })
}
</script>
