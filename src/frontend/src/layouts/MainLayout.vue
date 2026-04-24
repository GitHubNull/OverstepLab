<template>
  <div class="min-h-screen bg-slate-50 dark:bg-slate-900">
    <div class="flex">
      <!-- Sidebar -->
      <aside 
        class="fixed lg:static inset-y-0 left-0 z-50 w-72 bg-white dark:bg-slate-800 shadow-xl transition-transform duration-300"
        :class="{ '-translate-x-full lg:translate-x-0': !sidebarOpen }"
      >
        <!-- Logo -->
        <div class="p-6 border-b border-slate-200 dark:border-slate-700">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center shadow-lg">
              <i class="pi pi-shield text-xl text-white"></i>
            </div>
            <div>
              <h1 class="text-xl font-bold text-slate-800 dark:text-white">OverstepLab</h1>
              <p class="text-xs text-slate-500">越权漏洞测试靶场</p>
            </div>
          </div>
        </div>

        <!-- Navigation -->
        <nav class="p-4 space-y-1">
          <template v-for="item in menuItems" :key="item.label">
            <!-- Section Header -->
            <div v-if="item.isSection" class="px-3 py-2 text-xs font-semibold text-slate-400 uppercase tracking-wider mt-4 first:mt-0">
              {{ item.label }}
            </div>

            <!-- Menu Item with Children -->
            <div v-else-if="item.children" class="mb-2">
              <div 
                @click="item.expanded = !item.expanded"
                class="flex items-center justify-between px-3 py-2.5 rounded-xl cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors"
              >
                <div class="flex items-center gap-3">
                  <i :class="[item.icon, 'text-lg text-slate-500']"></i>
                  <span class="text-sm font-medium text-slate-700 dark:text-slate-300">{{ item.label }}</span>
                </div>
                <i 
                  class="pi text-xs text-slate-400 transition-transform"
                  :class="item.expanded ? 'pi-chevron-up' : 'pi-chevron-down'"
                ></i>
              </div>
              <div v-show="item.expanded" class="ml-4 mt-1 space-y-1">
                <router-link
                  v-for="child in item.children"
                  :key="child.label"
                  :to="child.to"
                  class="flex items-center gap-3 px-3 py-2 rounded-xl transition-colors"
                  :class="route.path === child.to 
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400' 
                    : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'"
                >
                  <i :class="[child.icon, 'text-sm']"></i>
                  <span class="text-sm">{{ child.label }}</span>
                </router-link>
              </div>
            </div>

            <!-- Single Menu Item -->
            <router-link
              v-else
              :to="item.to"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-colors"
              :class="route.path === item.to 
                ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 shadow-sm' 
                : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-700'"
            >
              <i :class="[item.icon, 'text-lg']"></i>
              <span class="text-sm font-medium">{{ item.label }}</span>
              <Badge v-if="item.badge" :value="item.badge" severity="danger" class="ml-auto" />
            </router-link>
          </template>
        </nav>

        <!-- Security Mode Toggle -->
        <div class="absolute bottom-0 left-0 right-0 p-4 border-t border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <i 
                class="pi text-lg"
                :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-500' : 'pi-lock text-green-500'"
              ></i>
              <span class="text-sm font-medium text-slate-700 dark:text-slate-300">
                {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
              </span>
            </div>
            <ToggleSwitch 
              v-model="isSecureMode"
              @change="toggleMode"
            />
          </div>
        </div>
      </aside>

      <!-- Mobile Sidebar Overlay -->
      <div 
        v-if="sidebarOpen" 
        @click="sidebarOpen = false"
        class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      ></div>

      <!-- Main Content -->
      <div class="flex-1 min-h-screen flex flex-col">
        <!-- Topbar -->
        <header class="bg-white dark:bg-slate-800 shadow-sm border-b border-slate-200 dark:border-slate-700 sticky top-0 z-30">
          <div class="flex items-center justify-between px-6 py-4">
            <div class="flex items-center gap-4">
              <Button 
                icon="pi pi-bars" 
                text 
                rounded
                class="lg:hidden"
                @click="sidebarOpen = !sidebarOpen"
              />
              <div>
                <h2 class="text-xl font-bold text-slate-800 dark:text-white">{{ pageTitle }}</h2>
                <p class="text-sm text-slate-500">{{ pageDescription }}</p>
              </div>
            </div>

            <div class="flex items-center gap-4">
              <!-- User Menu -->
              <div class="flex items-center gap-3">
                <Avatar 
                  :label="authStore.user?.username?.charAt(0).toUpperCase()" 
                  shape="circle"
                  class="bg-primary-100 text-primary-700"
                />
                <div class="hidden md:block">
                  <p class="text-sm font-medium text-slate-800 dark:text-white">{{ authStore.user?.username }}</p>
                  <p class="text-xs text-slate-500">{{ userRoleText }}</p>
                </div>
                <Button 
                  icon="pi pi-sign-out" 
                  text 
                  rounded
                  severity="secondary"
                  @click="handleLogout"
                  title="退出登录"
                />
              </div>
            </div>
          </div>
        </header>

        <!-- Page Content -->
        <main class="flex-1 p-6">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'
import Badge from 'primevue/badge'
import Avatar from 'primevue/avatar'
import ToggleSwitch from 'primevue/toggleswitch'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const sidebarOpen = ref(false)
const isSecureMode = ref(authStore.securityMode === 'secure')

const pageTitle = computed(() => {
  const names: Record<string, string> = {
    Dashboard: '仪表盘',
    VpsList: 'VPS 管理',
    VpsDetail: 'VPS 详情',
    Profile: '个人信息',
    Members: '企业成员',
    Orders: '订单管理',
    Bills: '账单管理',
    Tickets: '工单系统',
    ApiKeys: 'API Key 管理',
    AuditLogs: '审计日志',
    Challenges: '漏洞挑战',
    AdminUsers: '用户管理',
    AdminCompanies: '企业管理',
    AdminSystem: '系统设置',
  }
  return names[route.name as string] || 'OverstepLab'
})

const pageDescription = computed(() => {
  const descriptions: Record<string, string> = {
    Dashboard: '概览您的资源和挑战进度',
    VpsList: '管理您的虚拟服务器实例',
    Challenges: '发现并练习越权漏洞',
  }
  return descriptions[route.name as string] || ''
})

const userRoleText = computed(() => {
  if (authStore.isAdmin) return '平台管理员'
  if (authStore.isCompanyAdmin) return '企业管理员'
  if (authStore.isOperator) return '运维人员'
  if (authStore.isFinance) return '财务人员'
  if (authStore.isViewer) return '只读成员'
  return '个人用户'
})

const menuItems = computed(() => {
  const items: any[] = [
    { label: '仪表盘', icon: 'pi pi-home', to: '/' },
    { label: 'VPS 管理', icon: 'pi pi-server', to: '/vps' },
  ]

  if (authStore.user?.user_type === 'company' || authStore.isAdmin) {
    items.push({ label: '企业成员', icon: 'pi pi-users', to: '/members' })
  }

  items.push(
    { label: '财务管理', icon: 'pi pi-wallet', to: '/bills' },
    { label: '订单管理', icon: 'pi pi-list', to: '/orders' },
    { label: '工单系统', icon: 'pi pi-comments', to: '/tickets' },
    { label: 'API Key', icon: 'pi pi-key', to: '/apikeys' },
    { label: '审计日志', icon: 'pi pi-history', to: '/audit' },
    { 
      label: '漏洞挑战', 
      icon: 'pi pi-flag', 
      to: '/challenges',
      badge: '13'
    },
  )

  if (authStore.isAdmin) {
    items.push({ label: '平台管理', isSection: true })
    items.push({
      label: '管理后台',
      icon: 'pi pi-cog',
      expanded: false,
      children: [
        { label: '用户管理', icon: 'pi pi-users', to: '/admin/users' },
        { label: '企业管理', icon: 'pi pi-building', to: '/admin/companies' },
        { label: '系统设置', icon: 'pi pi-cog', to: '/admin/system' },
      ],
    })
  }

  return items
})

async function toggleMode() {
  await authStore.toggleSecurityMode()
  isSecureMode.value = authStore.securityMode === 'secure'
}

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>
