<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">Webhook 规则</h1>
        <p class="page-subtitle">配置仓库的 Webhook 触发规则</p>
      </div>
      <a-space>
        <a-select
          v-model:value="repoKey"
          placeholder="选择仓库"
          style="width: 220px"
          show-search
          :filter-option="filterRepoOption"
          @change="loadRules"
        >
          <a-select-option v-for="repo in repoStore.repos" :key="repo.key" :value="repo.key">
            {{ repo.name }}
          </a-select-option>
        </a-select>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          添加规则
        </a-button>
      </a-space>
    </div>

    <a-table
      :columns="columns"
      :data-source="webhookStore.rules"
      :loading="webhookStore.loading"
      row-key="id"
      :pagination="pagination"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'name'">
          <span class="rule-name">{{ record.name }}</span>
        </template>
        <template v-if="column.dataIndex === 'event_type'">
          <a-tag :color="eventTypeColor(record.event_type)">
            {{ record.event_type || '全部' }}
          </a-tag>
        </template>
        <template v-if="column.dataIndex === 'branch_pattern'">
          <span class="branch-tag">{{ record.branch_pattern || '*' }}</span>
        </template>
        <template v-if="column.dataIndex === 'enabled'">
          <a-switch
            :checked="record.enabled"
            checked-children="启用"
            un-checked-children="禁用"
            @change="(checked: boolean) => handleToggleEnabled(record, checked)"
          />
        </template>
        <template v-if="column.dataIndex === 'action'">
          <a-space :size="4">
            <a-button type="link" size="small" @click="openEdit(record)">
              <template #icon><EditOutlined /></template>
            </a-button>
            <a-popconfirm
              title="确定要删除该规则吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" size="small" danger>
                <template #icon><DeleteOutlined /></template>
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
      <template #emptyText>
        <a-empty description="暂无 Webhook 规则">
          <a-button type="primary" @click="openCreate">
            <template #icon><PlusOutlined /></template>
            添加第一条规则
          </a-button>
        </a-empty>
      </template>
    </a-table>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      :width="560"
      @ok="handleSubmit"
      ok-text="确定"
      cancel-text="取消"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入规则名称" />
        </a-form-item>

        <a-form-item label="仓库" required>
          <a-select
            v-model:value="formData.repo_key"
            placeholder="选择仓库"
            show-search
            :filter-option="filterRepoOption"
            style="width: 100%"
          >
            <a-select-option v-for="repo in repoStore.repos" :key="repo.key" :value="repo.key">
              {{ repo.name }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="事件类型">
              <a-select
                v-model:value="formData.event_type"
                style="width: 100%"
                :options="[
                  { label: '全部事件', value: '' },
                  { label: 'Push', value: 'push' },
                  { label: 'Merge Request', value: 'merge_request' },
                  { label: 'Tag', value: 'tag' }
                ]"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="触发动作">
              <a-select
                v-model:value="formData.action"
                style="width: 100%"
                :options="[{ label: '同步', value: 'sync' }]"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="分支过滤">
          <a-input v-model:value="formData.branch_pattern" placeholder="如 main,feature/*" />
          <div class="form-tip">
            <InfoCircleOutlined /> 支持 glob 模式，多个用逗号分隔。留空表示匹配所有分支
          </div>
        </a-form-item>

        <a-form-item label="关联同步任务">
          <a-select
            v-model:value="formData.sync_task_keys"
            mode="multiple"
            placeholder="选择关联的同步任务"
            style="width: 100%"
            :options="taskStore.tasks.map(t => ({ label: t.name, value: t.key }))"
          />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="最小间隔 (秒)">
              <a-input-number v-model:value="formData.min_interval" :min="0" :max="3600" style="width: 100%" />
              <div class="form-tip">
                <InfoCircleOutlined /> 防止短时间内重复触发，0 表示不限制
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="启用状态">
              <a-switch
                v-model:checked="formData.enabled"
                checked-children="启用"
                un-checked-children="禁用"
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'
import { useWebhookStore } from '@/stores/webhook'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import type { WebhookRule } from '@/types'
import { eventTypeColor } from '@/utils/dictionaries'
import { notifySuccess, notifyError, notifyWarning } from '@/utils/notify'

const webhookStore = useWebhookStore()
const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()
const repoKey = ref('')
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref(0)

const columns = [
  { title: '规则名称', dataIndex: 'name', width: 180, ellipsis: true },
  { title: '事件类型', dataIndex: 'event_type', width: 120, align: 'center' as const },
  { title: '分支过滤', dataIndex: 'branch_pattern', width: 140 },
  { title: '状态', dataIndex: 'enabled', width: 100, align: 'center' as const },
  { title: '操作', dataIndex: 'action', width: 100, align: 'center' as const },
]

const pagination = {
  pageSize: 10,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条规则`,
}

const formData = reactive({
  name: '',
  repo_key: '',
  event_type: '',
  branch_pattern: '',
  action: 'sync',
  sync_task_keys: [] as string[],
  min_interval: 0,
  enabled: true,
})

function filterRepoOption(input: string, option: any) {
  const repo = repoStore.repos.find(r => r.key === option.value)
  return repo?.name.toLowerCase().includes(input.toLowerCase()) || false
}

function loadRules() {
  if (repoKey.value) {
    webhookStore.fetchRules(repoKey.value).catch((e) => notifyError(e, '加载规则失败'))
  }
}

onMounted(async () => {
  try {
    await repoStore.fetchRepos()
    await taskStore.fetchTasks()
    if (repoStore.repos.length > 0) {
      repoKey.value = repoStore.repos[0].key
      loadRules()
    }
  } catch (e) {
    notifyError(e, '加载数据失败')
  }
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
    sync_task_keys: [],
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
    sync_task_keys: rule.sync_task_keys ? rule.sync_task_keys.split(',').filter(Boolean) : [],
    min_interval: rule.min_interval,
    enabled: rule.enabled,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.repo_key) {
    notifyWarning('请填写必填字段')
    return
  }
  const submitData = {
    ...formData,
    sync_task_keys: Array.isArray(formData.sync_task_keys) ? formData.sync_task_keys.join(',') : formData.sync_task_keys,
  }
  try {
    if (editingId.value) {
      await webhookStore.updateRule({ id: editingId.value, ...submitData })
      notifySuccess('更新规则成功')
    } else {
      await webhookStore.createRule(submitData)
      notifySuccess('创建规则成功')
    }
    dialogVisible.value = false
    loadRules()
  } catch (e) {
    notifyError(e, '保存规则失败')
  }
}

async function handleDelete(id: number) {
  try {
    await webhookStore.deleteRule(id)
    notifySuccess('删除规则成功')
    loadRules()
  } catch (e) {
    notifyError(e, '删除规则失败')
  }
}

async function handleToggleEnabled(rule: WebhookRule, checked: boolean) {
  try {
    await webhookStore.updateRule({ id: rule.id, enabled: checked })
    loadRules()
  } catch (e) {
    notifyError(e, '更新状态失败')
  }
}
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

.rule-name {
  font-weight: 500;
  color: $text-primary;
}

.form-tip {
  font-size: 12px;
  color: $text-secondary;
  margin-top: 4px;

  .anticon {
    margin-right: 4px;
  }
}
</style>
