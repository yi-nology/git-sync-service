<template>
  <a-layout class="app-layout">
    <AppSider :collapsed="collapsed" />
    <a-layout>
      <AppHeader :collapsed="collapsed" @toggle="collapsed = !collapsed" />
      <a-layout-content class="app-content">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppSider from './AppSider.vue'
import AppHeader from './AppHeader.vue'

// 侧边栏折叠状态(由顶部按钮控制)
const collapsed = ref(false)
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
}

.app-content {
  padding: 24px;
  min-height: 280px;
  background: #f5f5f5;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.15s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
