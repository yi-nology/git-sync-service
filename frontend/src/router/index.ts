import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: AppLayout,
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      { path: '/dashboard', component: () => import('@/views/dashboard/Dashboard.vue') },
      { path: '/sync', component: () => import('@/views/sync/SyncTaskList.vue') },
      { path: '/sync/history', component: () => import('@/views/sync/SyncHistory.vue') },
      { path: '/sync/new', component: () => import('@/views/sync/NewSyncTask.vue') },
      { path: '/webhook/rules', component: () => import('@/views/webhook/WebhookRules.vue') },
      { path: '/webhook/events', component: () => import('@/views/webhook/WebhookEvents.vue') },
      { path: '/repos', component: () => import('@/views/repos/RepoList.vue') },
      { path: '/repos/config/:id', component: () => import('@/views/repos/RepoUnifiedConfig.vue') },
      { path: '/local-repos/:id', component: () => import('@/views/repos/LocalRepoDetail.vue') },
      { path: '/logs/operations', component: () => import('@/views/logs/OperationLogs.vue') },
      { path: '/logs/sync', component: () => import('@/views/logs/SyncLogs.vue') },
      { path: '/logs/system', component: () => import('@/views/logs/SystemLogs.vue') },
      { path: '/settings', component: () => import('@/views/settings/Settings.vue') },
      { path: '/settings/platforms', component: () => import('@/views/settings/PlatformSettings.vue') },
      { path: '/settings/advanced', component: () => import('@/views/settings/AdvancedSettings.vue') },
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth !== false)

  if (requiresAuth && !authStore.isAuthenticated) {
    // Redirect to login if not authenticated
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    // Redirect to dashboard if already authenticated
    next('/dashboard')
  } else {
    next()
  }
})

export default router
