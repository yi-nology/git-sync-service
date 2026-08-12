<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">同步执行日志</h1>
        <p class="page-subtitle">查看同步任务的详细执行过程和输出</p>
      </div>
      <a-space>
        <a-button @click="fetchLogs">
          <template #icon><ReloadOutlined /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon" style="background: #E6F7FF; color: #1677FF;">
          <ThunderboltOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.today }}</div>
          <div class="stat-label">今日执行</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #F6FFED; color: #52C41A;">
          <CheckCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.successRate }}</div>
          <div class="stat-label">成功率</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #FFF7E6; color: #FAAD14;">
          <ClockCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.avgDuration }}</div>
          <div class="stat-label">平均耗时</div>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <a-select
        v-model:value="filters.taskKey"
        placeholder="选择任务"
        allow-clear
        show-search
        :filter-option="filterTaskOption"
        style="width: 200px"
        @change="handleFilterChange"
      >
        <a-select-option value="">全部任务</a-select-option>
        <a-select-option v-for="task in tasks" :key="task.key" :value="task.key">
          {{ task.name }}
        </a-select-option>
      </a-select>
      <a-select
        v-model:value="filters.status"
        placeholder="执行状态"
        allow-clear
        style="width: 140px"
        @change="handleFilterChange"
      >
        <a-select-option value="">全部状态</a-select-option>
        <a-select-option value="success">成功</a-select-option>
        <a-select-option value="failed">失败</a-select-option>
        <a-select-option value="running">运行中</a-select-option>
      </a-select>
      <a-range-picker
        v-model:value="filters.dateRange"
        :placeholder="['开始日期', '结束日期']"
        style="width: 260px"
        @change="handleDateChange"
      />
      <a-button type="primary" @click="fetchLogs">
        <template #icon><SearchOutlined /></template>
        查询
      </a-button>
      <a-button @click="resetFilters">
        <template #icon><UndoOutlined /></template>
        重置
      </a-button>
    </div>

    <!-- Log Table -->
    <a-table
      :columns="columns"
      :data-source="logs"
      :loading="loading"
      row-key="id"
      :scroll="{ x: 1000 }"
      :pagination="pagination"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'task_name'">
          <a class="task-link" @click="openDrawer(record)">{{ getTaskName(record.task_key) }}</a>
        </template>

        <template v-if="column.key === 'trigger_source'">
          <a-tag :color="triggerColor(record.trigger_source)">
            {{ triggerLabel(record.trigger_source) }}
          </a-tag>
        </template>

        <template v-if="column.key === 'status'">
          <StatusBadge :status="record.status" />
        </template>

        <template v-if="column.key === 'time'">
          <div class="time-cell">
            <div>{{ formatTime(record.start_time) }}</div>
            <div v-if="record.end_time" class="time-duration">
              耗时 {{ calcDuration(record.start_time, record.end_time) }}
            </div>
          </div>
        </template>

        <template v-if="column.key === 'repo'">
          <span class="repo-text">{{ getRepoDisplay(record.task_key) }}</span>
        </template>

        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openDrawer(record)">
              <template #icon><EyeOutlined /></template>
              查看日志
            </a-button>
            <a-button
              v-if="record.status === 'failed'"
              type="link"
              size="small"
              danger
              @click="handleRetry(record)"
            >
              <template #icon><RedoOutlined /></template>
              重试
            </a-button>
          </a-space>
        </template>
      </template>

      <template #emptyText>
        <a-empty description="暂无同步执行日志" />
      </template>
    </a-table>

    <!-- Log Detail Drawer -->
    <a-drawer
      v-model:open="drawerVisible"
      title="同步执行详情"
      :width="640"
      placement="right"
    >
      <template v-if="currentLog">
        <!-- Basic Info -->
        <a-descriptions :column="2" bordered size="small" class="detail-section">
          <a-descriptions-item label="任务名称" :span="2">
            {{ getTaskName(currentLog.task_key) }}
          </a-descriptions-item>
          <a-descriptions-item label="触发方式">
            <a-tag :color="triggerColor(currentLog.trigger_source)">
              {{ triggerLabel(currentLog.trigger_source) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="执行状态">
            <StatusBadge :status="currentLog.status" />
          </a-descriptions-item>
          <a-descriptions-item label="开始时间">
            {{ formatTime(currentLog.start_time) }}
          </a-descriptions-item>
          <a-descriptions-item label="结束时间">
            {{ formatTime(currentLog.end_time) }}
          </a-descriptions-item>
          <a-descriptions-item label="执行耗时" :span="2">
            {{ calcDuration(currentLog.start_time, currentLog.end_time) }}
          </a-descriptions-item>
        </a-descriptions>

        <!-- Sync Info -->
        <a-divider>同步信息</a-divider>
        <a-descriptions :column="2" bordered size="small" class="detail-section">
          <a-descriptions-item label="源仓库" :span="2">
            {{ getTaskSourceRepo(currentLog.task_key) }}
          </a-descriptions-item>
          <a-descriptions-item label="源分支">
            <span class="branch-tag">{{ getTaskSourceBranch(currentLog.task_key) }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="目标分支">
            <span class="branch-tag">{{ getTaskTargetBranch(currentLog.task_key) }}</span>
          </a-descriptions-item>
          <a-descriptions-item label="提交范围" :span="2">
            <span class="commit-range">{{ currentLog.commit_range || '-' }}</span>
          </a-descriptions-item>
        </a-descriptions>

        <!-- Execution Log -->
        <a-divider>执行日志</a-divider>
        <div class="execution-log">
          <div v-if="parsedDetails.length" class="log-steps">
            <div v-for="(step, index) in parsedDetails" :key="index" class="log-step">
              <div class="step-marker">
                <div class="step-dot" :class="step.type || 'info'"></div>
                <div v-if="index < parsedDetails.length - 1" class="step-line"></div>
              </div>
              <div class="step-content">
                <div class="step-time">{{ step.time || '' }}</div>
                <div class="step-message">{{ step.message }}</div>
              </div>
            </div>
          </div>
          <div v-else-if="currentLog.details" class="log-raw">
            <pre>{{ currentLog.details }}</pre>
          </div>
          <a-empty v-else :image="simpleImage" description="暂无执行日志" />
        </div>

        <!-- Error Info -->
        <template v-if="currentLog.status === 'failed' && currentLog.error_message">
          <a-divider>错误信息</a-divider>
          <a-alert
            type="error"
            :message="currentLog.error_message"
            show-icon
            class="error-alert"
          />
        </template>

        <!-- Actions -->
        <div class="drawer-actions">
          <a-space>
            <a-button
              v-if="currentLog.status === 'failed'"
              type="primary"
              danger
              @click="handleRetry(currentLog)"
            >
              <template #icon><RedoOutlined /></template>
              重试任务
            </a-button>
            <a-button @click="drawerVisible = false">关闭</a-button>
          </a-space>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import { message, Empty } from 'ant-design-vue'
import {
  ReloadOutlined,
  SearchOutlined,
  UndoOutlined,
  EyeOutlined,
  RedoOutlined,
  ThunderboltOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import { logApi, syncTaskApi } from '@/api'
import type { SyncLog, SyncLogRequest } from '@/api'
import type { SyncTask } from '@/types'
import StatusBadge from '@/components/common/StatusBadge.vue'
import dayjs, { type Dayjs } from 'dayjs'

const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE

const loading = ref(false)
const logs = ref<SyncLog[]>([])
const tasks = ref<SyncTask[]>([])
const drawerVisible = ref(false)
const currentLog = ref<SyncLog | null>(null)

const stats = reactive({
  today: 0,
  successRate: '-',
  avgDuration: '-',
})

const filters = reactive({
  taskKey: undefined as string | undefined,
  status: undefined as string | undefined,
  dateRange: null as [Dayjs, Dayjs] | null,
  startDate: '',
  endDate: '',
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条日志`,
})

const columns = [
  { title: '任务名称', key: 'task_name', dataIndex: 'task_key', width: 160, ellipsis: true },
  { title: '触发方式', key: 'trigger_source', width: 110, align: 'center' as const },
  { title: '状态', key: 'status', width: 90, align: 'center' as const },
  { title: '时间', key: 'time', width: 180 },
  { title: '仓库', key: 'repo', width: 200, ellipsis: true },
  { title: '操作', key: 'action', width: 170, fixed: 'right' as const },
]

// Task lookup map for quick access
const taskMap = computed(() => {
  const map: Record<string, SyncTask> = {}
  tasks.value.forEach(t => { map[t.key] = t })
  return map
})

function getTaskName(taskKey: string): string {
  return taskMap.value[taskKey]?.name || taskKey
}

function getTaskSourceRepo(taskKey: string): string {
  return taskMap.value[taskKey]?.source_repo_key || '-'
}

function getTaskSourceBranch(taskKey: string): string {
  return taskMap.value[taskKey]?.source_branch || 'main'
}

function getTaskTargetBranch(taskKey: string): string {
  return taskMap.value[taskKey]?.target_branch || 'main'
}

function getRepoDisplay(taskKey: string): string {
  const task = taskMap.value[taskKey]
  if (!task) return '-'
  return `${task.source_repo_key} -> ${task.target_repo_key}`
}

// Parse execution details into structured steps
const parsedDetails = computed(() => {
  if (!currentLog.value?.details) return []
  try {
    const parsed = JSON.parse(currentLog.value.details)
    if (Array.isArray(parsed)) return parsed
    return []
  } catch {
    // If not JSON, split by newlines into steps
    const lines = currentLog.value.details.split('\n').filter(l => l.trim())
    return lines.map(line => {
      const timeMatch = line.match(/^\[([^\]]+)\]\s*(.*)$/)
      if (timeMatch) {
        return { time: timeMatch[1], message: timeMatch[2], type: 'info' }
      }
      return { message: line, type: 'info' }
    })
  }
})

function triggerColor(trigger: string): string {
  const map: Record<string, string> = {
    manual: 'green',
    cron: 'blue',
    webhook: 'purple',
    retry: 'orange',
  }
  return map[trigger] || 'default'
}

function triggerLabel(trigger: string): string {
  const map: Record<string, string> = {
    manual: '手动',
    cron: '定时',
    webhook: 'Webhook',
    retry: '重试',
  }
  return map[trigger] || trigger
}

function formatTime(time: string): string {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

function calcDuration(start: string, end: string): string {
  if (!start || !end) return '-'
  try {
    const diff = new Date(end).getTime() - new Date(start).getTime()
    if (diff < 1000) return `${diff}ms`
    if (diff < 60000) return `${(diff / 1000).toFixed(1)}s`
    return `${Math.floor(diff / 60000)}m ${Math.floor((diff % 60000) / 1000)}s`
  } catch {
    return '-'
  }
}

function filterTaskOption(input: string, option: any) {
  if (!option.value) return true
  const task = tasks.value.find(t => t.key === option.value)
  return task?.name.toLowerCase().includes(input.toLowerCase()) || false
}

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    filters.startDate = dates[0].format('YYYY-MM-DD')
    filters.endDate = dates[1].format('YYYY-MM-DD')
  } else {
    filters.startDate = ''
    filters.endDate = ''
  }
}

function handleFilterChange() {
  pagination.current = 1
  fetchLogs()
}

function resetFilters() {
  filters.taskKey = undefined
  filters.status = undefined
  filters.dateRange = null
  filters.startDate = ''
  filters.endDate = ''
  pagination.current = 1
  fetchLogs()
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchLogs()
}

function openDrawer(record: SyncLog) {
  currentLog.value = record
  drawerVisible.value = true
}

async function handleRetry(record: SyncLog) {
  try {
    await syncTaskApi.run(record.task_key)
    message.success('重试任务已提交')
    // Refresh the log list after a short delay
    setTimeout(() => fetchLogs(), 1500)
  } catch (e: any) {
    message.error(e.error || '重试失败')
  }
}

function calculateStats(logList: SyncLog[]) {
  const today = dayjs().format('YYYY-MM-DD')
  const todayLogs = logList.filter(l => l.start_time?.startsWith(today))
  stats.today = todayLogs.length

  if (logList.length > 0) {
    const successCount = logList.filter(l => l.status === 'success').length
    stats.successRate = `${Math.round((successCount / logList.length) * 100)}%`

    const durations = logList
      .filter(l => l.start_time && l.end_time)
      .map(l => new Date(l.end_time).getTime() - new Date(l.start_time).getTime())
      .filter(d => d > 0)
    if (durations.length > 0) {
      const avg = durations.reduce((a, b) => a + b, 0) / durations.length
      if (avg < 1000) stats.avgDuration = `${Math.round(avg)}ms`
      else if (avg < 60000) stats.avgDuration = `${(avg / 1000).toFixed(1)}s`
      else stats.avgDuration = `${Math.floor(avg / 60000)}m ${Math.floor((avg % 60000) / 1000)}s`
    } else {
      stats.avgDuration = '-'
    }
  } else {
    stats.successRate = '-'
    stats.avgDuration = '-'
  }
}

async function fetchLogs() {
  loading.value = true
  try {
    // 使用同步历史 API
    const allLogs: any[] = []
    for (const task of tasks.value) {
      try {
        const resp = await syncTaskApi.history({ task_key: task.key, limit: 50 })
        const runs = resp.data?.runs || resp.runs || []
        allLogs.push(...runs.map((r: any) => ({ ...r, task_name: task.name })))
      } catch {
        // ignore individual failures
      }
    }

    // 按时间排序
    allLogs.sort((a, b) => {
      return new Date(b.start_time || 0).getTime() - new Date(a.start_time || 0).getTime()
    })

    // 应用筛选
    let filtered = allLogs
    if (filters.taskKey) {
      filtered = filtered.filter(log => log.task_key === filters.taskKey)
    }
    if (filters.status) {
      filtered = filtered.filter(log => log.status === filters.status)
    }

    logs.value = filtered
    pagination.total = filtered.length
    calculateStats(logs.value)
  } catch (e: any) {
    console.error('Failed to fetch logs:', e)
    logs.value = []
  } finally {
    loading.value = false
  }
}

async function fetchTasks() {
  try {
    const resp = await syncTaskApi.list({ page: 1, page_size: 100 })
    tasks.value = resp.data?.tasks || resp.tasks || []
  } catch {
    tasks.value = []
  }
}

function generateMockLogs(): SyncLog[] {
  const statuses = ['success', 'success', 'success', 'failed', 'success', 'running']
  const triggers = ['manual', 'cron', 'webhook', 'cron', 'manual', 'webhook']
  return Array.from({ length: 8 }, (_, i) => {
    const start = dayjs().subtract(i * 3, 'hour')
    const end = statuses[i % statuses.length] === 'running' ? '' : start.add(30 + i * 10, 'second').format('YYYY-MM-DD HH:mm:ss')
    return {
      id: i + 1,
      task_key: tasks.value[i % Math.max(tasks.value.length, 1)]?.key || 'task-1',
      trigger_source: triggers[i % triggers.length],
      status: statuses[i % statuses.length],
      start_time: start.format('YYYY-MM-DD HH:mm:ss'),
      end_time: end,
      commit_range: `abc${1000 + i}..def${2000 + i}`,
      details: i === 3 ? 'error: push failed: permission denied' : `[${start.format('HH:mm:ss')}] Starting sync\n[${start.add(5, 's').format('HH:mm:ss')}] Fetching source\n[${start.add(15, 's').format('HH:mm:ss')}] Pushing to target\n[${start.add(25, 's').format('HH:mm:ss')}] Done`,
      error_message: i === 3 ? 'push failed: permission denied to target repository' : '',
    }
  })
}

onMounted(async () => {
  await fetchTasks()
  fetchLogs()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
}

.page-header-bar {
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

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.stat-card {
  background: $card-background;
  border-radius: $radius-md;
  padding: $spacing-lg;
  display: flex;
  align-items: center;
  gap: $spacing-md;
  box-shadow: $shadow-card;
  transition: box-shadow 0.2s;

  &:hover {
    box-shadow: $shadow-card-hover;
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: $radius-md;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: $text-primary;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: $text-secondary;
  margin-top: 4px;
}

.filter-bar {
  display: flex;
  gap: $spacing-sm;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: $spacing-md;
  padding: $spacing-md;
  background: $card-background;
  border-radius: $radius-md;
  box-shadow: $shadow-card;
}

.task-link {
  color: $primary-color;
  font-weight: 500;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}

.time-cell {
  font-size: 13px;

  .time-duration {
    color: $text-secondary;
    font-size: 12px;
    margin-top: 2px;
  }
}

.repo-text {
  font-size: 13px;
  color: $text-primary;
}

// Drawer styles
.detail-section {
  margin-bottom: $spacing-md;
}

.execution-log {
  max-height: 400px;
  overflow-y: auto;
}

.log-steps {
  padding: $spacing-sm 0;
}

.log-step {
  display: flex;
  gap: $spacing-md;
  min-height: 40px;
}

.step-marker {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 16px;
  flex-shrink: 0;
}

.step-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 4px;
  flex-shrink: 0;

  &.info {
    background: $primary-color;
  }

  &.success {
    background: $success-color;
  }

  &.error {
    background: $error-color;
  }

  &.warning {
    background: $warning-color;
  }
}

.step-line {
  width: 2px;
  flex: 1;
  background: #E8E8E8;
  margin: 4px 0;
}

.step-content {
  flex: 1;
  padding-bottom: $spacing-md;
}

.step-time {
  font-size: 12px;
  color: $text-secondary;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
}

.step-message {
  font-size: 13px;
  color: $text-primary;
  margin-top: 2px;
  line-height: 1.5;
}

.log-raw {
  pre {
    background: #FAFAFA;
    border: 1px solid $border-color;
    border-radius: $radius-sm;
    padding: $spacing-md;
    font-size: 12px;
    font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
    color: $text-primary;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 300px;
    overflow-y: auto;
  }
}

.error-alert {
  margin-bottom: $spacing-md;
}

.drawer-actions {
  margin-top: $spacing-lg;
  padding-top: $spacing-md;
  border-top: 1px solid $border-color;
  display: flex;
  justify-content: flex-end;
}

.commit-range {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 12px;
  color: $text-secondary;
}

:deep(.ant-descriptions-item-label) {
  font-weight: 500;
  color: $text-secondary;
  width: 100px;
}
</style>
