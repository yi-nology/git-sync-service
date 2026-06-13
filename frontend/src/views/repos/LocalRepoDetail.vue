<template>
  <div class="page-container">
    <div class="repo-header">
      <div class="repo-left">
        <div class="repo-icon-lg">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
            <polyline points="13 2 13 9 20 9"/>
          </svg>
        </div>
        <div class="repo-info-lg">
          <div class="repo-name-lg">{{ repoInfo.name }}</div>
          <div class="repo-url-lg">{{ repoInfo.path }}</div>
        </div>
        <span class="badge-status" :class="repoInfo.status">{{ repoInfo.statusText }}</span>
      </div>
      <div class="header-actions-lg">
        <button class="btn-default-lg" @click="refreshRepo">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"/>
            <polyline points="1 20 1 14 7 14"/>
            <path d="M3.51 9a9 9 0 0 1 14.85 3.51L23 10"/>
            <path d="M21 14a9 9 0 0 1-14.85 3.51L1 14"/>
          </svg>
          刷新
        </button>
        <button class="btn-primary-lg" @click="syncNow">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          立即同步
        </button>
      </div>
    </div>

    <div class="tabs-bar">
      <button class="tab-btn" :class="{active: activeTab === 'sync'}" @click="changeTab('sync')">同步任务</button>
      <button class="tab-btn" :class="{active: activeTab === 'config'}" @click="changeTab('config')">同步配置</button>
      <button class="tab-btn" :class="{active: activeTab === 'webhook'}" @click="changeTab('webhook')">Webhook</button>
      <button class="tab-btn" :class="{active: activeTab === 'history'}" @click="changeTab('history')">同步历史</button>
    </div>

    <div v-if="activeTab === 'sync'" class="config-card">
      <div class="card-title">同步任务列表</div>
      <div class="task-list">
        <div class="task-item" v-for="task in syncTasks" :key="task.id">
          <div class="task-info">
            <div class="task-name">{{ task.name }}</div>
            <div class="task-meta">
              <span>{{ task.sourceBranch }}</span>
              <svg class="arrow-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 18 15 12 9 6"/>
              </svg>
              <span>{{ task.targetBranch }}</span>
            </div>
          </div>
          <div class="task-tags">
            <span class="tag" :class="task.mode">{{ task.modeText }}</span>
            <span class="tag" :class="task.status">{{ task.statusText }}</span>
            <span class="tag info">{{ task.lastRun }}</span>
          </div>
          <div class="task-actions">
            <button class="icon-btn" @click="runTask(task)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
            </button>
            <button class="icon-btn" @click="editTask(task)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
      <div class="card-footer">
        <button class="btn-default-light" @click="openNewTask">+ 新建同步任务</button>
      </div>
    </div>

    <div v-if="activeTab === 'config'" class="config-card">
      <div class="card-title">同步配置</div>
      <div class="switch-group">
        <div class="switch-item" v-for="item in configItems" :key="item.key">
          <div class="switch-info">
            <div class="switch-name">{{ item.name }}</div>
            <div class="switch-desc">{{ item.desc }}</div>
          </div>
          <div class="switch-btn" :class="{active: getConfigValue(item.key)}" @click="setConfigValue(item.key)">
            <span class="switch-dot"></span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'webhook'" class="config-card">
      <div class="card-title">Webhook 配置</div>
      <div class="webhook-info">
        <div class="webhook-url-wrap">
          <label>Webhook 地址</label>
          <div class="url-input-group">
            <input type="text" readonly :value="webhookUrl" class="url-input"/>
            <button class="copy-btn" @click="copyUrl">复制</button>
          </div>
        </div>
        <div class="webhook-secret">
          <label>密钥 (Secret)</label>
          <div class="url-input-group">
            <input type="text" readonly :value="webhookSecret" class="url-input"/>
            <button class="copy-btn" @click="generateSecret">重新生成</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'history'" class="config-card">
      <div class="card-title">同步历史记录</div>
      <div class="history-list">
        <div class="history-item" v-for="record in historyList" :key="record.id">
          <div class="history-icon" :class="record.status">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline v-if="record.status === 'success'" points="22 4 12 14.01 9 11.01"/>
              <circle v-if="record.status === 'error'" cx="12" cy="12" r="10"/>
              <line v-if="record.status === 'error'" x1="15" y1="9" x2="9" y2="15"/>
              <line v-if="record.status === 'error'" x1="9" y1="9" x2="15" y2="15"/>
              <polygon v-if="record.status === 'running'" points="5 3 19 12 5 21 5 3"/>
            </svg>
          </div>
          <div class="history-info">
            <div class="history-branch">{{ record.branch }}</div>
            <div class="history-meta">{{ record.source }} → {{ record.target }} · {{ record.duration }}</div>
          </div>
          <span class="history-status" :class="record.status">{{ record.statusText }}</span>
          <span class="history-time">{{ record.time }}</span>
        </div>
      </div>
    </div>

    <el-dialog v-model="taskDialogVisible" title="新建同步任务" width="500px">
      <el-form :model="taskForm" label-width="100px">
        <el-form-item label="任务名称">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称"/>
        </el-form-item>
        <el-form-item label="源分支">
          <el-input v-model="taskForm.sourceBranch" placeholder="例如: main"/>
        </el-form-item>
        <el-form-item label="目标分支">
          <el-input v-model="taskForm.targetBranch" placeholder="例如: main"/>
        </el-form-item>
        <el-form-item label="同步模式">
          <el-select v-model="taskForm.mode" style="width: 100%">
            <el-option label="单分支同步" value="single"/>
            <el-option label="全分支同步" value="full"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taskDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTask">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const repoId = ref('')
const activeTab = ref('sync')

const repoInfo = reactive({
  name: 'frontend-app',
  path: '/Users/zhangyi/projects/frontend-app',
  status: 'success',
  statusText: '已同步'
})

const config = reactive({
  incremental: true,
  distLock: true,
  syncTags: false,
  forcePush: false
})

const configItems = [
  { key: 'incremental', name: '启用增量同步', desc: '只同步变更的分支，提高效率' },
  { key: 'distLock', name: '启用分布式锁', desc: '多实例部署时防止并发冲突' },
  { key: 'syncTags', name: '同步所有标签', desc: '同步 Git tags 到目标仓库' },
  { key: 'forcePush', name: '强制推送', desc: '使用 --force 覆盖远程分支' }
]

const webhookUrl = ref('')
const webhookSecret = ref('abc123xyz789secretkey001')

const syncTasks = ref([
  { id: 1, name: 'main 分支同步', sourceBranch: 'main', targetBranch: 'main', mode: 'single', modeText: '单分支', status: 'success', statusText: '成功', lastRun: '2分钟前' },
  { id: 2, name: 'develop 分支同步', sourceBranch: 'develop', targetBranch: 'develop', mode: 'single', modeText: '单分支', status: 'running', statusText: '同步中', lastRun: '进行中' }
])

const historyList = ref([
  { id: 1, branch: 'main', source: 'GitHub', target: 'GitLab', status: 'success', statusText: '成功', duration: '45秒', time: '2分钟前' },
  { id: 2, branch: 'develop', source: 'GitHub', target: 'GitLab', status: 'success', statusText: '成功', duration: '38秒', time: '15分钟前' },
  { id: 3, branch: 'feature/auth', source: 'GitHub', target: 'GitLab', status: 'error', statusText: '失败', duration: '1分20秒', time: '1小时前' }
])

const taskDialogVisible = ref(false)
const taskForm = reactive({ name: '', sourceBranch: 'main', targetBranch: 'main', mode: 'single' })

function changeTab(tab: string) {
  activeTab.value = tab
  router.replace({ query: { tab } })
}

function syncNow() {
  ElMessage.info('开始同步...')
  repoInfo.status = 'running'
  repoInfo.statusText = '同步中'
}

function refreshRepo() {
  ElMessage.success('刷新成功')
}

function runTask(task: any) {
  ElMessage.info(`开始执行任务: ${task.name}`)
}

function editTask(task: any) {
  Object.assign(taskForm, task)
  taskDialogVisible.value = true
}

function openNewTask() {
  Object.assign(taskForm, { name: '', sourceBranch: 'main', targetBranch: 'main', mode: 'single' })
  taskDialogVisible.value = true
}

function saveTask() {
  if (!taskForm.name) {
    ElMessage.warning('请输入任务名称')
    return
  }
  syncTasks.value.unshift({
    id: Date.now(),
    ...taskForm,
    modeText: taskForm.mode === 'single' ? '单分支' : '全分支',
    status: 'stopped',
    statusText: '已停止',
    lastRun: '从未运行'
  })
  ElMessage.success('任务已创建')
  taskDialogVisible.value = false
}

function copyUrl() {
  navigator.clipboard.writeText(webhookUrl.value)
  ElMessage.success('已复制到剪贴板')
}

function generateSecret() {
  webhookSecret.value = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
  ElMessage.success('已重新生成密钥')
}

function getConfigValue(key: string): boolean {
  return (config as any)[key]
}

function setConfigValue(key: string): void {
  (config as any)[key] = !(config as any)[key]
}

onMounted(() => {
  repoId.value = route.params.id as string
  webhookUrl.value = `http://localhost:12345/webhook/git/${repoId.value}`
  
  if (route.query.tab) {
    activeTab.value = route.query.tab as string
  }
})
</script>

<style scoped lang="scss">
.page-container {
  background: #f0f2f5;
  min-height: 100%;
}

.repo-header {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .repo-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .repo-icon-lg {
    width: 56px;
    height: 56px;
    border-radius: 8px;
    background: #e6f7ff;
    color: #1890ff;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .repo-info-lg {
    .repo-name-lg {
      font-size: 18px;
      font-weight: 600;
      color: #262626;
    }
    .repo-url-lg {
      font-size: 13px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }

  .badge-status {
    padding: 5px 12px;
    border-radius: 4px;
    font-size: 13px;
    font-weight: 500;
    margin-left: 8px;

    &.success {
      background: #f6ffed;
      color: #52c41a;
    }
    &.running {
      background: #e6f7ff;
      color: #1890ff;
    }
    &.error {
      background: #fff2f0;
      color: #ff4d4f;
    }
  }

  .header-actions-lg {
    display: flex;
    gap: 12px;
  }

  .btn-default-lg {
    padding: 8px 16px;
    border-radius: 6px;
    background: #fff;
    border: 1px solid #d9d9d9;
    color: #595959;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;

    &:hover {
      color: #1890ff;
      border-color: #1890ff;
    }
  }

  .btn-primary-lg {
    padding: 8px 16px;
    border-radius: 6px;
    background: #1890ff;
    border: none;
    color: #fff;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;

    &:hover {
      background: #40a9ff;
    }
  }
}

.tabs-bar {
  display: flex;
  gap: 4px;
  margin-bottom: 16px;
  background: #fff;
  border-radius: 8px;
  padding: 4px;
  border: 1px solid #f0f0f0;
}

.tab-btn {
  padding: 8px 20px;
  border: none;
  background: transparent;
  border-radius: 6px;
  font-size: 14px;
  color: #595959;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    color: #1890ff;
  }

  &.active {
    background: #1890ff;
    color: #fff;
  }
}

.config-card {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;
  margin-bottom: 16px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.card-footer {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.btn-default-light {
  padding: 7px 14px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  color: #595959;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    color: #1890ff;
    border-color: #1890ff;
  }
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 6px;

  .task-info {
    .task-name {
      font-size: 14px;
      font-weight: 500;
      color: #262626;
    }
    .task-meta {
      font-size: 12px;
      color: #8c8c8c;
      margin-top: 4px;
      display: flex;
      gap: 8px;
      align-items: center;
    }
  }

  .task-tags {
    display: flex;
    gap: 8px;
    margin-right: 16px;

    .tag {
      padding: 4px 10px;
      border-radius: 4px;
      font-size: 12px;

      &.single {
        background: #fff7e6;
        color: #fa8c16;
      }
      &.full {
        background: #f6ffed;
        color: #52c41a;
      }
      &.success {
        background: #f6ffed;
        color: #52c41a;
      }
      &.running {
        background: #e6f7ff;
        color: #1890ff;
      }
      &.stopped {
        background: #f5f5f5;
        color: #8c8c8c;
      }
      &.info {
        background: #f5f5f5;
        color: #8c8c8c;
      }
    }
  }

  .task-actions {
    display: flex;
    gap: 8px;

    .icon-btn {
      width: 28px;
      height: 28px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;
      background: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      color: #595959;

      &:hover {
        color: #1890ff;
        border-color: #1890ff;
      }
    }
  }
}

.switch-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.switch-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;

  &:not(:last-child) {
    border-bottom: 1px solid #fafafa;
  }

  .switch-info {
    .switch-name {
      font-size: 14px;
      font-weight: 500;
      color: #262626;
    }
    .switch-desc {
      font-size: 12px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }
}

.switch-btn {
  width: 44px;
  height: 24px;
  border-radius: 12px;
  background: #d9d9d9;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  &.active {
    background: #52c41a;
  }

  .switch-dot {
    position: absolute;
    left: 2px;
    top: 2px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #fff;
    transition: all 0.2s;
  }

  &.active .switch-dot {
    transform: translateX(20px);
  }
}

.webhook-info {
  display: flex;
  flex-direction: column;
  gap: 16px;

  label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #595959;
    margin-bottom: 8px;
  }

  .url-input-group {
    display: flex;
    gap: 8px;

    .url-input {
      flex: 1;
      height: 36px;
      padding: 0 12px;
      border: 1px solid #d9d9d9;
      border-radius: 6px;
      font-size: 13px;
      color: #262626;
      background: #fafafa;
    }

    .copy-btn {
      padding: 0 16px;
      height: 36px;
      border-radius: 6px;
      background: #fff;
      border: 1px solid #d9d9d9;
      color: #595959;
      font-size: 13px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        color: #1890ff;
        border-color: #1890ff;
      }
    }
  }
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 6px;

  .history-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    &.success {
      background: #f6ffed;
      color: #52c41a;
    }
    &.error {
      background: #fff2f0;
      color: #ff4d4f;
    }
    &.running {
      background: #e6f7ff;
      color: #1890ff;
    }
  }

  .history-info {
    flex: 1;

    .history-branch {
      font-size: 14px;
      font-weight: 500;
      color: #262626;
    }
    .history-meta {
      font-size: 12px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }

  .history-status {
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 12px;

    &.success {
      background: #f6ffed;
      color: #52c41a;
    }
    &.error {
      background: #fff2f0;
      color: #ff4d4f;
    }
    &.running {
      background: #e6f7ff;
      color: #1890ff;
    }
  }

  .history-time {
    font-size: 12px;
    color: #8c8c8c;
  }
}
</style>
