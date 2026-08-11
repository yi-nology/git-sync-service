<template>
  <header class="top-nav">
    <div class="nav-container">
      <!-- Logo -->
      <div class="nav-logo" @click="router.push('/dashboard')">
        <div class="logo-icon">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          </svg>
        </div>
        <span class="logo-text">Git Sync</span>
      </div>

      <!-- Navigation Menu -->
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="horizontal"
        theme="light"
        @click="handleMenuClick"
        class="nav-menu"
      >
        <a-menu-item key="/dashboard">
          <template #icon><DashboardOutlined /></template>
          <span>仪表盘</span>
        </a-menu-item>

        <a-sub-menu key="sync">
          <template #icon><SyncOutlined /></template>
          <template #title>同步管理</template>
          <a-menu-item key="/sync">同步任务</a-menu-item>
          <a-menu-item key="/sync/history">同步历史</a-menu-item>
          <a-menu-item key="/sync/new">
            <PlusCircleOutlined style="margin-right: 4px; font-size: 12px;" />
            新建任务
          </a-menu-item>
        </a-sub-menu>

        <a-menu-item key="/repos">
          <template #icon><FolderOutlined /></template>
          <span>仓库管理</span>
        </a-menu-item>

        <a-sub-menu key="webhook">
          <template #icon><ApiOutlined /></template>
          <template #title>触发配置</template>
          <a-menu-item key="/webhook/rules">Webhook 规则</a-menu-item>
          <a-menu-item key="/webhook/events">事件日志</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="logs">
          <template #icon><FileTextOutlined /></template>
          <template #title>日志管理</template>
          <a-menu-item key="/logs/operations">操作日志</a-menu-item>
          <a-menu-item key="/logs/sync">同步执行日志</a-menu-item>
          <a-menu-item key="/logs/system">系统运行日志</a-menu-item>
        </a-sub-menu>

        <a-sub-menu key="system">
          <template #icon><SettingOutlined /></template>
          <template #title>系统</template>
          <a-menu-item key="/settings">系统设置</a-menu-item>
          <a-menu-item key="/settings/platforms">
            <CloudServerOutlined style="margin-right: 4px; font-size: 12px;" />
            平台管理
          </a-menu-item>
          <a-menu-item key="/settings/advanced">高级配置</a-menu-item>
        </a-sub-menu>
      </a-menu>

      <!-- Right Actions -->
      <div class="nav-actions">
        <a-tooltip title="刷新页面">
          <a-button type="text" @click="handleRefresh">
            <template #icon><ReloadOutlined /></template>
          </a-button>
        </a-tooltip>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  DashboardOutlined,
  SyncOutlined,
  ApiOutlined,
  FolderOutlined,
  SettingOutlined,
  ReloadOutlined,
  PlusCircleOutlined,
  CloudServerOutlined,
  FileTextOutlined,
} from '@ant-design/icons-vue'

const route = useRoute()
const router = useRouter()

const selectedKeys = ref<string[]>([])
const openKeys = ref<string[]>([])

// Submenu mapping for auto-opening
const submenuMap: Record<string, string> = {
  '/dashboard': '',
  '/sync': 'sync',
  '/sync/history': 'sync',
  '/sync/new': 'sync',
  '/repos': '',
  '/webhook/rules': 'webhook',
  '/webhook/events': 'webhook',
  '/logs/operations': 'logs',
  '/logs/sync': 'logs',
  '/logs/system': 'logs',
  '/settings': 'system',
  '/settings/platforms': 'system',
  '/settings/advanced': 'system',
}

// Update selected keys and open keys based on route
watch(() => route.path, (path) => {
  // For dynamic routes, highlight the parent
  let highlightPath = path
  if (path.startsWith('/repos/config/') || path.startsWith('/local-repos/')) {
    highlightPath = '/repos'
  }

  selectedKeys.value = [highlightPath]
  const submenu = submenuMap[highlightPath]
  if (submenu) {
    openKeys.value = [submenu]
  }
}, { immediate: true })

const handleMenuClick = ({ key }: { key: string }) => {
  router.push(key)
}

const handleRefresh = () => {
  router.go(0)
}
</script>

<style scoped lang="scss">
.top-nav {
  position: sticky;
  top: 0;
  z-index: 100;
  height: 56px;
  background: #FFFFFF;
  border-bottom: 1px solid #F0F0F0;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.04);
}

.nav-container {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 24px;
  max-width: 100%;
}

.nav-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: opacity 0.2s;
  flex-shrink: 0;
  margin-right: 32px;

  &:hover {
    opacity: 0.8;
  }

  .logo-icon {
    color: #1677FF;
    display: flex;
    align-items: center;
  }

  .logo-text {
    color: #1677FF;
    font-size: 18px;
    font-weight: 700;
    white-space: nowrap;
    letter-spacing: -0.5px;
  }
}

.nav-menu {
  flex: 1;
  border-bottom: none;
  line-height: 56px;

  :deep(.ant-menu-item) {
    height: 56px;
    line-height: 56px;
    padding: 0 16px;

    &.ant-menu-item-selected {
      color: #1677FF;
      font-weight: 500;
    }
  }

  :deep(.ant-menu-submenu-title) {
    height: 56px;
    line-height: 56px;
    padding: 0 16px;
  }
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  margin-left: 16px;
}
</style>
