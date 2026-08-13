<template>
  <!-- 不渲染可见内容,仅作为 <a-app> 的子组件来获取上下文并注入 message 实例 -->
  <span aria-hidden="true" style="display: none" />
</template>

<script setup lang="ts">
import { App as AntApp } from 'ant-design-vue'
import { registerMessage } from '@/utils/notify'

// 必须在 <a-app> 内部调用,才能拿到与 App 上下文绑定的 message 实例
const { message } = AntApp.useApp()
registerMessage(message)
</script>
