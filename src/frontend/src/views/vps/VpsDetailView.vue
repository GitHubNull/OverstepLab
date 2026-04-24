<template>
  <div v-if="vpsStore.currentVps" class="space-y-5">
    <!-- Breadcrumb and Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <Button
          icon="pi pi-arrow-left"
          text
          rounded
          size="small"
          class="text-slate-400"
          @click="$router.push('/vps')"
        />
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-lg font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.name }}</h2>
            <Tag
              :severity="vpsStore.currentVps.status === 'running' ? 'success' : 'danger'"
              :value="vpsStore.currentVps.status === 'running' ? '运行中' : '已停止'"
              class="text-[10px]"
            />
          </div>
          <code class="text-xs text-slate-400">{{ vpsStore.currentVps.ip_address }}</code>
        </div>
      </div>

      <div class="flex gap-2 flex-wrap">
        <Button
          v-if="vpsStore.currentVps.status === 'stopped'"
          label="启动"
          icon="pi pi-play"
          severity="success"
          size="small"
          @click="handleStart"
        />
        <Button
          v-if="vpsStore.currentVps.status === 'running'"
          label="停止"
          icon="pi pi-stop"
          severity="warn"
          size="small"
          @click="handleStop"
        />
        <Button
          label="重启"
          icon="pi pi-refresh"
          size="small"
          outlined
          @click="handleRestart"
        />
        <Button
          label="重装系统"
          icon="pi pi-download"
          severity="info"
          size="small"
          outlined
          @click="showReinstallDialog = true"
        />
        <Button
          label="删除"
          icon="pi pi-trash"
          severity="danger"
          size="small"
          outlined
          @click="deleteVps"
        />
      </div>
    </div>

    <!-- Info Cards -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Configuration -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-cog"></i>
            <span>配置信息</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-2 gap-3">
            <div class="spec-box">
              <div class="flex items-center gap-1.5 mb-2">
                <i class="pi pi-microchip text-indigo-400 text-xs"></i>
                <span class="text-[10px] text-slate-400 font-medium uppercase">CPU</span>
              </div>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.cpu }} <span class="text-xs font-normal text-slate-400">核</span></p>
            </div>
            <div class="spec-box">
              <div class="flex items-center gap-1.5 mb-2">
                <i class="pi pi-database text-indigo-400 text-xs"></i>
                <span class="text-[10px] text-slate-400 font-medium uppercase">内存</span>
              </div>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.memory }} <span class="text-xs font-normal text-slate-400">MB</span></p>
            </div>
            <div class="spec-box">
              <div class="flex items-center gap-1.5 mb-2">
                <i class="pi pi-hdd text-indigo-400 text-xs"></i>
                <span class="text-[10px] text-slate-400 font-medium uppercase">磁盘</span>
              </div>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.disk }} <span class="text-xs font-normal text-slate-400">GB</span></p>
            </div>
            <div class="spec-box">
              <div class="flex items-center gap-1.5 mb-2">
                <i class="pi pi-wifi text-indigo-400 text-xs"></i>
                <span class="text-[10px] text-slate-400 font-medium uppercase">带宽</span>
              </div>
              <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.bandwidth }} <span class="text-xs font-normal text-slate-400">Mbps</span></p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Network and System -->
      <Card class="shadow-none">
        <template #title>
          <div class="section-title">
            <i class="pi pi-globe"></i>
            <span>网络与系统</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-blue-50 dark:bg-blue-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-desktop text-blue-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-[10px] text-slate-400 uppercase">操作系统</p>
                  <p class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ vpsStore.currentVps.os_image }}</p>
                </div>
              </div>
            </div>
            <div class="info-row">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-map-marker text-emerald-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-[10px] text-slate-400 uppercase">IP 地址</p>
                  <code class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ vpsStore.currentVps.ip_address }}</code>
                </div>
              </div>
            </div>
            <div class="info-row">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-purple-50 dark:bg-purple-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-calendar text-purple-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-[10px] text-slate-400 uppercase">到期时间</p>
                  <p class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ formatDate(vpsStore.currentVps.expire_at) }}</p>
                </div>
              </div>
            </div>
            <div class="info-row">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 bg-orange-50 dark:bg-orange-900/20 rounded-lg flex items-center justify-center">
                  <i class="pi pi-clock text-orange-500 text-xs"></i>
                </div>
                <div>
                  <p class="text-[10px] text-slate-400 uppercase">创建时间</p>
                  <p class="text-sm font-semibold text-slate-700 dark:text-slate-200">{{ formatDate(vpsStore.currentVps.created_at) }}</p>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Console Access -->
    <Card class="shadow-none">
      <template #title>
        <div class="section-title">
          <i class="pi pi-terminal"></i>
          <span>控制台访问</span>
        </div>
      </template>
      <template #content>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-slate-100 dark:bg-slate-700 rounded-xl flex items-center justify-center">
              <i class="pi pi-desktop text-slate-400"></i>
            </div>
            <div>
              <p class="font-semibold text-sm text-slate-700 dark:text-slate-200">Web 控制台</p>
              <p class="text-[11px] text-slate-400">通过浏览器访问 VPS 终端</p>
            </div>
          </div>
          <Button
            label="打开控制台"
            icon="pi pi-external-link"
            size="small"
            outlined
          />
        </div>
      </template>
    </Card>

    <!-- Reinstall Dialog -->
    <Dialog v-model:visible="showReinstallDialog" header="重装系统" modal :style="{ width: '400px' }">
      <div class="space-y-3">
        <p class="text-sm text-slate-500">选择要安装的操作系统镜像：</p>
        <Select v-model="reinstallOs" :options="osOptions" placeholder="选择操作系统" class="w-full" />
        <p class="text-xs text-red-500">注意：重装系统将清除所有数据，此操作不可恢复！</p>
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
import Card from 'primevue/card'
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
const osOptions = ['Ubuntu 22.04', 'Ubuntu 20.04', 'CentOS 8', 'Debian 11', 'Debian 12']

onMounted(() => {
  vpsStore.fetchDetail(Number(route.params.id))
})

async function handleStart() {
  try {
    await vpsStore.start(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已启动', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  }
}

async function handleStop() {
  try {
    await vpsStore.stop(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已停止', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  }
}

async function handleRestart() {
  try {
    await vpsStore.restart(vpsStore.currentVps!.id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已重启', life: 2000 })
    vpsStore.fetchDetail(Number(route.params.id))
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  }
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
