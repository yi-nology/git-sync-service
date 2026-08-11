<template>
  <div class="page-container">
    <PageHeader title="同步任务">
      <template #actions>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          创建任务
        </a-button>
      </template>
    </PageHeader>

    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-num blue">{{ taskStore.total }}</div>
        <div class="stat-name">总任务数</div>
      </div>
      <div class="stat-card">
        <div class="stat-num green">{{ taskStore.tasks.filter(t => t.last_status === 'success').length }}</div>
        <div class="stat-name">成功</div>
      </div>
      <div class="stat-card">
        <div class="stat-num orange">{{ taskStore.tasks.filter(t => t.last_status === 'running').length }}</div>
        <div class="stat-name">运行中</div>
      </div>
      <div class="stat-card">
        <div class="stat-num red">{{ taskStore.tasks.filter(t => t.last_status === 'failed').length }}</div>
        <div class="stat-name">失败</div>
      </div>
    </div>

    <div class="content-card">
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="taskStore.tasks"
          :loading="taskStore.loading"
          :pagination="pagination"
          row-key="key"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <span class="task-name">{{ record.name }}</span>
            </template>
            <template v-if="column.key === 'branches'">
              <a-space :size="4">
                <span class="branch-tag">{{ record.source_branch }}</span>
                <ArrowRightOutlined style="color: #8C8C8C; font-size: 12px;" />
                <span class="branch-tag">{{ record.target_branch }}</span>
              </a-space>
            </template>
            <template v-if="column.key === 'mode'">
              <a-tag :color="record.sync_mode === 'all' ? 'blue' : 'default'">
                {{ record.sync_mode === 'all' ? '全分支' : '单分支' }}
              </a-tag>
            </template>
            <template v-if="column.key === 'status'">
              <StatusBadge :status="record.last_status" />
            </template>
            <template v-if="column.key === 'last_run'">
              <span class="time-text">{{ record.last_run_at || '未运行' }}</span>
            </template>
            <template v-if="column.key === 'actions'">
              <a-space :size="8">
                <a-button type="link" size="small" @click="handleRun(record.key)">
                  <template #icon><PlayCircleOutlined /></template>
                  运行
                </a-button>
                <a-button type="link" size="small" @click="openEdit(record)">
                  <template #icon><EditOutlined /></template>
                  编辑
                </a-button>
                <a-popconfirm
                  title="确定要删除该任务吗？"
                  ok-text="确定"
                  cancel-text="取消"
                  @confirm="handleDelete(record.key)"
                >
                  <a-button type="link" danger size="small">
                    <template #icon><DeleteOutlined /></template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </div>
    </div>

    <a-modal v-model:open="dialogVisible" :title="dialogTitle" :width="600" @ok="handleSubmit" okText="确定" cancelText="取消">
      <a-form :model="formData" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-form-item label="源仓库 Key" required>
          <a-input v-model:value="formData.source_repo_key" placeholder="请输入源仓库 Key" />
        </a-form-item>
        <a-form-item label="源分支">
          <a-input v-model:value="formData.source_branch" placeholder="main" />
        </a-form-item>
        <a-form-item label="目标仓库 Key" required>
          <a-input v-model:value="formData.target_repo_key" placeholder="请输入目标仓库 Key" />
        </a-form-item>
        <a-form-item label="目标分支">
          <a-input v-model:value="formData.target_branch" placeholder="main" />
        </a-form-item>
        <a-form-item label="Cron 表达式">
          <a-input v-model:value="formData.cron" placeholder="可选，如 0 */5 * * * *" />
        </a-form-item>
        <a-form-item label="同步模式">
          <a-select v-model:value="formData.sync_mode" style="width: 100%"
            :options="[{label: '单分支', value: 'single'}, {label: '全分支', value: 'all'}]" />
        </a-form-item>
        <a-form-item label="选项">
          <a-space direction="vertical" :size="8">
            <a-checkbox v-model:checked="formData.git_tags">同步 Tags</a-checkbox>
            <a-checkbox v-model:checked="formData.git_force">强制推送</a-checkbox>
            <a-checkbox v-model:checked="formData.git_prune">Prune</a-checkbox>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  PlayCircleOutlined,
  EditOutlined,
  DeleteOutlined,
  ArrowRightOutlined,
} from '@ant-design/icons-vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import type { SyncTask } from '@/types'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'

const taskStore = useSyncTaskStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingKey = ref('')

const formData = reactive({
  name: '',
  source_repo_key: '',
  source_branch: 'main',
  target_repo_key: '',
  target_branch: 'main',
  sync_mode: 'single',
  cron: '',
  git_tags: false,
  git_force: false,
  git_prune: false,
})

const columns = [
  { title: '任务名称', key: 'name', dataIndex: 'name', width: 200 },
  { title: '分支', key: 'branches', width: 220 },
  { title: '同步模式', key: 'mode', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '最后运行', key: 'last_run', width: 160 },
  { title: '操作', key: 'actions', width: 200, fixed: 'right' as const },
]

const pagination = {
  pageSize: 10,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

onMounted(() => {
  taskStore.fetchTasks()
})

function openCreate() {
  dialogTitle.value = '创建任务'
  editingKey.value = ''
  Object.assign(formData, {
    name: '', source_repo_key: '', source_branch: 'main',
    target_repo_key: '', target_branch: 'main', sync_mode: 'single',
    cron: '', git_tags: false, git_force: false, git_prune: false,
  })
  dialogVisible.value = true
}

function openEdit(task: SyncTask) {
  dialogTitle.value = '编辑任务'
  editingKey.value = task.key
  Object.assign(formData, {
    name: task.name,
    source_repo_key: task.source_repo_key,
    source_branch: task.source_branch,
    target_repo_key: task.target_repo_key,
    target_branch: task.target_branch,
    sync_mode: task.sync_mode,
    cron: task.cron,
    git_tags: task.git_tags,
    git_force: task.git_force,
    git_prune: task.git_prune,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.source_repo_key || !formData.target_repo_key) {
    message.warning('请填写必填字段')
    return
  }
  if (editingKey.value) {
    await taskStore.updateTask({ key: editingKey.value, ...formData })
  } else {
    await taskStore.createTask(formData)
  }
  dialogVisible.value = false
}

async function handleDelete(key: string) {
  await taskStore.deleteTask(key)
}

async function handleRun(key: string) {
  await taskStore.runTask(key)
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.stat-card {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: 24px;
  box-shadow: $shadow-card;
  text-align: center;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;

  &.blue { color: $primary-color; }
  &.green { color: $success-color; }
  &.orange { color: $warning-color; }
  &.red { color: $error-color; }
}

.stat-name {
  font-size: 13px;
  color: $text-secondary;
}

.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  overflow: hidden;
}

.card-body {
  padding: 0;
}

.task-name {
  font-weight: 500;
  color: $text-primary;
}

.branch-tag {
  display: inline-block;
  padding: 2px 8px;
  background: #F5F5F5;
  border-radius: $border-radius-sm;
  font-size: 12px;
  color: $text-primary;
}

.time-text {
  color: $text-secondary;
  font-size: 13px;
}
</style>
