<template>
  <span class="status-badge" :class="statusClass">
    <span class="status-dot" :class="statusClass"></span>
    {{ statusText }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string | undefined
}>()

const statusClass = computed(() => {
  const status = props.status || 'idle'
  const classMap: Record<string, string> = {
    success: 'success',
    running: 'running',
    failed: 'failed',
    received: 'running',
    processed: 'success',
    active: 'success',
    idle: 'idle',
    stopped: 'idle',
  }
  return classMap[status] || 'idle'
})

const statusText = computed(() => {
  const status = props.status || 'idle'
  const textMap: Record<string, string> = {
    success: '成功',
    running: '运行中',
    failed: '失败',
    received: '已接收',
    processed: '已处理',
    active: '活跃',
    idle: '未运行',
    stopped: '已停止',
  }
  return textMap[status] || status || '-'
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: $border-radius-sm;
  font-size: 12px;
  font-weight: 500;

  &.success {
    background: #F6FFED;
    color: $success-color;
  }

  &.running {
    background: #E6F7FF;
    color: $primary-color;
  }

  &.failed {
    background: #FFF2F0;
    color: $error-color;
  }

  &.idle {
    background: #F5F5F5;
    color: $text-secondary;
  }
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;

  &.success { background: $success-color; }
  &.running { background: $primary-color; }
  &.failed { background: $error-color; }
  &.idle { background: $text-secondary; }
}
</style>
