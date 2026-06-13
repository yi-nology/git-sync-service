import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// 同步任务 API
export const syncTaskApi = {
  list: (params?: any) => api.get('/sync/tasks', { params }),
  get: (id: number) => api.get(`/sync/tasks/${id}`),
  create: (data: any) => api.post('/sync/tasks', data),
  update: (id: number, data: any) => api.put(`/sync/tasks/${id}`, data),
  delete: (id: number) => api.delete(`/sync/tasks/${id}`),
  run: (id: number) => api.post(`/sync/tasks/${id}/run`),
  history: (params?: any) => api.get('/sync/history', { params }),
}

// Webhook API
export const webhookApi = {
  rules: (params?: any) => api.get('/webhook/rules', { params }),
  createRule: (data: any) => api.post('/webhook/rules', data),
  updateRule: (id: number, data: any) => api.put(`/webhook/rules/${id}`, data),
  deleteRule: (id: number) => api.delete(`/webhook/rules/${id}`),
  events: (params?: any) => api.get('/webhook/events', { params }),
}

// 系统设置 API
export const settingsApi = {
  get: () => api.get('/settings'),
  update: (data: any) => api.put('/settings', data),
  credentials: () => api.get('/settings/credentials'),
  createCredential: (data: any) => api.post('/settings/credentials', data),
  updateCredential: (id: number, data: any) => api.put(`/settings/credentials/${id}`, data),
  deleteCredential: (id: number) => api.delete(`/settings/credentials/${id}`),
  testCredential: (id: number) => api.post(`/settings/credentials/${id}/test`),
}

// 仪表盘 API
export const dashboardApi = {
  stats: () => api.get('/dashboard/stats'),
  recentSyncs: () => api.get('/dashboard/recent-syncs'),
  chartData: () => api.get('/dashboard/chart-data'),
}

export default api
