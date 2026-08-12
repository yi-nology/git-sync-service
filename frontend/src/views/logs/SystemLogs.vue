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
        <div class="stat-icon" style="background: #FFF2F0; color: #FF4D4F;">
          <CloseCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-error">{{ stats.errorCount }}</div>
          <div class="stat-label">ERROR</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #FFF7E6; color: #FAAD14;">
          <WarningOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-warn">{{ stats.warnCount }}</div>
          <div class="stat-label">WARN</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #E6F7FF; color: #1677FF;">
          <InfoCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-value stat-info">{{ stats.infoCount }}</div>
          <div class="stat-label">INFO</div>
        </div>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
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
          <span class="log-time">{{ formatTime(record.time) }}</span>
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
        <!-- Basic Info -->
        <a-descriptions :column="2" bordered size="small" class="detail-section">
          <a-descriptions-item label="时间" :span="2">
            {{ formatTime(currentLog.time) }}
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

        <!-- Message -->
        <a-divider>日志消息</a-divider>
        <div class="detail-message">
          <pre>{{ currentLog.message }}</pre>
        </div>

        <!-- Details -->
        <template v-if="currentLog.details">
          <a-divider>详细信息</a-divider>
          <div class="detail-extra">
            <pre>{{ currentLog.details }}</pre>
          </div>
        </template>

        <!-- Stack Trace -->
        <template v-if="getStackTrace(currentLog)">
          <a-divider>堆栈跟踪</a-divider>
          <div class="detail-stacktrace">
            <pre>{{ getStackTrace(currentLog) }}</pre>
          </div>
        </template>

        <!-- Actions -->
        <div class="drawer-actions">
          <a-button @click="drawerVisible = false">关闭</a-button>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, reactive, computed } from 'vue'
import { Empty } from 'ant-design-vue'
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
import type { SystemLog, SystemLogRequest } from '@/api'
import dayjs, { type Dayjs } from 'dayjs'

const loading = ref(false)
const logs = ref<SystemLog[]>([])
const drawerVisible = ref(false)
const currentLog = ref<SystemLog | null>(null)
let refreshTimer: ReturnType<typeof setInterval> | null = null

const autoRefreshInterval = ref(0)

const stats = reactive({
  errorCount: 0,
  warnCount: 0,
  infoCount: 0,
})

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

// Extract module from details or message
function getModule(record: SystemLog): string {
  if (!record.details) return '-'
  try {
    const parsed = JSON.parse(record.details)
    if (parsed.module) return parsed.module
  } catch {
    // Try to extract from message prefix like [module] or [server]
    const match = record.message?.match(/^\[([a-zA-Z]+)\]/)
    if (match) return match[1]
  }
  return '-'
}

// Extract stack trace from details
function getStackTrace(record: SystemLog): string | null {
  if (!record.details) return null
  try {
    const parsed = JSON.parse(record.details)
    if (parsed.stack_trace || parsed.stackTrace) {
      return parsed.stack_trace || parsed.stackTrace
    }
  } catch {
    // If details contains "Stack Trace:" or "goroutine" patterns
    const stackMatch = record.details.match(/(?:Stack Trace:|goroutine \d+)([\s\S]*)/i)
    if (stackMatch) return stackMatch[1].trim()
  }
  return null
}

// Client-side filtered logs
const filteredLogs = computed(() => {
  let result = logs.value

  if (filters.module) {
    result = result.filter(log => getModule(log) === filters.module)
  }

  if (filters.search) {
    const searchLower = filters.search.toLowerCase()
    result = result.filter(log =>
      log.message?.toLowerCase().includes(searchLower) ||
      log.details?.toLowerCase().includes(searchLower)
    )
  }

  if (filters.startDate && filters.endDate) {
    result = result.filter(log => {
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

function formatTime(time: string): string {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
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
  filters.level = undefined
  filters.module = undefined
  filters.search = ''
  filters.dateRange = null
  filters.startDate = ''
  filters.endDate = ''
  pagination.current = 1
  fetchLogs()
}

function handleTableChange(pag: any) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
}

function openDrawer(record: SystemLog) {
  currentLog.value = record
  drawerVisible.value = true
}

function calculateStats(logList: SystemLog[]) {
  stats.errorCount = logList.filter(l => l.level === 'ERROR').length
  stats.warnCount = logList.filter(l => l.level === 'WARN').length
  stats.infoCount = logList.filter(l => l.level === 'INFO').length
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
    const params: SystemLogRequest = {
      page: pagination.current,
      page_size: pagination.pageSize,
    }
    if (filters.level) params.level = filters.level

    const resp: any = await logApi.listSystem(params)
    // 响应信封为 { code, message, data: { list, pagination } }
    const data = resp?.data ?? resp
    logs.value = data.list || []
    pagination.total = data.pagination?.total || 0
    calculateStats(logs.value)
  } catch (e: any) {
    // 接口不可用时如实清空，不再展示假数据
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
  line-height: 1.2;

  &.stat-error {
    color: #FF4D4F;
  }

  &.stat-warn {
    color: #FAAD14;
  }

  &.stat-info {
    color: #1677FF;
  }
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

// Drawer styles
.detail-section {
  margin-bottom: $spacing-md;
}

.detail-message,
.detail-extra,
.detail-stacktrace {
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
    margin: 0;
  }
}

.detail-stacktrace {
  pre {
    background: #FFF2F0;
    border-color: #FFCCC7;
    color: #CF1322;
  }
}

.drawer-actions {
  margin-top: $spacing-lg;
  padding-top: $spacing-md;
  border-top: 1px solid $border-color;
  display: flex;
  justify-content: flex-end;
}

:deep(.ant-descriptions-item-label) {
  font-weight: 500;
  color: $text-secondary;
  width: 100px;
}
</style>
