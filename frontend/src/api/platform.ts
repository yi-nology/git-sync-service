import api from './index'

export interface Platform {
  id: number
  key: string
  name: string
  type: string
  instance_url: string
  api_url: string
  access_token?: string
  skip_tls_verify: boolean
  ca_cert_path: string
  proxy_url: string
  is_default: boolean
  status: string
  repo_count: number
  last_test_at?: string
  last_test_result?: string
  created_at: string
  updated_at: string
}

export interface CreatePlatformRequest {
  name: string
  type: string
  instance_url?: string
  api_url: string
  access_token: string
  skip_tls_verify?: boolean
  ca_cert_path?: string
  proxy_url?: string
  is_default?: boolean
}

export interface UpdatePlatformRequest {
  key: string
  name?: string
  instance_url?: string
  api_url?: string
  access_token?: string
  skip_tls_verify?: boolean
  ca_cert_path?: string
  proxy_url?: string
  is_default?: boolean
}

export const platformApi = {
  // 列出所有平台
  list: () =>
    api.get<any, { platforms: Platform[]; total: number }>('/platforms'),

  // 获取单个平台
  get: (key: string) =>
    api.get<any, { platform: Platform }>('/platform', { params: { key } }),

  // 创建平台
  create: (data: CreatePlatformRequest) =>
    api.post<any, { platform: Platform }>('/platform', data),

  // 更新平台
  update: (data: UpdatePlatformRequest) =>
    api.put<any, { platform: Platform }>('/platform', data),

  // 删除平台
  delete: (key: string) =>
    api.delete<any, void>('/platform', { params: { key } }),

  // 测试连接
  test: (key: string) =>
    api.post<any, { result: any }>('/platform/test', null, { params: { key } }),

  // 设置默认
  setDefault: (key: string) =>
    api.post<any, { message: string }>('/platform/set-default', null, { params: { key } }),

  // 列出平台上的仓库
  listRepos: (key: string) =>
    api.get<any, { repos: any[]; total: number }>('/platform/repos', { params: { key } }),

  // 同步仓库到本地
  syncRepos: (key: string) =>
    api.post<any, { message: string; syncedCount: number }>('/platform/sync-repos', null, { params: { key } }),
}
