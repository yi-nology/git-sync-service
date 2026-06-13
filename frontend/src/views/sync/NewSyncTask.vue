<template><div class="page-container">
  <div class="page-header"><h1 class="page-title-light">新建同步任务</h1></div>
  <div class="steps-bar">
    <div class="step" :class="{active: step===1, done: step>1}"><span class="step-num">1</span><span class="step-label">基本配置</span></div>
    <div class="step-line" :class="{active: step>1}"></div>
    <div class="step" :class="{active: step===2, done: step>2}"><span class="step-num">2</span><span class="step-label">同步规则</span></div>
    <div class="step-line" :class="{active: step>2}"></div>
    <div class="step" :class="{active: step===3}"><span class="step-num">3</span><span class="step-label">高级选项</span></div>
  </div>
  <div class="form-card" style="margin-top:24px;">
    <div class="form-body" style="padding:24px;">
      <div v-if="step===1">
        <div class="form-item"><label>任务名称</label><input type="text" class="form-input" v-model="form.name" placeholder="请输入任务名称"></div>
        <div class="form-item"><label>源 Git 平台</label><select class="form-select" v-model="form.source"><option>GitHub</option><option>GitLab</option><option>Gitee</option></select></div>
        <div class="form-item"><label>目标 Git 平台</label><select class="form-select" v-model="form.target"><option>GitLab</option><option>GitHub</option><option>Gitee</option></select></div>
        <div class="form-item"><label>同步模式</label><div class="radio-group"><label class="radio"><input type="radio" v-model="form.mode" value="single"> 单分支同步</label><label class="radio"><input type="radio" v-model="form.mode" value="all"> 全分支同步</label></div></div>
      </div>
      <div v-if="step===2"><div class="form-item"><label>源分支</label><input type="text" class="form-input" v-model="form.sourceBranch" value="main"></div><div class="form-item"><label>目标分支</label><input type="text" class="form-input" v-model="form.targetBranch" value="main"></div></div>
      <div v-if="step===3"><div class="form-item"><label>同步所有标签</label><div class="switch-btn" :class="{active: form.syncTags}" @click="form.syncTags=!form.syncTags"><span class="switch-dot"></span></div></div><div class="form-item"><label>强制推送</label><div class="switch-btn" :class="{active: form.force}" @click="form.force=!form.force"><span class="switch-dot"></span></div></div></div>
    </div>
    <div class="form-footer"><button class="btn-default-light" v-if="step>1" @click="step--">上一步</button><button class="btn-primary-light" v-if="step<3" @click="step++">下一步</button><button class="btn-primary-light" style="background:#52c41a;" v-if="step===3" @click="submit">创建任务</button></div>
  </div>
</div></template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
const step = ref(1)
const form = reactive({ name: '', source: 'GitHub', target: 'GitLab', mode: 'single', sourceBranch: 'main', targetBranch: 'main', syncTags: false, force: false })
function submit() { alert('创建成功!') }
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; }
.page-header { margin-bottom: 24px; }
.page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
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
.form-input { width: 100%; height: 40px; padding: 0 12px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 14px; }
.form-select { width: 100%; height: 40px; padding: 0 12px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 14px; background: #fff; }
.btn-default-light { padding: 7px 20px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #262626; font-size: 14px; cursor: pointer; }
.btn-primary-light { padding: 7px 20px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 14px; cursor: pointer; }
.radio-group { display: flex; gap: 24px; }
.radio { display: flex; align-items: center; gap: 8px; font-size: 14px; color: #262626; }
.switch-btn { width: 44px; height: 24px; border-radius: 12px; background: #d9d9d9; cursor: pointer; transition: all 0.2s; position: relative; &.active { background: #1890ff; } .switch-dot { position: absolute; left: 2px; top: 2px; width: 20px; height: 20px; border-radius: 50%; background: #fff; transition: all 0.2s; } &.active .switch-dot { transform: translateX(20px); } }
</style>
