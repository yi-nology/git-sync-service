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
            row-key="key"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'branch'">
                {{ record.source_branch }} → {{ record.target_branch }}
              </template>
              <template v-if="column.dataIndex === 'last_status'">
                <StatusBadge :status="record.last_status" />
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
            row-key="key"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'platform'">
                <a-tag color="blue">{{ record.platform }}</a-tag>
              </template>
              <template v-if="column.dataIndex === 'status'">
                <StatusBadge :status="record.status" />
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
  { title: '分支', dataIndex: 'branch', key: 'branch' },
  { title: '状态', dataIndex: 'last_status', key: 'last_status' },
]

const repoColumns = [
  { title: '仓库名称', dataIndex: 'name', key: 'name' },
  { title: '平台', dataIndex: 'platform', key: 'platform' },
  { title: '状态', dataIndex: 'status', key: 'status' },
]

onMounted(() => {
  repoStore.fetchRepos()
  taskStore.fetchTasks()
})
</script>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

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

  &.blue { background: #E6F7FF; color: $primary-color; }
  &.green { background: #F6FFED; color: $success-color; }
  &.orange { background: #FFF7E6; color: $warning-color; }
  &.red { background: #FFF2F0; color: $error-color; }
}

.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: $text-primary;
}

.stat-name {
  font-size: 13px;
  color: $text-secondary;
  margin-top: $spacing-xs;
}

.grid-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-md;
}

.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid $border-color;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
}

.card-body {
  padding: $spacing-md $spacing-lg;
}
</style>
