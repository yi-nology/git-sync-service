<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">仪表盘</h1>
    </div>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon blue">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
            <polyline points="13 2 13 9 20 9"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ repoStore.total }}</div>
          <div class="stat-name">仓库总数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ taskStore.total }}</div>
          <div class="stat-name">同步任务</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ runningCount }}</div>
          <div class="stat-name">运行中</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon red">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ failedCount }}</div>
          <div class="stat-name">失败任务</div>
        </div>
      </div>
    </div>

    <div class="grid-row">
      <div class="card">
        <div class="card-title">最近同步任务</div>
        <div v-if="taskStore.tasks.length === 0" class="empty-text">暂无任务</div>
        <div v-else class="recent-list">
          <div class="recent-item" v-for="task in recentTasks" :key="task.key">
            <div class="recent-info">
              <div class="recent-name">{{ task.name }}</div>
              <div class="recent-meta">{{ task.source_branch }} → {{ task.target_branch }}</div>
            </div>
            <span class="status-badge" :class="task.last_status || 'idle'">{{ statusText(task.last_status) }}</span>
          </div>
        </div>
        <div class="card-footer">
          <router-link to="/sync" class="link">查看全部 →</router-link>
        </div>
      </div>

      <div class="card">
        <div class="card-title">仓库列表</div>
        <div v-if="repoStore.repos.length === 0" class="empty-text">暂无仓库</div>
        <div v-else class="recent-list">
          <div class="recent-item" v-for="repo in recentRepos" :key="repo.key">
            <div class="recent-info">
              <div class="recent-name">{{ repo.name }}</div>
              <div class="recent-meta">{{ repo.platform }} | {{ repo.clone_url }}</div>
            </div>
            <span class="badge" :class="repo.status">{{ repo.status === 'active' ? '活跃' : '停用' }}</span>
          </div>
        </div>
        <div class="card-footer">
          <router-link to="/repos" class="link">查看全部 →</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import { statusText } from '@/utils'

const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()

const runningCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'running').length)
const failedCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'failed').length)
const recentTasks = computed(() => taskStore.tasks.slice(0, 5))
const recentRepos = computed(() => repoStore.repos.slice(0, 5))

onMounted(() => {
  repoStore.fetchRepos()
  taskStore.fetchTasks()
})
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 24px; }
.stat-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px; display: flex; align-items: center; gap: 16px; }
.stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; &.blue { background: #e6f7ff; color: #1890ff; } &.green { background: #f6ffed; color: #52c41a; } &.orange { background: #fff7e6; color: #fa8c16; } &.red { background: #fff2f0; color: #ff4d4f; } }
.stat-num { font-size: 28px; font-weight: 700; color: #262626; }
.stat-name { font-size: 13px; color: #8c8c8c; margin-top: 4px; }
.grid-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px; }
.card-title { font-size: 15px; font-weight: 600; color: #262626; margin-bottom: 16px; padding-bottom: 12px; border-bottom: 1px solid #f0f0f0; }
.empty-text { font-size: 13px; color: #8c8c8c; text-align: center; padding: 24px; }
.recent-list { display: flex; flex-direction: column; gap: 12px; }
.recent-item { display: flex; justify-content: space-between; align-items: center; padding: 10px 0; border-bottom: 1px solid #f5f5f5; &:last-child { border-bottom: none; } }
.recent-info .recent-name { font-size: 14px; font-weight: 500; color: #262626; }
.recent-info .recent-meta { font-size: 12px; color: #8c8c8c; margin-top: 4px; }
.status-badge { padding: 4px 10px; border-radius: 4px; font-size: 12px; &.success { background: #f6ffed; color: #52c41a; } &.running { background: #e6f7ff; color: #1890ff; } &.failed { background: #fff2f0; color: #ff4d4f; } &.idle { background: #f5f5f5; color: #8c8c8c; } }
.badge { padding: 4px 10px; border-radius: 4px; font-size: 12px; background: #f5f5f5; color: #8c8c8c; &.active { background: #f6ffed; color: #52c41a; } }
.card-footer { margin-top: 12px; padding-top: 12px; border-top: 1px solid #f0f0f0; text-align: right; }
.link { font-size: 13px; color: #1890ff; text-decoration: none; &:hover { text-decoration: underline; } }
</style>
