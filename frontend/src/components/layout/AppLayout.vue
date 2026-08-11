<template>
  <a-layout class="app-layout">
    <!-- Top Navigation -->
    <TopNav />

    <!-- Main Content Area -->
    <a-layout class="main-layout">
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
import TopNav from './TopNav.vue'
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
}

.main-layout {
  min-height: calc(100vh - 56px);
}

.app-content {
  padding: 24px;
  min-height: 280px;
}

// Page transition
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.15s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
