<template>
  <div v-if="vpsStore.currentVps" class="space-y-6">
    <!-- Breadcrumb and Actions -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center gap-2">
        <Button 
          icon="pi pi-arrow-left" 
          text 
          rounded
          @click="$router.push('/vps')"
        />
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.name }}</h2>
            <Tag 
              :severity="vpsStore.currentVps.status === 'running' ? 'success' : 'danger'"
              :value="vpsStore.currentVps.status === 'running' ? '运行中' : '已停止'"
            />
          </div>
          <code class="text-sm text-slate-500">{{ vpsStore.currentVps.ip_address }}</code>
        </div>
      </div>

      <div class="flex gap-2">
        <Button 
          v-if="vpsStore.currentVps.status === 'stopped'"
          label="启动" 
          icon="pi pi-play" 
          severity="success"
          @click="vpsStore.start(vpsStore.currentVps!.id)"
        />
        <Button 
          v-if="vpsStore.currentVps.status === 'running'"
          label="停止" 
          icon="pi pi-stop" 
          severity="warning"
          @click="vpsStore.stop(vpsStore.currentVps!.id)"
        />
        <Button 
          label="重启" 
          icon="pi pi-refresh" 
          @click="vpsStore.restart(vpsStore.currentVps!.id)"
        />
        <Button 
          label="重装系统" 
          icon="pi pi-download" 
          severity="info"
        />
        <Button 
          label="删除" 
          icon="pi pi-trash" 
          severity="danger"
          outlined
          @click="deleteVps"
        />
      </div>
    </div>

    <!-- Info Cards -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Configuration -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-cog text-primary-500"></i>
            <span class="font-bold">配置信息</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-microchip text-primary-500"></i>
                  <span class="text-sm text-slate-500">CPU</span>
                </div>
                <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.cpu }} <span class="text-sm font-normal text-slate-500">核</span></p>
              </div>

              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-database text-primary-500"></i>
                  <span class="text-sm text-slate-500">内存</span>
                </div>
                <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.memory }} <span class="text-sm font-normal text-slate-500">MB</span></p>
              </div>

              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-hdd text-primary-500"></i>
                  <span class="text-sm text-slate-500">磁盘</span>
                </div>
                <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.disk }} <span class="text-sm font-normal text-slate-500">GB</span></p>
              </div>

              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-wifi text-primary-500"></i>
                  <span class="text-sm text-slate-500">带宽</span>
                </div>
                <p class="text-2xl font-bold text-slate-800 dark:text-white">{{ vpsStore.currentVps.bandwidth }} <span class="text-sm font-normal text-slate-500">Mbps</span></p>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Network and System -->
      <Card class="shadow-sm">
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-globe text-primary-500"></i>
            <span class="font-bold">网络与系统</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-desktop text-blue-600"></i>
                </div>
                <div>
                  <p class="text-sm text-slate-500">操作系统</p>
                  <p class="font-semibold text-slate-800 dark:text-white">{{ vpsStore.currentVps.os_image }}</p>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-map-marker text-green-600"></i>
                </div>
                <div>
                  <p class="text-sm text-slate-500">IP 地址</p>
                  <code class="font-semibold text-slate-800 dark:text-white">{{ vpsStore.currentVps.ip_address }}</code>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-calendar text-purple-600"></i>
                </div>
                <div>
                  <p class="text-sm text-slate-500">到期时间</p>
                  <p class="font-semibold text-slate-800 dark:text-white">{{ formatDate(vpsStore.currentVps.expire_at) }}</p>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-700/50 rounded-xl">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-orange-100 dark:bg-orange-900/30 rounded-lg flex items-center justify-center">
                  <i class="pi pi-clock text-orange-600"></i>
                </div>
                <div>
                  <p class="text-sm text-slate-500">创建时间</p>
                  <p class="font-semibold text-slate-800 dark:text-white">{{ formatDate(vpsStore.currentVps.created_at) }}</p>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Console Access -->
    <Card class="shadow-sm">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-terminal text-primary-500"></i>
          <span class="font-bold">控制台访问</span>
        </div>
      </template>
      <template #content>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-slate-100 dark:bg-slate-700 rounded-xl flex items-center justify-center">
              <i class="pi pi-desktop text-2xl text-slate-400"></i>
            </div>
            <div>
              <p class="font-semibold text-slate-800 dark:text-white">Web 控制台</p>
              <p class="text-sm text-slate-500">通过浏览器访问 VPS 终端</p>
            </div>
          </div>
          <Button 
            label="打开控制台" 
            icon="pi pi-external-link"
            severity="primary"
          />
        </div>
      </template>
    </Card>
  </div>

  <!-- Loading State -->
  <div v-else class="flex items-center justify-center py-20">
    <ProgressSpinner />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useVpsStore } from '@/stores/vps'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import ProgressSpinner from 'primevue/progressspinner'
import { formatDate } from '@/utils/date'

const route = useRoute()
const router = useRouter()
const vpsStore = useVpsStore()
const confirm = useConfirm()
const toast = useToast()

onMounted(() => {
  vpsStore.fetchDetail(Number(route.params.id))
})

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
