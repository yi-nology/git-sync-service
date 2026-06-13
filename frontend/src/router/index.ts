import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'

const routes = [
  {
    path: '/',
    component: AppLayout,
    redirect: '/sync',
    children: [
      { path: '/sync', component: () => import('@/views/sync/SyncTaskList.vue') },
      { path: '/sync/history', component: () => import('@/views/sync/SyncHistory.vue') },
      { path: '/sync/new', component: () => import('@/views/sync/NewSyncTask.vue') },
      { path: '/webhook/rules', component: () => import('@/views/webhook/WebhookRules.vue') },
      { path: '/webhook/events', component: () => import('@/views/webhook/WebhookEvents.vue') },
        { path: '/settings', component: () => import('@/views/settings/Settings.vue') },
        { path: '/settings/advanced', component: () => import('@/views/settings/AdvancedSettings.vue') },
        { path: '/repos', component: () => import('@/views/repos/RepoList.vue') },
        { path: '/repos/config/:id', component: () => import('@/views/repos/RepoUnifiedConfig.vue') },
        { path: '/local-repos/:id', component: () => import('@/views/repos/LocalRepoDetail.vue') },
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })
export default router
