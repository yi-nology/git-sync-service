import { defineStore } from 'pinia'
import { ref } from 'vue'
import { syncTaskApi } from '@/api'
import type { SyncTask, SyncRun, CreateTaskRequest, UpdateTaskRequest, Pagination } from '@/types'
import { ElMessage } from 'element-plus'

export const useSyncTaskStore = defineStore('syncTask', () => {
  const tasks = ref<SyncTask[]>([])
  const total = ref(0)
  const loading = ref(false)
  const history = ref<SyncRun[]>([])

  async function fetchTasks(params?: { repo_key?: string } & Pagination) {
    loading.value = true
    try {
      const data = await syncTaskApi.list(params)
      tasks.value = data.tasks || []
      total.value = data.total || 0
    } catch (e: any) {
      ElMessage.error(e.error || '获取任务列表失败')
    } finally {
      loading.value = false
    }
  }

  async function getTask(key: string): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.get(key)
      return data.task
    } catch (e: any) {
      ElMessage.error(e.error || '获取任务详情失败')
      return null
    }
  }

  async function createTask(req: CreateTaskRequest): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.create(req)
      ElMessage.success('创建任务成功')
      await fetchTasks()
      return data.task
    } catch (e: any) {
      ElMessage.error(e.error || '创建任务失败')
      return null
    }
  }

  async function updateTask(req: UpdateTaskRequest): Promise<SyncTask | null> {
    try {
      const data = await syncTaskApi.update(req)
      ElMessage.success('更新任务成功')
      await fetchTasks()
      return data.task
    } catch (e: any) {
      ElMessage.error(e.error || '更新任务失败')
      return null
    }
  }

  async function deleteTask(key: string) {
    try {
      await syncTaskApi.delete(key)
      ElMessage.success('删除任务成功')
      await fetchTasks()
    } catch (e: any) {
      ElMessage.error(e.error || '删除任务失败')
    }
  }

  async function runTask(key: string) {
    try {
      await syncTaskApi.run(key)
      ElMessage.success('任务已启动')
    } catch (e: any) {
      ElMessage.error(e.error || '启动任务失败')
    }
  }

  async function fetchHistory(taskKey: string, limit = 50) {
    try {
      const data = await syncTaskApi.history({ task_key: taskKey, limit })
      history.value = data.runs || []
    } catch (e: any) {
      ElMessage.error(e.error || '获取历史记录失败')
    }
  }

  return {
    tasks, total, loading, history,
    fetchTasks, getTask, createTask, updateTask, deleteTask, runTask, fetchHistory,
  }
})
