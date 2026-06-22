<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">高级配置</h1>
    </div>

    <div class="config-card">
      <div class="section-title">基本配置</div>
      <div class="form-row">
        <div class="form-item">
          <label>启用增量同步</label>
          <el-switch v-model="config.incrementalSync"/>
        </div>
        <div class="form-item">
          <label>启用分布式锁</label>
          <el-switch v-model="config.distributedLock"/>
        </div>
      </div>
    </div>

    <div class="config-card">
      <div class="section-title">同步配置</div>
      <div class="form-row">
        <div class="form-item">
          <label>同步超时 (秒)</label>
          <el-input-number v-model="config.syncTimeout" :min="10" :max="600"/>
        </div>
        <div class="form-item">
          <label>最大并发同步数</label>
          <el-input-number v-model="config.maxConcurrent" :min="1" :max="10"/>
        </div>
        <div class="form-item">
          <label>重试次数</label>
          <el-input-number v-model="config.retryCount" :min="0" :max="5"/>
        </div>
      </div>
      <div class="form-row">
        <div class="form-item full-width">
          <label>排除的分支 (正则)</label>
          <el-input v-model="config.excludeBranches" placeholder="例如: ^feature/.*$"/>
        </div>
      </div>
    </div>

    <div class="config-card">
      <div class="section-title">高级选项</div>
      <div class="form-row">
        <div class="form-item">
          <label>启用 Git LFS</label>
          <el-switch v-model="config.gitLfs"/>
        </div>
        <div class="form-item">
          <label>验证提交签名</label>
          <el-switch v-model="config.verifySignature"/>
        </div>
        <div class="form-item">
          <label>启用审计日志</label>
          <el-switch v-model="config.auditLog"/>
        </div>
      </div>
      <div class="form-row">
        <div class="form-item full-width">
          <label>历史记录保留天数</label>
          <el-input-number v-model="config.retentionDays" :min="1" :max="365"/>
        </div>
      </div>
    </div>

    <div class="action-footer">
      <el-button @click="resetConfig">重置</el-button>
      <el-button type="primary" @click="saveConfig">保存配置</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { ElMessage } from 'element-plus'

const config = reactive({
  incrementalSync: true,
  distributedLock: true,
  syncTimeout: 300,
  maxConcurrent: 3,
  retryCount: 2,
  excludeBranches: '',
  gitLfs: false,
  verifySignature: false,
  auditLog: true,
  retentionDays: 30,
})

function saveConfig() {
  ElMessage.success('配置已保存')
}

function resetConfig() {
  Object.assign(config, {
    incrementalSync: true,
    distributedLock: true,
    syncTimeout: 300,
    maxConcurrent: 3,
    retryCount: 2,
    excludeBranches: '',
    gitLfs: false,
    verifySignature: false,
    auditLog: true,
    retentionDays: 30,
  })
  ElMessage.info('配置已重置')
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.config-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; padding: 20px 24px; margin-bottom: 16px; }
.section-title { font-size: 15px; font-weight: 600; color: #262626; margin-bottom: 20px; padding-bottom: 12px; border-bottom: 1px solid #f0f0f0; }
.form-row { display: flex; gap: 24px; margin-bottom: 20px; &:last-child { margin-bottom: 0; } }
.form-item { display: flex; flex-direction: column; gap: 8px; min-width: 200px; label { font-size: 14px; font-weight: 500; color: #595959; } &.full-width { flex: 1; width: 100%; } }
.action-footer { display: flex; justify-content: flex-end; gap: 12px; padding-top: 8px; }
</style>
