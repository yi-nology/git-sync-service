<template>
  <div class="page-container">
    <PageHeader title="事件日志">
      <template #actions>
        <a-input
          v-model:value="repoKey"
          placeholder="输入仓库 Key"
          style="width: 200px;"
          @pressEnter="loadEvents"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-button @click="loadEvents">
          <template #icon><ReloadOutlined /></template>
          刷新
        </a-button>
      </template>
    </PageHeader>

    <a-tabs v-model:activeKey="activeTab" style="margin-bottom: 16px;">
      <a-tab-pane key="all" tab="全部" />
      <a-tab-pane key="received" tab="已接收" />
      <a-tab-pane key="processed" tab="已处理" />
    </a-tabs>

    <a-table
      :columns="columns"
      :data-source="filteredEvents"
      :loading="webhookStore.loading"
      row-key="id"
      :pagination="false"
      size="middle"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'event_id'">
          <a-tooltip :title="record.event_id">
            <span class="truncated-id">{{ record.event_id }}</span>
          </a-tooltip>
        </template>
        <template v-if="column.dataIndex === 'event_type'">
          <a-tag color="blue">{{ record.event_type }}</a-tag>
        </template>
        <template v-if="column.dataIndex === 'source'">
          <span>{{ record.source || record.repo_key || '-' }}</span>
        </template>
        <template v-if="column.dataIndex === 'actor_name'">
          <span>{{ record.actor_name || '-' }}</span>
        </template>
        <template v-if="column.dataIndex === 'branch'">
          <span>{{ record.branch || '-' }}</span>
        </template>
        <template v-if="column.dataIndex === 'status'">
          <StatusBadge :status="record.status" />
        </template>
        <template v-if="column.dataIndex === 'processed_at'">
          <span style="color: #8c8c8c;">{{ record.processed_at || record.created_at || '-' }}</span>
        </template>
        <template v-if="column.dataIndex === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="showDetail(record)">详情</a-button>
            <a-button
              v-if="record.status === 'failed' || record.status === 'received'"
              type="link"
              size="small"
              @click="handleRetry(record.id)"
            >
              重试
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="detailVisible"
      title="事件详情"
      :width="600"
      :footer="null"
    >
      <a-descriptions
        v-if="currentEvent"
        :column="1"
        bordered
        size="small"
      >
        <a-descriptions-item label="事件 ID">
          <span>{{ currentEvent.event_id }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="触发时间">
          {{ currentEvent.created_at }}
        </a-descriptions-item>
        <a-descriptions-item label="事件类型">
          <a-tag color="blue">{{ currentEvent.event_type }}</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="仓库">
          {{ currentEvent.repo_key }}
        </a-descriptions-item>
        <a-descriptions-item label="来源">
          {{ currentEvent.source || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="分支">
          {{ currentEvent.branch || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="触发者">
          {{ currentEvent.actor_name || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="Commit">
          {{ currentEvent.commit_sha || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <StatusBadge :status="currentEvent.status" />
        </a-descriptions-item>
        <a-descriptions-item v-if="currentEvent.error_message" label="错误信息">
          <a-typography-text type="danger">{{ currentEvent.error_message }}</a-typography-text>
        </a-descriptions-item>
      </a-descriptions>
      <div style="text-align: right; margin-top: 16px;">
        <a-button @click="detailVisible = false">关闭</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { SearchOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { useWebhookStore } from '@/stores/webhook'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import type { WebhookEvent } from '@/types'

const webhookStore = useWebhookStore()
const repoKey = ref('')
const activeTab = ref('all')
const detailVisible = ref(false)
const currentEvent = ref<WebhookEvent | null>(null)

const columns = [
  { title: '事件 ID', dataIndex: 'event_id', width: 180 },
  { title: '事件类型', dataIndex: 'event_type', width: 120 },
  { title: '来源', dataIndex: 'source', width: 120 },
  { title: '操作者', dataIndex: 'actor_name', width: 100 },
  { title: '分支', dataIndex: 'branch', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '处理时间', dataIndex: 'processed_at', width: 160 },
  { title: '操作', dataIndex: 'action', width: 120, align: 'center' as const },
]

const filteredEvents = computed(() => {
  if (activeTab.value === 'all') return webhookStore.events
  return webhookStore.events.filter(e => e.status === activeTab.value)
})

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
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

.truncated-id {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
}
</style>
