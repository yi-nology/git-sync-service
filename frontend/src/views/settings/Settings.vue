<template><div class="page-container">
  <div class="page-header"><h1 class="page-title-light">系统设置</h1></div>
  <div class="settings-card" style="margin-top:24px;">
    <div class="setting-item">
      <div class="setting-info"><div class="setting-name">启用增量同步</div><div class="setting-desc">只同步变更的分支，提高效率</div></div>
      <div class="switch-btn" :class="{active: settings.incremental}" @click="settings.incremental=!settings.incremental"><span class="switch-dot"></span></div>
    </div>
    <div class="setting-item">
      <div class="setting-info"><div class="setting-name">启用分布式锁</div><div class="setting-desc">多实例部署时防止并发冲突</div></div>
      <div class="switch-btn" :class="{active: settings.distLock}" @click="settings.distLock=!settings.distLock"><span class="switch-dot"></span></div>
    </div>
    <div class="setting-item">
      <div class="setting-info"><div class="setting-name">同步超时时间</div><div class="setting-desc">单位：秒</div></div>
      <input type="number" class="setting-input" v-model="settings.timeout" style="width:120px;">
    </div>
    <div class="setting-item">
      <div class="setting-info"><div class="setting-name">最大并发同步数</div><div class="setting-desc">同时执行的同步任务数量</div></div>
      <input type="number" class="setting-input" v-model="settings.maxConcurrent" style="width:120px;">
    </div>
    <div class="setting-item">
      <div class="setting-info"><div class="setting-name">历史记录保留天数</div><div class="setting-desc">超过此天数的历史记录将被清理</div></div>
      <input type="number" class="setting-input" v-model="settings.retentionDays" style="width:120px;">
    </div>
  </div>
  <div style="margin-top:24px;"><button class="btn-primary-light">保存设置</button></div>
</div></template>

<script setup lang="ts">
import { reactive } from 'vue'
const settings = reactive({ incremental: true, distLock: true, timeout: 300, maxConcurrent: 3, retentionDays: 30 })
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; }
.page-header { margin-bottom: 24px; }
.page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.settings-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; }
.setting-item { display: flex; justify-content: space-between; align-items: center; padding: 16px 24px; border-bottom: 1px solid #f0f0f0; &:last-child { border-bottom: none; } }
.setting-info { .setting-name { font-size: 14px; font-weight: 500; color: #262626; } .setting-desc { font-size: 12px; color: #8c8c8c; margin-top: 4px; } }
.setting-input { height: 32px; padding: 0 12px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 14px; text-align: center; }
.btn-primary-light { padding: 7px 20px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 14px; cursor: pointer; }
.switch-btn { width: 44px; height: 24px; border-radius: 12px; background: #d9d9d9; cursor: pointer; transition: all 0.2s; position: relative; &.active { background: #1890ff; } .switch-dot { position: absolute; left: 2px; top: 2px; width: 20px; height: 20px; border-radius: 50%; background: #fff; transition: all 0.2s; } &.active .switch-dot { transform: translateX(20px); } }
</style>
