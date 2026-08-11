<template>
  <div class="page-container">
    <PageHeader title="系统设置" />

    <div class="content-card">
      <div class="card-body">
        <a-form layout="vertical" :model="settings">
          <a-divider orientation="left">通用设置</a-divider>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="服务名称">
                <a-input v-model:value="settings.serviceName" placeholder="Git Sync Service" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="日志级别">
                <a-select v-model:value="settings.logLevel">
                  <a-select-option value="debug">Debug</a-select-option>
                  <a-select-option value="info">Info</a-select-option>
                  <a-select-option value="warn">Warn</a-select-option>
                  <a-select-option value="error">Error</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>

          <a-divider orientation="left">通知设置</a-divider>

          <a-form-item label="Webhook URL">
            <a-input v-model:value="settings.webhookUrl" placeholder="https://..." />
          </a-form-item>

          <a-form-item>
            <a-button type="primary" @click="handleSave">保存设置</a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import PageHeader from '@/components/common/PageHeader.vue'

const settings = reactive({
  serviceName: 'Git Sync Service',
  logLevel: 'info',
  webhookUrl: '',
})

function handleSave() {
  message.success('设置已保存')
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  padding: $spacing-lg;
}

.card-body {
  max-width: 800px;
}
</style>
