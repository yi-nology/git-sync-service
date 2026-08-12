<template>
  <span class="status-badge" :class="kind">
    <span class="status-dot" :class="kind"></span>
    {{ text }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { statusLabel, statusKind } from '@/constants/status'

const props = defineProps<{
  status: string | undefined
}>()

const kind = computed(() => statusKind(props.status))
const text = computed(() => statusLabel(props.status))
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 2px 10px;
  border-radius: $radius-sm;
  font-size: 12px;
  font-weight: 500;

  &.success { background: #f6ffed; color: $success; }
  &.running { background: #e6f7ff; color: $primary; }
  &.failed  { background: #fff2f0; color: $error; }
  &.warning { background: #fff7e6; color: $warning; }
  &.idle    { background: #f5f5f5; color: $text-secondary; }
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;

  &.success { background: $success; }
  &.running { background: $primary; }
  &.failed  { background: $error; }
  &.warning { background: $warning; }
  &.idle    { background: $text-secondary; }
}
</style>
