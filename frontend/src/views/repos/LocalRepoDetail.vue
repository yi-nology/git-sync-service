<template>
  <div class="page-container">
    <div v-if="loading" class="loading-state">加载中...</div>
    <template v-else-if="repo">
      <div class="repo-header">
        <div class="repo-left">
          <div class="repo-icon-lg">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
              <polyline points="13 2 13 9 20 9"/>
            </svg>
          </div>
          <div class="repo-info-lg">
            <div class="repo-name-lg">{{ repo.name }}</div>
            <div class="repo-url-lg">{{ repo.clone_url }}</div>
          </div>
          <span class="badge" :class="repo.status">{{ repo.status === 'active' ? '活跃' : '停用' }}</span>
        </div>
        <div class="header-actions">
          <button class="btn-default" @click="refresh">刷新</button>
          <button class="btn-primary" @click="syncNow">立即同步</button>
        </div>
      </div>

      <div class="tabs-bar">
        <button class="tab-btn" :class="{active: activeTab === 'tasks'}" @click="changeTab('tasks')">同步任务</button>
        <button class="tab-btn" :class="{active: activeTab === 'history'}" @click="changeTab('history')">同步历史</button>
        <button class="tab-btn" :class="{active: activeTab === 'webhook'}" @click="changeTab('webhook')">Webhook</button>
      </div>

      <div v-if="activeTab === 'tasks'" class="config-card">
        <div class="card-title">同步任务列表</div>
        <div v-if="tasks.length === 0" class="empty-text">暂无同步任务</div>
        <div v-else class="task-list">
          <div class="task-item" v-for="task in tasks" :key="task.key">
            <div class="task-info">
              <div class="task-name">{{ task.name }}</div>
              <div class="task-meta">
                <span>{{ task.source_branch }}</span>
                <svg class="arrow-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9 18 15 12 9 6"/>
                </svg>
                <span>{{ task.target_branch }}</span>
              </div>
            </div>
            <div class="task-tags">
              <span class="tag" :class="task.last_status || 'idle'">{{ statusText(task.last_status) }}</span>
            </div>
            <div class="task-actions">
              <button class="action-btn run" @click="runTask(task.key)">运行</button>
              <button class="action-btn edit" @click="editTask(task.key)">编辑</button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'history'" class="config-card">
        <div class="card-title">同步历史</div>
        <div v-if="history.length === 0" class="empty-text">暂无历史记录</div>
        <div v-else class="history-list">
          <div class="history-item" v-for="run in history" :key="run.id">
            <div class="history-icon" :class="run.status">
              <svg v-if="run.status === 'success'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            </div>
            <div class="history-info">
              <div class="history-task">{{ run.task_key }}</div>
              <div class="history-meta">{{ run.trigger_source }} | {{ run.start_time }}</div>
            </div>
            <span class="status-badge" :class="run.status">{{ statusText(run.status) }}</span>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'webhook'" class="config-card">
        <div class="card-title">Webhook 配置</div>
        <div class="webhook-info">
          <label>Webhook 地址</label>
          <div class="url-input-group">
            <input type="text" readonly :value="webhookUrl" class="url-input"/>
            <button class="copy-btn" @click="copyUrl">复制</button>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="empty-state">仓库不存在</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import { statusText, copyToClipboard } from '@/utils'
import type { Repo, SyncTask, SyncRun } from '@/types'

const route = useRoute()
const router = useRouter()
const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()

const loading = ref(true)
const repo = ref<Repo | null>(null)
const tasks = ref<SyncTask[]>([])
const history = ref<SyncRun[]>([])
const activeTab = ref('tasks')

const repoKey = computed(() => route.params.id as string)
const webhookUrl = computed(() => `${window.location.origin}/api/v1/webhook/receive/${repoKey.value}`)

onMounted(async () => {
  try {
    repo.value = await repoStore.getRepo(repoKey.value)
    if (repo.value) {
      await taskStore.fetchTasks({ repo_key: repoKey.value })
      tasks.value = taskStore.tasks
    }
    if (route.query.tab) {
      activeTab.value = route.query.tab as string
    }
  } finally {
    loading.value = false
  }
})

function changeTab(tab: string) {
  activeTab.value = tab
  router.replace({ query: { tab } })
}

function refresh() {
  ElMessage.success('刷新成功')
}

function syncNow() {
  if (tasks.value.length > 0) {
    taskStore.runTask(tasks.value[0].key)
  } else {
    ElMessage.warning('暂无同步任务')
  }
}

function runTask(key: string) {
  taskStore.runTask(key)
}

function editTask(key: string) {
  router.push(`/sync?edit=${key}`)
}

function copyUrl() {
  copyToClipboard(webhookUrl.value)
  ElMessage.success('已复制到剪贴板')
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.repo-header { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px 24px; display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.repo-left { display: flex; align-items: center; gap: 12px; }
.repo-icon-lg { width: 56px; height: 56px; border-radius: 8px; background: #e6f7ff; color: #1890ff; display: flex; align-items: center; justify-content: center; }
.repo-info-lg .repo-name-lg { font-size: 18px; font-weight: 600; color: #262626; }
.repo-info-lg .repo-url-lg { font-size: 13px; color: #8c8c8c; margin-top: 4px; }
.badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; background: #f5f5f5; color: #8c8c8c; margin-left: 8px; &.active { background: #f6ffed; color: #52c41a; } }
.header-actions { display: flex; gap: 12px; }
.btn-default { padding: 7px 14px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
.btn-primary { padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; cursor: pointer; &:hover { background: #40a9ff; } }
.tabs-bar { display: flex; gap: 4px; margin-bottom: 16px; background: #fff; border-radius: 8px; padding: 4px; border: 1px solid #f0f0f0; }
.tab-btn { padding: 8px 20px; border: none; background: transparent; border-radius: 6px; font-size: 14px; color: #595959; cursor: pointer; transition: all 0.2s; &:hover { color: #1890ff; } &.active { background: #1890ff; color: #fff; } }
.config-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px 24px; margin-bottom: 16px; }
.card-title { font-size: 15px; font-weight: 600; color: #262626; margin-bottom: 16px; padding-bottom: 12px; border-bottom: 1px solid #f0f0f0; }
.empty-text { font-size: 13px; color: #8c8c8c; text-align: center; padding: 24px; }
.task-list { display: flex; flex-direction: column; gap: 12px; }
.task-item { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; background: #fafafa; border-radius: 6px; }
.task-info .task-name { font-size: 14px; font-weight: 500; color: #262626; }
.task-info .task-meta { font-size: 12px; color: #8c8c8c; margin-top: 4px; display: flex; gap: 8px; align-items: center; }
.arrow-icon { color: #bfbfbf; }
.task-tags { display: flex; gap: 8px; margin-right: 16px; }
.tag { padding: 4px 10px; border-radius: 4px; font-size: 12px; &.success { background: #f6ffed; color: #52c41a; } &.running { background: #e6f7ff; color: #1890ff; } &.failed { background: #fff2f0; color: #ff4d4f; } &.idle { background: #f5f5f5; color: #8c8c8c; } }
.task-actions { display: flex; gap: 8px; }
.action-btn { padding: 4px 8px; border-radius: 4px; border: none; cursor: pointer; font-size: 12px; transition: all 0.2s; &.run { background: #e6f7ff; color: #1890ff; &:hover { background: #bae7ff; } } &.edit { background: #f6ffed; color: #52c41a; &:hover { background: #d9f7be; } } }
.history-list { display: flex; flex-direction: column; gap: 12px; }
.history-item { display: flex; align-items: center; gap: 12px; padding: 12px 16px; background: #fafafa; border-radius: 6px; }
.history-icon { width: 32px; height: 32px; border-radius: 50%; display: flex; align-items: center; justify-content: center; &.success { background: #f6ffed; color: #52c41a; } &.failed { background: #fff2f0; color: #ff4d4f; } &.running { background: #e6f7ff; color: #1890ff; } }
.history-info { flex: 1; }
.history-task { font-size: 14px; font-weight: 500; color: #262626; }
.history-meta { font-size: 12px; color: #8c8c8c; margin-top: 4px; }
.status-badge { padding: 4px 10px; border-radius: 4px; font-size: 12px; &.success { background: #f6ffed; color: #52c41a; } &.running { background: #e6f7ff; color: #1890ff; } &.failed { background: #fff2f0; color: #ff4d4f; } }
.webhook-info { label { display: block; font-size: 13px; font-weight: 500; color: #595959; margin-bottom: 8px; } }
.url-input-group { display: flex; gap: 8px; }
.url-input { flex: 1; height: 36px; padding: 0 12px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 13px; color: #8c8c8c; background: #fafafa; }
.copy-btn { padding: 0 16px; height: 36px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
</style>
