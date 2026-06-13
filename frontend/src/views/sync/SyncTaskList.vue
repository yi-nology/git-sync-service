<template><div class="page-container">
   <div class="page-header">
     <h1 class="page-title-light">同步任务</h1>
     <div class="header-actions">
       <button class="btn-primary-light" @click="openCreate">
         <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
           <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
         </svg>创建任务
       </button>
     </div>
   </div>

  <div class="stats-row-light">
    <div class="stat-card-light" v-for="stat in stats" :key="stat.name">
      <div class="stat-card-inner">
        <div class="stat-content-light">
          <div class="stat-num-light" :style="{color: stat.color}">{{ stat.num }}</div>
          <div class="stat-name-light">{{ stat.name }}</div>
        </div>
        <div class="stat-icon-wrap" :style="{background: stat.bg}">
          <svg v-html="stat.icon" width="24" height="24" viewBox="0 0 24 24" fill="none" :stroke="stat.color" stroke-width="2"></svg>
        </div>
      </div>
    </div>
  </div>

  <div class="filter-bar-light">
    <div class="filter-input-light" style="width:280px;">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#8c8c8c" stroke-width="2">
        <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <input type="text" placeholder="搜索任务名称..." v-model="filters.keyword">
    </div>
    <div class="filter-select-light"><span>状态</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8c8c8c" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg></div>
    <div class="filter-select-light"><span>仓库</span><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#8c8c8c" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg></div>
  </div>

  <div class="table-card">
    <table class="sync-table">
      <thead><tr>
        <th style="text-align:left;width:180px;">任务名称</th>
        <th style="text-align:left;">源分支</th>
        <th style="text-align:left;">目标分支</th>
        <th style="text-align:left;">同步模式</th>
        <th style="text-align:left;width:100px;">状态</th>
        <th style="text-align:left;width:120px;">最后运行</th>
        <th style="text-align:center;width:120px;">操作</th>
      </tr></thead>
      <tbody><tr v-for="task in tasks" :key="task.id">
        <td><span class="task-name-text">{{ task.name }}</span></td>
        <td><span class="branch-tag">{{ task.source }}</span></td>
        <td><span class="branch-tag">{{ task.target }}</span></td>
        <td><span class="mode-text">{{ task.mode }}</span></td>
        <td><span class="status-dot-light" :class="task.status"></span>{{ task.statusText }}</td>
        <td class="time-col">{{ task.lastRun }}</td>
       <td class="action-col">
            <button class="action-btn edit" title="立即同步" @click="handleSync(task)">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>
              <span>运行</span>
            </button>
            <button class="action-btn run" title="编辑" @click="openEdit(task)">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              <span>编辑</span>
            </button>
            <button class="action-btn delete" title="删除" @click="handleDelete(task.id)">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
              <span>删除</span>
            </button>
          </td>
      </tr></tbody>
    </table>
    <div class="table-footer">
      <span class="total-text">共 12 条任务</span>
      <div class="pager">
        <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg></button>
        <button class="page-btn active">1</button><button class="page-btn">2</button>
        <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg></button>
      </div>
     </div>
   </div>

    <DeleteConfirmModal 
      v-model="deleteModal.visible" 
      :task-name="deleteModal.taskName"
      @confirm="confirmDelete"
    />
    
    <EditTaskModal
      v-model="editModal.visible"
      :task="editModal.task"
      @confirm="handleEditSubmit"
    />
 </div></template>

  <script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import DeleteConfirmModal from '@/components/DeleteConfirmModal.vue'
import EditTaskModal from '@/components/EditTaskModal.vue'

const filters = reactive({ keyword: '' })

const tasks = ref([
 { id: 1, name: 'frontend-sync', source: 'GitHub / main', target: 'GitLab / main', mode: '单分支', status: 'success', statusText: '成功', lastRun: '2 分钟前' },
 { id: 2, name: 'backend-api-sync', source: 'GitLab / develop', target: 'Gitee / develop', mode: '单分支', status: 'running', statusText: '同步中', lastRun: '运行中...' },
 { id: 3, name: 'docs-mirror', source: 'GitHub / master', target: 'Gitee / master', mode: '单分支', status: 'error', statusText: '失败', lastRun: '15 分钟前' },
 { id: 4, name: 'config-backup', source: 'GitLab / main', target: 'GitHub / main', mode: '单分支', status: 'stopped', statusText: '已停止', lastRun: '1 小时前' },
])

const stats = ref([
   { name: '总任务数', num: 12, color: '#1890ff', bg: '#e6f7ff', icon: '<circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>' },
   { name: '运行中', num: 8, color: '#52c41a', bg: '#f6ffed', icon: '<polygon points="5 3 19 12 5 21 5 3"/>' },
   { name: '今日同步', num: 156, color: '#faad14', bg: '#fffbe6', icon: '<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>' },
   { name: '失败任务', num: 1, color: '#ff4d4f', bg: '#fff2f0', icon: '<circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>' },
 ])

// 删除弹窗
const deleteModal = reactive({
  visible: false,
  taskName: '',
  deletingId: null as number | null,
})

// 编辑弹窗
const editModal = reactive({
  visible: false,
  task: null as any,
})

function openCreate() {
  editModal.task = null
  editModal.visible = true
}

function openEdit(task: any) {
  editModal.task = task
  editModal.visible = true
}

function showDeleteModal(task: any) {
  deleteModal.taskName = task.name
  deleteModal.deletingId = task.id
  deleteModal.visible = true
}

function confirmDelete() {
  if (deleteModal.deletingId) {
    tasks.value = tasks.value.filter(t => t.id !== deleteModal.deletingId)
    ElMessage.success('删除成功')
  }
  deleteModal.visible = false
  deleteModal.deletingId = null
}

function handleSync(task: any) {
  ElMessage.info(`开始同步: ${task.name}`)
}

function handleEditSubmit(data: any) {
  if (editModal.task?.id) {
    const index = tasks.value.findIndex(t => t.id === editModal.task.id)
    if (index > -1) {
      tasks.value[index] = { ...tasks.value[index], ...data }
    }
    ElMessage.success('更新成功')
  } else {
    tasks.value.unshift({
      id: Date.now(),
      ...data,
      status: 'stopped',
      statusText: '已停止',
      lastRun: '未运行'
    })
    ElMessage.success('创建成功')
  }
}
 </script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.btn-primary-light { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.2s; &:hover { background: #40a9ff; } }
.stats-row-light { display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px; margin-bottom: 20px; }
.stat-card-light { background: #fff; border-radius: 8px; height: 100px; padding: 0 20px; display: flex; align-items: center; }
.stat-card-inner { display: flex; align-items: center; justify-content: space-between; width: 100%; }
.stat-num-light { font-size: 28px; font-weight: 700; line-height: 1.2; margin-bottom: 4px; }
.stat-name-light { font-size: 13px; color: #8c8c8c; }
.stat-icon-wrap { width: 48px; height: 48px; border-radius: 24px; display: flex; align-items: center; justify-content: center; }
.filter-bar-light { display: flex; gap: 12px; margin-bottom: 16px; }
.filter-input-light { background: #fff; border: 1px solid #d9d9d9; border-radius: 6px; padding: 0 12px; height: 32px; display: flex; align-items: center; gap: 8px; input { background: transparent; border: none; outline: none; font-size: 13px; color: #262626; width: 100%; } }
.filter-select-light { background: #fff; border: 1px solid #d9d9d9; border-radius: 6px; padding: 0 12px; height: 32px; display: flex; align-items: center; gap: 8px; font-size: 13px; color: #262626; }
.table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
.sync-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
.task-name-text { font-weight: 500; }
.branch-tag { padding: 5px 12px; background: #f5f5f5; border-radius: 4px; font-size: 12px; }
.mode-text { color: #8c8c8c; }
.status-dot-light { display: inline-block; width: 6px; height: 6px; border-radius: 50%; margin-right: 8px; &.success { background: #52c41a; } &.running { background: #1890ff; } &.error { background: #ff4d4f; } &.stopped { background: #d9d9d9; } }
.time-col { color: #8c8c8c; }
.action-col { display: flex; justify-content: center; gap: 8px; }
.action-btn { 
  display: flex; 
  align-items: center; 
  gap: 4px; 
  padding: 4px 8px; 
  border-radius: 4px; 
  border: none; 
  cursor: pointer; 
  font-size: 12px; 
  font-family: 'Inter';
  transition: all 0.2s;
  
  &.edit { 
    background: #e6f7ff; 
    color: #1890ff; 
    &:hover { background: #bae7ff; }
  }
  &.run { 
    background: #f6ffed; 
    color: #52c41a; 
    &:hover { background: #d9f7be; }
  }
  &.delete { 
    background: #fff2f0; 
    color: #ff4d4f; 
    &:hover { background: #ffccc7; }
  }
}
.table-footer { display: flex; justify-content: space-between; align-items: center; padding: 0 24px; height: 48px; border-top: 1px solid #f0f0f0; }
.total-text { font-size: 13px; color: #8c8c8c; }
.pager { display: flex; gap: 8px; }
.page-btn { width: 28px; height: 28px; border-radius: 4px; border: 1px solid #d9d9d9; background: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; color: #595959; &.active { background: #1890ff; border-color: #1890ff; color: #fff; } }
</style>
