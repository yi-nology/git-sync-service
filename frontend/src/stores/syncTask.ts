import { defineStore } from 'pinia'
import { ref } from 'vue'
import { syncTaskApi } from '@/api'
import type { SyncTask, SyncRun, CreateTaskRequest, UpdateTaskRequest, Pagination } from '@/types'
import { message } from 'ant-design-vue'

export interface TaskListParams extends Partial<Pagination> {
  repo_key?: string
  status?: string
  search?: string
}

export const useSyncTaskStore = defineStore('syncTask', () => {
  const tasks = ref<SyncTask[]>([])
  const total = ref(0)
  const loading = ref(false)
  const history = ref<SyncRun[]>([])
  // Keep the last-used params so callers can refresh with the same filters
  const lastParams = ref<TaskListParams>({})

  async function fetchTasks(params?: TaskListParams) {
    loading.value = true
    if (params) lastParams.value = { ...params }
    try {
      const data = await syncTaskApi.list(params as any)
      tasks.value = data.tasks || []
      total.value = data.total || 0
    } catch (e: any) {
      message.error(e.error || '获取任务列表失败')
    } finally {
      loading.value = false
    }
  }

  async function refreshTasks() {
    await fetchTasks(lastParams.value)
  }

  async function getTask(key: string): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.get(key)
      return data.task
    } catch (e: any) {
      message.error(e.error || '获取任务详情失败')
      return null
    }
  }

  async function createTask(req: CreateTaskRequest): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.create(req)
      message.success('创建任务成功')
      await fetchTasks()
      return data.task
    } catch (e: any) {
      message.error(e.error || '创建任务失败')
      return null
    }
  }

  async function updateTask(req: UpdateTaskRequest): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.update(req)
      message.success('更新任务成功')
      await fetchTasks()
      return data.task
    } catch (e: any) {
      message.error(e.error || '更新任务失败')
      return null
    }
  }

  async function deleteTask(key: string) {
    try {
      await syncTaskApi.delete(key)
      message.success('删除任务成功')
      await fetchTasks()
    } catch (e: any) {
      message.error(e.error || '删除任务失败')
    }
  }

  async function runTask(key: string) {
    try {
      await syncTaskApi.run(key)
      message.success('任务已启动')
    } catch (e: any) {
      message.error(e.error || '启动任务失败')
    }
  }

  async function fetchHistory(taskKey: string, limit = 50) {
    try {
      const data = await syncTaskApi.history({ task_key: taskKey, limit })
      history.value = data.runs || []
    } catch (e: any) {
      message.error(e.error || '获取历史记录失败')
    }
  }

  async function batchDelete(keys: string[]) {
    try {
      const results = await Promise.allSettled(keys.map(k => syncTaskApi.delete(k)))
      const succeeded = results.filter(r => r.status === 'fulfilled').length
      const failed = results.length - succeeded
      if (succeeded > 0) message.success(`成功删除 ${succeeded} 个任务`)
      if (failed > 0) message.error(`${failed} 个任务删除失败`)
      await refreshTasks()
    } catch (e: any) {
      message.error('批量删除失败')
    }
  }

  return {
    tasks, total, loading, history, lastParams,
    fetchTasks, refreshTasks, getTask, createTask, updateTask, deleteTask, runTask, fetchHistory, batchDelete,
  }
})
