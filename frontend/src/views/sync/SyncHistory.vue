<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">同步历史</h1>
        <p class="page-subtitle">查看所有同步任务的执行记录</p>
      </div>
      <a-space>
        <a-select
          v-model:value="selectedTask"
          placeholder="筛选任务"
          style="width: 240px"
          allow-clear
          show-search
          :filter-option="filterTaskOption"
          @change="handleTaskChange"
        >
          <a-select-option value="">全部任务</a-select-option>
          <a-select-option v-for="task in taskStore.tasks" :key="task.key" :value="task.key">
            {{ task.name }}
          </a-select-option>
        </a-select>
        <a-button @click="refresh">
          <template #icon><ReloadOutlined /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- 统计概览:总执行次数 / 成功(含成功率) / 失败 / 平均耗时 -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon blue"><UnorderedListOutlined /></div>
        <div class="stat-content">
          <div class="stat-num">{{ stats.total }}</div>
          <div class="stat-name">总执行次数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green"><CheckCircleOutlined /></div>
        <div class="stat-content">
          <div class="stat-num">{{ stats.success }}</div>
          <div class="stat-name">成功 ({{ stats.successRate }})</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon red"><CloseCircleOutlined /></div>
        <div class="stat-content">
          <div class="stat-num">{{ stats.failed }}</div>
          <div class="stat-name">失败</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange"><ClockCircleOutlined /></div>
        <div class="stat-content">
          <div class="stat-num">{{ stats.avgDuration }}</div>
          <div class="stat-name">平均耗时</div>
        </div>
      </div>
    </div>

    <!-- 筛选:任务名/分支搜索 + 状态 -->
    <div class="filter-bar">
      <a-input
        v-model:value="filters.search"
        placeholder="搜索任务名称、分支..."
        allow-clear
        class="filter-input"
      >
        <template #prefix><SearchOutlined style="color: #8c8c8c" /></template>
      </a-input>
      <a-select
        v-model:value="filters.status"
        class="filter-select"
        placeholder="所有状态"
        allow-clear
      >
        <a-select-option value="success">成功</a-select-option>
        <a-select-option value="failed">失败</a-select-option>
        <a-select-option value="running">运行中</a-select-option>
      </a-select>
      <span class="filter-bar-spacer" />
    </div>

    <a-table
      :columns="columns"
      :data-source="filteredData"
      :loading="taskStore.loading"
      :pagination="pagination"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'time'">
          <div class="time-cell">
            <div>{{ record.start_time }}</div>
            <div v-if="record.end_time" class="time-end">→ {{ record.end_time }}</div>
          </div>
        </template>
        <template v-if="column.key === 'task'">
          <span class="task-name">{{ taskNameOf(record.task_key) }}</span>
        </template>
        <template v-if="column.key === 'flow'">
          <div class="flow-cell">
            <span class="flow-repo">{{ repoShort(sourceRepoOf(record.task_key)) }}</span>
            <span class="flow-branch">{{ sourceBranchOf(record.task_key) }}</span>
            <ArrowRightOutlined class="flow-arrow" />
            <span class="flow-repo">{{ repoShort(targetRepoOf(record.task_key)) }}</span>
            <span class="flow-branch">{{ targetBranchOf(record.task_key) }}</span>
          </div>
        </template>
        <template v-if="column.key === 'trigger'">
          <a-tag :color="triggerColor(record.trigger_source)">
            {{ triggerLabel(record.trigger_source) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <StatusBadge :status="record.status" />
        </template>
        <template v-if="column.key === 'duration'">
          <span class="duration-text">{{ formatDuration(record.duration_ms) }}</span>
        </template>
        <template v-if="column.key === 'actions'">
          <a-tooltip title="重新运行">
            <a-button
              type="text"
              size="small"
              :loading="runningKeys[record.task_key]"
              @click="rerun(record)"
            >
              <template #icon><PlayCircleOutlined /></template>
            </a-button>
          </a-tooltip>
        </template>
      </template>
      <template #emptyText>
        <a-empty description="暂无同步历史">
          <template #description>
            <span>{{ selectedTask ? '该任务暂无执行记录' : '暂无同步执行记录' }}</span>
          </template>
        </a-empty>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import {
  ReloadOutlined,
  SearchOutlined,
  PlayCircleOutlined,
  ArrowRightOutlined,
  UnorderedListOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import { syncTaskApi } from '@/api'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { triggerColor, triggerLabel } from '@/utils/dictionaries'
import { notifyError, notifySuccess } from '@/utils/notify'

const taskStore = useSyncTaskStore()
const selectedTask = ref<string>('')

const columns = [
  { title: '时间', key: 'time', width: 210 },
  { title: '任务名称', key: 'task', width: 200, ellipsis: true },
  { title: '源 → 目标', key: 'flow', minWidth: 260 },
  { title: '触发方式', key: 'trigger', width: 100, align: 'center' as const },
  { title: '状态', key: 'status', width: 90, align: 'center' as const },
  { title: '耗时', key: 'duration', width: 90 },
  { title: '操作', key: 'actions', width: 80, align: 'center' as const },
]

const pagination = {
  pageSize: 20,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条记录`,
}

// All history data - when no task selected, show all tasks' history
const allHistory = ref<any[]>([])
const filters = reactive({ search: '', status: undefined as string | undefined })
const runningKeys = ref<Record<string, boolean>>({})

const displayData = computed(() => {
  if (selectedTask.value) {
    return taskStore.history
  }
  return allHistory.value
})

const taskMap = computed(() => {
  const map: Record<string, any> = {}
  for (const t of taskStore.tasks) map[t.key] = t
  return map
})

const filteredData = computed(() => {
  let rows = displayData.value
  if (filters.status) {
    rows = rows.filter(r => r.status === filters.status)
  }
  if (filters.search) {
    const q = filters.search.toLowerCase()
    rows = rows.filter(r => {
      const task = taskMap.value[r.task_key]
      const hay = [
        task?.name || '',
        task?.source_repo_key || '',
        task?.target_repo_key || '',
        task?.source_branch || '',
        task?.target_branch || '',
      ].join(' ').toLowerCase()
      return hay.includes(q)
    })
  }
  return rows
})

const stats = computed(() => {
  const rows = displayData.value
  const total = rows.length
  const success = rows.filter(r => r.status === 'success').length
  const failed = rows.filter(r => r.status === 'failed').length
  const durations = rows.map(r => r.duration_ms || 0).filter(d => d > 0)
  const avg = durations.length ? durations.reduce((a, b) => a + b, 0) / durations.length : 0
  return {
    total,
    success,
    failed,
    successRate: total ? `${Math.round((success / total) * 100)}%` : '-',
    avgDuration: avg ? formatDuration(Math.round(avg)) : '-',
  }
})

function taskNameOf(key: string) {
  return taskMap.value[key]?.name || key
}

function sourceRepoOf(key: string) {
  return taskMap.value[key]?.source_repo_key || '-'
}

function targetRepoOf(key: string) {
  return taskMap.value[key]?.target_repo_key || '-'
}

function sourceBranchOf(key: string) {
  return taskMap.value[key]?.source_branch || ''
}

function targetBranchOf(key: string) {
  return taskMap.value[key]?.target_branch || ''
}

// 仓库全名较长,显示 owner/repo 的 repo 短名,悬浮有完整名
function repoShort(full: string) {
  if (!full || full === '-') return '-'
  const parts = full.split('/')
  return parts[parts.length - 1] || full
}

function formatDuration(ms?: number) {
  if (!ms || ms <= 0) return '-'
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${Math.floor(ms / 60000)}m ${Math.floor((ms % 60000) / 1000)}s`
}

function filterTaskOption(input: string, option: any) {
  if (!option.value) return true
  const task = taskStore.tasks.find(t => t.key === option.value)
  return task?.name.toLowerCase().includes(input.toLowerCase()) || false
}

function handleTaskChange(key: string) {
  if (key) {
    taskStore.fetchHistory(key).catch((e) => notifyError(e, '加载历史失败'))
  } else {
    loadAllHistory()
  }
}

async function loadAllHistory() {
  // Fetch all tasks' history
  allHistory.value = []
  for (const task of taskStore.tasks) {
    try {
      await taskStore.fetchHistory(task.key, 10)
      allHistory.value.push(...taskStore.history.map(h => ({ ...h })))
    } catch {
      // ignore individual failures
    }
  }
  // Sort by start_time descending
  allHistory.value.sort((a, b) => {
    return new Date(b.start_time || 0).getTime() - new Date(a.start_time || 0).getTime()
  })
}

function refresh() {
  if (selectedTask.value) {
    taskStore.fetchHistory(selectedTask.value).catch((e) => notifyError(e, '加载历史失败'))
  } else {
    loadAllHistory()
  }
}

async function rerun(record: any) {
  runningKeys.value[record.task_key] = true
  try {
    await syncTaskApi.run(record.task_key)
    notifySuccess('任务已触发,请稍后刷新查看结果')
  } catch (e) {
    notifyError(e, '触发任务失败')
  } finally {
    runningKeys.value[record.task_key] = false
  }
}

onMounted(async () => {
  try {
    await taskStore.fetchTasks()
  } catch (e) {
    notifyError(e, '加载任务失败')
  }
  loadAllHistory()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.filter-input {
  width: 280px;
}

.filter-select {
  width: 140px;
}

.filter-bar-spacer {
  flex: 1;
}

.task-name {
  font-weight: 500;
}

.time-cell {
  .time-end {
    color: $text-secondary;
    font-size: 12px;
  }
}

.flow-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;

  .flow-repo {
    font-weight: 500;
  }

  .flow-branch {
    color: $text-secondary;
    font-size: 12px;
  }

  .flow-arrow {
    color: $text-secondary;
    font-size: 12px;
  }
}

.duration-text {
  font-variant-numeric: tabular-nums;
}

@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-input,
  .filter-select {
    width: 100%;
  }
}
</style>
