<template>
  <div class="app-layout" :class="{ 'dark-theme': isDark }">
    <el-container class="layout-container">
      <el-aside :width="220" class="sidebar">
        <div class="sidebar-logo">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
          </svg>
          <span class="logo-text">Git Sync</span>
        </div>
        <nav class="sidebar-nav">
          <div class="nav-label">同步管理</div>
          <router-link to="/sync" class="nav-item" :class="{ active: $route.path === '/sync' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line>
            </svg>
            <span>同步任务</span>
          </router-link>
          <router-link to="/sync/history" class="nav-item" :class="{ active: $route.path === '/sync/history' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
            </svg>
            <span>同步历史</span>
          </router-link>
          <div class="nav-sep"></div>
          <div class="nav-label">触发配置</div>
          <router-link to="/webhook/rules" class="nav-item" :class="{ active: $route.path === '/webhook/rules' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            <span>Webhook 规则</span>
          </router-link>
          <router-link to="/webhook/events" class="nav-item" :class="{ active: $route.path === '/webhook/events' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
            </svg>
            <span>事件日志</span>
          </router-link>
          <div class="nav-sep"></div>
           <div class="nav-label">仓库管理</div>
           <router-link to="/repos" class="nav-item" :class="{ active: $route.path.startsWith('/repos') }">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
               <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
               <polyline points="13 2 13 9 20 9"></polyline>
             </svg>
             <span>仓库列表</span>
           </router-link>
           <div class="nav-sep"></div>
            <div class="nav-label">系统</div>
           <router-link to="/settings" class="nav-item" :class="{ active: $route.path === '/settings' }">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
               <circle cx="12" cy="12" r="3"></circle><path d="M12 1v4M12 19v4M4.22 4.22l2.83 2.83M16.95 16.95l2.83 2.83M1 12h4M19 12h4M4.22 19.78l2.83-2.83M16.95 7.05l2.83-2.83"></path>
             </svg>
             <span>系统设置</span>
           </router-link>
           <router-link to="/settings/advanced" class="nav-item" :class="{ active: $route.path === '/settings/advanced' }">
             <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
               <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path>
             </svg>
             <span>高级配置</span>
           </router-link>
        </nav>
      </el-aside>
      <el-container class="main-container">
        <el-main class="app-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isDark = computed(() => route.path === '/sync/history')
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
  background: #f0f2f5;

  &.dark-theme {
    background: #0f172a;
  }
}

.layout-container {
  height: 100vh;
}

.sidebar {
  background: #001529;
  display: flex;
  flex-direction: column;
  width: 220px;
  flex-shrink: 0;

  .sidebar-logo {
    height: 56px;
    display: flex;
    align-items: center;
    padding: 0 20px;
    gap: 10px;

    .logo-text {
      color: #ffffff;
      font-size: 16px;
      font-weight: 700;
    }
  }

  .sidebar-nav {
    flex: 1;
    padding-top: 16px;

    .nav-label {
      font-size: 11px;
      font-weight: 600;
      color: #666666;
      padding: 0 24px;
      margin-bottom: 2px;
      letter-spacing: 0.5px;
    }

    .nav-sep {
      height: 24px;
    }

    .nav-item {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 0 24px;
      height: 44px;
      color: rgba(255, 255, 255, 0.65);
      text-decoration: none;
      font-size: 14px;
      transition: all 0.2s;

      &:hover {
        background: rgba(255, 255, 255, 0.08);
        color: #ffffff;
      }

      &.active {
        background: rgba(255, 255, 255, 0.08);
        color: #ffffff;
      }
    }
  }
}

.main-container {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-main {
  padding: 24px 32px;
  overflow-y: auto;
}
</style>
