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
          row-key="id"
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

<style scoped lang="scss">
@import '@/styles/variables.scss';

.page-container {
  background: $background-color;
  min-height: 100%;
  padding: $spacing-lg;
}

.content-card {
  background: $card-background;
  border-radius: $border-radius-md;
  box-shadow: $shadow-card;
  overflow: hidden;
}

.card-header {
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid $border-color;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-body {
  padding: 0;
}
</style>
