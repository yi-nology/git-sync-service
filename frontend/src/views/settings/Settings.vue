<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
    </div>

    <div class="settings-card">
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">启用增量同步</div>
          <div class="setting-desc">只同步变更的分支，提高效率</div>
        </div>
        <a-switch v-model:checked="settings.incremental"/>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">启用分布式锁</div>
          <div class="setting-desc">多实例部署时防止并发冲突</div>
        </div>
        <a-switch v-model:checked="settings.distLock"/>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">同步超时时间</div>
          <div class="setting-desc">单位：秒</div>
        </div>
        <a-input-number v-model:value="settings.timeout" :min="10" :max="600" style="width: 120px;"/>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">最大并发同步数</div>
          <div class="setting-desc">同时执行的同步任务数量</div>
        </div>
        <a-input-number v-model:value="settings.maxConcurrent" :min="1" :max="10" style="width: 120px;"/>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">历史记录保留天数</div>
          <div class="setting-desc">超过此天数的历史记录将被清理</div>
        </div>
        <a-input-number v-model:value="settings.retentionDays" :min="1" :max="365" style="width: 120px;"/>
      </div>
    </div>

    <div style="margin-top: 24px;">
      <a-button type="primary" @click="saveSettings">保存设置</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { message } from 'ant-design-vue'

const settings = reactive({
  incremental: true,
  distLock: true,
  timeout: 300,
  maxConcurrent: 3,
  retentionDays: 30,
})

function saveSettings() {
  message.success('设置已保存')
}
</script>

<style scoped lang="scss">
.page-container { background: #f0f2f5; min-height: 100%; padding: 24px; }
.page-header { margin-bottom: 24px; }
.page-title { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
.settings-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; }
.setting-item { display: flex; justify-content: space-between; align-items: center; padding: 16px 24px; border-bottom: 1px solid #f0f0f0; &:last-child { border-bottom: none; } }
.setting-info { .setting-name { font-size: 14px; font-weight: 500; color: #262626; } .setting-desc { font-size: 12px; color: #8c8c8c; margin-top: 4px; } }
</style>
