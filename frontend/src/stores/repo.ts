import { defineStore } from 'pinia'
import { ref } from 'vue'
import { repoApi } from '@/api'
import type { Repo, CreateRepoRequest, UpdateRepoRequest } from '@/types'
import type { TestConnectionData } from '@/types/api'

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
      // 无参调用(任务向导、下拉选择等)需要全量列表;
      // 后端默认 page_size=10,不传参会只拿到最新 10 条,导致下拉选不到大部分仓库
      const query = { page: 1, page_size: 500, ...params }
      const data = await repoApi.list(query)
      repos.value = data.list
      total.value = data.pagination?.total ?? 0
    } finally {
      loading.value = false
    }
  }

  async function getRepo(key: string): Promise<Repo> {
    const data = await repoApi.get(key)
    return data.repo
  }

  async function createRepo(req: CreateRepoRequest): Promise<Repo> {
    const data = await repoApi.create(req)
    await fetchRepos()
    return data.repo
  }

  async function updateRepo(req: UpdateRepoRequest): Promise<Repo> {
    const data = await repoApi.update(req)
    await fetchRepos()
    return data.repo
  }

  async function deleteRepo(key: string) {
    await repoApi.delete(key)
    await fetchRepos()
  }

  /** 批量删除(后端 /repos/batch 支持) */
  async function batchDelete(keys: string[]): Promise<{ succeeded: number; failed: number }> {
    const data = await repoApi.batchDelete(keys)
    await fetchRepos()
    return { succeeded: data.success, failed: data.failed }
  }

  async function testConnection(key: string): Promise<TestConnectionData> {
    return repoApi.testConnection(key)
  }

  async function listBranches(key: string): Promise<string[]> {
    const data = await repoApi.listBranches(key)
    return data.branches || []
  }

  return {
    repos, total, loading,
    fetchRepos, getRepo, createRepo, updateRepo, deleteRepo, batchDelete,
    testConnection, listBranches,
  }
})
