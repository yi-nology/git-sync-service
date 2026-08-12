<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">系统设置</h1>
        <p class="page-subtitle">配置系统的基本运行参数</p>
      </div>
    </div>

    <a-alert
      type="warning"
      show-icon
      message="功能开发中"
      description="本页设置项暂未接入后端，保存不会生效。相关能力正在建设中。"
      style="max-width: 800px; margin-bottom: 16px;"
    />

    <div class="settings-card">
      <a-form layout="vertical" :model="settings">
        <a-divider orientation="left">
          <span class="divider-title">通用设置</span>
        </a-divider>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="服务名称">
              <a-input v-model:value="settings.serviceName" placeholder="Git Sync Service" />
              <div class="form-tip">
                <InfoCircleOutlined /> 显示在页面标题和通知中的服务名称
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="日志级别">
              <a-select v-model:value="settings.logLevel" style="width: 100%;">
                <a-select-option value="debug">
                  <span style="font-weight: 500;">Debug</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">输出所有日志，适合调试</span>
                </a-select-option>
                <a-select-option value="info">
                  <span style="font-weight: 500;">Info</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">常规运行日志（推荐）</span>
                </a-select-option>
                <a-select-option value="warn">
                  <span style="font-weight: 500;">Warn</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">仅警告和错误</span>
                </a-select-option>
                <a-select-option value="error">
                  <span style="font-weight: 500;">Error</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">仅错误日志</span>
                </a-select-option>
              </a-select>
              <div class="form-tip">
                <InfoCircleOutlined /> 日志级别越高，输出越少。生产环境建议使用 Info 或 Warn
              </div>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">
          <span class="divider-title">通知设置</span>
        </a-divider>

        <a-form-item label="Webhook 回调 URL">
          <a-input v-model:value="settings.webhookUrl" placeholder="https://your-domain.com/webhook/callback" />
          <div class="form-tip">
            <InfoCircleOutlined /> 同步完成后会向此 URL 发送通知。支持 HTTP/HTTPS，留空则不发送通知
          </div>
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">
              <template #icon><SaveOutlined /></template>
              保存设置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { InfoCircleOutlined, SaveOutlined } from '@ant-design/icons-vue'
import { notifyInfo } from '@/utils/notify'

const settings = reactive({
  serviceName: 'Git Sync Service',
  logLevel: 'info',
  webhookUrl: '',
})

function handleSave() {
  notifyInfo('该功能正在开发中，暂未接入后端')
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
}

.page-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: $spacing-lg;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: $text-primary;
  margin: 0;
  line-height: 1.3;
}

.page-subtitle {
  font-size: 14px;
  color: $text-secondary;
  margin: 4px 0 0 0;
}

.settings-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  padding: $spacing-lg $spacing-xl;
  max-width: 800px;
}

.divider-title {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
}

.form-tip {
  font-size: 12px;
  color: $text-secondary;
  margin-top: 4px;

  .anticon {
    margin-right: 4px;
  }
}
</style>
