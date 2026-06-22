import { defineStore } from 'pinia'
import { ref } from 'vue'
import { repoApi } from '@/api'
import type { Repo, CreateRepoRequest, UpdateRepoRequest, TestConnectionResp, Pagination } from '@/types'
import { ElMessage } from 'element-plus'

export const useRepoStore = defineStore('repo', () => {
  const repos = ref<Repo[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchRepos(params?: Pagination) {
    loading.value = true
    try {
      const data = await repoApi.list(params)
      repos.value = data.repos || []
      total.value = data.total || 0
    } catch (e: any) {
      ElMessage.error(e.error || '获取仓库列表失败')
    } finally {
      loading.value = false
    }
  }

  async function getRepo(key: string): Promise<Repo | null> {
    try {
      const data = await repoApi.get(key)
      return data.repo
    } catch (e: any) {
      ElMessage.error(e.error || '获取仓库详情失败')
      return null
    }
  }

  async function createRepo(req: CreateRepoRequest): Promise<Repo | null> {
    try {
      const data = await repoApi.create(req)
      ElMessage.success('创建仓库成功')
      await fetchRepos()
      return data.repo
    } catch (e: any) {
      ElMessage.error(e.error || '创建仓库失败')
      return null
    }
  }

  async function updateRepo(req: UpdateRepoRequest): Promise<Repo | null> {
    try {
      const data = await repoApi.update(req)
      ElMessage.success('更新仓库成功')
      await fetchRepos()
      return data.repo
    } catch (e: any) {
      ElMessage.error(e.error || '更新仓库失败')
      return null
    }
  }

  async function deleteRepo(key: string) {
    try {
      await repoApi.delete(key)
      ElMessage.success('删除仓库成功')
      await fetchRepos()
    } catch (e: any) {
      ElMessage.error(e.error || '删除仓库失败')
    }
  }

  async function testConnection(key: string): Promise<TestConnectionResp | null> {
    try {
      const data = await repoApi.testConnection(key)
      return data
    } catch (e: any) {
      ElMessage.error(e.error || '测试连接失败')
      return null
    }
  }

  async function listBranches(key: string): Promise<string[]> {
    try {
      const data = await repoApi.listBranches(key)
      return data.branches || []
    } catch (e: any) {
      ElMessage.error(e.error || '获取分支列表失败')
      return []
    }
  }

  return {
    repos, total, loading,
    fetchRepos, getRepo, createRepo, updateRepo, deleteRepo,
    testConnection, listBranches,
  }
})
