<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">仓库管理</h1>
      <div class="header-actions">
        <button class="btn-primary" @click="openCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加仓库
        </button>
      </div>
    </div>

    <div v-if="repoStore.loading" class="loading-state">加载中...</div>
    <div v-else-if="repoStore.repos.length === 0" class="empty-state">暂无仓库数据</div>
    <div v-else class="repo-list">
      <div class="repo-card" v-for="repo in repoStore.repos" :key="repo.key">
        <div class="repo-header">
          <div class="repo-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
              <polyline points="13 2 13 9 20 9"/>
            </svg>
          </div>
          <div class="repo-info">
            <div class="repo-name">{{ repo.name }}</div>
            <div class="repo-url">{{ repo.clone_url }}</div>
          </div>
          <span class="badge" :class="repo.status">
            {{ repo.status === 'active' ? '活跃' : '停用' }}
          </span>
        </div>
        <div class="repo-configs">
          <div class="config-row"><div class="config-label">平台</div><div class="config-value">{{ repo.platform }}</div></div>
          <div class="config-row"><div class="config-label">所有者</div><div class="config-value">{{ repo.platform_owner }}</div></div>
          <div class="config-row"><div class="config-label">仓库</div><div class="config-value">{{ repo.platform_repo }}</div></div>
          <div class="config-row"><div class="config-label">默认分支</div><div class="config-value">{{ repo.default_branch }}</div></div>
        </div>
        <div class="repo-actions">
          <button class="btn-default" @click="testConn(repo.key)">测试连接</button>
          <button class="btn-default" @click="openEdit(repo)">编辑</button>
          <button class="btn-danger" @click="handleDelete(repo.key)">删除</button>
        </div>
      </div>
    </div>

    <a-modal v-model:open="dialogVisible" :title="dialogTitle" :width="500" @ok="handleSubmit" okText="确定" cancelText="取消">
      <a-form :model="formData" :label-col="{ span: 6 }">
        <a-form-item label="仓库名称">
          <a-input v-model:value="formData.name" placeholder="请输入仓库名称"/>
        </a-form-item>
        <a-form-item label="仓库地址">
          <a-input v-model:value="formData.remote_url" placeholder="请输入仓库地址"/>
        </a-form-item>
        <a-form-item label="访问令牌">
          <a-input-password v-model:value="formData.access_token" placeholder="请输入访问令牌"/>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRepoStore } from '@/stores/repo'
import type { Repo } from '@/types'

const repoStore = useRepoStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingKey = ref('')

const formData = reactive({
  name: '',
  remote_url: '',
  access_token: '',
})

onMounted(() => {
  repoStore.fetchRepos()
})

function openCreate() {
  dialogTitle.value = '添加仓库'
  editingKey.value = ''
  Object.assign(formData, { name: '', remote_url: '', access_token: '' })
  dialogVisible.value = true
}

function openEdit(repo: Repo) {
  dialogTitle.value = '编辑仓库'
  editingKey.value = repo.key
  Object.assign(formData, { name: repo.name, remote_url: repo.clone_url, access_token: '' })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.remote_url) {
    message.warning('请填写仓库名称和地址')
    return
  }
  if (editingKey.value) {
    await repoStore.updateRepo({ key: editingKey.value, name: formData.name, access_token: formData.access_token })
  } else {
    await repoStore.createRepo(formData)
  }
  dialogVisible.value = false
}

async function handleDelete(key: string) {
  await new Promise<void>((resolve, reject) => {
    Modal.confirm({
      title: '提示',
      content: '确定要删除该仓库吗？',
      okText: '确定',
      cancelText: '取消',
      onOk: () => resolve(),
      onCancel: () => reject(new Error('cancelled')),
    })
  })
  await repoStore.deleteRepo(key)
}

async function testConn(key: string) {
  const result = await repoStore.testConnection(key)
  if (result) {
    message[result.success ? 'success' : 'error'](result.message)
  }
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.btn-primary { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; cursor: pointer; &:hover { background: #40a9ff; } }
.btn-default { padding: 6px 12px; border-radius: 4px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 12px; cursor: pointer; &:hover { color: #1890ff; border-color: #1890ff; } }
.btn-danger { padding: 6px 12px; border-radius: 4px; background: #fff; border: 1px solid #d9d9d9; color: #ff4d4f; font-size: 12px; cursor: pointer; &:hover { border-color: #ff4d4f; } }
.loading-state, .empty-state { text-align: center; padding: 48px; color: #8c8c8c; font-size: 14px; }
.repo-list { display: flex; flex-direction: column; gap: 16px; }
.repo-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px; }
.repo-header { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid #f0f0f0; }
.repo-icon { width: 40px; height: 40px; border-radius: 8px; background: #e6f7ff; color: #1890ff; display: flex; align-items: center; justify-content: center; }
.repo-info { flex: 1; }
.repo-name { font-size: 15px; font-weight: 600; color: #262626; }
.repo-url { font-size: 12px; color: #8c8c8c; margin-top: 4px; }
.badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; background: #f5f5f5; color: #8c8c8c; &.active { background: #f6ffed; color: #52c41a; } }
.repo-configs { margin-bottom: 16px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.config-row { display: flex; gap: 8px; }
.config-label { font-size: 13px; color: #8c8c8c; width: 80px; }
.config-value { font-size: 13px; color: #262626; font-weight: 500; }
.repo-actions { display: flex; gap: 8px; padding-top: 16px; border-top: 1px solid #f0f0f0; }
</style>
