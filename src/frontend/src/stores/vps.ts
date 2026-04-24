import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getVpsList, getVpsDetail, startVps, stopVps, restartVps, reinstallVps, deleteVps, createVps } from '@/api/vps'
import type { VPSInstance } from '@/types'

export const useVpsStore = defineStore('vps', () => {
  const vpsList = ref<VPSInstance[]>([])
  const currentVps = ref<VPSInstance | null>(null)
  const loading = ref(false)

  async function fetchList() {
    loading.value = true
    try {
      const response = await getVpsList()
      vpsList.value = response.data.data!
    } finally {
      loading.value = false
    }
  }

  async function fetchDetail(id: number) {
    loading.value = true
    try {
      const response = await getVpsDetail(id)
      currentVps.value = response.data.data!
    } finally {
      loading.value = false
    }
  }

  async function start(id: number) {
    await startVps(id)
    await fetchList()
  }

  async function stop(id: number) {
    await stopVps(id)
    await fetchList()
  }

  async function restart(id: number) {
    await restartVps(id)
    await fetchList()
  }

  async function reinstall(id: number, osImage: string) {
    await reinstallVps(id, osImage)
    await fetchDetail(id)
  }

  async function remove(id: number) {
    await deleteVps(id)
    await fetchList()
  }

  async function create(data: { name: string; cpu: number; memory: number; disk: number; bandwidth: number; os_image: string }) {
    await createVps(data)
    await fetchList()
  }

  return {
    vpsList, currentVps, loading,
    fetchList, fetchDetail, start, stop, restart, reinstall, remove, create,
  }
})
