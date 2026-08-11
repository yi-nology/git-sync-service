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

    <a-table
      :columns="columns"
      :data-source="displayData"
      :loading="taskStore.loading"
      :pagination="pagination"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'task'">
          <span class="task-name">{{ record.task_key }}</span>
        </template>
        <template v-if="column.key === 'status'">
          <StatusBadge :status="record.status" />
        </template>
        <template v-if="column.key === 'trigger'">
          <a-tag :color="triggerColor(record.trigger_source)">
            {{ triggerLabel(record.trigger_source) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'duration'">
          <span class="duration-text">{{ calcDuration(record.start_time, record.end_time) }}</span>
        </template>
        <template v-if="column.key === 'time'">
          <div class="time-cell">
            <div>{{ record.start_time }}</div>
            <div v-if="record.end_time" class="time-end">→ {{ record.end_time }}</div>
          </div>
        </template>
        <template v-if="column.key === 'commits'">
          <span class="commit-text">{{ record.commit_range || '-' }}</span>
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
import { ref, computed, onMounted } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import StatusBadge from '@/components/common/StatusBadge.vue'

const taskStore = useSyncTaskStore()
const selectedTask = ref<string>('')

const columns = [
  { title: '任务', key: 'task', dataIndex: 'task_key', width: 150, ellipsis: true },
  { title: '触发方式', key: 'trigger', width: 110, align: 'center' as const },
  { title: '状态', key: 'status', width: 90, align: 'center' as const },
  { title: '耗时', key: 'duration', width: 100 },
  { title: '时间范围', key: 'time', width: 280 },
  { title: '提交范围', key: 'commits', ellipsis: true },
]

const pagination = {
  pageSize: 20,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条记录`,
}

// All history data - when no task selected, show all tasks' history
const allHistory = ref<any[]>([])

const displayData = computed(() => {
  if (selectedTask.value) {
    return taskStore.history
  }
  return allHistory.value
})

function filterTaskOption(input: string, option: any) {
  if (!option.value) return true
  const task = taskStore.tasks.find(t => t.key === option.value)
  return task?.name.toLowerCase().includes(input.toLowerCase()) || false
}

const triggerColor = (trigger: string) => {
  const map: Record<string, string> = {
    manual: 'green',
    cron: 'blue',
    webhook: 'purple',
  }
  return map[trigger] || 'default'
}

const triggerLabel = (trigger: string) => {
  const map: Record<string, string> = {
    manual: '手动',
    cron: '定时',
    webhook: 'Webhook',
  }
  return map[trigger] || trigger
}

function calcDuration(start: string, end: string) {
  if (!start || !end) return '-'
  try {
    const startTime = new Date(start).getTime()
    const endTime = new Date(end).getTime()
    const diff = endTime - startTime
    if (diff < 1000) return `${diff}ms`
    if (diff < 60000) return `${(diff / 1000).toFixed(1)}s`
    return `${Math.floor(diff / 60000)}m ${Math.floor((diff % 60000) / 1000)}s`
  } catch {
    return '-'
  }
}

function handleTaskChange(key: string) {
  if (key) {
    taskStore.fetchHistory(key)
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
    taskStore.fetchHistory(selectedTask.value)
  } else {
    loadAllHistory()
  }
}

onMounted(async () => {
  await taskStore.fetchTasks()
  loadAllHistory()
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

.task-name {
  font-weight: 500;
  color: $text-primary;
}

.duration-text {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 13px;
  color: $text-primary;
}

.time-cell {
  font-size: 13px;

  .time-end {
    color: $text-secondary;
    font-size: 12px;
    margin-top: 2px;
  }
}

.commit-text {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 12px;
  color: $text-secondary;
}
</style>
