<template>
  <a-layout class="app-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      :width="240"
      :collapsed-width="80"
      class="app-sidebar"
    >
      <!-- Logo -->
      <div class="sidebar-logo">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#1677FF" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
        <span v-if="!collapsed" class="logo-text">Git Sync</span>
      </div>

      <!-- Navigation Menu -->
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        theme="light"
        @click="handleMenuClick"
      >
        <a-sub-menu key="overview">
          <template #icon><DashboardOutlined /></template>
          <template #title>概览</template>
          <a-menu-item key="/dashboard">仪表盘</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="sync">
          <template #icon><SyncOutlined /></template>
          <template #title>同步管理</template>
          <a-menu-item key="/sync">同步任务</a-menu-item>
          <a-menu-item key="/sync/history">同步历史</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="webhook">
          <template #icon><ApiOutlined /></template>
          <template #title>触发配置</template>
          <a-menu-item key="/webhook/rules">Webhook 规则</a-menu-item>
          <a-menu-item key="/webhook/events">事件日志</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="repos">
          <template #icon><FolderOutlined /></template>
          <template #title>仓库管理</template>
          <a-menu-item key="/repos">仓库列表</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="system">
          <template #icon><SettingOutlined /></template>
          <template #title>系统</template>
          <a-menu-item key="/settings">系统设置</a-menu-item>
          <a-menu-item key="/settings/advanced">高级配置</a-menu-item>
        </a-sub-menu>
      </a-menu>

      <!-- Collapse Trigger -->
      <div class="sidebar-trigger" @click="collapsed = !collapsed">
        <MenuFoldOutlined v-if="!collapsed" />
        <MenuUnfoldOutlined v-else />
      </div>
    </a-layout-sider>

    <a-layout class="main-layout">
      <a-layout-header class="app-header">
        <!-- Breadcrumb -->
        <a-breadcrumb>
          <a-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">
            <router-link v-if="item.path" :to="item.path">{{ item.label }}</router-link>
            <span v-else>{{ item.label }}</span>
          </a-breadcrumb-item>
        </a-breadcrumb>

        <!-- Refresh button -->
        <a-button type="text" @click="handleRefresh">
          <template #icon><ReloadOutlined /></template>
        </a-button>
      </a-layout-header>

      <a-layout-content class="app-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DashboardOutlined,
  SyncOutlined,
  ApiOutlined,
  FolderOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ReloadOutlined
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])

// Breadcrumb mapping
const breadcrumbMap: Record<string, { label: string; path?: string }[]> = {
  '/dashboard': [{ label: '仪表盘' }],
  '/sync': [{ label: '同步管理', path: '/sync' }, { label: '同步任务' }],
  '/sync/history': [{ label: '同步管理', path: '/sync' }, { label: '同步历史' }],
  '/webhook/rules': [{ label: '触发配置', path: '/webhook/rules' }, { label: 'Webhook 规则' }],
  '/webhook/events': [{ label: '触发配置', path: '/webhook/events' }, { label: '事件日志' }],
  '/repos': [{ label: '仓库管理' }],
  '/settings': [{ label: '系统', path: '/settings' }, { label: '系统设置' }],
  '/settings/advanced': [{ label: '系统', path: '/settings' }, { label: '高级配置' }]
}

// Submenu mapping for auto-opening
const submenuMap: Record<string, string> = {
  '/dashboard': 'overview',
  '/sync': 'sync',
  '/sync/history': 'sync',
  '/webhook/rules': 'webhook',
  '/webhook/events': 'webhook',
  '/repos': 'repos',
  '/settings': 'system',
  '/settings/advanced': 'system'
}

const breadcrumbs = computed(() => {
  return breadcrumbMap[route.path] || [{ label: '首页' }]
})

// Update selected keys and open keys based on route
watch(() => route.path, (path) => {
  selectedKeys.value = [path]
  const submenu = submenuMap[path]
  if (submenu && !collapsed.value) {
    openKeys.value = [submenu]
  }
}, { immediate: true })

// When sidebar collapses, close all submenus
watch(collapsed, (isCollapsed) => {
  if (isCollapsed) {
    openKeys.value = []
  } else {
    // Re-open the current submenu
    const submenu = submenuMap[route.path]
    if (submenu) {
      openKeys.value = [submenu]
    }
  }
})

const handleMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}

const handleRefresh = () => {
  router.go(0)
}
</script>

<style scoped lang="scss">
.app-layout {
  min-height: 100vh;
}

.app-sidebar {
  background: #FFFFFF !important;
  border-right: 1px solid #F0F0F0;
  box-shadow: 2px 0 8px 0 rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;

  :deep(.ant-menu) {
    border-right: none;
    flex: 1;
  }

  :deep(.ant-menu-item) {
    margin: 4px 8px;
    border-radius: 8px;

    &.ant-menu-item-selected {
      background: #E6F7FF;
      color: #1677FF;
    }
  }

  :deep(.ant-menu-submenu-title) {
    margin: 4px 8px;
    border-radius: 8px;
  }
}

.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 10px;
  border-bottom: 1px solid #F0F0F0;

  .logo-text {
    color: #1677FF;
    font-size: 18px;
    font-weight: 700;
    white-space: nowrap;
  }
}

.sidebar-trigger {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-top: 1px solid #F0F0F0;
  color: #666666;
  transition: all 0.2s;

  &:hover {
    color: #1677FF;
    background: #F5F5F5;
  }
}

.main-layout {
  min-height: 100vh;
}

.app-header {
  background: #FFFFFF !important;
  padding: 0 24px !important;
  height: 56px;
  line-height: 56px;
  border-bottom: 1px solid #F0F0F0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-content {
  margin: 24px;
  padding: 24px;
  background: #FFFFFF;
  border-radius: 8px;
  min-height: 280px;
}
</style>
