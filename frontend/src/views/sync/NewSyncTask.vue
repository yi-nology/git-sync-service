<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">新建同步任务</h1>
    </div>

    <div class="steps-bar">
      <div class="step" :class="{active: step===1, done: step>1}"><span class="step-num">1</span><span class="step-label">基本配置</span></div>
      <div class="step-line" :class="{active: step>1}"></div>
      <div class="step" :class="{active: step===2, done: step>2}"><span class="step-num">2</span><span class="step-label">同步规则</span></div>
      <div class="step-line" :class="{active: step>2}"></div>
      <div class="step" :class="{active: step===3}"><span class="step-num">3</span><span class="step-label">高级选项</span></div>
    </div>

    <div class="form-card">
      <div class="form-body">
        <div v-if="step===1">
          <div class="form-item">
            <label>任务名称</label>
            <el-input v-model="form.name" placeholder="请输入任务名称"/>
          </div>
          <div class="form-item">
            <label>源仓库 Key</label>
            <el-input v-model="form.source_repo_key" placeholder="请输入源仓库 Key"/>
          </div>
          <div class="form-item">
            <label>目标仓库 Key</label>
            <el-input v-model="form.target_repo_key" placeholder="请输入目标仓库 Key"/>
          </div>
          <div class="form-item">
            <label>同步模式</label>
            <el-radio-group v-model="form.sync_mode">
              <el-radio value="single">单分支同步</el-radio>
              <el-radio value="all">全分支同步</el-radio>
            </el-radio-group>
          </div>
        </div>

        <div v-if="step===2">
          <div class="form-item">
            <label>源分支</label>
            <el-input v-model="form.source_branch" placeholder="main"/>
          </div>
          <div class="form-item">
            <label>目标分支</label>
            <el-input v-model="form.target_branch" placeholder="main"/>
          </div>
          <div class="form-item">
            <label>Cron 表达式（可选）</label>
            <el-input v-model="form.cron" placeholder="如 0 */5 * * * *"/>
          </div>
        </div>

        <div v-if="step===3">
          <div class="form-item">
            <el-checkbox v-model="form.git_tags">同步 Tags</el-checkbox>
          </div>
          <div class="form-item">
            <el-checkbox v-model="form.git_force">强制推送</el-checkbox>
          </div>
          <div class="form-item">
            <el-checkbox v-model="form.git_prune">Prune 远程已删除分支</el-checkbox>
          </div>
        </div>
      </div>

      <div class="form-footer">
        <el-button v-if="step > 1" @click="step--">上一步</el-button>
        <el-button v-if="step < 3" type="primary" @click="step++">下一步</el-button>
        <el-button v-if="step === 3" type="primary" @click="submit" :loading="submitting">创建任务</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useSyncTaskStore } from '@/stores/syncTask'

const router = useRouter()
const taskStore = useSyncTaskStore()
const step = ref(1)
const submitting = ref(false)

const form = reactive({
  name: '',
  source_repo_key: '',
  target_repo_key: '',
  sync_mode: 'single',
  source_branch: 'main',
  target_branch: 'main',
  cron: '',
  git_tags: false,
  git_force: false,
  git_prune: false,
})

async function submit() {
  if (!form.name || !form.source_repo_key || !form.target_repo_key) {
    ElMessage.warning('请填写必填字段')
    return
  }
  submitting.value = true
  const result = await taskStore.createTask(form)
  submitting.value = false
  if (result) {
    router.push('/sync')
  }
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.steps-bar { display: flex; align-items: center; justify-content: center; gap: 0; margin-bottom: 24px; }
.step { display: flex; align-items: center; gap: 8px; }
.step-num { width: 32px; height: 32px; border-radius: 50%; background: #f5f5f5; border: 2px solid #d9d9d9; color: #8c8c8c; display: flex; align-items: center; justify-content: center; font-weight: 600; font-size: 14px; }
.step.active .step-num { background: #1890ff; border-color: #1890ff; color: #fff; }
.step.done .step-num { background: #52c41a; border-color: #52c41a; color: #fff; }
.step-label { font-size: 14px; color: #8c8c8c; }
.step.active .step-label { color: #1890ff; font-weight: 500; }
.step-line { width: 60px; height: 2px; background: #d9d9d9; margin: 0 8px; }
.step-line.active { background: #52c41a; }
.form-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; }
.form-body { padding: 24px; }
.form-footer { padding: 16px 24px; border-top: 1px solid #f0f0f0; display: flex; justify-content: center; gap: 12px; }
.form-item { margin-bottom: 24px; label { display: block; font-size: 14px; font-weight: 500; color: #262626; margin-bottom: 8px; } }
</style>
