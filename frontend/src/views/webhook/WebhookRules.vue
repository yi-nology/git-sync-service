<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Webhook 规则管理</h1>
      <div class="header-actions">
        <a-input v-model:value="repoKey" placeholder="输入仓库 Key" style="width: 200px;" @keyup.enter="loadRules"/>
        <button class="btn-primary" @click="openCreate">添加规则</button>
      </div>
    </div>

    <div v-if="webhookStore.loading" class="loading-state">加载中...</div>
    <div v-else-if="webhookStore.rules.length === 0" class="empty-state">暂无规则数据</div>
    <div v-else class="table-card">
      <table class="data-table">
        <thead><tr>
          <th style="width:200px;">规则名称</th>
          <th>仓库</th>
          <th>事件类型</th>
          <th>分支过滤</th>
          <th>状态</th>
          <th style="width:120px;text-align:center;">操作</th>
        </tr></thead>
        <tbody><tr v-for="rule in webhookStore.rules" :key="rule.id">
          <td class="task-name">{{ rule.name }}</td>
          <td>{{ rule.repo_key }}</td>
          <td><span class="badge event-tag">{{ rule.event_type || '全部' }}</span></td>
          <td><span class="badge filter-tag">{{ rule.branch_pattern || '全部' }}</span></td>
          <td><span class="status-badge" :class="rule.enabled ? 'enabled' : 'disabled'">{{ rule.enabled ? '已启用' : '已停用' }}</span></td>
          <td class="action-col">
            <button class="action-btn edit" @click="openEdit(rule)">编辑</button>
            <button class="action-btn delete" @click="handleDelete(rule.id)">删除</button>
          </td>
        </tr></tbody>
      </table>
    </div>

    <a-modal v-model:open="dialogVisible" :title="dialogTitle" :width="500" @ok="handleSubmit" okText="确定" cancelText="取消">
      <a-form :model="formData" :label-col="{ span: 6 }">
        <a-form-item label="规则名称">
          <a-input v-model:value="formData.name" placeholder="请输入规则名称"/>
        </a-form-item>
        <a-form-item label="仓库 Key">
          <a-input v-model:value="formData.repo_key" placeholder="请输入仓库 Key"/>
        </a-form-item>
        <a-form-item label="事件类型">
          <a-select v-model:value="formData.event_type" style="width: 100%"
            :options="[{label: '全部', value: ''}, {label: 'push', value: 'push'}, {label: 'merge_request', value: 'merge_request'}, {label: 'tag', value: 'tag'}]"/>
        </a-form-item>
        <a-form-item label="分支过滤">
          <a-input v-model:value="formData.branch_pattern" placeholder="如 main,feature/*"/>
        </a-form-item>
        <a-form-item label="触发动作">
          <a-select v-model:value="formData.action" style="width: 100%"
            :options="[{label: '同步', value: 'sync'}]"/>
        </a-form-item>
        <a-form-item label="关联任务">
          <a-input v-model:value="formData.sync_task_keys" placeholder="任务 Key，逗号分隔"/>
        </a-form-item>
        <a-form-item label="最小间隔">
          <a-input-number v-model:value="formData.min_interval" :min="0" :max="3600"/>
          <span style="margin-left: 8px; color: #8c8c8c;">秒</span>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="formData.enabled"/>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useWebhookStore } from '@/stores/webhook'
import type { WebhookRule } from '@/types'

const webhookStore = useWebhookStore()
const repoKey = ref('')
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref(0)

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
  Object.assign(formData, { name: '', repo_key: repoKey.value, event_type: '', branch_pattern: '', action: 'sync', sync_task_keys: '', min_interval: 0, enabled: true })
  dialogVisible.value = true
}

function openEdit(rule: WebhookRule) {
  dialogTitle.value = '编辑规则'
  editingId.value = rule.id
  Object.assign(formData, { name: rule.name, repo_key: rule.repo_key, event_type: rule.event_type, branch_pattern: rule.branch_pattern, action: rule.action, sync_task_keys: rule.sync_task_keys, min_interval: rule.min_interval, enabled: rule.enabled })
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
  await new Promise<void>((resolve, reject) => {
    Modal.confirm({
      title: '提示',
      content: '确定要删除该规则吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => resolve(),
      onCancel: () => reject(new Error('cancelled')),
    })
  })
  await webhookStore.deleteRule(id)
  loadRules()
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; align-items: center; }
.btn-primary { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; cursor: pointer; &:hover { background: #40a9ff; } }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
.task-name { font-weight: 500; }
.badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block; }
.event-tag { background: #e6f7ff; color: #1890ff; }
.filter-tag { background: #fff7e6; color: #fa8c16; }
.status-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; &.enabled { background: #f6ffed; color: #52c41a; } &.disabled { background: #f5f5f5; color: #8c8c8c; } }
.action-col { display: flex; justify-content: center; gap: 8px; }
.action-btn { padding: 4px 8px; border-radius: 4px; border: none; cursor: pointer; font-size: 12px; transition: all 0.2s;
  &.edit { background: #e6f7ff; color: #1890ff; &:hover { background: #bae7ff; } }
  &.delete { background: #fff2f0; color: #ff4d4f; &:hover { background: #ffccc7; } }
}
</style>
