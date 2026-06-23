import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
      },
      {
        path: 'vps',
        name: 'VpsList',
        component: () => import('@/views/vps/VpsListView.vue'),
      },
      {
        path: 'vps/:id',
        name: 'VpsDetail',
        component: () => import('@/views/vps/VpsDetailView.vue'),
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/user/ProfileView.vue'),
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/user/CompanyMembersView.vue'),
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/finance/OrdersView.vue'),
      },
      {
        path: 'bills',
        name: 'Bills',
        component: () => import('@/views/finance/BillsView.vue'),
      },
      {
        path: 'tickets',
        name: 'Tickets',
        component: () => import('@/views/tickets/TicketListView.vue'),
      },
      {
        path: 'tickets/:id',
        name: 'TicketDetail',
        component: () => import('@/views/tickets/TicketDetailView.vue'),
      },
      {
        path: 'apikeys',
        name: 'ApiKeys',
        component: () => import('@/views/apikeys/ApiKeysView.vue'),
      },
      {
        path: 'audit',
        name: 'AuditLogs',
        component: () => import('@/views/audit/AuditLogsView.vue'),
      },
      {
        path: 'challenges',
        name: 'Challenges',
        component: () => import('@/views/challenges/ChallengesView.vue'),
      },
      {
        path: 'admin/users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/AdminUsersView.vue'),
        meta: { requiresAdmin: true },
      },
      {
        path: 'admin/companies',
        name: 'AdminCompanies',
        component: () => import('@/views/admin/AdminCompaniesView.vue'),
        meta: { requiresAdmin: true },
      },
      {
        path: 'admin/system',
        name: 'AdminSystem',
        component: () => import('@/views/admin/AdminSystemView.vue'),
        meta: { requiresAdmin: true },
      },
      {
        path: 'admin/announcements',
        name: 'AdminAnnouncements',
        component: () => import('@/views/admin/AnnouncementsView.vue'),
        meta: { requiresAdmin: true },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token) {
    return '/login'
  }
  if (to.meta.requiresAdmin) {
    const user = localStorage.getItem('user')
    if (user) {
      try {
        const parsed = JSON.parse(user)
        if (parsed.user_type !== 'platform_admin') {
          return '/'
        }
      } catch {
        return '/login'
      }
    } else {
      return '/login'
    }
  }
  if (to.path === '/login' && token) {
    return '/'
  }
})

router.onError((error: any) => {
  // Handle dynamic import failures (e.g. Vite HMR timeout, network error)
  if (error.message?.includes('Failed to fetch dynamically imported module')) {
    console.warn('[Router] Dynamic import failed, reloading page...')
    window.location.href = router.currentRoute.value.fullPath
  }
})

export default router
