<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">VPS 列表</h2>
        <p class="text-slate-500">管理您的虚拟服务器实例</p>
      </div>
      <Button 
        label="购买 VPS" 
        icon="pi pi-plus" 
        severity="primary"
        @click="showCreateDialog = true"
      />
    </div>

    <!-- VPS Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <Card 
        v-for="vps in vpsStore.vpsList" 
        :key="vps.id"
        class="shadow-sm hover:shadow-lg transition-all duration-300 group"
      >
        <template #title>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div 
                class="w-10 h-10 rounded-xl flex items-center justify-center"
                :class="vps.status === 'running' ? 'bg-green-100 dark:bg-green-900/30' : 'bg-red-100 dark:bg-red-900/30'"
              >
                <i 
                  class="pi text-lg"
                  :class="vps.status === 'running' ? 'pi-server text-green-600' : 'pi-server text-red-600'"
                ></i>
              </div>
              <div>
                <h3 class="font-bold text-slate-800 dark:text-white">{{ vps.name }}</h3>
                <Tag 
                  :severity="vps.status === 'running' ? 'success' : 'danger'"
                  :value="vps.status === 'running' ? '运行中' : '已停止'"
                  class="text-xs mt-1"
                />
              </div>
            </div>
            <Button 
              icon="pi pi-ellipsis-v" 
              text 
              rounded
              @click="(e) => toggleMenu(e, vps)"
            />
          </div>
        </template>

        <template #content>
          <div class="space-y-3">
            <!-- IP Address -->
            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-500">IP 地址</span>
              <code class="bg-slate-100 dark:bg-slate-700 px-2 py-1 rounded text-sm font-mono">{{ vps.ip_address }}</code>
            </div>

            <!-- Specs -->
            <div class="grid grid-cols-3 gap-2">
              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-lg p-2 text-center">
                <i class="pi pi-microchip text-slate-400 mb-1"></i>
                <p class="text-lg font-bold text-slate-800 dark:text-white">{{ vps.cpu }}</p>
                <p class="text-xs text-slate-500">核 CPU</p>
              </div>
              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-lg p-2 text-center">
                <i class="pi pi-database text-slate-400 mb-1"></i>
                <p class="text-lg font-bold text-slate-800 dark:text-white">{{ vps.memory }}</p>
                <p class="text-xs text-slate-500">MB 内存</p>
              </div>
              <div class="bg-slate-50 dark:bg-slate-700/50 rounded-lg p-2 text-center">
                <i class="pi pi-hdd text-slate-400 mb-1"></i>
                <p class="text-lg font-bold text-slate-800 dark:text-white">{{ vps.disk }}</p>
                <p class="text-xs text-slate-500">GB 磁盘</p>
              </div>
            </div>

            <!-- OS -->
            <div class="flex items-center justify-between">
              <span class="text-sm text-slate-500">操作系统</span>
              <div class="flex items-center gap-2">
                <i class="pi pi-desktop text-slate-400"></i>
                <span class="text-sm">{{ vps.os_image }}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-2">
              <Button 
                label="查看详情" 
                icon="pi pi-eye" 
                text 
                size="small"
                class="flex-1"
                @click="$router.push(`/vps/${vps.id}`)"
              />
              <Button 
                v-if="vps.status === 'stopped'"
                label="启动" 
                icon="pi pi-play" 
                severity="success"
                size="small"
                class="flex-1"
                @click="vpsStore.start(vps.id)"
              />
              <Button 
                v-if="vps.status === 'running'"
                label="停止" 
                icon="pi pi-stop" 
                severity="warning"
                size="small"
                class="flex-1"
                @click="vpsStore.stop(vps.id)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-if="vpsStore.vpsList.length === 0 && !vpsStore.loading" class="text-center py-16">
      <div class="w-20 h-20 bg-slate-100 dark:bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-4">
        <i class="pi pi-server text-4xl text-slate-400"></i>
      </div>
      <h3 class="text-lg font-semibold text-slate-700 dark:text-slate-300 mb-2">暂无 VPS 实例</h3>
      <p class="text-slate-500 mb-4">点击上方按钮购买您的第一个 VPS</p>
      <Button label="购买 VPS" icon="pi pi-plus" @click="showCreateDialog = true" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useVpsStore } from '@/stores/vps'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'

const router = useRouter()
const vpsStore = useVpsStore()
const showCreateDialog = ref(false)

onMounted(() => {
  vpsStore.fetchList()
})

function toggleMenu(event: Event, vps: any) {
  // TODO: Implement context menu
  console.log('Toggle menu for', vps.name)
}
</script>
