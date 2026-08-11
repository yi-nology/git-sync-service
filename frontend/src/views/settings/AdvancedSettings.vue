<template>
  <div class="page-container">
    <PageHeader title="高级配置" />

    <div class="content-card">
      <div class="card-body">
        <a-form layout="vertical" :model="config">
          <a-divider orientation="left">数据库配置</a-divider>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="数据库驱动">
                <a-select v-model:value="config.dbDriver">
                  <a-select-option value="sqlite">SQLite</a-select-option>
                  <a-select-option value="mysql">MySQL</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="数据库路径">
                <a-input v-model:value="config.dbSource" placeholder="data/git-sync.db" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-divider orientation="left">Git 配置</a-divider>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="默认分支">
                <a-input v-model:value="config.defaultBranch" placeholder="main" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="超时时间 (秒)">
                <a-input-number v-model:value="config.timeout" :min="10" :max="600" style="width: 100%" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="最大并发同步数">
                <a-input-number v-model:value="config.maxConcurrent" :min="1" :max="10" style="width: 100%" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="历史记录保留天数">
                <a-input-number v-model:value="config.retentionDays" :min="1" :max="365" style="width: 100%" />
              </a-form-item>
            </a-col>
          </a-row>

          <a-form-item>
            <a-space>
              <a-button type="primary" @click="handleSave">保存配置</a-button>
              <a-button danger @click="handleReset">重置为默认</a-button>
            </a-space>
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
