<template>
  <div class="page-container">
    <div v-if="loading" class="loading-state">加载中...</div>
    <template v-else-if="repo">
      <div class="config-header">
        <div class="config-header-left">
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
        <div class="header-actions-lg">
          <button class="btn-default" @click="testConn">测试连接</button>
          <button class="btn-primary" @click="runSync">立即同步</button>
        </div>
      </div>

      <div class="config-sections">
        <div class="config-card">
          <div class="card-title">仓库信息</div>
          <div class="info-grid">
            <div class="info-item"><span class="info-label">平台</span><span class="info-value">{{ repo.platform }}</span></div>
            <div class="info-item"><span class="info-label">所有者</span><span class="info-value">{{ repo.platform_owner }}</span></div>
            <div class="info-item"><span class="info-label">仓库名</span><span class="info-value">{{ repo.platform_repo }}</span></div>
            <div class="info-item"><span class="info-label">默认分支</span><span class="info-value">{{ repo.default_branch }}</span></div>
            <div class="info-item"><span class="info-label">SSH URL</span><span class="info-value">{{ repo.ssh_url || '-' }}</span></div>
            <div class="info-item"><span class="info-label">创建时间</span><span class="info-value">{{ repo.created_at }}</span></div>
          </div>
        </div>

        <div class="config-card">
          <div class="card-title">关联同步任务</div>
          <div v-if="tasks.length === 0" class="empty-text">暂无关联任务</div>
          <div v-else class="task-list">
            <div class="task-item" v-for="task in tasks" :key="task.key">
              <div class="task-info">
                <div class="task-name">{{ task.name }}</div>
                <div class="task-meta">{{ task.source_branch }} → {{ task.target_branch }}</div>
              </div>
              <span class="status-badge" :class="task.last_status || 'idle'">{{ statusText(task.last_status) }}</span>
            </div>
          </div>
        </div>

        <div class="config-card">
          <div class="card-title">Webhook 配置</div>
          <div class="webhook-info">
            <div class="webhook-url-wrap">
              <label>Webhook 地址</label>
              <div class="url-input-group">
                <input type="text" readonly :value="webhookUrl" class="url-input"/>
                <button class="copy-btn" @click="copyUrl">复制</button>
              </div>
            </div>
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
import { message } from 'ant-design-vue'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import { statusText, copyToClipboard } from '@/utils'
import type { Repo, SyncTask } from '@/types'

const route = useRoute()
const router = useRouter()
const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()

const loading = ref(true)
const repo = ref<Repo | null>(null)
const tasks = ref<SyncTask[]>([])

const repoKey = computed(() => route.params.id as string)
const webhookUrl = computed(() => `${window.location.origin}/api/v1/webhook/receive/${repoKey.value}`)

onMounted(async () => {
  try {
    repo.value = await repoStore.getRepo(repoKey.value)
    if (repo.value) {
      await taskStore.fetchTasks({ repo_key: repoKey.value })
      tasks.value = taskStore.tasks
    }
  } finally {
    loading.value = false
  }
})

async function testConn() {
  const result = await repoStore.testConnection(repoKey.value)
  if (result) {
    message[result.success ? 'success' : 'error'](result.message)
  }
}

function runSync() {
  if (tasks.value.length > 0) {
    taskStore.runTask(tasks.value[0].key)
  } else {
    message.warning('暂无关联任务')
  }
}

function copyUrl() {
  copyToClipboard(webhookUrl.value)
  message.success('已复制到剪贴板')
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.config-header { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px 24px; display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.config-header-left { display: flex; align-items: center; gap: 12px; }
.repo-icon-lg { width: 56px; height: 56px; border-radius: 8px; background: #e6f7ff; color: #1890ff; display: flex; align-items: center; justify-content: center; }
.repo-info-lg .repo-name-lg { font-size: 18px; font-weight: 600; color: #262626; }
.repo-info-lg .repo-url-lg { font-size: 13px; color: #8c8c8c; margin-top: 4px; }
.badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; background: #f5f5f5; color: #8c8c8c; margin-left: 8px; &.active { background: #f6ffed; color: #52c41a; } }
.header-actions-lg { display: flex; gap: 12px; }
.btn-default { padding: 7px 14px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
.btn-primary { padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; cursor: pointer; &:hover { background: #40a9ff; } }
.config-sections { display: flex; flex-direction: column; gap: 16px; }
.config-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px 24px; }
.card-title { font-size: 15px; font-weight: 600; color: #262626; margin-bottom: 16px; padding-bottom: 12px; border-bottom: 1px solid #f0f0f0; }
.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.info-item { display: flex; gap: 8px; }
.info-label { font-size: 13px; color: #8c8c8c; width: 80px; flex-shrink: 0; }
.info-value { font-size: 13px; color: #262626; font-weight: 500; word-break: break-all; }
.empty-text { font-size: 13px; color: #8c8c8c; text-align: center; padding: 24px; }
.task-list { display: flex; flex-direction: column; gap: 12px; }
.task-item { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; background: #fafafa; border-radius: 6px; }
.task-info .task-name { font-size: 14px; font-weight: 500; color: #262626; }
.task-info .task-meta { font-size: 12px; color: #8c8c8c; margin-top: 4px; }
.status-badge { padding: 4px 10px; border-radius: 4px; font-size: 12px; &.success { background: #f6ffed; color: #52c41a; } &.running { background: #e6f7ff; color: #1890ff; } &.failed { background: #fff2f0; color: #ff4d4f; } &.idle { background: #f5f5f5; color: #8c8c8c; } }
.webhook-info label { display: block; font-size: 13px; font-weight: 500; color: #595959; margin-bottom: 8px; }
.url-input-group { display: flex; gap: 8px; }
.url-input { flex: 1; height: 36px; padding: 0 12px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 13px; color: #8c8c8c; background: #fafafa; }
.copy-btn { padding: 0 16px; height: 36px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
</style>
