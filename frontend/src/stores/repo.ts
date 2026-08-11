import { defineStore } from 'pinia'
import { ref } from 'vue'
import { repoApi } from '@/api'
import type { Repo, CreateRepoRequest, UpdateRepoRequest, TestConnectionResp } from '@/types'
import { message } from 'ant-design-vue'

export interface RepoListParams {
  page?: number
  page_size?: number
  search?: string
  platform?: string
  status?: string
}

export const useRepoStore = defineStore('repo', () => {
  const repos = ref<Repo[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchRepos(params?: RepoListParams) {
    loading.value = true
    try {
      const data = await repoApi.list(params)
      repos.value = data.repos || []
      total.value = data.total || 0
    } catch (e: any) {
      message.error(e.error || '获取仓库列表失败')
    } finally {
      loading.value = false
    }
  }

  async function getRepo(key: string): Promise<Repo | null> {
    try {
      const data = await repoApi.get(key)
      return data.repo
    } catch (e: any) {
      message.error(e.error || '获取仓库详情失败')
      return null
    }
  }

  async function createRepo(req: CreateRepoRequest): Promise<Repo | null> {
    try {
      const data = await repoApi.create(req)
      message.success('创建仓库成功')
      await fetchRepos()
      return data.repo
    } catch (e: any) {
      message.error(e.error || '创建仓库失败')
      return null
    }
  }

  async function updateRepo(req: UpdateRepoRequest): Promise<Repo | null> {
    try {
      const data = await repoApi.update(req)
      message.success('更新仓库成功')
      await fetchRepos()
      return data.repo
    } catch (e: any) {
      message.error(e.error || '更新仓库失败')
      return null
    }
  }

  async function deleteRepo(key: string) {
    try {
      await repoApi.delete(key)
      message.success('删除仓库成功')
      await fetchRepos()
    } catch (e: any) {
      message.error(e.error || '删除仓库失败')
    }
  }

  async function testConnection(key: string): Promise<TestConnectionResp | null> {
    try {
      const data = await repoApi.testConnection(key)
      return data
    } catch (e: any) {
      message.error(e.error || '测试连接失败')
      return null
    }
  }

  async function listBranches(key: string): Promise<string[]> {
    try {
      const data = await repoApi.listBranches(key)
      return data.branches || []
    } catch (e: any) {
      message.error(e.error || '获取分支列表失败')
      return []
    }
  }

  return {
    repos, total, loading,
    fetchRepos, getRepo, createRepo, updateRepo, deleteRepo,
    testConnection, listBranches,
  }
})
