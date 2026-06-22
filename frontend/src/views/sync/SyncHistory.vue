<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">同步历史</h1>
      <div class="header-actions">
        <el-input v-model="taskKey" placeholder="输入任务 Key" style="width: 200px;" @keyup.enter="loadHistory"/>
        <button class="btn-default" @click="loadHistory">查询</button>
      </div>
    </div>

    <div v-if="taskStore.history.length === 0" class="empty-state">暂无历史数据</div>
    <div v-else class="history-list">
      <div class="history-card" v-for="run in taskStore.history" :key="run.id">
        <div class="card-header">
          <div class="header-left">
            <div class="status-icon" :class="run.status">
              <svg v-if="run.status === 'success'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            </div>
            <div class="title-wrap">
              <div class="task-key">{{ run.task_key }}</div>
              <div class="task-trigger">{{ run.trigger_source }} | {{ run.start_time }}</div>
            </div>
          </div>
          <span class="status-badge" :class="run.status">{{ statusText(run.status) }}</span>
        </div>
        <div class="card-body" v-if="run.details">
          <pre class="details-pre">{{ run.details }}</pre>
        </div>
        <div class="card-body" v-if="run.error_message">
          <div class="error-message">{{ run.error_message }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSyncTaskStore } from '@/stores/syncTask'

const taskStore = useSyncTaskStore()
const taskKey = ref('')

function statusText(status: string) {
  const map: Record<string, string> = { success: '成功', running: '运行中', failed: '失败' }
  return map[status] || status
}

function loadHistory() {
  if (taskKey.value) {
    taskStore.fetchHistory(taskKey.value)
  }
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; align-items: center; }
.btn-default { padding: 7px 14px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
.empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.history-list { display: flex; flex-direction: column; gap: 16px; }
.history-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
.card-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid #f0f0f0; }
.header-left { display: flex; align-items: center; gap: 12px; }
.status-icon { width: 36px; height: 36px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  &.success { background: #f6ffed; color: #52c41a; }
  &.failed { background: #fff2f0; color: #ff4d4f; }
  &.running { background: #e6f7ff; color: #1890ff; }
}
.task-key { font-weight: 600; color: #262626; font-size: 14px; }
.task-trigger { font-size: 12px; color: #8c8c8c; margin-top: 2px; }
.status-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px;
  &.success { background: #f6ffed; color: #52c41a; }
  &.failed { background: #fff2f0; color: #ff4d4f; }
  &.running { background: #e6f7ff; color: #1890ff; }
}
.card-body { padding: 16px 20px; }
.details-pre { background: #f5f5f5; border-radius: 6px; padding: 12px; font-size: 12px; line-height: 1.6; color: #595959; margin: 0; max-height: 200px; overflow-y: auto; white-space: pre-wrap; }
.error-message { background: #fff2f0; border: 1px solid #ffccc7; border-radius: 6px; padding: 12px; font-size: 13px; color: #ff4d4f; }
</style>
