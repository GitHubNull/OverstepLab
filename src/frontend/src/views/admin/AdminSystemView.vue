<template>
  <div>
    <PageHeader title="系统设置" description="管理系统配置和数据" />

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- System Config -->
      <Card>
        <template #title>
          <div class="section-title">
            <i class="pi pi-sliders-h"></i>
            <span>系统配置</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-3">
            <div v-for="cfg in configFields" :key="cfg.key" class="flex items-center justify-between gap-3">
              <div class="flex-1 min-w-0">
                <label class="text-xs font-semibold text-[var(--text-secondary)]">{{ cfg.label }}</label>
              </div>
              <div class="flex items-center gap-2">
                <InputText
                  v-if="cfg.type === 'text'"
                  v-model="configForm[cfg.key]"
                  size="small"
                  class="w-36 !text-xs"
                  @change="saveConfig(cfg.key)"
                />
                <Select
                  v-else-if="cfg.type === 'select'"
                  v-model="configForm[cfg.key]"
                  :options="cfg.options"
                  size="small"
                  class="w-36 !text-xs"
                  @change="saveConfig(cfg.key)"
                />
                <InputNumber
                  v-else-if="cfg.type === 'number'"
                  v-model="configForm[cfg.key]"
                  size="small"
                  class="w-36 !text-xs"
                  @input="saveConfig(cfg.key)"
                />
                <Button
                  icon="pi pi-save"
                  text
                  rounded
                  size="small"
                  :loading="savingKey === cfg.key"
                  class="text-[var(--text-tertiary)] flex-shrink-0"
                  @click="saveConfig(cfg.key)"
                />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Security Mode -->
      <Card>
        <template #title>
          <div class="section-title">
            <i class="pi pi-shield"></i>
            <span>安全模式</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-4">
            <div
              class="rounded-xl p-4 border transition-all duration-300"
              :class="authStore.securityMode === 'vulnerable' ? 'status-pulse-danger' : ''"
              style="border-width: 1px;"
              :style="authStore.securityMode === 'vulnerable'
                ? { background: 'var(--danger-subtle)', borderColor: 'var(--danger-subtle)' }
                : { background: 'var(--success-subtle)', borderColor: 'var(--success-subtle)' }"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-11 h-11 rounded-xl flex items-center justify-center"
                  :style="authStore.securityMode === 'vulnerable'
                    ? { background: 'var(--danger-subtle)' }
                    : { background: 'var(--success-subtle)' }"
                >
                  <i
                    class="pi text-xl"
                    :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock' : 'pi-lock'"
                    :style="{ color: authStore.securityMode === 'vulnerable' ? 'var(--danger)' : 'var(--success)' }"
                  ></i>
                </div>
                <div>
                  <p class="font-semibold text-sm text-[var(--text-primary)]">
                    当前模式: {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
                  </p>
                  <p class="text-[10px] text-[var(--text-tertiary)] mt-0.5">
                    {{ authStore.securityMode === 'vulnerable'
                      ? '所有越权漏洞均可触发'
                      : '所有漏洞已被修复，权限校验严格' }}
                  </p>
                </div>
              </div>
            </div>

            <div class="info-row">
              <span class="text-xs text-[var(--text-secondary)]">切换安全模式</span>
              <ToggleSwitch
                v-model="isSecureMode"
                @change="toggleMode"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- Database Management -->
      <Card>
        <template #title>
          <div class="section-title">
            <i class="pi pi-database"></i>
            <span>数据库管理</span>
          </div>
        </template>
        <template #content>
          <div class="relative rounded-xl p-4 overflow-hidden" style="background: var(--danger-subtle); border: 1px solid var(--danger-subtle);">
            <div class="absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full" style="background: var(--danger);"></div>
            <div class="pl-3 space-y-4">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-lg flex items-center justify-center" style="background: var(--danger-subtle);">
                  <i class="pi pi-exclamation-triangle text-sm" style="color: var(--danger);"></i>
                </div>
                <div>
                  <p class="font-semibold text-sm text-[var(--text-primary)]">重置数据库</p>
                  <p class="text-[10px] text-[var(--text-tertiary)]">清除所有数据并恢复初始状态</p>
                </div>
              </div>
              <p class="text-xs text-[var(--text-secondary)] leading-relaxed">
                此操作将删除所有用户数据、VPS 实例、订单记录等，并将数据库恢复到初始种子状态。<strong style="color: var(--danger);">此操作不可撤销！</strong>
              </p>
              <Button
                label="重置数据库"
                icon="pi pi-refresh"
                severity="danger"
                size="small"
                @click="handleReset"
              />
            </div>
          </div>
        </template>
      </Card>

      <!-- System Info -->
      <Card>
        <template #title>
          <div class="section-title">
            <i class="pi pi-info-circle"></i>
            <span>系统信息</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row">
              <span class="text-xs text-[var(--text-tertiary)]">版本</span>
              <span class="text-xs font-medium text-[var(--text-primary)]">v1.0.0</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-[var(--text-tertiary)]">技术栈</span>
              <span class="text-xs font-medium text-[var(--text-primary)]">Go + Vue3 + PrimeVue</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-[var(--text-tertiary)]">数据库</span>
              <span class="text-xs font-medium text-[var(--text-primary)]">SQLite3</span>
            </div>
            <div class="info-row">
              <span class="text-xs text-[var(--text-tertiary)]">框架</span>
              <span class="text-xs font-medium text-[var(--text-primary)]">Gin + GORM</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Quick Links -->
      <Card>
        <template #title>
          <div class="section-title">
            <i class="pi pi-directions"></i>
            <span>快速导航</span>
          </div>
        </template>
        <template #content>
          <div class="space-y-2">
            <div class="info-row cursor-pointer group" @click="$router.push('/admin/users')">
              <div class="flex items-center gap-2">
                <i class="pi pi-users text-[var(--primary)] text-xs"></i>
                <span class="text-xs font-medium text-[var(--text-secondary)] group-hover:text-[var(--primary)] transition-colors">用户管理</span>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/admin/companies')">
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-[var(--primary)] text-xs"></i>
                <span class="text-xs font-medium text-[var(--text-secondary)] group-hover:text-[var(--primary)] transition-colors">企业管理</span>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/challenges')">
              <div class="flex items-center gap-2">
                <i class="pi pi-flag text-[var(--primary)] text-xs"></i>
                <span class="text-xs font-medium text-[var(--text-secondary)] group-hover:text-[var(--primary)] transition-colors">漏洞挑战</span>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-[10px]"></i>
            </div>
            <div class="info-row cursor-pointer group" @click="$router.push('/audit')">
              <div class="flex items-center gap-2">
                <i class="pi pi-history text-[var(--primary)] text-xs"></i>
                <span class="text-xs font-medium text-[var(--text-secondary)] group-hover:text-[var(--primary)] transition-colors">审计日志</span>
              </div>
              <i class="pi pi-chevron-right text-[var(--text-tertiary)] text-[10px]"></i>
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import * as api from '@/api'
import PageHeader from '@/components/PageHeader.vue'
import Card from 'primevue/card'
import Button from 'primevue/button'
import ToggleSwitch from 'primevue/toggleswitch'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'

const authStore = useAuthStore()
const toast = useToast()
const confirm = useConfirm()
const isSecureMode = ref(authStore.securityMode === 'secure')
const savingKey = ref<string | null>(null)

const configForm = ref<Record<string, any>>({
  site_name: 'CloudNest',
  allow_registration: 'true',
  default_vps_expire_days: '30',
  max_vps_per_user: '10',
})

const configFields = [
  { key: 'site_name', label: '站点名称', type: 'text' },
  { key: 'allow_registration', label: '允许注册', type: 'select', options: ['true', 'false'] },
  { key: 'default_vps_expire_days', label: 'VPS 默认过期天数', type: 'number' },
  { key: 'max_vps_per_user', label: '每用户最大 VPS 数', type: 'number' },
]

onMounted(async () => {
  try {
    const res = await api.getSystemConfig()
    if (res.data.data) {
      Object.assign(configForm.value, res.data.data)
    }
  } catch {}
})

async function saveConfig(key: string) {
  savingKey.value = key
  try {
    await api.updateSystemConfig({ key, value: String(configForm.value[key]) })
    toast.add({ severity: 'success', summary: '成功', detail: '配置已保存', life: 2000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '保存失败', life: 3000 })
  } finally {
    savingKey.value = null
  }
}

async function toggleMode() {
  await authStore.toggleSecurityMode()
  isSecureMode.value = authStore.securityMode === 'secure'
}

function handleReset() {
  confirm.require({
    message: '确定要重置数据库吗？所有数据将被清除并恢复初始状态，此操作不可撤销！',
    header: '确认重置数据库',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.adminResetDb()
        toast.add({ severity: 'success', summary: '成功', detail: '数据库已重置', life: 3000 })
        window.location.reload()
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: '重置失败', life: 3000 })
      }
    },
  })
}
</script>

<style scoped>
@keyframes pulse-danger {
  0% { box-shadow: 0 0 0 0 rgba(244, 63, 94, 0.3); }
  70% { box-shadow: 0 0 0 8px rgba(244, 63, 94, 0); }
  100% { box-shadow: 0 0 0 0 rgba(244, 63, 94, 0); }
}

.status-pulse-danger {
  animation: pulse-danger 2.5s infinite;
}
</style>
