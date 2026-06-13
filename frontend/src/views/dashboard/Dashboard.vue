<template>
  <div class="dashboard">
    <div class="stats-row">
      <div class="stat-card">
        <div class="stat-icon primary">
          <el-icon><Odometer /></el-icon>
        </div>
        <div class="stat-value">{{ stats.totalTasks }}</div>
        <div class="stat-label">同步任务总数</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon success">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-value">{{ stats.runningTasks }}</div>
        <div class="stat-label">运行中任务</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon warning">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-value">{{ stats.todaySyncs }}</div>
        <div class="stat-label">今日同步次数</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon danger">
          <el-icon><Warning /></el-icon>
        </div>
        <div class="stat-value">{{ stats.failedTasks }}</div>
        <div class="stat-label">失败任务</div>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card class="chart-card">
          <template #header>
            <span class="card-title">同步趋势（最近7天）</span>
          </template>
          <div ref="chartRef" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="chart-card">
          <template #header>
            <span class="card-title">同步状态分布</span>
          </template>
          <div ref="pieChartRef" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="recent-activities" style="margin-top: 20px;">
      <template #header>
        <span class="card-title">最近同步记录</span>
        <el-button type="primary" link @click="$router.push('/sync/history')">查看全部</el-button>
      </template>
      <el-table :data="recentSyncs" stripe>
        <el-table-column prop="taskName" label="任务名称" width="200" />
        <el-table-column prop="source" label="源仓库" width="180">
          <template #default="{ row }">
            <span class="repo-badge source">{{ row.source }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="target" label="目标仓库" width="180">
          <template #default="{ row }">
            <span class="repo-badge target">{{ row.target }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'success'" type="success" size="small">成功</el-tag>
            <el-tag v-else-if="row.status === 'running'" type="info" size="small">同步中</el-tag>
            <el-tag v-else type="danger" size="small">失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时" width="100" align="center" />
        <el-table-column prop="time" label="时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import * as echarts from 'echarts'
import { Odometer, CircleCheck, Clock, Warning } from '@element-plus/icons-vue'

const chartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()

const stats = reactive({
  totalTasks: 12,
  runningTasks: 3,
  todaySyncs: 47,
  failedTasks: 1,
})

const recentSyncs = ref([
  { taskName: 'main分支同步', source: 'GitHub', target: 'GitLab', status: 'success', duration: '2.3s', time: '2024-05-16 15:30:00' },
  { taskName: 'develop分支同步', source: 'GitHub', target: 'Gitee', status: 'success', duration: '1.8s', time: '2024-05-16 14:20:00' },
  { taskName: '全分支镜像', source: 'GitLab', target: 'GitHub', status: 'running', duration: '-', time: '2024-05-16 13:15:00' },
  { taskName: 'tag同步', source: 'GitHub', target: 'Gitee', status: 'success', duration: '0.5s', time: '2024-05-16 12:00:00' },
  { taskName: 'feature分支同步', source: 'Gitee', target: 'GitLab', status: 'failed', duration: '3.2s', time: '2024-05-16 11:30:00' },
])

onMounted(() => {
  initChart()
  initPieChart()
})

function initChart() {
  if (!chartRef.value) return
  const chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['同步次数', '成功次数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['5/10', '5/11', '5/12', '5/13', '5/14', '5/15', '5/16'],
    },
    yAxis: { type: 'value' },
    series: [
      {
        name: '同步次数',
        type: 'line',
        smooth: true,
        data: [32, 45, 38, 52, 41, 44, 47],
        itemStyle: { color: '#409eff' },
      },
      {
        name: '成功次数',
        type: 'line',
        smooth: true,
        data: [30, 44, 37, 50, 40, 43, 46],
        itemStyle: { color: '#67c23a' },
      },
    ],
  })
}

function initPieChart() {
  if (!pieChartRef.value) return
  const chart = echarts.init(pieChartRef.value)
  chart.setOption({
    tooltip: { trigger: 'item' },
    series: [
      {
        name: '同步状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: { label: { show: true, fontSize: 20, fontWeight: 'bold' } },
        labelLine: { show: false },
        data: [
          { value: 45, name: '成功', itemStyle: { color: '#67c23a' } },
          { value: 2, name: '失败', itemStyle: { color: '#f56c6c' } },
          { value: 3, name: '进行中', itemStyle: { color: '#409eff' } },
        ],
      },
    ],
  })
}
</script>

<style scoped lang="scss">
.dashboard {
  .card-title {
    font-weight: 600;
    font-size: 14px;
  }

  .repo-badge {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;

    &.source {
      background: rgba(64, 158, 255, 0.1);
      color: #409eff;
    }

    &.target {
      background: rgba(103, 194, 58, 0.1);
      color: #67c23a;
    }
  }

  .recent-activities {
    .el-card__header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
