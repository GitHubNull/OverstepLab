<template>
  <div class="min-h-screen bg-[#f0f2f5] dark:bg-slate-900">
    <div class="flex">
      <!-- Sidebar -->
      <aside
        class="fixed lg:static inset-y-0 left-0 z-50 w-[260px] bg-white dark:bg-slate-800 border-r border-slate-200/80 dark:border-slate-700/50 transition-transform duration-300 flex flex-col shadow-sm"
        :class="{ '-translate-x-full lg:translate-x-0': !sidebarOpen }"
      >
        <!-- Logo -->
        <div class="p-5 flex items-center gap-3">
          <div class="w-9 h-9 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-md shadow-indigo-200">
            <i class="pi pi-shield text-base text-white"></i>
          </div>
          <div>
            <h1 class="text-base font-bold text-slate-800 dark:text-white leading-tight">OverstepLab</h1>
            <p class="text-[10px] text-slate-400 font-medium tracking-wide">越权漏洞测试靶场</p>
          </div>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 px-3 py-2 space-y-0.5 overflow-y-auto">
          <template v-for="item in menuItems" :key="item.label">
            <!-- Section Header -->
            <div v-if="item.isSection" class="px-3 py-2 text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-5 first:mt-0">
              {{ item.label }}
            </div>

            <!-- Menu Item with Children -->
            <div v-else-if="item.children">
              <div
                @click="item.expanded = !item.expanded"
                class="flex items-center justify-between px-3 py-2 rounded-lg cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors group"
              >
                <div class="flex items-center gap-2.5">
                  <i :class="[item.icon, 'text-sm text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-300 transition-colors']"></i>
                  <span class="text-[13px] font-medium text-slate-600 dark:text-slate-300">{{ item.label }}</span>
                </div>
                <i
                  class="pi text-[10px] text-slate-400 transition-transform duration-200"
                  :class="item.expanded ? 'pi-chevron-up' : 'pi-chevron-down'"
                ></i>
              </div>
              <div v-show="item.expanded" class="ml-3 mt-0.5 space-y-0.5 pl-3 border-l border-slate-200 dark:border-slate-700">
                <router-link
                  v-for="child in item.children"
                  :key="child.label"
                  :to="child.to"
                  class="flex items-center gap-2.5 px-3 py-1.5 rounded-lg transition-all duration-150"
                  :class="route.path === child.to
                    ? 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 font-medium'
                    : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-700/50 hover:text-slate-700 dark:hover:text-slate-300'"
                >
                  <i :class="[child.icon, 'text-xs']"></i>
                  <span class="text-[13px]">{{ child.label }}</span>
                </router-link>
              </div>
            </div>

            <!-- Single Menu Item -->
            <router-link
              v-else
              :to="item.to"
              class="flex items-center gap-2.5 px-3 py-2 rounded-lg transition-all duration-150 group"
              :class="route.path === item.to
                ? 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 font-medium'
                : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-700/50 hover:text-slate-700 dark:hover:text-slate-300'"
            >
              <i :class="[item.icon, 'text-sm']"></i>
              <span class="text-[13px] font-medium">{{ item.label }}</span>
              <Badge v-if="item.badge" :value="item.badge" severity="danger" class="ml-auto text-[10px]" />
            </router-link>
          </template>
        </nav>

        <!-- Security Mode Toggle -->
        <div class="p-3 border-t border-slate-200/80 dark:border-slate-700/50">
          <div
            class="flex items-center justify-between px-3 py-2.5 rounded-xl transition-colors"
            :class="authStore.securityMode === 'vulnerable'
              ? 'bg-red-50 dark:bg-red-900/10'
              : 'bg-emerald-50 dark:bg-emerald-900/10'"
          >
            <div class="flex items-center gap-2">
              <div
                class="w-6 h-6 rounded-md flex items-center justify-center"
                :class="authStore.securityMode === 'vulnerable' ? 'bg-red-100 dark:bg-red-900/30' : 'bg-emerald-100 dark:bg-emerald-900/30'"
              >
                <i
                  class="pi text-xs"
                  :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock text-red-500' : 'pi-lock text-emerald-500'"
                ></i>
              </div>
              <span class="text-xs font-semibold" :class="authStore.securityMode === 'vulnerable' ? 'text-red-600 dark:text-red-400' : 'text-emerald-600 dark:text-emerald-400'">
                {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
              </span>
            </div>
            <ToggleSwitch
              v-model="isSecureMode"
              @change="toggleMode"
              class="scale-90"
            />
          </div>
        </div>
      </aside>

      <!-- Mobile Sidebar Overlay -->
      <div
        v-if="sidebarOpen"
        @click="sidebarOpen = false"
        class="fixed inset-0 bg-black/30 backdrop-blur-sm z-40 lg:hidden"
      ></div>

      <!-- Main Content -->
      <div class="flex-1 min-h-screen flex flex-col min-w-0">
        <!-- Topbar -->
        <header class="bg-white/80 dark:bg-slate-800/80 backdrop-blur-md border-b border-slate-200/60 dark:border-slate-700/40 sticky top-0 z-30">
          <div class="flex items-center justify-between px-6 h-14">
            <div class="flex items-center gap-3">
              <Button
                icon="pi pi-bars"
                text
                rounded
                size="small"
                class="lg:hidden text-slate-500"
                @click="sidebarOpen = !sidebarOpen"
              />
              <div>
                <h2 class="text-base font-bold text-slate-800 dark:text-white leading-tight">{{ pageTitle }}</h2>
                <p v-if="pageDescription" class="text-[11px] text-slate-400 mt-0.5">{{ pageDescription }}</p>
              </div>
            </div>

            <div class="flex items-center gap-3">
              <!-- User Menu -->
              <div class="flex items-center gap-2.5 pl-3 border-l border-slate-200 dark:border-slate-700">
                <Avatar
                  :label="authStore.user?.username?.charAt(0).toUpperCase()"
                  shape="circle"
                  size="small"
                  class="bg-indigo-100 text-indigo-600 text-xs"
                />
                <div class="hidden md:block">
                  <p class="text-xs font-semibold text-slate-700 dark:text-slate-200 leading-tight">{{ authStore.user?.username }}</p>
                  <p class="text-[10px] text-slate-400">{{ userRoleText }}</p>
                </div>
                <Button
                  icon="pi pi-sign-out"
                  text
                  rounded
                  size="small"
                  severity="secondary"
                  class="text-slate-400 hover:text-red-500"
                  @click="handleLogout"
                  title="退出登录"
                />
              </div>
            </div>
          </div>
        </header>

        <!-- Page Content -->
        <main class="flex-1 p-5 lg:p-6">
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
    Profile: '管理您的账户信息',
    Bills: '查看消费记录和账户余额',
    Orders: '查看购买和续费记录',
    Tickets: '提交和管理技术支持请求',
    ApiKeys: '管理您的 API 访问凭证',
    AuditLogs: '查看系统操作记录',
    Members: '管理企业内的用户和权限',
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
    { label: '仪表盘', icon: 'pi pi-th-large', to: '/' },
    { label: 'VPS 管理', icon: 'pi pi-server', to: '/vps' },
  ]

  if (authStore.user?.user_type === 'company' || authStore.isAdmin) {
    items.push({ label: '企业成员', icon: 'pi pi-users', to: '/members' })
  }

  items.push(
    { label: '财务管理', icon: 'pi pi-wallet', to: '/bills' },
    { label: '订单管理', icon: 'pi pi-receipt', to: '/orders' },
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
    items.push({ label: '管理', isSection: true })
    items.push({
      label: '管理后台',
      icon: 'pi pi-cog',
      expanded: false,
      children: [
        { label: '用户管理', icon: 'pi pi-user', to: '/admin/users' },
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
