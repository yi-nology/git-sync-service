<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title-light">仓库管理</h1>
      <div class="header-actions">
        <button class="btn-primary-light" @click="openCreate">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加仓库
        </button>
      </div>
    </div>

    <div class="repo-list">
      <div class="repo-card" v-for="repo in dataList" :key="repo.id">
        <div class="repo-header">
          <div class="repo-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
              <polyline points="13 2 13 9 20 9"/>
            </svg>
          </div>
          <div class="repo-info">
            <div class="repo-name">{{ repo.name }}</div>
            <div class="repo-url">{{ repo.url }}</div>
          </div>
          <span class="badge" :class="{ active: repo.syncEnabled }">
            {{ repo.syncEnabled ? '已同步' : '未同步' }}
          </span>
        </div>
        <div class="repo-configs">
          <div class="config-row"><div class="config-label">源平台</div><div class="config-value">{{ repo.source }}</div></div>
          <div class="config-row"><div class="config-label">目标平台</div><div class="config-value">{{ repo.target }}</div></div>
          <div class="config-row"><div class="config-label">同步模式</div><div class="config-value">{{ repo.mode }}</div></div>
          <div class="config-row"><div class="config-label">上次同步</div><div class="config-value">{{ repo.lastSync }}</div></div>
        </div>
       <div class="repo-actions">
         <router-link :to="`/local-repos/${repo.id}?tab=sync`" class="btn-default-light">查看同步</router-link>
         <router-link :to="`/repos/config/${repo.id}`" class="btn-default-light">查看配置</router-link>
         <button class="btn-default-light" @click="openEdit(repo)">编辑</button>
          <button class="btn-default-light danger" @click="handleDelete(repo.id!)">删除</button>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" label-width="100px">
        <el-form-item label="仓库名称">
          <el-input v-model="formData.name" placeholder="请输入仓库名称"/>
        </el-form-item>
        <el-form-item label="仓库地址">
          <el-input v-model="formData.url" placeholder="请输入仓库地址"/>
        </el-form-item>
        <el-form-item label="源平台">
          <el-select v-model="formData.source" style="width: 100%">
            <el-option label="GitHub" value="GitHub"/>
            <el-option label="GitLab" value="GitLab"/>
            <el-option label="Gitee" value="Gitee"/>
            <el-option label="GitLab 自建" value="GitLab 自建"/>
            <el-option disabled label="Bitbucket (开发中)" value="Bitbucket"/>
            <el-option disabled label="Azure DevOps (开发中)" value="Azure DevOps"/>
            <el-option disabled label="Coding.net (开发中)" value="Coding.net"/>
          </el-select>
        </el-form-item>
        <el-form-item label="目标平台">
          <el-select v-model="formData.target" style="width: 100%">
            <el-option label="GitLab" value="GitLab"/>
            <el-option label="GitHub" value="GitHub"/>
            <el-option label="Gitee" value="Gitee"/>
            <el-option label="GitLab 自建" value="GitLab 自建"/>
            <el-option disabled label="Bitbucket (开发中)" value="Bitbucket"/>
            <el-option disabled label="Azure DevOps (开发中)" value="Azure DevOps"/>
            <el-option disabled label="Coding.net (开发中)" value="Coding.net"/>
          </el-select>
        </el-form-item>
        <el-form-item label="同步模式">
          <el-select v-model="formData.mode" style="width: 100%">
            <el-option label="单分支同步" value="单分支同步"/>
            <el-option label="全分支同步" value="全分支同步"/>
          </el-select>
        </el-form-item>
        <el-form-item label="启用同步">
          <el-switch v-model="formData.syncEnabled"/>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeDialog">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCrud } from '@/composables/useCrud'

interface Repo {
  id?: number
  name: string
  url: string
  source: string
  target: string
  mode: string
  syncEnabled: boolean
  lastSync: string
}

const {
  dataList,
  dialogVisible,
  dialogTitle,
  formData,
  openCreate,
  openEdit,
  handleDelete,
  handleSubmit,
  closeDialog
 } = useCrud<Repo>({
   name: '仓库',
   fetchApi: async () => [
     { id: 1, name: 'owner/repo1', url: 'https://github.com/owner/repo1', source: 'GitHub', target: 'GitLab', mode: '单分支同步', syncEnabled: true, lastSync: '2分钟前' },
     { id: 2, name: 'group/project', url: 'https://gitlab.com/group/project', source: 'GitLab', target: 'Gitee', mode: '全分支同步', syncEnabled: true, lastSync: '15分钟前' },
     { id: 3, name: 'org/docs', url: 'https://gitee.com/org/docs', source: 'Gitee', target: 'GitHub', mode: '单分支同步', syncEnabled: true, lastSync: '1小时前' },
     { id: 4, name: 'internal/app', url: 'https://gitlab.internal.com/group/app', source: 'GitLab 自建', target: 'Gitee', mode: '单分支同步', syncEnabled: true, lastSync: '30分钟前' },
   ]
 })

onMounted(() => {
  // fetchData()
})
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.btn-primary-light { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.2s; &:hover { background: #40a9ff; } }
.btn-default-light { padding: 6px 12px; border-radius: 4px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 12px; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; &:hover { color: #1890ff; border-color: #1890ff; } &.danger:hover { color: #ff4d4f; border-color: #ff4d4f; } }
.repo-list { display: flex; flex-direction: column; gap: 16px; }
.repo-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px; }
.repo-header { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid #f0f0f0; }
.repo-icon { width: 40px; height: 40px; border-radius: 8px; background: #e6f7ff; color: #1890ff; display: flex; align-items: center; justify-content: center; }
.repo-info { flex: 1; }
.repo-name { font-size: 15px; font-weight: 600; color: #262626; }
.repo-url { font-size: 12px; color: #8c8c8c; margin-top: 4px; }
.badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block; background: #f5f5f5; color: #8c8c8c; &.active { background: #f6ffed; color: #52c41a; } }
.repo-configs { margin-bottom: 16px; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.config-row { display: flex; gap: 8px; }
.config-label { font-size: 13px; color: #8c8c8c; width: 80px; }
.config-value { font-size: 13px; color: #262626; font-weight: 500; }
.repo-actions { display: flex; gap: 8px; padding-top: 16px; border-top: 1px solid #f0f0f0; }
</style>
