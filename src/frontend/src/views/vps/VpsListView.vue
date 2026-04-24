<template>
  <div>
    <PageHeader title="VPS 管理" description="管理您的虚拟服务器实例">
      <template #actions>
        <Button label="购买 VPS" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
      </template>
    </PageHeader>

    <!-- VPS Cards -->
    <div v-if="vpsStore.vpsList.length > 0" class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-5">
      <div
        v-for="vps in vpsStore.vpsList"
        :key="vps.id"
        class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden transition-all duration-200 hover:border-[var(--border-strong)] hover:shadow-sm"
      >
        <div class="p-5">
          <!-- Header: Name + Status -->
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-3">
              <div
                class="w-10 h-10 rounded-xl flex items-center justify-center"
                :style="{ backgroundColor: vps.status === 'running' ? 'var(--success-subtle)' : 'var(--danger-subtle)' }"
              >
                <i
                  class="pi pi-server text-base"
                  :style="{ color: vps.status === 'running' ? 'var(--success)' : 'var(--danger)' }"
                ></i>
              </div>
              <div>
                <h3 class="font-semibold text-[15px] text-[var(--text-primary)] leading-tight">{{ vps.name }}</h3>
                <span
                  class="inline-flex items-center gap-1 mt-0.5 text-[11px] font-medium"
                  :style="{ color: vps.status === 'running' ? 'var(--success)' : 'var(--danger)' }"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :style="{ backgroundColor: vps.status === 'running' ? 'var(--success)' : 'var(--danger)' }"></span>
                  {{ vps.status === 'running' ? '运行中' : '已停止' }}
                </span>
              </div>
            </div>
          </div>

          <!-- IP -->
          <div class="flex items-center justify-between py-2.5 border-t border-[var(--border-subtle)]">
            <span class="text-[12px] text-[var(--text-tertiary)]">IP 地址</span>
            <code class="text-[12px] text-[var(--text-primary)] mono font-medium">{{ vps.ip_address }}</code>
          </div>

          <!-- Specs Row -->
          <div class="grid grid-cols-3 gap-3 py-3">
            <div class="text-center">
              <p class="text-lg font-bold text-[var(--text-primary)] mono leading-none">{{ vps.cpu }}</p>
              <p class="text-[11px] text-[var(--text-tertiary)] mt-1">CPU 核心</p>
            </div>
            <div class="text-center border-x border-[var(--border-subtle)]">
              <p class="text-lg font-bold text-[var(--text-primary)] mono leading-none">{{ vps.memory }}</p>
              <p class="text-[11px] text-[var(--text-tertiary)] mt-1">MB 内存</p>
            </div>
            <div class="text-center">
              <p class="text-lg font-bold text-[var(--text-primary)] mono leading-none">{{ vps.disk }}</p>
              <p class="text-[11px] text-[var(--text-tertiary)] mt-1">GB 磁盘</p>
            </div>
          </div>

          <!-- OS -->
          <div class="flex items-center justify-between py-2.5 border-t border-[var(--border-subtle)]">
            <span class="text-[12px] text-[var(--text-tertiary)]">操作系统</span>
            <span class="text-[12px] text-[var(--text-secondary)] font-medium">{{ vps.os_image }}</span>
          </div>

          <!-- Actions -->
          <div class="flex gap-2 mt-4 pt-3 border-t border-[var(--border-default)]">
            <Button label="详情" icon="pi pi-eye" text size="small" class="flex-1" @click="$router.push(`/vps/${vps.id}`)" />
            <Button
              v-if="vps.status === 'stopped'"
              label="启动" icon="pi pi-play" severity="success" size="small" class="flex-1"
              @click="handleStart(vps.id)"
            />
            <Button
              v-if="vps.status === 'running'"
              label="停止" icon="pi pi-stop" severity="warn" size="small" class="flex-1"
              @click="handleStop(vps.id)"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!vpsStore.loading" class="empty-state">
      <div class="empty-state-icon">
        <i class="pi pi-server text-3xl text-[var(--text-tertiary)]"></i>
      </div>
      <h3 class="text-base font-semibold text-[var(--text-secondary)] mb-1">暂无 VPS 实例</h3>
      <p class="text-[var(--text-tertiary)] text-sm mb-4">点击上方按钮购买您的第一个 VPS</p>
      <Button label="购买 VPS" icon="pi pi-plus" size="small" @click="showCreateDialog = true" />
    </div>

    <!-- Create VPS Dialog -->
    <Dialog v-model:visible="showCreateDialog" header="购买 VPS" modal :style="{ width: '450px' }">
      <div class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">实例名称</label>
          <InputText v-model="createForm.name" class="w-full" placeholder="例如: my-web-server" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">CPU (核)</label>
            <Select v-model="createForm.cpu" :options="cpuOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">内存 (MB)</label>
            <Select v-model="createForm.memory" :options="memoryOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">磁盘 (GB)</label>
            <Select v-model="createForm.disk" :options="diskOptions" placeholder="选择" class="w-full" />
          </div>
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">带宽 (Mbps)</label>
            <Select v-model="createForm.bandwidth" :options="bandwidthOptions" placeholder="选择" class="w-full" />
          </div>
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">操作系统</label>
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
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import PageHeader from '@/components/PageHeader.vue'

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
