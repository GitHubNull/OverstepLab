<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="page-header">
        <h2>VPS 管理</h2>
        <p>管理您的虚拟服务器实例</p>
      </div>
      <Button
        label="购买 VPS"
        icon="pi pi-plus"
        size="small"
        @click="showCreateDialog = true"
      />
    </div>

    <!-- VPS Cards -->
    <div v-if="vpsStore.vpsList.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <Card
        v-for="vps in vpsStore.vpsList"
        :key="vps.id"
        class="shadow-none hover:shadow-md group"
      >
        <template #content>
          <div class="space-y-3.5">
            <!-- Header -->
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <div
                  class="w-9 h-9 rounded-lg flex items-center justify-center"
                  :class="vps.status === 'running' ? 'bg-emerald-50 dark:bg-emerald-900/20' : 'bg-red-50 dark:bg-red-900/20'"
                >
                  <i
                    class="pi pi-server text-sm"
                    :class="vps.status === 'running' ? 'text-emerald-500' : 'text-red-400'"
                  ></i>
                </div>
                <div>
                  <h3 class="font-semibold text-sm text-slate-800 dark:text-white">{{ vps.name }}</h3>
                  <Tag
                    :severity="vps.status === 'running' ? 'success' : 'danger'"
                    :value="vps.status === 'running' ? '运行中' : '已停止'"
                    class="text-[10px]"
                  />
                </div>
              </div>
            </div>

            <!-- IP Address -->
            <div class="flex items-center justify-between text-xs">
              <span class="text-slate-400">IP 地址</span>
              <code class="bg-slate-50 dark:bg-slate-700/50 px-1.5 py-0.5 rounded text-slate-600 dark:text-slate-300 font-mono">{{ vps.ip_address }}</code>
            </div>

            <!-- Specs -->
            <div class="grid grid-cols-3 gap-2">
              <div class="spec-box !py-2">
                <i class="pi pi-microchip text-slate-300 text-xs"></i>
                <p class="text-base font-bold text-slate-700 dark:text-white mt-0.5">{{ vps.cpu }}</p>
                <p class="text-[10px] text-slate-400">核 CPU</p>
              </div>
              <div class="spec-box !py-2">
                <i class="pi pi-database text-slate-300 text-xs"></i>
                <p class="text-base font-bold text-slate-700 dark:text-white mt-0.5">{{ vps.memory }}</p>
                <p class="text-[10px] text-slate-400">MB 内存</p>
              </div>
              <div class="spec-box !py-2">
                <i class="pi pi-hdd text-slate-300 text-xs"></i>
                <p class="text-base font-bold text-slate-700 dark:text-white mt-0.5">{{ vps.disk }}</p>
                <p class="text-[10px] text-slate-400">GB 磁盘</p>
              </div>
            </div>

            <!-- OS -->
            <div class="flex items-center justify-between text-xs">
              <span class="text-slate-400">操作系统</span>
              <span class="text-slate-600 dark:text-slate-300 flex items-center gap-1">
                <i class="pi pi-desktop text-slate-300 text-[10px]"></i>
                {{ vps.os_image }}
              </span>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-1 border-t border-slate-100 dark:border-slate-700/50">
              <Button
                label="详情"
                icon="pi pi-eye"
                text
                size="small"
                class="flex-1 !text-xs"
                @click="$router.push(`/vps/${vps.id}`)"
              />
              <Button
                v-if="vps.status === 'stopped'"
                label="启动"
                icon="pi pi-play"
                severity="success"
                size="small"
                class="flex-1 !text-xs"
                @click="handleStart(vps.id)"
              />
              <Button
                v-if="vps.status === 'running'"
                label="停止"
                icon="pi pi-stop"
                severity="warn"
                size="small"
                class="flex-1 !text-xs"
                @click="handleStop(vps.id)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-else-if="!vpsStore.loading" class="empty-state">
      <div class="empty-state-icon">
        <i class="pi pi-server text-3xl text-slate-300"></i>
      </div>
      <h3 class="text-base font-semibold text-slate-600 dark:text-slate-300 mb-1">暂无 VPS 实例</h3>
      <p class="text-slate-400 text-sm mb-4">点击上方按钮购买您的第一个 VPS</p>
      <Button label="购买 VPS" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
    </div>

    <!-- Create VPS Dialog -->
    <Dialog
      v-model:visible="showCreateDialog"
      header="购买 VPS"
      modal
      :style="{ width: '450px' }"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">实例名称</label>
          <InputText v-model="createForm.name" class="w-full" placeholder="例如: my-web-server" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">CPU (核)</label>
            <Select v-model="createForm.cpu" :options="cpuOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">内存 (MB)</label>
            <Select v-model="createForm.memory" :options="memoryOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">磁盘 (GB)</label>
            <Select v-model="createForm.disk" :options="diskOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">带宽 (Mbps)</label>
            <Select v-model="createForm.bandwidth" :options="bandwidthOptions" placeholder="选择" class="w-full" />
          </div>
        </div>
        <div>
          <label class="block text-xs font-semibold text-slate-500 mb-1.5 uppercase tracking-wider">操作系统</label>
          <Select v-model="createForm.os_image" :options="osOptions" placeholder="选择操作系统" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showCreateDialog = false" />
        <Button label="确认购买" icon="pi pi-check" size="small" :loading="creating" @click="handleCreate" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useVpsStore } from '@/stores/vps'
import { useToast } from 'primevue/usetoast'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'

const vpsStore = useVpsStore()
const toast = useToast()
const showCreateDialog = ref(false)
const creating = ref(false)

const cpuOptions = [1, 2, 4, 8]
const memoryOptions = [512, 1024, 2048, 4096, 8192]
const diskOptions = [20, 40, 80, 160, 320]
const bandwidthOptions = [1, 5, 10, 20, 50, 100]
const osOptions = ['Ubuntu 22.04', 'Ubuntu 20.04', 'CentOS 8', 'Debian 11', 'Debian 12']

const createForm = ref({
  name: '',
  cpu: 2,
  memory: 1024,
  disk: 40,
  bandwidth: 10,
  os_image: 'Ubuntu 22.04',
})

onMounted(() => {
  vpsStore.fetchList()
})

async function handleStart(id: number) {
  try {
    await vpsStore.start(id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已启动', life: 2000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  }
}

async function handleStop(id: number) {
  try {
    await vpsStore.stop(id)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已停止', life: 2000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  }
}

async function handleCreate() {
  if (!createForm.value.name) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请输入实例名称', life: 2000 })
    return
  }
  creating.value = true
  try {
    await vpsStore.create(createForm.value)
    toast.add({ severity: 'success', summary: '成功', detail: 'VPS 已创建', life: 2000 })
    showCreateDialog.value = false
    createForm.value = { name: '', cpu: 2, memory: 1024, disk: 40, bandwidth: 10, os_image: 'Ubuntu 22.04' }
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '创建失败', life: 3000 })
  } finally {
    creating.value = false
  }
}
</script>
