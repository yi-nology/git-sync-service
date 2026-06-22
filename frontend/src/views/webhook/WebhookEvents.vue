<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Webhook 事件日志</h1>
      <div class="header-actions">
        <el-input v-model="repoKey" placeholder="输入仓库 Key" style="width: 200px;" @keyup.enter="loadEvents"/>
        <button class="btn-default" @click="loadEvents">刷新</button>
      </div>
    </div>

    <div class="tabs-bar">
      <button class="tab-btn" :class="{active: activeTab === 'all'}" @click="activeTab = 'all'">全部</button>
      <button class="tab-btn" :class="{active: activeTab === 'received'}" @click="activeTab = 'received'">已接收</button>
      <button class="tab-btn" :class="{active: activeTab === 'processed'}" @click="activeTab = 'processed'">已处理</button>
    </div>

    <div v-if="webhookStore.loading" class="loading-state">加载中...</div>
    <div v-else-if="filteredEvents.length === 0" class="empty-state">暂无事件数据</div>
    <div v-else class="table-card">
      <table class="data-table">
        <thead><tr>
          <th style="width:180px;">时间</th>
          <th>事件类型</th>
          <th>仓库</th>
          <th>分支</th>
          <th>触发者</th>
          <th>状态</th>
          <th style="width:120px;text-align:center;">操作</th>
        </tr></thead>
        <tbody><tr v-for="e in filteredEvents" :key="e.id">
          <td style="color:#8c8c8c;">{{ e.created_at }}</td>
          <td><span class="event-badge" :class="e.event_type">{{ e.event_type }}</span></td>
          <td>{{ e.repo_key }}</td>
          <td>{{ e.branch || '-' }}</td>
          <td>{{ e.actor_name || '-' }}</td>
          <td><span class="status-badge" :class="e.status">{{ statusText(e.status) }}</span></td>
          <td class="action-col">
            <button class="action-btn view" @click="showDetail(e)">详情</button>
            <button v-if="e.status === 'received'" class="action-btn retry" @click="handleRetry(e.id)">重试</button>
          </td>
        </tr></tbody>
      </table>
    </div>

    <el-dialog v-model="detailVisible" title="事件详情" width="600px">
      <div class="detail-content" v-if="currentEvent">
        <div class="detail-row"><span class="detail-label">事件ID</span><span class="detail-value">{{ currentEvent.event_id }}</span></div>
        <div class="detail-row"><span class="detail-label">触发时间</span><span class="detail-value">{{ currentEvent.created_at }}</span></div>
        <div class="detail-row"><span class="detail-label">事件类型</span><span class="event-badge" :class="currentEvent.event_type">{{ currentEvent.event_type }}</span></div>
        <div class="detail-row"><span class="detail-label">仓库</span><span class="detail-value">{{ currentEvent.repo_key }}</span></div>
        <div class="detail-row"><span class="detail-label">分支</span><span class="detail-value">{{ currentEvent.branch || '-' }}</span></div>
        <div class="detail-row"><span class="detail-label">触发者</span><span class="detail-value">{{ currentEvent.actor_name || '-' }}</span></div>
        <div class="detail-row"><span class="detail-label">Commit</span><span class="detail-value">{{ currentEvent.commit_sha || '-' }}</span></div>
        <div class="detail-row"><span class="detail-label">状态</span><span class="status-badge" :class="currentEvent.status">{{ statusText(currentEvent.status) }}</span></div>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useWebhookStore } from '@/stores/webhook'
import type { WebhookEvent } from '@/types'

const webhookStore = useWebhookStore()
const repoKey = ref('')
const activeTab = ref('all')
const detailVisible = ref(false)
const currentEvent = ref<WebhookEvent | null>(null)

const filteredEvents = computed(() => {
  if (activeTab.value === 'all') return webhookStore.events
  return webhookStore.events.filter(e => e.status === activeTab.value)
})

function statusText(status: string) {
  const map: Record<string, string> = { received: '已接收', processed: '已处理', failed: '失败' }
  return map[status] || status
}

function loadEvents() {
  if (repoKey.value) {
    webhookStore.fetchEvents(repoKey.value)
  }
}

onMounted(() => {
  repoKey.value = 'default'
  loadEvents()
})

function showDetail(e: WebhookEvent) {
  currentEvent.value = e
  detailVisible.value = true
}

async function handleRetry(id: number) {
  await webhookStore.retryEvent(id)
  loadEvents()
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; align-items: center; }
.btn-default { padding: 7px 14px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
.tabs-bar { display: flex; gap: 4px; margin-bottom: 16px; }
.tab-btn { padding: 6px 16px; border: none; background: transparent; border-radius: 6px; font-size: 13px; color: #595959; cursor: pointer; transition: all 0.2s; &:hover { color: #1890ff; } &.active { background: #1890ff; color: #fff; } }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
.event-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block;
  &.push { background: #e6f7ff; color: #1890ff; }
  &.merge_request { background: #f6ffed; color: #52c41a; }
  &.tag { background: #fff7e6; color: #fa8c16; }
}
.status-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block;
  &.received { background: #e6f7ff; color: #1890ff; }
  &.processed { background: #f6ffed; color: #52c41a; }
  &.failed { background: #fff2f0; color: #ff4d4f; }
}
.action-col { display: flex; justify-content: center; gap: 8px; }
.action-btn { padding: 4px 8px; border-radius: 4px; border: none; cursor: pointer; font-size: 12px; transition: all 0.2s;
  &.view { background: #e6f7ff; color: #1890ff; &:hover { background: #bae7ff; } }
  &.retry { background: #fff7e6; color: #fa8c16; &:hover { background: #ffe7ba; } }
}
.detail-content { display: flex; flex-direction: column; gap: 16px; }
.detail-row { display: flex; align-items: center; gap: 12px; }
.detail-label { min-width: 80px; font-size: 13px; color: #8c8c8c; }
.detail-value { font-size: 13px; color: #262626; flex: 1; }
</style>
