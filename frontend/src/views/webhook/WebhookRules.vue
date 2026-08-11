<template>
  <div class="page-container">
    <PageHeader title="Webhook 规则管理">
      <template #actions>
        <a-input
          v-model:value="repoKey"
          placeholder="输入仓库 Key"
          style="width: 200px;"
          @keyup.enter="loadRules"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
        </a-input>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          添加规则
        </a-button>
      </template>
    </PageHeader>

    <a-table
      :columns="columns"
      :data-source="webhookStore.rules"
      :loading="webhookStore.loading"
      row-key="id"
      :pagination="false"
      size="middle"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'name'">
          <span class="task-name">{{ record.name }}</span>
        </template>
        <template v-if="column.dataIndex === 'event_type'">
          <a-tag color="blue">{{ record.event_type || '全部' }}</a-tag>
        </template>
        <template v-if="column.dataIndex === 'branch_pattern'">
          <a-tag color="orange">{{ record.branch_pattern || '全部' }}</a-tag>
        </template>
        <template v-if="column.dataIndex === 'enabled'">
          <a-switch :checked="record.enabled" disabled size="small" />
        </template>
        <template v-if="column.dataIndex === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="openEdit(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除该规则吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      :width="500"
      @ok="handleSubmit"
      ok-text="确定"
      cancel-text="取消"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="仓库 Key" required>
          <a-input v-model:value="formData.repo_key" placeholder="请输入仓库 Key" />
        </a-form-item>
        <a-form-item label="事件类型">
          <a-select
            v-model:value="formData.event_type"
            style="width: 100%"
            :options="[
              { label: '全部', value: '' },
              { label: 'push', value: 'push' },
              { label: 'merge_request', value: 'merge_request' },
              { label: 'tag', value: 'tag' }
            ]"
          />
        </a-form-item>
        <a-form-item label="分支过滤">
          <a-input v-model:value="formData.branch_pattern" placeholder="如 main,feature/*" />
        </a-form-item>
        <a-form-item label="触发动作">
          <a-select
            v-model:value="formData.action"
            style="width: 100%"
            :options="[{ label: '同步', value: 'sync' }]"
          />
        </a-form-item>
        <a-form-item label="关联任务">
          <a-input v-model:value="formData.sync_task_keys" placeholder="任务 Key，逗号分隔" />
        </a-form-item>
        <a-form-item label="最小间隔">
          <a-input-number v-model:value="formData.min_interval" :min="0" :max="3600" />
          <span style="margin-left: 8px; color: #8c8c8c;">秒</span>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="formData.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { useWebhookStore } from '@/stores/webhook'
import PageHeader from '@/components/common/PageHeader.vue'
import type { WebhookRule } from '@/types'

const webhookStore = useWebhookStore()
const repoKey = ref('')
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref(0)

const columns = [
  { title: '规则名称', dataIndex: 'name', width: 200 },
  { title: '仓库', dataIndex: 'repo_key' },
  { title: '事件类型', dataIndex: 'event_type' },
  { title: '分支过滤', dataIndex: 'branch_pattern' },
  { title: '状态', dataIndex: 'enabled', width: 100 },
  { title: '操作', dataIndex: 'action', width: 120, align: 'center' as const },
]

const formData = reactive({
  name: '',
  repo_key: '',
  event_type: '',
  branch_pattern: '',
  action: 'sync',
  sync_task_keys: '',
  min_interval: 0,
  enabled: true,
})

function loadRules() {
  if (repoKey.value) {
    webhookStore.fetchRules(repoKey.value)
  }
}

onMounted(() => {
  repoKey.value = 'default'
  loadRules()
})

function openCreate() {
  dialogTitle.value = '添加规则'
  editingId.value = 0
  Object.assign(formData, {
    name: '',
    repo_key: repoKey.value,
    event_type: '',
    branch_pattern: '',
    action: 'sync',
    sync_task_keys: '',
    min_interval: 0,
    enabled: true,
  })
  dialogVisible.value = true
}

function openEdit(rule: WebhookRule) {
  dialogTitle.value = '编辑规则'
  editingId.value = rule.id
  Object.assign(formData, {
    name: rule.name,
    repo_key: rule.repo_key,
    event_type: rule.event_type,
    branch_pattern: rule.branch_pattern,
    action: rule.action,
    sync_task_keys: rule.sync_task_keys,
    min_interval: rule.min_interval,
    enabled: rule.enabled,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.repo_key) {
    message.warning('请填写必填字段')
    return
  }
  if (editingId.value) {
    await webhookStore.updateRule({ id: editingId.value, ...formData })
  } else {
    await webhookStore.createRule(formData)
  }
  dialogVisible.value = false
  loadRules()
}

async function handleDelete(id: number) {
  await webhookStore.deleteRule(id)
  loadRules()
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

.task-name {
  font-weight: 500;
}
</style>
