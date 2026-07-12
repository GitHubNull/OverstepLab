<template>
  <div class="min-h-screen flex bg-[var(--bg-base)]">
    <!-- Sidebar -->
    <aside
      class="fixed lg:static inset-y-0 left-0 z-50 w-[260px] bg-[var(--bg-surface)] border-r border-[var(--border-default)] transition-transform duration-300 flex flex-col"
      :class="{ '-translate-x-full lg:translate-x-0': !sidebarOpen }"
    >
      <!-- Logo -->
      <div class="h-16 flex items-center gap-3 px-5 border-b border-[var(--border-default)] flex-shrink-0">
        <div class="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0" style="background: linear-gradient(135deg, var(--primary) 0%, var(--primary-hover) 100%);">
          <i class="pi pi-shield text-sm text-white"></i>
        </div>
        <div>
          <h1 class="text-[15px] font-bold text-[var(--text-primary)] leading-tight">OverstepLab</h1>
          <p class="text-[10px] text-[var(--text-tertiary)] font-medium tracking-wider">越权漏洞测试靶场</p>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-3 py-4 overflow-y-auto">
        <template v-for="(group, gi) in menuGroups" :key="gi">
          <div v-if="group.items.length > 0">
            <!-- Group label -->
            <div v-if="group.label" class="px-3 mb-2 text-[11px] font-semibold text-[var(--text-tertiary)] uppercase tracking-wider">
              {{ group.label }}
            </div>
            <!-- Group items -->
            <div class="space-y-0.5 mb-4">
              <template v-for="item in group.items" :key="item.label">
                <!-- Item with children (expandable) -->
                <div v-if="item.children">
                  <div
                    @click="item.expanded = !item.expanded"
                    class="flex items-center justify-between px-3 h-10 rounded-lg cursor-pointer transition-colors"
                    :class="isChildActive(item.children) ? 'text-[var(--primary)] bg-[var(--primary-subtle)]' : 'text-[var(--text-secondary)] hover:bg-[var(--bg-surface-hover)]'"
                  >
                    <div class="flex items-center gap-3">
                      <i :class="[item.icon, 'text-[14px]']"></i>
                      <span class="text-[13px] font-medium">{{ item.label }}</span>
                    </div>
                    <i
                      class="pi text-[10px] transition-transform duration-200"
                      :class="item.expanded ? 'pi-angle-up' : 'pi-angle-down'"
                    ></i>
                  </div>
                  <div v-show="item.expanded" class="mt-0.5 ml-4 pl-4 border-l-2 border-[var(--border-default)] space-y-0.5">
                    <router-link
                      v-for="child in item.children"
                      :key="child.label"
                      :to="child.to || ''"
                      class="flex items-center gap-2.5 px-3 h-9 rounded-lg transition-colors text-[13px]"
                      :class="route.path === child.to
                        ? 'text-[var(--primary)] font-medium bg-[var(--primary-subtle)]'
                        : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)]'"
                      @click="closeMobileSidebar"
                    >
                      <span>{{ child.label }}</span>
                    </router-link>
                  </div>
                </div>

                <!-- Single item -->
                <router-link
                  v-else
                  :to="item.to || ''"
                  class="flex items-center gap-3 px-3 h-10 rounded-lg transition-colors relative"
                  :class="route.path === item.to
                    ? 'text-[var(--primary)] font-medium bg-[var(--primary-subtle)]'
                    : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)]'"
                  @click="closeMobileSidebar"
                >
                  <div
                    v-if="route.path === item.to"
                    class="absolute left-0 top-2 bottom-2 w-[3px] rounded-r-full"
                    style="background: var(--primary);"
                  ></div>
                  <i :class="[item.icon, 'text-[14px]']"></i>
                  <span class="text-[13px] font-medium">{{ item.label }}</span>
                  <Badge v-if="item.badge" :value="item.badge" severity="danger" class="ml-auto" />
                </router-link>
              </template>
            </div>
          </div>
        </template>
      </nav>

      <!-- Security Mode Panel -->
      <div class="p-3 border-t border-[var(--border-default)] flex-shrink-0">
        <div
          class="rounded-xl px-3 py-3 transition-all"
          :style="authStore.securityMode === 'vulnerable'
            ? 'background: var(--danger-subtle); border: 1px solid color-mix(in srgb, var(--danger) 20%, transparent);'
            : 'background: var(--success-subtle); border: 1px solid color-mix(in srgb, var(--success) 20%, transparent);'"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <div
                class="w-7 h-7 rounded-lg flex items-center justify-center"
                :style="authStore.securityMode === 'vulnerable' ? 'background: color-mix(in srgb, var(--danger) 15%, transparent);' : 'background: color-mix(in srgb, var(--success) 15%, transparent);'"
              >
                <i
                  class="pi text-xs"
                  :class="authStore.securityMode === 'vulnerable' ? 'pi-unlock' : 'pi-lock'"
                  :style="{ color: authStore.securityMode === 'vulnerable' ? 'var(--danger)' : 'var(--success)' }"
                ></i>
              </div>
              <div>
                <span
                  class="text-[12px] font-semibold block leading-tight"
                  :style="{ color: authStore.securityMode === 'vulnerable' ? 'var(--danger)' : 'var(--success)' }"
                >
                  {{ authStore.securityMode === 'vulnerable' ? '漏洞模式' : '安全模式' }}
                </span>
              </div>
            </div>
            <ToggleSwitch
              v-model="isSecureMode"
              @change="toggleMode"
              class="scale-[0.8]"
            />
          </div>
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
      <header
        class="sticky top-0 z-30 h-14 flex items-center justify-between px-4 lg:px-6 border-b border-[var(--border-default)]"
        style="background: color-mix(in srgb, var(--bg-surface) 85%, transparent); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px);"
      >
        <div class="flex items-center gap-3 min-w-0">
          <Button
            icon="pi pi-bars"
            text
            rounded
            size="small"
            class="lg:hidden text-[var(--text-secondary)] flex-shrink-0"
            @click="sidebarOpen = !sidebarOpen"
          />
          <!-- Breadcrumb -->
          <nav v-if="breadcrumbs.length > 0" class="flex items-center gap-1.5 text-sm min-w-0">
            <template v-for="(crumb, index) in breadcrumbs" :key="crumb.path">
              <router-link
                v-if="index < breadcrumbs.length - 1"
                :to="crumb.path"
                class="text-[var(--text-tertiary)] hover:text-[var(--text-primary)] transition-colors text-xs font-medium truncate"
              >
                {{ crumb.label }}
              </router-link>
              <span v-else class="text-[var(--text-primary)] font-semibold text-xs truncate">
                {{ crumb.label }}
              </span>
              <i v-if="index < breadcrumbs.length - 1" class="pi pi-chevron-right text-[8px] text-[var(--text-tertiary)] flex-shrink-0"></i>
            </template>
          </nav>
        </div>

        <div class="flex items-center gap-1.5 flex-shrink-0">
          <!-- Active Encoding Challenge Indicator -->
          <div
            v-if="encodingStore.isActive"
            class="hidden md:flex items-center gap-2 px-2.5 py-1 rounded-lg bg-[var(--primary-subtle)] border border-[var(--primary)]/20 cursor-pointer"
            @click="$router.push('/challenges')"
          >
            <i class="pi pi-lock text-[10px] text-[var(--primary)]"></i>
            <span class="text-[11px] font-medium text-[var(--primary)]">
              {{ encodingStore.activeChallenge?.id }}
            </span>
          </div>

          <ThemeToggle />

          <Button
            icon="pi pi-bell"
            text
            rounded
            size="small"
            class="text-[var(--text-secondary)] hidden sm:flex"
            @click="toast.add({ severity: 'info', summary: '提示', detail: '暂无新通知', life: 2000 })"
          />

          <div class="w-px h-5 bg-[var(--border-default)] mx-1 hidden sm:block"></div>

          <!-- User Menu -->
          <div class="relative">
            <button
              @click="userMenuOpen = !userMenuOpen"
              class="flex items-center gap-2 pl-2 pr-2.5 py-1.5 rounded-lg hover:bg-[var(--bg-surface-hover)] transition-colors"
            >
              <Avatar
                :label="authStore.user?.username?.charAt(0).toUpperCase()"
                shape="circle"
                size="small"
                class="bg-[var(--primary-subtle)] text-[var(--primary)] text-xs w-7 h-7"
              />
              <span class="hidden md:block text-xs font-medium text-[var(--text-primary)]">{{ authStore.user?.username }}</span>
              <i class="pi pi-chevron-down text-[9px] text-[var(--text-tertiary)] hidden md:block"></i>
            </button>

            <!-- User Dropdown -->
            <div
              v-if="userMenuOpen"
              v-click-outside="() => userMenuOpen = false"
              class="absolute right-0 top-full mt-1.5 w-52 bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl shadow-lg py-1 z-50"
            >
              <div class="px-3.5 py-2.5 border-b border-[var(--border-default)]">
                <p class="text-sm font-semibold text-[var(--text-primary)]">{{ authStore.user?.username }}</p>
                <p class="text-[11px] text-[var(--text-tertiary)] mt-0.5">{{ authStore.user?.email }}</p>
              </div>
              <div class="py-1">
                <router-link
                  to="/profile"
                  class="flex items-center gap-2.5 px-3.5 py-2 text-[13px] text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-surface-hover)] transition-colors"
                  @click="userMenuOpen = false"
                >
                  <i class="pi pi-user text-xs"></i>
                  个人资料
                </router-link>
                <button
                  class="w-full flex items-center gap-2.5 px-3.5 py-2 text-[13px] text-[var(--danger)] hover:bg-[var(--danger-subtle)] transition-colors text-left"
                  @click="handleLogout"
                >
                  <i class="pi pi-sign-out text-xs"></i>
                  退出登录
                </button>
              </div>
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
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useEncodingChallengeStore } from '@/stores/encodingChallenge'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Badge from 'primevue/badge'
import Avatar from 'primevue/avatar'
import ToggleSwitch from 'primevue/toggleswitch'
import ThemeToggle from '@/components/ThemeToggle.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const encodingStore = useEncodingChallengeStore()
const toast = useToast()

const sidebarOpen = ref(false)
const userMenuOpen = ref(false)
const isSecureMode = ref(authStore.securityMode === 'secure')

function closeMobileSidebar() {
  if (window.innerWidth < 1024) {
    sidebarOpen.value = false
  }
}

function isChildActive(children: any[]) {
  return children.some(c => route.path === c.to || route.path.startsWith(c.to + '/'))
}

const breadcrumbs = computed(() => {
  const crumbs: { label: string; path: string }[] = []
  const name = route.name as string

  if (route.path.startsWith('/admin')) {
    crumbs.push({ label: '管理后台', path: '/admin/users' })
    if (name === 'AdminUsers') crumbs.push({ label: '用户管理', path: '/admin/users' })
    else if (name === 'AdminCompanies') crumbs.push({ label: '企业管理', path: '/admin/companies' })
    else if (name === 'AdminSystem') crumbs.push({ label: '系统设置', path: '/admin/system' })
  } else if (route.path.startsWith('/vps/')) {
    crumbs.push({ label: 'VPS 管理', path: '/vps' })
    crumbs.push({ label: '实例详情', path: route.path })
  } else if (route.path.startsWith('/tickets/')) {
    crumbs.push({ label: '工单系统', path: '/tickets' })
    crumbs.push({ label: '工单详情', path: route.path })
  } else {
    const labels: Record<string, string> = {
      Dashboard: '仪表盘',
      VpsList: 'VPS 管理',
      Profile: '个人资料',
      Members: '企业成员',
      Bills: '账单管理',
      Orders: '订单管理',
      Tickets: '工单系统',
      ApiKeys: 'API Key 管理',
      AuditLogs: '审计日志',
      Challenges: '漏洞挑战',
      CryptoTools: '编码工具',
    }
    if (labels[name]) {
      crumbs.push({ label: labels[name], path: route.path })
    }
  }

  return crumbs
})

interface MenuItem {
  label: string
  icon?: string
  to?: string
  badge?: string
  expanded?: boolean
  children?: MenuItem[]
}

interface MenuGroup {
  label?: string
  items: MenuItem[]
}

const menuGroups = computed<MenuGroup[]>(() => {
  const groups: MenuGroup[] = []

  // Core - no group label
  groups.push({
    label: undefined,
    items: [
      { label: '仪表盘', icon: 'pi pi-th-large', to: '/' },
      { label: 'VPS 管理', icon: 'pi pi-server', to: '/vps' },
    ]
  })

  // Account
  const accountItems: MenuItem[] = [
    { label: '个人资料', icon: 'pi pi-user', to: '/profile' },
  ]
  if (authStore.user?.user_type === 'company') {
    accountItems.push({ label: '企业成员', icon: 'pi pi-users', to: '/members' })
  }
  accountItems.push({ label: 'API Key', icon: 'pi pi-key', to: '/apikeys' })
  groups.push({ label: '账户', items: accountItems })

  // Finance
  groups.push({
    label: '财务',
    items: [
      { label: '账单管理', icon: 'pi pi-wallet', to: '/bills' },
      { label: '订单管理', icon: 'pi pi-receipt', to: '/orders' },
    ]
  })

  // Support
  groups.push({
    label: '支持',
    items: [
      { label: '工单系统', icon: 'pi pi-comments', to: '/tickets' },
      { label: '审计日志', icon: 'pi pi-history', to: '/audit' },
    ]
  })

  // Challenges
  groups.push({
    label: undefined,
    items: [
      { label: '编码工具', icon: 'pi pi-wrench', to: '/tools' },
      { label: '漏洞挑战', icon: 'pi pi-flag', to: '/challenges', badge: '21' },
    ]
  })

  // Admin
  if (authStore.isAdmin) {
    groups.push({
      label: '管理',
      items: [
        {
          label: '管理后台',
          icon: 'pi pi-cog',
          expanded: route.path.startsWith('/admin'),
          children: [
            { label: '用户管理', to: '/admin/users' },
            { label: '企业管理', to: '/admin/companies' },
            { label: '公告管理', to: '/admin/announcements' },
            { label: '系统设置', to: '/admin/system' },
          ],
        },
      ]
    })
  }

  return groups
})

async function toggleMode() {
  await authStore.toggleSecurityMode()
  isSecureMode.value = authStore.securityMode === 'secure'
}

async function handleLogout() {
  userMenuOpen.value = false
  await authStore.logout()
  router.push('/login')
}

// Click outside directive
const vClickOutside = {
  mounted(el: HTMLElement, binding: any) {
    (el as any)._clickOutside = (event: Event) => {
      if (!(el === event.target || el.contains(event.target as Node))) {
        binding.value()
      }
    }
    document.addEventListener('click', (el as any)._clickOutside, true)
  },
  unmounted(el: HTMLElement) {
    document.removeEventListener('click', (el as any)._clickOutside, true)
  }
}
</script>
