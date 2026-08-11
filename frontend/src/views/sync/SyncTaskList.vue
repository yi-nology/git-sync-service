<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">同步任务</h1>
      <div class="header-actions">
        <button class="btn-primary" @click="openCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建任务
        </button>
      </div>
    </div>

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

    <div v-if="taskStore.loading" class="loading-state">加载中...</div>
    <div v-else-if="taskStore.tasks.length === 0" class="empty-state">暂无任务数据</div>
    <div v-else class="table-card">
      <table class="data-table">
        <thead><tr>
          <th style="text-align:left;">任务名称</th>
          <th style="text-align:left;">源分支</th>
          <th style="text-align:left;">目标分支</th>
          <th style="text-align:left;">同步模式</th>
          <th style="text-align:left;width:100px;">状态</th>
          <th style="text-align:left;width:120px;">最后运行</th>
          <th style="text-align:center;width:180px;">操作</th>
        </tr></thead>
        <tbody><tr v-for="task in taskStore.tasks" :key="task.key">
          <td><span class="task-name">{{ task.name }}</span></td>
          <td><span class="branch-tag">{{ task.source_branch }}</span></td>
          <td><span class="branch-tag">{{ task.target_branch }}</span></td>
          <td><span class="mode-text">{{ task.sync_mode || 'single' }}</span></td>
          <td>
            <span class="status-dot" :class="task.last_status || 'idle'"></span>
            {{ statusText(task.last_status) }}
          </td>
          <td class="time-col">{{ task.last_run_at || '未运行' }}</td>
          <td class="action-col">
            <button class="action-btn run" @click="handleRun(task.key)">运行</button>
            <button class="action-btn edit" @click="openEdit(task)">编辑</button>
            <button class="action-btn delete" @click="handleDelete(task.key)">删除</button>
          </td>
        </tr></tbody>
      </table>
    </div>

    <a-modal v-model:open="dialogVisible" :title="dialogTitle" :width="600" @ok="handleSubmit" okText="确定" cancelText="取消">
      <a-form :model="formData" :label-col="{ span: 6 }">
        <a-form-item label="任务名称">
          <a-input v-model:value="formData.name" placeholder="请输入任务名称"/>
        </a-form-item>
        <a-form-item label="源仓库 Key">
          <a-input v-model:value="formData.source_repo_key" placeholder="请输入源仓库 Key"/>
        </a-form-item>
        <a-form-item label="源分支">
          <a-input v-model:value="formData.source_branch" placeholder="main"/>
        </a-form-item>
        <a-form-item label="目标仓库 Key">
          <a-input v-model:value="formData.target_repo_key" placeholder="请输入目标仓库 Key"/>
        </a-form-item>
        <a-form-item label="目标分支">
          <a-input v-model:value="formData.target_branch" placeholder="main"/>
        </a-form-item>
        <a-form-item label="Cron 表达式">
          <a-input v-model:value="formData.cron" placeholder="可选，如 0 */5 * * * *"/>
        </a-form-item>
        <a-form-item label="同步模式">
          <a-select v-model:value="formData.sync_mode" style="width: 100%"
            :options="[{label: '单分支', value: 'single'}, {label: '全分支', value: 'all'}]"/>
        </a-form-item>
        <a-form-item label="选项">
          <a-checkbox v-model:checked="formData.git_tags">同步 Tags</a-checkbox>
          <a-checkbox v-model:checked="formData.git_force">强制推送</a-checkbox>
          <a-checkbox v-model:checked="formData.git_prune">Prune</a-checkbox>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import type { SyncTask } from '@/types'

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

onMounted(() => {
  taskStore.fetchTasks()
})

function statusText(status: string) {
  const map: Record<string, string> = { success: '成功', running: '运行中', failed: '失败', idle: '未运行' }
  return map[status] || '未运行'
}

function openCreate() {
  dialogTitle.value = '创建任务'
  editingKey.value = ''
  Object.assign(formData, { name: '', source_repo_key: '', source_branch: 'main', target_repo_key: '', target_branch: 'main', sync_mode: 'single', cron: '', git_tags: false, git_force: false, git_prune: false })
  dialogVisible.value = true
}

function openEdit(task: SyncTask) {
  dialogTitle.value = '编辑任务'
  editingKey.value = task.key
  Object.assign(formData, {
    name: task.name, source_repo_key: task.source_repo_key, source_branch: task.source_branch,
    target_repo_key: task.target_repo_key, target_branch: task.target_branch,
    sync_mode: task.sync_mode, cron: task.cron,
    git_tags: task.git_tags, git_force: task.git_force, git_prune: task.git_prune,
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
  await new Promise<void>((resolve, reject) => {
    Modal.confirm({
      title: '提示',
      content: '确定要删除该任务吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => resolve(),
      onCancel: () => reject(new Error('cancelled')),
    })
  })
  await taskStore.deleteTask(key)
}

async function handleRun(key: string) {
  await taskStore.runTask(key)
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.btn-primary { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; cursor: pointer; &:hover { background: #40a9ff; } }
.stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 20px; }
.stat-card { background: #fff; border-radius: 8px; padding: 20px; text-align: center; }
.stat-num { font-size: 28px; font-weight: 700; margin-bottom: 4px; &.blue { color: #1890ff; } &.green { color: #52c41a; } &.orange { color: #faad14; } &.red { color: #ff4d4f; } }
.stat-name { font-size: 13px; color: #8c8c8c; }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; background: #fff; border-radius: 8px; }
.table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
.data-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
.task-name { font-weight: 500; }
.branch-tag { padding: 5px 12px; background: #f5f5f5; border-radius: 4px; font-size: 12px; }
.mode-text { color: #8c8c8c; }
.status-dot { display: inline-block; width: 6px; height: 6px; border-radius: 50%; margin-right: 8px; background: #d9d9d9; &.success { background: #52c41a; } &.running { background: #1890ff; } &.failed { background: #ff4d4f; } }
.time-col { color: #8c8c8c; }
.action-col { display: flex; justify-content: center; gap: 8px; }
.action-btn { display: flex; align-items: center; gap: 4px; padding: 4px 8px; border-radius: 4px; border: none; cursor: pointer; font-size: 12px; transition: all 0.2s;
  &.run { background: #e6f7ff; color: #1890ff; &:hover { background: #bae7ff; } }
  &.edit { background: #f6ffed; color: #52c41a; &:hover { background: #d9f7be; } }
  &.delete { background: #fff2f0; color: #ff4d4f; &:hover { background: #ffccc7; } }
}
</style>
