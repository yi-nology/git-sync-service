<template>
  <div class="page-container">
    <div class="dashboard-header">
      <div>
        <h1 class="page-title">仪表盘</h1>
        <p class="page-subtitle">系统运行概览</p>
      </div>
      <a-space>
        <a-button @click="router.push('/sync/new')">
          <template #icon><PlusOutlined /></template>
          新建同步任务
        </a-button>
        <a-button type="primary" @click="router.push('/repos')">
          <template #icon><FolderAddOutlined /></template>
          添加仓库
        </a-button>
      </a-space>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <div class="stat-card" @click="router.push('/repos')">
        <div class="stat-icon blue">
          <FolderOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ repoStore.total }}</div>
          <div class="stat-name">仓库总数</div>
        </div>
        <div class="stat-arrow"><RightOutlined /></div>
      </div>
      <div class="stat-card" @click="router.push('/sync')">
        <div class="stat-icon green">
          <SyncOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ taskStore.total }}</div>
          <div class="stat-name">同步任务</div>
        </div>
        <div class="stat-arrow"><RightOutlined /></div>
      </div>
      <div class="stat-card" @click="router.push('/sync')">
        <div class="stat-icon orange">
          <PlayCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ runningCount }}</div>
          <div class="stat-name">运行中</div>
        </div>
        <div class="stat-arrow"><RightOutlined /></div>
      </div>
      <div class="stat-card" @click="router.push('/sync')">
        <div class="stat-icon red">
          <CloseCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ failedCount }}</div>
          <div class="stat-name">失败任务</div>
        </div>
        <div class="stat-arrow"><RightOutlined /></div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="quick-actions">
      <div class="section-title">快捷操作</div>
      <div class="action-grid">
        <div class="action-item" @click="router.push('/sync/new')">
          <div class="action-icon blue"><SyncOutlined /></div>
          <span>创建同步任务</span>
        </div>
        <div class="action-item" @click="router.push('/repos')">
          <div class="action-icon green"><FolderAddOutlined /></div>
          <span>添加仓库</span>
        </div>
        <div class="action-item" @click="router.push('/webhook/rules')">
          <div class="action-icon purple"><ApiOutlined /></div>
          <span>配置 Webhook</span>
        </div>
        <div class="action-item" @click="router.push('/sync/history')">
          <div class="action-icon cyan"><HistoryOutlined /></div>
          <span>查看日志</span>
        </div>
      </div>
    </div>

    <!-- Content Grid -->
    <div class="grid-row">
      <div class="content-card">
        <div class="card-header">
          <span class="card-title">
            <SyncOutlined style="margin-right: 8px; color: #1677FF;" />
            最近同步任务
          </span>
          <router-link to="/sync">
            <a-button type="link" size="small">查看全部 <RightOutlined /></a-button>
          </router-link>
        </div>
        <div class="card-body">
          <a-table
            :columns="taskColumns"
            :data-source="recentTasks"
            :pagination="false"
            size="small"
            row-key="key"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'name'">
                <a class="task-link" @click="router.push(`/sync`)">{{ record.name }}</a>
              </template>
              <template v-if="column.dataIndex === 'branch'">
                <span class="branch-tag">{{ record.source_branch }}</span>
                <ArrowRightOutlined style="margin: 0 4px; color: #BFBFBF; font-size: 12px;" />
                <span class="branch-tag">{{ record.target_branch }}</span>
              </template>
              <template v-if="column.dataIndex === 'last_status'">
                <StatusBadge :status="record.last_status" />
              </template>
              <template v-if="column.dataIndex === 'last_run_at'">
                <span class="time-text">{{ record.last_run_at || '未运行' }}</span>
              </template>
            </template>
            <template #emptyText>
              <a-empty description="暂无同步任务" :image-style="{ height: '60px' }">
                <a-button type="primary" size="small" @click="router.push('/sync/new')">创建任务</a-button>
              </a-empty>
            </template>
          </a-table>
        </div>
      </div>

      <div class="content-card">
        <div class="card-header">
          <span class="card-title">
            <FolderOutlined style="margin-right: 8px; color: #52C41A;" />
            仓库列表
          </span>
          <router-link to="/repos">
            <a-button type="link" size="small">查看全部 <RightOutlined /></a-button>
          </router-link>
        </div>
        <div class="card-body">
          <a-table
            :columns="repoColumns"
            :data-source="recentRepos"
            :pagination="false"
            size="small"
            row-key="key"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'name'">
                <a class="task-link" @click="router.push(`/local-repos/${record.key}`)">{{ record.name }}</a>
              </template>
              <template v-if="column.dataIndex === 'platform'">
                <a-tag :color="platformColor(record.platform)">{{ platformLabel(record.platform) }}</a-tag>
              </template>
              <template v-if="column.dataIndex === 'status'">
                <StatusBadge :status="record.status" />
              </template>
            </template>
            <template #emptyText>
              <a-empty description="暂无仓库" :image-style="{ height: '60px' }">
                <a-button type="primary" size="small" @click="router.push('/repos')">添加仓库</a-button>
              </a-empty>
            </template>
          </a-table>
        </div>
      </div>
    </div>

    <!-- System Status -->
    <div class="system-status-section">
      <div class="section-title">
        <DashboardOutlined style="margin-right: 8px;" />
        系统状态
      </div>
      <div class="system-status-card">
        <div class="status-items">
          <div class="status-item">
            <div class="status-item-label">
              <CheckCircleOutlined v-if="systemStatus?.status === 'running'" style="color: $success; margin-right: 6px;" />
              <CloseCircleOutlined v-else style="color: $error; margin-right: 6px;" />
              服务状态
            </div>
            <a-tag :color="systemStatus?.status === 'running' ? 'success' : 'error'">
              {{ systemStatus?.status === 'running' ? '运行中' : '已停止' }}
            </a-tag>
          </div>
          <a-divider type="vertical" style="height: 40px;" />
          <div class="status-item">
            <div class="status-item-label">
              <ClockCircleOutlined style="color: $primary; margin-right: 6px;" />
              最后同步时间
            </div>
            <span class="status-value">{{ systemStatus?.lastSyncAt || '暂无' }}</span>
          </div>
          <a-divider type="vertical" style="height: 40px;" />
          <div class="status-item">
            <div class="status-item-label">
              <InfoCircleOutlined style="color: #722ED1; margin-right: 6px;" />
              版本号
            </div>
            <a-tag color="blue">{{ systemStatus?.version || 'v0.0.0' }}</a-tag>
          </div>
          <a-divider type="vertical" style="height: 40px;" />
          <div class="status-item">
            <div class="status-item-label">
              <CloudServerOutlined style="color: $warning; margin-right: 6px;" />
              运行时长
            </div>
            <span class="status-value">{{ formatUptime(systemStatus?.uptime) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import { systemApi } from '@/api'
import type { SystemStatusResp } from '@/api'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  FolderOutlined,
  SyncOutlined,
  PlayCircleOutlined,
  CloseCircleOutlined,
  PlusOutlined,
  FolderAddOutlined,
  RightOutlined,
  ArrowRightOutlined,
  ApiOutlined,
  HistoryOutlined,
  DashboardOutlined,
  ClockCircleOutlined,
  InfoCircleOutlined,
  CloudServerOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()

const systemStatus = ref<SystemStatusResp | null>(null)

const runningCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'running').length)
const failedCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'failed').length)
const recentTasks = computed(() => taskStore.tasks.slice(0, 5))
const recentRepos = computed(() => repoStore.repos.slice(0, 5))

const formatUptime = (seconds?: number) => {
  if (!seconds) return '--'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  if (days > 0) return `${days}天 ${hours}小时`
  if (hours > 0) return `${hours}小时 ${mins}分钟`
  return `${mins}分钟`
}

const platformColor = (platform: string) => {
  const map: Record<string, string> = {
    github: '#24292E',
    gitlab: '#FC6D26',
    gitea: '#609926',
    gitee: '#C71D23',
    gitcode: '#0066FF',
    atomgit: '#0084FF',
    tencent_code: '#00B4D8',
  }
  return map[platform] || 'blue'
}

const platformLabel = (platform: string) => {
  const map: Record<string, string> = {
    github: 'GitHub',
    gitlab: 'GitLab',
    gitea: 'Gitea',
    gitee: 'Gitee',
    gitcode: 'GitCode',
    atomgit: 'AtomGit',
    tencent_code: '腾讯工蜂',
  }
  return map[platform] || platform
}

const taskColumns = [
  { title: '任务名称', dataIndex: 'name', key: 'name', ellipsis: true },
  { title: '分支', dataIndex: 'branch', key: 'branch', width: 180 },
  { title: '状态', dataIndex: 'last_status', key: 'last_status', width: 90, align: 'center' as const },
  { title: '最后运行', dataIndex: 'last_run_at', key: 'last_run_at', width: 140 },
]

const repoColumns = [
  { title: '仓库名称', dataIndex: 'name', key: 'name', ellipsis: true },
  { title: '平台', dataIndex: 'platform', key: 'platform', width: 100, align: 'center' as const },
  { title: '状态', dataIndex: 'status', key: 'status', width: 90, align: 'center' as const },
]

const fetchSystemStatus = async () => {
  try {
    const resp = await systemApi.status()
    // Handle nested response format
    systemStatus.value = resp.data || resp
  } catch (e) {
    console.error('Failed to fetch system status:', e)
  }
}

onMounted(() => {
  repoStore.fetchRepos()
  taskStore.fetchTasks()
  fetchSystemStatus()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: $spacing-lg;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: $text-primary;
  margin: 0;
  line-height: 1.3;
}

.page-subtitle {
  font-size: 14px;
  color: $text-secondary;
  margin: 4px 0 0 0;
}

// Stats
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.stat-card {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: 20px;
  box-shadow: $shadow-card;
  display: flex;
  align-items: center;
  gap: $spacing-md;
  transition: all 0.2s ease;
  cursor: pointer;
  border: 1px solid transparent;

  &:hover {
    box-shadow: $shadow-card-hover;
    border-color: #E6F7FF;

    .stat-arrow {
      opacity: 1;
      transform: translateX(0);
    }
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: $border-radius-lg;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;

  &.blue { background: #E6F7FF; color: $primary-color; }
  &.green { background: #F6FFED; color: $success-color; }
  &.orange { background: #FFF7E6; color: $warning-color; }
  &.red { background: #FFF2F0; color: $error-color; }
}

.stat-content {
  flex: 1;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: $text-primary;
  line-height: 1.2;
}

.stat-name {
  font-size: 13px;
  color: $text-secondary;
  margin-top: 2px;
}

.stat-arrow {
  color: #BFBFBF;
  font-size: 12px;
  opacity: 0;
  transform: translateX(-4px);
  transition: all 0.2s ease;
}

// Quick Actions
.quick-actions {
  margin-bottom: $spacing-lg;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-md;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
}

.action-item {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: $shadow-card;
  border: 1px solid transparent;

  &:hover {
    border-color: #E6F7FF;
    box-shadow: $shadow-card-hover;
    transform: translateY(-1px);
  }

  span {
    font-size: 14px;
    font-weight: 500;
    color: $text-primary;
  }
}

.action-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;

  &.blue { background: #E6F7FF; color: $primary-color; }
  &.green { background: #F6FFED; color: $success-color; }
  &.purple { background: #F9F0FF; color: #722ED1; }
  &.cyan { background: #E6FFFB; color: #13C2C2; }
}

// Content Grid
.grid-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-md;
}

.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px $spacing-lg;
  border-bottom: 1px solid $border-color;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
  display: flex;
  align-items: center;
}

.card-body {
  padding: $spacing-md $spacing-lg;
}

.task-link {
  color: $text-primary;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.2s;

  &:hover {
    color: $primary-color;
  }
}

.time-text {
  color: $text-secondary;
  font-size: 13px;
}

// System Status
.system-status-section {
  margin-top: $spacing-lg;
}

.system-status-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  padding: 20px $spacing-lg;
}

.status-items {
  display: flex;
  align-items: center;
  gap: $spacing-lg;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 140px;
}

.status-item-label {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: $text-secondary;
}

.status-value {
  font-size: 14px;
  font-weight: 500;
  color: $text-primary;
}

// Responsive
@media (max-width: 1200px) {
  .stats-row,
  .action-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .grid-row {
    grid-template-columns: 1fr;
  }

  .status-items {
    flex-wrap: wrap;
    gap: $spacing-md;
  }
}

@media (max-width: 768px) {
  .status-items {
    flex-direction: column;
    align-items: flex-start;
  }

  .status-items .ant-divider-vertical {
    display: none;
  }
}
</style>
