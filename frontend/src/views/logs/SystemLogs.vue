<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">系统运行日志</h1>
        <p class="page-subtitle">查看系统的运行状态、错误、警告等信息</p>
      </div>
      <a-space>
        <a-select
          v-model:value="autoRefreshInterval"
          style="width: 120px"
          @change="handleAutoRefreshChange"
        >
          <a-select-option :value="0">自动刷新</a-select-option>
          <a-select-option :value="5000">5 秒</a-select-option>
          <a-select-option :value="10000">10 秒</a-select-option>
          <a-select-option :value="30000">30 秒</a-select-option>
        </a-select>
        <a-button @click="fetchLogs">
          <template #icon><ReloadOutlined /></template>
          刷新
        </a-button>
      </a-space>
    </div>

    <!-- Stats Cards -->
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon" style="background: #fff2f0; color: #ff4d4f;">
          <CloseCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-error">{{ stats.errorCount }}</div>
          <div class="stat-label">ERROR</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #fff7e6; color: #faad14;">
          <WarningOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-warn">{{ stats.warnCount }}</div>
          <div class="stat-label">WARN</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #e6f7ff; color: #1677ff;">
          <InfoCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-info">{{ stats.infoCount }}</div>
          <div class="stat-label">INFO</div>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar carded">
      <a-select
        v-model:value="filters.level"
        placeholder="日志级别"
        allow-clear
        style="width: 140px"
        @change="handleFilterChange"
      >
        <a-select-option value="">全部级别</a-select-option>
        <a-select-option value="ERROR">ERROR</a-select-option>
        <a-select-option value="WARN">WARN</a-select-option>
        <a-select-option value="INFO">INFO</a-select-option>
        <a-select-option value="DEBUG">DEBUG</a-select-option>
      </a-select>
      <a-select
        v-model:value="filters.module"
        placeholder="所属模块"
        allow-clear
        style="width: 140px"
        @change="handleFilterChange"
      >
        <a-select-option value="">全部模块</a-select-option>
        <a-select-option value="server">server</a-select-option>
        <a-select-option value="database">database</a-select-option>
        <a-select-option value="sync">sync</a-select-option>
        <a-select-option value="git">git</a-select-option>
        <a-select-option value="scheduler">scheduler</a-select-option>
      </a-select>
      <a-input
        v-model:value="filters.search"
        placeholder="搜索日志内容..."
        allow-clear
        style="width: 220px"
        @pressEnter="fetchLogs"
      >
        <template #prefix><SearchOutlined /></template>
      </a-input>
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
      :data-source="filteredLogs"
      :loading="loading"
      row-key="id"
      :scroll="{ x: 1000 }"
      :pagination="pagination"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'time'">
          <span class="log-time">{{ formatDate(record.time) }}</span>
        </template>

        <template v-if="column.key === 'level'">
          <a-tag :color="levelColor(record.level)">
            {{ record.level }}
          </a-tag>
        </template>

        <template v-if="column.key === 'module'">
          <span class="log-module">{{ getModule(record) }}</span>
        </template>

        <template v-if="column.key === 'message'">
          <span class="log-message">{{ record.message }}</span>
        </template>

        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="openDrawer(record)">
            <template #icon><EyeOutlined /></template>
            详情
          </a-button>
        </template>
      </template>

      <template #emptyText>
        <a-empty description="暂无系统日志" />
      </template>
    </a-table>

    <!-- Log Detail Drawer -->
    <a-drawer
      v-model:open="drawerVisible"
      title="日志详情"
      :width="640"
      placement="right"
    >
      <template v-if="currentLog">
        <a-descriptions :column="2" bordered size="small" class="detail-section">
          <a-descriptions-item label="时间" :span="2">
            {{ formatDate(currentLog.time) }}
          </a-descriptions-item>
          <a-descriptions-item label="级别">
            <a-tag :color="levelColor(currentLog.level)">
              {{ currentLog.level }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="模块">
            {{ getModule(currentLog) }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>日志消息</a-divider>
        <div class="detail-message">
          <pre>{{ currentLog.message }}</pre>
        </div>

        <template v-if="currentLog.details">
          <a-divider>详细信息</a-divider>
          <div class="detail-extra">
            <pre>{{ currentLog.details }}</pre>
          </div>
        </template>

        <template v-if="getStackTrace(currentLog)">
          <a-divider>堆栈跟踪</a-divider>
          <div class="detail-stacktrace">
            <pre>{{ getStackTrace(currentLog) }}</pre>
          </div>
        </template>

        <div class="drawer-actions">
          <a-button @click="drawerVisible = false">关闭</a-button>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, reactive, computed } from 'vue'
import type { TablePaginationConfig } from 'ant-design-vue'
import {
  ReloadOutlined,
  SearchOutlined,
  UndoOutlined,
  EyeOutlined,
  CloseCircleOutlined,
  WarningOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'
import { logApi } from '@/api'
import type { SystemLog, SystemLogParams } from '@/types/api'
import dayjs, { type Dayjs } from 'dayjs'
import { formatDate } from '@/utils'

const loading = ref(false)
const logs = ref<SystemLog[]>([])
const drawerVisible = ref(false)
const currentLog = ref<SystemLog | null>(null)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const autoRefreshInterval = ref(0)

const stats = reactive({ errorCount: 0, warnCount: 0, infoCount: 0 })

const filters = reactive({
  level: undefined as string | undefined,
  module: undefined as string | undefined,
  search: '',
  dateRange: null as [Dayjs, Dayjs] | null,
  startDate: '',
  endDate: '',
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条日志`,
})

const columns = [
  { title: '时间', key: 'time', width: 180 },
  { title: '级别', key: 'level', width: 100, align: 'center' as const },
  { title: '模块', key: 'module', width: 120 },
  { title: '消息', key: 'message', ellipsis: true },
  { title: '操作', key: 'action', width: 100, fixed: 'right' as const },
]

// 缓存 record.id → module 映射,避免模板循环中反复 JSON.parse
const moduleCache = new Map<number, string>()
function getModule(record: SystemLog): string {
  if (moduleCache.has(record.id)) return moduleCache.get(record.id)!
  let result = '-'
  if (record.details) {
    try {
      const parsed = JSON.parse(record.details)
      if (parsed.module) result = parsed.module
    } catch {
      const match = record.message?.match(/^\[([a-zA-Z]+)\]/)
      if (match) result = match[1]
    }
  }
  moduleCache.set(record.id, result)
  return result
}

function getStackTrace(record: SystemLog): string | null {
  if (!record.details) return null
  try {
    const parsed = JSON.parse(record.details)
    if (parsed.stack_trace || parsed.stackTrace) {
      return parsed.stack_trace || parsed.stackTrace
    }
  } catch {
    const stackMatch = record.details.match(/(?:Stack Trace:|goroutine \d+)([\s\S]*)/i)
    if (stackMatch) return stackMatch[1].trim()
  }
  return null
}

const filteredLogs = computed(() => {
  let result = logs.value

  if (filters.module) {
    result = result.filter((log) => getModule(log) === filters.module)
  }

  if (filters.search) {
    const searchLower = filters.search.toLowerCase()
    result = result.filter(
      (log) =>
        log.message?.toLowerCase().includes(searchLower) ||
        log.details?.toLowerCase().includes(searchLower),
    )
  }

  if (filters.startDate && filters.endDate) {
    result = result.filter((log) => {
      if (!log.time) return false
      const logDate = dayjs(log.time).format('YYYY-MM-DD')
      return logDate >= filters.startDate && logDate <= filters.endDate
    })
  }

  return result
})

function levelColor(level: string): string {
  const map: Record<string, string> = {
    ERROR: 'red',
    WARN: 'orange',
    INFO: 'blue',
    DEBUG: 'default',
  }
  return map[level] || 'default'
}

// formatTime 已移至 @/utils 统一管理,使用 formatDate 替代

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
  filters.level = undefined
  filters.module = undefined
  filters.search = ''
  filters.dateRange = null
  filters.startDate = ''
  filters.endDate = ''
  pagination.current = 1
  fetchLogs()
}

function handleTableChange(pag: TablePaginationConfig) {
  pagination.current = pag.current || 1
  pagination.pageSize = pag.pageSize || 20
}

function openDrawer(record: SystemLog) {
  currentLog.value = record
  drawerVisible.value = true
}

function calculateStats(logList: SystemLog[]) {
  let errors = 0, warns = 0, infos = 0
  for (const l of logList) {
    if (l.level === 'ERROR') errors++
    else if (l.level === 'WARN') warns++
    else if (l.level === 'INFO') infos++
  }
  stats.errorCount = errors
  stats.warnCount = warns
  stats.infoCount = infos
}

function handleAutoRefreshChange(value: number) {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
  if (value > 0) {
    refreshTimer = setInterval(() => {
      fetchLogs()
    }, value)
  }
}

async function fetchLogs() {
  loading.value = true
  try {
    const params: SystemLogParams = {
      page: pagination.current,
      page_size: pagination.pageSize,
    }
    if (filters.level) params.level = filters.level

    const data = await logApi.listSystem(params)
    logs.value = data.list || []
    pagination.total = data.pagination?.total || 0
    calculateStats(logs.value)
  } catch (e) {
    // 接口不可用时如实清空,不展示假数据
    console.error('Failed to fetch system logs:', e)
    logs.value = []
    pagination.total = 0
    calculateStats(logs.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLogs()
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.stats-row {
  grid-template-columns: repeat(3, 1fr);
}

.filter-bar.carded {
  padding: $spacing-md;
  background: $bg-primary;
  border-radius: $radius-md;
  box-shadow: $shadow-card;
}

.stat-value {
  &.stat-error { color: $error; }
  &.stat-warn  { color: $warning; }
  &.stat-info  { color: $primary; }
}

.log-time {
  color: $text-secondary;
  font-size: 13px;
}

.log-module {
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 13px;
  color: $text-primary;
}

.log-message {
  color: $text-primary;
  font-size: 13px;
}

// 抽屉详情
.detail-section {
  margin-bottom: $spacing-md;
}

.detail-message,
.detail-extra,
.detail-stacktrace {
  pre {
    background: #fafafa;
    border: 1px solid $border;
    border-radius: $radius-sm;
    padding: $spacing-md;
    font-size: 12px;
    font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
    color: $text-primary;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 300px;
    overflow-y: auto;
    margin: 0;
  }
}

.detail-stacktrace {
  pre {
    background: #fff2f0;
    border-color: #ffccc7;
    color: #cf1322;
  }
}

.drawer-actions {
  margin-top: $spacing-lg;
  padding-top: $spacing-md;
  border-top: 1px solid $border;
  display: flex;
  justify-content: flex-end;
}

:deep(.ant-descriptions-item-label) {
  font-weight: 500;
  color: $text-secondary;
  width: 100px;
}
</style>
