import axios from 'axios'
import type {
  ListReposResp, ListTasksResp, ListHistoryResp, ListRulesResp, ListEventsResp,
  TestConnectionResp, PreviewSyncResp,
  CreateRepoRequest, UpdateRepoRequest,
  CreateTaskRequest, UpdateTaskRequest,
  CreateRuleRequest, UpdateRuleRequest,
  Pagination,
} from '@/types'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

// Add request interceptor to include API key
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    const apiKey = authStore.getApiKey()
    if (apiKey) {
      config.headers['X-API-Key'] = apiKey
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error.response?.data || error.message)
    // If 401, clear API key and redirect to login
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.clearApiKey()
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

export const repoApi = {
  list: (params?: Pagination) =>
    api.get<any, ListReposResp>('/repos', { params }),
  get: (key: string) =>
    api.get<any, { repo: any }>('/repo', { params: { key } }),
  create: (data: CreateRepoRequest) =>
    api.post<any, { repo: any }>('/repo/create', data),
  update: (data: UpdateRepoRequest) =>
    api.post<any, { repo: any }>('/repo/update', data),
  delete: (key: string) =>
    api.post<any, { success: boolean }>('/repo/delete', null, { params: { key } }),
  testConnection: (key: string) =>
    api.post<any, TestConnectionResp>('/repo/test', null, { params: { key } }),
  listBranches: (key: string) =>
    api.get<any, { branches: string[] }>('/repo/branches', { params: { key } }),
}

export const syncTaskApi = {
  list: (params?: { repo_key?: string } & Pagination) =>
    api.get<any, ListTasksResp>('/sync/tasks', { params }),
  get: (key: string) =>
    api.get<any, { task: any }>('/sync/task', { params: { key } }),
  create: (data: CreateTaskRequest) =>
    api.post<any, { task: any }>('/sync/task/create', data),
  update: (data: UpdateTaskRequest) =>
    api.post<any, { task: any }>('/sync/task/update', data),
  delete: (key: string) =>
    api.post<any, { success: boolean }>('/sync/task/delete', null, { params: { key } }),
  run: (key: string) =>
    api.post<any, { success: boolean; message: string }>('/sync/task/run', null, { params: { key } }),
  preview: (data: { source_repo_key: string; source_branch?: string; target_repo_key: string; target_branch?: string }) =>
    api.post<any, PreviewSyncResp>('/sync/preview', data),
  history: (params: { task_key: string; limit?: number }) =>
    api.get<any, ListHistoryResp>('/sync/history', { params }),
}

export const webhookApi = {
  listRules: (repo_key: string) =>
    api.get<any, ListRulesResp>('/webhook/rules', { params: { repo_key } }),
  getRule: (id: number) =>
    api.get<any, { rule: any }>('/webhook/rule', { params: { id } }),
  createRule: (data: CreateRuleRequest) =>
    api.post<any, { rule: any }>('/webhook/rule/create', data),
  updateRule: (data: UpdateRuleRequest) =>
    api.post<any, { rule: any }>('/webhook/rule/update', data),
  deleteRule: (id: number) =>
    api.post<any, { success: boolean }>('/webhook/rule/delete', null, { params: { id } }),
  listEvents: (params: { repo_key: string; limit?: number }) =>
    api.get<any, ListEventsResp>('/webhook/events', { params }),
  retryEvent: (id: number) =>
    api.post<any, { success: boolean; message: string }>('/webhook/event/retry', null, { params: { id } }),
}

export interface OperationLogRequest {
  page?: number
  page_size?: number
  search?: string
  action?: string
  user?: string
  start_date?: string
  end_date?: string
}

export interface OperationLog {
  id: number
  time: string
  user: string
  action: string
  resource: string
  details: string
  ip: string
}

export interface OperationLogResp {
  list: OperationLog[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

export interface SyncLogRequest {
  page?: number
  page_size?: number
  task_key?: string
  status?: string
}

export interface SyncLog {
  id: number
  task_key: string
  trigger_source: string
  status: string
  start_time: string
  end_time: string
  commit_range: string
  details: string
  error_message: string
}

export interface SyncLogResp {
  list: SyncLog[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

export interface SystemLogRequest {
  page?: number
  page_size?: number
  level?: string
}

export interface SystemLog {
  id: number
  time: string
  level: string
  message: string
  details: string
}

export interface SystemLogResp {
  list: SystemLog[]
  pagination: {
    page: number
    page_size: number
    total: number
  }
}

export const logApi = {
  listOperations: (params?: OperationLogRequest) =>
    api.get<any, OperationLogResp>('/logs/operations', { params }),
  listSync: (params?: SyncLogRequest) =>
    api.get<any, SyncLogResp>('/logs/sync', { params }),
  listSystem: (params?: SystemLogRequest) =>
    api.get<any, SystemLogResp>('/logs/system', { params }),
}

export interface SystemStatusResp {
  status: string
  version: string
  uptime: number
  repoCount: number
  taskCount: number
  runningTask: number
  lastSyncAt: string
}

export interface HealthResp {
  status: string
  database: {
    status: string
    size: number
  }
}

export const systemApi = {
  status: () =>
    api.get<any, SystemStatusResp>('/system/status'),
  health: () =>
    api.get<any, HealthResp>('/system/health'),
}

export default api
