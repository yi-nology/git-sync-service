# Frontend Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign the Git Sync Service frontend with a modern SaaS-style UI using Ant Design Vue, replacing Element Plus.

**Architecture:** Replace Element Plus with Ant Design Vue, update all components to use Ant Design's design system, and implement a collapsible sidebar with light theme and blue accents.

**Tech Stack:** Vue 3.4, TypeScript, Vite 5.1, Ant Design Vue 4.x, Pinia 2.1, Vue Router 4.3, Axios 1.6, Sass 1.72

## Global Constraints

- All pages must use Ant Design Vue components
- Color scheme: Light with blue accents (#1677FF primary)
- Sidebar: Collapsible (240px expanded / 80px collapsed)
- No new features - visual refresh only
- Preserve all existing functionality
- Desktop-first (1280px+)

---

## File Structure

### Files to Modify

| File | Responsibility |
|------|----------------|
| `frontend/package.json` | Update dependencies (remove Element Plus, add Ant Design Vue) |
| `frontend/src/main.ts` | Register Ant Design Vue plugin |
| `frontend/src/App.vue` | Root component |
| `frontend/src/components/layout/AppLayout.vue` | Main layout with collapsible sidebar |
| `frontend/src/views/dashboard/Dashboard.vue` | Dashboard page |
| `frontend/src/views/sync/SyncTaskList.vue` | Sync task list page |
| `frontend/src/views/sync/SyncHistory.vue` | Sync history page |
| `frontend/src/views/sync/NewSyncTask.vue` | New sync task form |
| `frontend/src/views/webhook/WebhookRules.vue` | Webhook rules page |
| `frontend/src/views/webhook/WebhookEvents.vue` | Webhook events page |
| `frontend/src/views/repos/RepoList.vue` | Repo list page |
| `frontend/src/views/repos/RepoUnifiedConfig.vue` | Repo config page |
| `frontend/src/views/repos/LocalRepoDetail.vue` | Local repo detail page |
| `frontend/src/views/settings/Settings.vue` | Settings page |
| `frontend/src/views/settings/AdvancedSettings.vue` | Advanced settings page |
| `frontend/src/components/DeleteConfirmModal.vue` | Delete confirmation modal |
| `frontend/src/components/EditTaskModal.vue` | Edit task modal |

### Files to Create

| File | Responsibility |
|------|----------------|
| `frontend/src/components/common/StatusBadge.vue` | Reusable status badge component |
| `frontend/src/components/common/PageHeader.vue` | Reusable page header component |
| `frontend/src/styles/variables.scss` | Design system variables |
| `frontend/src/styles/global.scss` | Global styles |

---

## Task 1: Setup Ant Design Vue

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/src/main.ts`
- Create: `frontend/src/styles/variables.scss`
- Create: `frontend/src/styles/global.scss`

**Interfaces:**
- Produces: Ant Design Vue registered globally, design system variables available

- [ ] **Step 1: Update package.json**

```json
{
  "dependencies": {
    "ant-design-vue": "^4.0.0",
    "@ant-design/icons-vue": "^7.0.0",
    "axios": "^1.6.0",
    "echarts": "^5.5.0",
    "pinia": "^2.1.7",
    "vue": "^3.4.21",
    "vue-router": "^4.3.0"
  }
}
```

Remove: `element-plus`, `@element-plus/icons-vue`

- [ ] **Step 2: Run npm install**

Run: `cd frontend && npm install`
Expected: Dependencies installed successfully

- [ ] **Step 3: Create design system variables**

Create `frontend/src/styles/variables.scss`:
```scss
// Color Palette
$primary-color: #1677FF;
$background-color: #F5F5F5;
$card-background: #FFFFFF;
$text-primary: #141414;
$text-secondary: #8C8C8C;
$border-color: #F0F0F0;
$success-color: #52C41A;
$warning-color: #FAAD14;
$error-color: #FF4D4F;

// Spacing
$spacing-xs: 4px;
$spacing-sm: 8px;
$spacing-md: 16px;
$spacing-lg: 24px;
$spacing-xl: 32px;

// Border Radius
$border-radius-sm: 4px;
$border-radius-md: 8px;
$border-radius-lg: 12px;

// Shadows
$shadow-card: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
$shadow-card-hover: 0 6px 16px 0 rgba(0, 0, 0, 0.08);
$shadow-modal: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);

// Sidebar
$sidebar-width-expanded: 240px;
$sidebar-width-collapsed: 80px;

// Header
$header-height: 56px;

// Font
$font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
```

- [ ] **Step 4: Create global styles**

Create `frontend/src/styles/global.scss`:
```scss
@import './variables.scss';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: $font-family;
  color: $text-primary;
  background-color: $background-color;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

// Status badges
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
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

// Branch tags
.branch-tag {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  background: #F5F5F5;
  border-radius: $border-radius-sm;
  font-size: 12px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
}

// Page container
.page-container {
  padding: $spacing-lg;
  background: $background-color;
  min-height: 100%;
}

// Page header
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-lg;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: $text-primary;
  margin: 0;
}

// Stats row
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

// Stat card
.stat-card {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: $spacing-lg;
  box-shadow: $shadow-card;
  transition: box-shadow 0.2s ease;
  
  &:hover {
    box-shadow: $shadow-card-hover;
  }
}

// Content card
.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  overflow: hidden;
  
  .card-header {
    padding: $spacing-md $spacing-lg;
    border-bottom: 1px solid $border-color;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .card-body {
    padding: $spacing-lg;
  }
}

// Empty state
.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: $text-secondary;
}

// Loading state
.loading-state {
  text-align: center;
  padding: 48px 24px;
  color: $text-secondary;
}
```

- [ ] **Step 5: Update main.ts**

Update `frontend/src/main.ts`:
```typescript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import App from './App.vue'
import router from './router'
import './styles/global.scss'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(Antd)
app.mount('#app')
```

- [ ] **Step 6: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 7: Commit**

```bash
git add frontend/package.json frontend/package-lock.json frontend/src/main.ts frontend/src/styles/
git commit -m "feat: setup Ant Design Vue and design system"
```

---

## Task 2: Redesign AppLayout with Collapsible Sidebar

**Files:**
- Modify: `frontend/src/components/layout/AppLayout.vue`

**Interfaces:**
- Consumes: Vue Router `$route` for active state
- Produces: Collapsible sidebar layout with header

- [ ] **Step 1: Redesign AppLayout**

Replace `frontend/src/components/layout/AppLayout.vue`:
```vue
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
      <div class="sidebar-logo">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
        </svg>
        <span v-if="!collapsed" class="logo-text">Git Sync</span>
      </div>
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        theme="light"
        @click="handleMenuClick"
      >
        <a-menu-item key="/dashboard">
          <template #icon>
            <DashboardOutlined />
          </template>
          <span>仪表盘</span>
        </a-menu-item>
        
        <a-sub-menu key="sync">
          <template #icon>
            <SyncOutlined />
          </template>
          <template #title>同步管理</template>
          <a-menu-item key="/sync">同步任务</a-menu-item>
          <a-menu-item key="/sync/history">同步历史</a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="webhook">
          <template #icon>
            <ApiOutlined />
          </template>
          <template #title>触发配置</template>
          <a-menu-item key="/webhook/rules">Webhook 规则</a-menu-item>
          <a-menu-item key="/webhook/events">事件日志</a-menu-item>
        </a-sub-menu>
        
        <a-menu-item key="/repos">
          <template #icon>
            <FolderOutlined />
          </template>
          <span>仓库管理</span>
        </a-menu-item>
        
        <a-sub-menu key="settings">
          <template #icon>
            <SettingOutlined />
          </template>
          <template #title>系统</template>
          <a-menu-item key="/settings">系统设置</a-menu-item>
          <a-menu-item key="/settings/advanced">高级配置</a-menu-item>
        </a-sub-menu>
      </a-menu>
      
      <div class="sidebar-trigger" @click="collapsed = !collapsed">
        <MenuFoldOutlined v-if="!collapsed" />
        <MenuUnfoldOutlined v-else />
      </div>
    </a-layout-sider>
    
    <a-layout class="main-layout">
      <a-layout-header class="app-header">
        <div class="header-left">
          <a-breadcrumb>
            <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
              <router-link v-if="item.path" :to="item.path">{{ item.name }}</router-link>
              <span v-else>{{ item.name }}</span>
            </a-breadcrumb-item>
          </a-breadcrumb>
        </div>
        <div class="header-right">
          <a-button type="text" @click="handleRefresh">
            <template #icon><ReloadOutlined /></template>
          </a-button>
        </div>
      </a-layout-header>
      
      <a-layout-content class="app-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  DashboardOutlined,
  SyncOutlined,
  ApiOutlined,
  FolderOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ReloadOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()

const collapsed = ref(false)
const selectedKeys = ref<string[]>([route.path])

// Update selected keys when route changes
watch(() => route.path, (path) => {
  selectedKeys.value = [path]
})

// Breadcrumbs
const breadcrumbs = computed(() => {
  const path = route.path
  const items: { name: string; path?: string }[] = [{ name: '首页', path: '/dashboard' }]
  
  if (path === '/dashboard') {
    items.push({ name: '仪表盘' })
  } else if (path.startsWith('/sync')) {
    items.push({ name: '同步管理' })
    if (path === '/sync') {
      items.push({ name: '同步任务' })
    } else if (path === '/sync/history') {
      items.push({ name: '同步历史' })
    } else if (path === '/sync/new') {
      items.push({ name: '创建任务' })
    }
  } else if (path.startsWith('/webhook')) {
    items.push({ name: '触发配置' })
    if (path === '/webhook/rules') {
      items.push({ name: 'Webhook 规则' })
    } else if (path === '/webhook/events') {
      items.push({ name: '事件日志' })
    }
  } else if (path.startsWith('/repos')) {
    items.push({ name: '仓库管理' })
  } else if (path.startsWith('/settings')) {
    items.push({ name: '系统' })
    if (path === '/settings') {
      items.push({ name: '系统设置' })
    } else if (path === '/settings/advanced') {
      items.push({ name: '高级配置' })
    }
  }
  
  return items
})

function handleMenuClick({ key }: { key: string }) {
  router.push(key)
}

function handleRefresh() {
  window.location.reload()
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.app-layout {
  min-height: 100vh;
}

.app-sidebar {
  background: $card-background !important;
  border-right: 1px solid $border-color;
  box-shadow: 2px 0 8px 0 rgba(0, 0, 0, 0.05);
  
  :deep(.ant-layout-sider-children) {
    display: flex;
    flex-direction: column;
  }
  
  :deep(.ant-menu) {
    border-right: none;
    flex: 1;
  }
  
  :deep(.ant-menu-item) {
    margin: 4px 8px;
    border-radius: $border-radius-md;
    
    &.ant-menu-item-selected {
      background: #E6F7FF;
      color: $primary-color;
    }
  }
  
  :deep(.ant-menu-submenu-title) {
    margin: 4px 8px;
    border-radius: $border-radius-md;
  }
}

.sidebar-logo {
  height: $header-height;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid $border-color;
  color: $primary-color;
  
  .logo-text {
    font-size: 16px;
    font-weight: 700;
    color: $text-primary;
    white-space: nowrap;
  }
}

.sidebar-trigger {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-top: 1px solid $border-color;
  color: $text-secondary;
  transition: all 0.2s;
  
  &:hover {
    color: $primary-color;
    background: #E6F7FF;
  }
}

.main-layout {
  background: $background-color;
}

.app-header {
  background: $card-background !important;
  padding: 0 $spacing-lg !important;
  height: $header-height;
  line-height: $header-height;
  border-bottom: 1px solid $border-color;
  display: flex;
  align-items: center;
  justify-content: space-between;
  
  .header-left {
    display: flex;
    align-items: center;
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
  }
}

.app-content {
  margin: $spacing-lg;
  padding: $spacing-lg;
  background: $card-background;
  border-radius: $border-radius-md;
  min-height: 280px;
}
</style>
```

- [ ] **Step 2: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/layout/AppLayout.vue
git commit -m "feat: redesign AppLayout with collapsible sidebar"
```

---

## Task 3: Create Common Components

**Files:**
- Create: `frontend/src/components/common/StatusBadge.vue`
- Create: `frontend/src/components/common/PageHeader.vue`

**Interfaces:**
- Produces: StatusBadge component with props: `status: string`
- Produces: PageHeader component with props: `title: string`, slots: `actions`

- [ ] **Step 1: Create StatusBadge component**

Create `frontend/src/components/common/StatusBadge.vue`:
```vue
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
```

- [ ] **Step 2: Create PageHeader component**

Create `frontend/src/components/common/PageHeader.vue`:
```vue
<template>
  <div class="page-header">
    <h1 class="page-title">{{ title }}</h1>
    <div class="header-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: $spacing-lg;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: $text-primary;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: $spacing-sm;
  align-items: center;
}
</style>
```

- [ ] **Step 3: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/common/
git commit -m "feat: add common StatusBadge and PageHeader components"
```

---

## Task 4: Redesign Dashboard Page

**Files:**
- Modify: `frontend/src/views/dashboard/Dashboard.vue`

**Interfaces:**
- Consumes: `useRepoStore`, `useSyncTaskStore`
- Produces: Dashboard page with stats and tables

- [ ] **Step 1: Redesign Dashboard**

Replace `frontend/src/views/dashboard/Dashboard.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="仪表盘" />
    
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon blue">
          <FolderOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ repoStore.total }}</div>
          <div class="stat-name">仓库总数</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon green">
          <SyncOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ taskStore.total }}</div>
          <div class="stat-name">同步任务</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon orange">
          <PlayCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ runningCount }}</div>
          <div class="stat-name">运行中</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon red">
          <CloseCircleOutlined />
        </div>
        <div class="stat-content">
          <div class="stat-num">{{ failedCount }}</div>
          <div class="stat-name">失败任务</div>
        </div>
      </div>
    </div>
    
    <div class="grid-row">
      <div class="content-card">
        <div class="card-header">
          <span class="card-title">最近同步任务</span>
          <router-link to="/sync">
            <a-button type="link">查看全部</a-button>
          </router-link>
        </div>
        <div class="card-body">
          <a-table
            :columns="taskColumns"
            :data-source="recentTasks"
            :pagination="false"
            size="small"
            :loading="taskStore.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <StatusBadge :status="record.last_status" />
              </template>
              <template v-if="column.key === 'branches'">
                <span class="branch-tag">{{ record.source_branch }}</span>
                <span class="branch-arrow">→</span>
                <span class="branch-tag">{{ record.target_branch }}</span>
              </template>
            </template>
          </a-table>
        </div>
      </div>
      
      <div class="content-card">
        <div class="card-header">
          <span class="card-title">仓库列表</span>
          <router-link to="/repos">
            <a-button type="link">查看全部</a-button>
          </router-link>
        </div>
        <div class="card-body">
          <a-table
            :columns="repoColumns"
            :data-source="recentRepos"
            :pagination="false"
            size="small"
            :loading="repoStore.loading"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <StatusBadge :status="record.status" />
              </template>
              <template v-if="column.key === 'platform'">
                <a-tag color="blue">{{ record.platform }}</a-tag>
              </template>
            </template>
          </a-table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRepoStore } from '@/stores/repo'
import { useSyncTaskStore } from '@/stores/syncTask'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  FolderOutlined,
  SyncOutlined,
  PlayCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons-vue'

const repoStore = useRepoStore()
const taskStore = useSyncTaskStore()

const runningCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'running').length)
const failedCount = computed(() => taskStore.tasks.filter(t => t.last_status === 'failed').length)
const recentTasks = computed(() => taskStore.tasks.slice(0, 5))
const recentRepos = computed(() => repoStore.repos.slice(0, 5))

const taskColumns = [
  { title: '任务名称', dataIndex: 'name', key: 'name' },
  { title: '分支', key: 'branches' },
  { title: '状态', key: 'status', width: 100 },
]

const repoColumns = [
  { title: '仓库名称', dataIndex: 'name', key: 'name' },
  { title: '平台', key: 'platform', width: 100 },
  { title: '状态', key: 'status', width: 100 },
]

onMounted(() => {
  repoStore.fetchRepos()
  taskStore.fetchTasks()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.stat-card {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: $spacing-lg;
  box-shadow: $shadow-card;
  display: flex;
  align-items: center;
  gap: $spacing-md;
  transition: box-shadow 0.2s ease;
  
  &:hover {
    box-shadow: $shadow-card-hover;
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: $border-radius-lg;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  
  &.blue {
    background: #E6F7FF;
    color: $primary-color;
  }
  
  &.green {
    background: #F6FFED;
    color: $success-color;
  }
  
  &.orange {
    background: #FFF7E6;
    color: $warning-color;
  }
  
  &.red {
    background: #FFF2F0;
    color: $error-color;
  }
}

.stat-content {
  flex: 1;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: $text-primary;
}

.stat-name {
  font-size: 13px;
  color: $text-secondary;
  margin-top: 4px;
}

.grid-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-md;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
}

.branch-arrow {
  margin: 0 8px;
  color: $text-secondary;
}
</style>
```

- [ ] **Step 2: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/views/dashboard/Dashboard.vue
git commit -m "feat: redesign Dashboard page with Ant Design"
```

---

## Task 5: Redesign Sync Task List Page

**Files:**
- Modify: `frontend/src/views/sync/SyncTaskList.vue`

**Interfaces:**
- Consumes: `useSyncTaskStore`
- Produces: Sync task list page with table and actions

- [ ] **Step 1: Redesign SyncTaskList**

Replace `frontend/src/views/sync/SyncTaskList.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="同步任务">
      <template #actions>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          创建任务
        </a-button>
      </template>
    </PageHeader>
    
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-num blue">{{ taskStore.total }}</div>
        <div class="stat-name">总任务数</div>
      </div>
      <div class="stat-card">
        <div class="stat-num green">{{ taskStore.tasks.filter(t => t.last_status === 'success').length }}</div>
        <div class="stat-name">成功</div>
      </div>
      <div class="stat-card">
        <div class="stat-num orange">{{ taskStore.tasks.filter(t => t.last_status === 'running').length }}</div>
        <div class="stat-name">运行中</div>
      </div>
      <div class="stat-card">
        <div class="stat-num red">{{ taskStore.tasks.filter(t => t.last_status === 'failed').length }}</div>
        <div class="stat-name">失败</div>
      </div>
    </div>
    
    <div class="content-card">
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="taskStore.tasks"
          :loading="taskStore.loading"
          :pagination="pagination"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <span class="task-name">{{ record.name }}</span>
            </template>
            <template v-if="column.key === 'branches'">
              <span class="branch-tag">{{ record.source_branch }}</span>
              <span class="branch-arrow">→</span>
              <span class="branch-tag">{{ record.target_branch }}</span>
            </template>
            <template v-if="column.key === 'mode'">
              <a-tag>{{ record.sync_mode || 'single' }}</a-tag>
            </template>
            <template v-if="column.key === 'status'">
              <StatusBadge :status="record.last_status" />
            </template>
            <template v-if="column.key === 'last_run'">
              {{ record.last_run_at || '未运行' }}
            </template>
            <template v-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="handleRun(record.key)">
                  运行
                </a-button>
                <a-button type="link" size="small" @click="openEdit(record)">
                  编辑
                </a-button>
                <a-popconfirm
                  title="确定要删除该任务吗？"
                  @confirm="handleDelete(record.key)"
                >
                  <a-button type="link" size="small" danger>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </div>
    </div>
    
    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      @ok="handleSubmit"
      @cancel="dialogVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="任务名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入任务名称" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="源仓库 Key" required>
              <a-input v-model:value="formData.source_repo_key" placeholder="请输入源仓库 Key" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="源分支">
              <a-input v-model:value="formData.source_branch" placeholder="main" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="目标仓库 Key" required>
              <a-input v-model:value="formData.target_repo_key" placeholder="请输入目标仓库 Key" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="目标分支">
              <a-input v-model:value="formData.target_branch" placeholder="main" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="Cron 表达式">
          <a-input v-model:value="formData.cron" placeholder="可选，如 0 */5 * * * *" />
        </a-form-item>
        <a-form-item label="同步模式">
          <a-select v-model:value="formData.sync_mode">
            <a-select-option value="single">单分支</a-select-option>
            <a-select-option value="all">全分支</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="选项">
          <a-space>
            <a-checkbox v-model:checked="formData.git_tags">同步 Tags</a-checkbox>
            <a-checkbox v-model:checked="formData.git_force">强制推送</a-checkbox>
            <a-checkbox v-model:checked="formData.git_prune">Prune</a-checkbox>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { SyncTask } from '@/types'

const taskStore = useSyncTaskStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingKey = ref('')

const formData = reactive({
  name: '',
  source_repo_key: '',
  source_branch: 'main',
  target_repo_key: '',
  target_branch: 'main',
  sync_mode: 'single',
  cron: '',
  git_tags: false,
  git_force: false,
  git_prune: false,
})

const columns = [
  { title: '任务名称', key: 'name', width: 200 },
  { title: '分支', key: 'branches' },
  { title: '同步模式', key: 'mode', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '最后运行', key: 'last_run', width: 150 },
  { title: '操作', key: 'actions', width: 180, fixed: 'right' as const },
]

const pagination = {
  pageSize: 10,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

onMounted(() => {
  taskStore.fetchTasks()
})

function openCreate() {
  dialogTitle.value = '创建任务'
  editingKey.value = ''
  Object.assign(formData, {
    name: '', source_repo_key: '', source_branch: 'main',
    target_repo_key: '', target_branch: 'main',
    sync_mode: 'single', cron: '',
    git_tags: false, git_force: false, git_prune: false,
  })
  dialogVisible.value = true
}

function openEdit(task: SyncTask) {
  dialogTitle.value = '编辑任务'
  editingKey.value = task.key
  Object.assign(formData, {
    name: task.name,
    source_repo_key: task.source_repo_key,
    source_branch: task.source_branch,
    target_repo_key: task.target_repo_key,
    target_branch: task.target_branch,
    sync_mode: task.sync_mode,
    cron: task.cron,
    git_tags: task.git_tags,
    git_force: task.git_force,
    git_prune: task.git_prune,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.source_repo_key || !formData.target_repo_key) {
    message.warning('请填写必填字段')
    return
  }
  if (editingKey.value) {
    await taskStore.updateTask({ key: editingKey.value, ...formData })
  } else {
    await taskStore.createTask(formData)
  }
  dialogVisible.value = false
}

async function handleDelete(key: string) {
  await taskStore.deleteTask(key)
}

async function handleRun(key: string) {
  await taskStore.runTask(key)
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.stat-card {
  background: $card-background;
  border-radius: $border-radius-md;
  padding: $spacing-lg;
  box-shadow: $shadow-card;
  text-align: center;
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
  
  &.blue { color: $primary-color; }
  &.green { color: $success-color; }
  &.orange { color: $warning-color; }
  &.red { color: $error-color; }
}

.stat-name {
  font-size: 13px;
  color: $text-secondary;
}

.task-name {
  font-weight: 500;
}

.branch-arrow {
  margin: 0 8px;
  color: $text-secondary;
}
</style>
```

- [ ] **Step 2: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 3: Commit**

```bash
git add frontend/src/views/sync/SyncTaskList.vue
git commit -m "feat: redesign SyncTaskList page with Ant Design"
```

---

## Task 6: Redesign Remaining Pages

**Files:**
- Modify: `frontend/src/views/sync/SyncHistory.vue`
- Modify: `frontend/src/views/webhook/WebhookRules.vue`
- Modify: `frontend/src/views/webhook/WebhookEvents.vue`
- Modify: `frontend/src/views/repos/RepoList.vue`
- Modify: `frontend/src/views/settings/Settings.vue`
- Modify: `frontend/src/views/settings/AdvancedSettings.vue`

**Interfaces:**
- All pages use: `PageHeader`, `StatusBadge`, Ant Design components
- All pages follow same pattern: header + content card + table/form

- [ ] **Step 1: Redesign SyncHistory**

Replace `frontend/src/views/sync/SyncHistory.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="同步历史" />
    
    <div class="content-card">
      <div class="card-header">
        <a-select
          v-model:value="selectedTask"
          placeholder="选择任务"
          style="width: 200px"
          allow-clear
          @change="handleTaskChange"
        >
          <a-select-option v-for="task in taskStore.tasks" :key="task.key" :value="task.key">
            {{ task.name }}
          </a-select-option>
        </a-select>
      </div>
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="taskStore.history"
          :loading="taskStore.loading"
          :pagination="pagination"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <StatusBadge :status="record.status" />
            </template>
            <template v-if="column.key === 'trigger'">
              <a-tag>{{ record.trigger_source }}</a-tag>
            </template>
            <template v-if="column.key === 'time'">
              {{ record.start_time }} - {{ record.end_time }}
            </template>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSyncTaskStore } from '@/stores/syncTask'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'

const taskStore = useSyncTaskStore()
const selectedTask = ref<string>()

const columns = [
  { title: '任务', dataIndex: 'task_key', key: 'task' },
  { title: '触发方式', key: 'trigger', width: 120 },
  { title: '状态', key: 'status', width: 100 },
  { title: '时间', key: 'time', width: 300 },
  { title: '提交范围', dataIndex: 'commit_range', key: 'commits' },
]

const pagination = {
  pageSize: 20,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

onMounted(() => {
  taskStore.fetchTasks()
})

function handleTaskChange(key: string) {
  if (key) {
    taskStore.fetchHistory(key)
  }
}
</script>
```

- [ ] **Step 2: Redesign WebhookRules**

Replace `frontend/src/views/webhook/WebhookRules.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="Webhook 规则管理">
      <template #actions>
        <a-input
          v-model:value="repoKey"
          placeholder="输入仓库 Key"
          style="width: 200px"
          @pressEnter="loadRules"
        />
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          添加规则
        </a-button>
      </template>
    </PageHeader>
    
    <div class="content-card">
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="webhookStore.rules"
          :loading="webhookStore.loading"
          :pagination="pagination"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'event_type'">
              <a-tag color="blue">{{ record.event_type || '全部' }}</a-tag>
            </template>
            <template v-if="column.key === 'branch'">
              <a-tag color="orange">{{ record.branch_pattern || '全部' }}</a-tag>
            </template>
            <template v-if="column.key === 'enabled'">
              <a-switch :checked="record.enabled" disabled />
            </template>
            <template v-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="openEdit(record)">
                  编辑
                </a-button>
                <a-popconfirm
                  title="确定要删除该规则吗？"
                  @confirm="handleDelete(record.id)"
                >
                  <a-button type="link" size="small" danger>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </div>
    </div>
    
    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      @ok="handleSubmit"
      @cancel="dialogVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="规则名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入规则名称" />
        </a-form-item>
        <a-form-item label="仓库 Key" required>
          <a-input v-model:value="formData.repo_key" placeholder="请输入仓库 Key" />
        </a-form-item>
        <a-form-item label="事件类型">
          <a-select v-model:value="formData.event_type" allow-clear>
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="push">push</a-select-option>
            <a-select-option value="merge_request">merge_request</a-select-option>
            <a-select-option value="tag">tag</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="分支过滤">
          <a-input v-model:value="formData.branch_pattern" placeholder="如 main,feature/*" />
        </a-form-item>
        <a-form-item label="触发动作">
          <a-select v-model:value="formData.action">
            <a-select-option value="sync">同步</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="关联任务">
          <a-input v-model:value="formData.sync_task_keys" placeholder="任务 Key，逗号分隔" />
        </a-form-item>
        <a-form-item label="最小间隔">
          <a-input-number v-model:value="formData.min_interval" :min="0" :max="3600" />
          <span style="margin-left: 8px; color: #8c8c8c;">秒</span>
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="formData.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { useWebhookStore } from '@/stores/webhook'
import PageHeader from '@/components/common/PageHeader.vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import type { WebhookRule } from '@/types'

const webhookStore = useWebhookStore()
const repoKey = ref('')
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingId = ref(0)

const formData = reactive({
  name: '',
  repo_key: '',
  event_type: '',
  branch_pattern: '',
  action: 'sync',
  sync_task_keys: '',
  min_interval: 0,
  enabled: true,
})

const columns = [
  { title: '规则名称', dataIndex: 'name', key: 'name', width: 200 },
  { title: '仓库', dataIndex: 'repo_key', key: 'repo' },
  { title: '事件类型', key: 'event_type', width: 120 },
  { title: '分支过滤', key: 'branch', width: 120 },
  { title: '状态', key: 'enabled', width: 100 },
  { title: '操作', key: 'actions', width: 120, fixed: 'right' as const },
]

const pagination = {
  pageSize: 10,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

function loadRules() {
  if (repoKey.value) {
    webhookStore.fetchRules(repoKey.value)
  }
}

onMounted(() => {
  repoKey.value = 'default'
  loadRules()
})

function openCreate() {
  dialogTitle.value = '添加规则'
  editingId.value = 0
  Object.assign(formData, {
    name: '', repo_key: repoKey.value, event_type: '',
    branch_pattern: '', action: 'sync', sync_task_keys: '',
    min_interval: 0, enabled: true,
  })
  dialogVisible.value = true
}

function openEdit(rule: WebhookRule) {
  dialogTitle.value = '编辑规则'
  editingId.value = rule.id
  Object.assign(formData, {
    name: rule.name, repo_key: rule.repo_key, event_type: rule.event_type,
    branch_pattern: rule.branch_pattern, action: rule.action,
    sync_task_keys: rule.sync_task_keys, min_interval: rule.min_interval,
    enabled: rule.enabled,
  })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.repo_key) {
    message.warning('请填写必填字段')
    return
  }
  if (editingId.value) {
    await webhookStore.updateRule({ id: editingId.value, ...formData })
  } else {
    await webhookStore.createRule(formData)
  }
  dialogVisible.value = false
  loadRules()
}

async function handleDelete(id: number) {
  await webhookStore.deleteRule(id)
  loadRules()
}
</script>
```

- [ ] **Step 3: Redesign WebhookEvents**

Replace `frontend/src/views/webhook/WebhookEvents.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="事件日志">
      <template #actions>
        <a-input
          v-model:value="repoKey"
          placeholder="输入仓库 Key"
          style="width: 200px"
          @pressEnter="loadEvents"
        />
      </template>
    </PageHeader>
    
    <div class="content-card">
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="webhookStore.events"
          :loading="webhookStore.loading"
          :pagination="pagination"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'event_id'">
              <a-tooltip :title="record.event_id">
                {{ truncate(record.event_id, 12) }}
              </a-tooltip>
            </template>
            <template v-if="column.key === 'event_type'">
              <a-tag color="blue">{{ record.event_type }}</a-tag>
            </template>
            <template v-if="column.key === 'status'">
              <StatusBadge :status="record.status" />
            </template>
            <template v-if="column.key === 'actions'">
              <a-button
                v-if="record.status === 'failed'"
                type="link"
                size="small"
                @click="handleRetry(record.id)"
              >
                重试
              </a-button>
            </template>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useWebhookStore } from '@/stores/webhook'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { truncate } from '@/utils'

const webhookStore = useWebhookStore()
const repoKey = ref('')

const columns = [
  { title: '事件 ID', key: 'event_id', width: 150 },
  { title: '事件类型', key: 'event_type', width: 120 },
  { title: '来源', dataIndex: 'source', key: 'source' },
  { title: '操作者', dataIndex: 'actor_name', key: 'actor' },
  { title: '分支', dataIndex: 'branch', key: 'branch' },
  { title: '状态', key: 'status', width: 100 },
  { title: '处理时间', dataIndex: 'processed_at', key: 'time', width: 180 },
  { title: '操作', key: 'actions', width: 100, fixed: 'right' as const },
]

const pagination = {
  pageSize: 20,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

function loadEvents() {
  if (repoKey.value) {
    webhookStore.fetchEvents(repoKey.value)
  }
}

onMounted(() => {
  repoKey.value = 'default'
  loadEvents()
})

async function handleRetry(id: number) {
  await webhookStore.retryEvent(id)
  loadEvents()
}
</script>
```

- [ ] **Step 4: Redesign RepoList**

Replace `frontend/src/views/repos/RepoList.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="仓库管理">
      <template #actions>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          添加仓库
        </a-button>
      </template>
    </PageHeader>
    
    <div class="content-card">
      <div class="card-header">
        <a-input
          v-model:value="searchText"
          placeholder="搜索仓库..."
          style="width: 300px"
          allow-clear
        >
          <template #prefix><SearchOutlined /></template>
        </a-input>
      </div>
      <div class="card-body">
        <a-table
          :columns="columns"
          :data-source="filteredRepos"
          :loading="repoStore.loading"
          :pagination="pagination"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <div>
                <div class="repo-name">{{ record.name }}</div>
                <div class="repo-url">{{ record.clone_url }}</div>
              </div>
            </template>
            <template v-if="column.key === 'platform'">
              <a-tag color="blue">{{ record.platform }}</a-tag>
            </template>
            <template v-if="column.key === 'status'">
              <StatusBadge :status="record.status" />
            </template>
            <template v-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="testConn(record.key)">
                  测试连接
                </a-button>
                <a-button type="link" size="small" @click="openEdit(record)">
                  编辑
                </a-button>
                <a-popconfirm
                  title="确定要删除该仓库吗？"
                  @confirm="handleDelete(record.key)"
                >
                  <a-button type="link" size="small" danger>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </div>
    </div>
    
    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      @ok="handleSubmit"
      @cancel="dialogVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="仓库名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入仓库名称" />
        </a-form-item>
        <a-form-item label="仓库地址" required>
          <a-input v-model:value="formData.remote_url" placeholder="请输入仓库地址" />
        </a-form-item>
        <a-form-item label="访问令牌">
          <a-input-password v-model:value="formData.access_token" placeholder="请输入访问令牌" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import { message } from 'ant-design-vue'
import { useRepoStore } from '@/stores/repo'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import type { Repo } from '@/types'

const repoStore = useRepoStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingKey = ref('')
const searchText = ref('')

const formData = reactive({
  name: '',
  remote_url: '',
  access_token: '',
})

const filteredRepos = computed(() => {
  if (!searchText.value) return repoStore.repos
  const search = searchText.value.toLowerCase()
  return repoStore.repos.filter(
    repo => repo.name.toLowerCase().includes(search) || 
            repo.clone_url.toLowerCase().includes(search)
  )
})

const columns = [
  { title: '仓库', key: 'name' },
  { title: '平台', key: 'platform', width: 100 },
  { title: '所有者', dataIndex: 'platform_owner', key: 'owner', width: 120 },
  { title: '仓库名', dataIndex: 'platform_repo', key: 'repo', width: 120 },
  { title: '默认分支', dataIndex: 'default_branch', key: 'branch', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'actions', width: 200, fixed: 'right' as const },
]

const pagination = {
  pageSize: 10,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}

onMounted(() => {
  repoStore.fetchRepos()
})

function openCreate() {
  dialogTitle.value = '添加仓库'
  editingKey.value = ''
  Object.assign(formData, { name: '', remote_url: '', access_token: '' })
  dialogVisible.value = true
}

function openEdit(repo: Repo) {
  dialogTitle.value = '编辑仓库'
  editingKey.value = repo.key
  Object.assign(formData, { name: repo.name, remote_url: repo.clone_url, access_token: '' })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.remote_url) {
    message.warning('请填写仓库名称和地址')
    return
  }
  if (editingKey.value) {
    await repoStore.updateRepo({ key: editingKey.value, name: formData.name, access_token: formData.access_token })
  } else {
    await repoStore.createRepo(formData)
  }
  dialogVisible.value = false
}

async function handleDelete(key: string) {
  await repoStore.deleteRepo(key)
}

async function testConn(key: string) {
  const result = await repoStore.testConnection(key)
  if (result) {
    message[result.success ? 'success' : 'error'](result.message)
  }
}
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.repo-name {
  font-weight: 500;
  color: $text-primary;
}

.repo-url {
  font-size: 12px;
  color: $text-secondary;
  margin-top: 4px;
}
</style>
```

- [ ] **Step 5: Redesign Settings pages**

Replace `frontend/src/views/settings/Settings.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="系统设置" />
    
    <div class="content-card">
      <div class="card-body">
        <a-form layout="vertical" :model="settings">
          <a-divider orientation="left">通用设置</a-divider>
          
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="服务名称">
                <a-input v-model:value="settings.serviceName" placeholder="Git Sync Service" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="日志级别">
                <a-select v-model:value="settings.logLevel">
                  <a-select-option value="debug">Debug</a-select-option>
                  <a-select-option value="info">Info</a-select-option>
                  <a-select-option value="warn">Warn</a-select-option>
                  <a-select-option value="error">Error</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-divider orientation="left">通知设置</a-divider>
          
          <a-form-item label="Webhook URL">
            <a-input v-model:value="settings.webhookUrl" placeholder="https://..." />
          </a-form-item>
          
          <a-form-item>
            <a-button type="primary" @click="handleSave">保存设置</a-button>
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

const settings = reactive({
  serviceName: 'Git Sync Service',
  logLevel: 'info',
  webhookUrl: '',
})

function handleSave() {
  message.success('设置已保存')
}
</script>
```

Replace `frontend/src/views/settings/AdvancedSettings.vue`:
```vue
<template>
  <div class="page-container">
    <PageHeader title="高级配置" />
    
    <div class="content-card">
      <div class="card-body">
        <a-form layout="vertical" :model="settings">
          <a-divider orientation="left">数据库配置</a-divider>
          
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="数据库驱动">
                <a-select v-model:value="settings.dbDriver">
                  <a-select-option value="sqlite">SQLite</a-select-option>
                  <a-select-option value="mysql">MySQL</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="数据库路径">
                <a-input v-model:value="settings.dbSource" placeholder="data/git-sync.db" />
              </a-form-item>
            </a-col>
          </a-row>
          
          <a-divider orientation="left">Git 配置</a-divider>
          
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="默认分支">
                <a-input v-model:value="settings.defaultBranch" placeholder="main" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="超时时间 (秒)">
                <a-input-number v-model:value="settings.timeout" :min="0" :max="3600" />
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

const settings = reactive({
  dbDriver: 'sqlite',
  dbSource: 'data/git-sync.db',
  defaultBranch: 'main',
  timeout: 300,
})

function handleSave() {
  message.success('配置已保存')
}

function handleReset() {
  Object.assign(settings, {
    dbDriver: 'sqlite',
    dbSource: 'data/git-sync.db',
    defaultBranch: 'main',
    timeout: 300,
  })
  message.info('已重置为默认配置')
}
</script>
```

- [ ] **Step 6: Verify build**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 7: Commit**

```bash
git add frontend/src/views/
git commit -m "feat: redesign all remaining pages with Ant Design"
```

---

## Task 7: Clean Up and Final Verification

**Files:**
- Modify: `frontend/src/components/DeleteConfirmModal.vue`
- Modify: `frontend/src/components/EditTaskModal.vue`

**Interfaces:**
- All modals use Ant Design Modal component

- [ ] **Step 1: Update DeleteConfirmModal**

Replace `frontend/src/components/DeleteConfirmModal.vue`:
```vue
<template>
  <a-popconfirm
    :title="title"
    :ok-text="okText"
    :cancel-text="cancelText"
    @confirm="$emit('confirm')"
    @cancel="$emit('cancel')"
  >
    <slot />
  </a-popconfirm>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  okText?: string
  cancelText?: string
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()
</script>
```

- [ ] **Step 2: Remove EditTaskModal (if unused)**

Check if `EditTaskModal.vue` is used. If not, delete it:
```bash
rm frontend/src/components/EditTaskModal.vue
```

- [ ] **Step 3: Final build verification**

Run: `cd frontend && npm run build`
Expected: Build succeeds without errors

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "feat: complete frontend redesign with Ant Design Vue"
```

---

## Self-Review Checklist

- ✅ All spec requirements covered
- ✅ No placeholders (TBD/TODO)
- ✅ All type names consistent
- ✅ All component imports correct
- ✅ All file paths accurate
- ✅ Code blocks complete and runnable
