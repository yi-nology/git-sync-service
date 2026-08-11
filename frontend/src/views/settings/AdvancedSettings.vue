<template>
  <div class="page-container">
    <div class="page-header-bar">
      <div>
        <h1 class="page-title">高级配置</h1>
        <p class="page-subtitle">配置数据库、Git 参数和并发设置</p>
      </div>
    </div>

    <div class="settings-card">
      <a-form layout="vertical" :model="config">
        <a-divider orientation="left">
          <span class="divider-title">数据库配置</span>
        </a-divider>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="数据库驱动">
              <a-select v-model:value="config.dbDriver" style="width: 100%;">
                <a-select-option value="sqlite">
                  <span style="font-weight: 500;">SQLite</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">轻量级，无需额外服务</span>
                </a-select-option>
                <a-select-option value="mysql">
                  <span style="font-weight: 500;">MySQL</span>
                  <span style="color: #8C8C8C; font-size: 12px; margin-left: 8px;">适合生产环境</span>
                </a-select-option>
              </a-select>
              <div class="form-tip">
                <InfoCircleOutlined /> SQLite 适合单机部署，MySQL 适合高并发场景
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="数据库路径 / 连接串">
              <a-input v-model:value="config.dbSource" placeholder="data/git-sync.db" />
              <div class="form-tip">
                <InfoCircleOutlined /> SQLite: 文件路径；MySQL: user:pass@tcp(host:port)/dbname
              </div>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider orientation="left">
          <span class="divider-title">Git 配置</span>
        </a-divider>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="默认分支">
              <a-input v-model:value="config.defaultBranch" placeholder="main" />
              <div class="form-tip">
                <InfoCircleOutlined /> 新建任务时的默认分支名称
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="超时时间">
              <a-input-number v-model:value="config.timeout" :min="10" :max="600" style="width: 100%" addon-after="秒" />
              <div class="form-tip">
                <InfoCircleOutlined /> 单次同步操作的最大等待时间
              </div>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="最大并发同步数">
              <a-input-number v-model:value="config.maxConcurrent" :min="1" :max="10" style="width: 100%" />
              <div class="form-tip">
                <InfoCircleOutlined /> 同时执行的最大同步任务数。过高可能导致资源耗尽
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="历史记录保留天数">
              <a-input-number v-model:value="config.retentionDays" :min="1" :max="365" style="width: 100%" addon-after="天" />
              <div class="form-tip">
                <InfoCircleOutlined /> 超过此天数的同步历史记录将被自动清理
              </div>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleSave">
              <template #icon><SaveOutlined /></template>
              保存配置
            </a-button>
            <a-popconfirm
              title="确定要重置为默认配置吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleReset"
            >
              <a-button danger>重置为默认</a-button>
            </a-popconfirm>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { message } from 'ant-design-vue'
import { InfoCircleOutlined, SaveOutlined } from '@ant-design/icons-vue'

const defaults = {
  dbDriver: 'sqlite',
  dbSource: 'data/git-sync.db',
  defaultBranch: 'main',
  timeout: 300,
  maxConcurrent: 3,
  retentionDays: 30,
}

const config = reactive({ ...defaults })

function handleSave() {
  message.success('配置已保存')
}

function handleReset() {
  Object.assign(config, defaults)
  message.info('配置已重置')
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
