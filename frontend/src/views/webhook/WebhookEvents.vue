<template><div class="page-container">
   <div class="page-header"><h1 class="page-title-light">Webhook 事件日志</h1></div>
   <div class="tabs-bar">
     <button class="tab-btn" :class="{active: activeTab === 'all'}" @click="activeTab = 'all'">全部</button>
     <button class="tab-btn" :class="{active: activeTab === 'success'}" @click="activeTab = 'success'">成功</button>
     <button class="tab-btn" :class="{active: activeTab === 'failed'}" @click="activeTab = 'failed'">失败</button>
   </div>
   <div class="filter-bar-light" style="margin-bottom:16px;">
     <div class="filter-input-light" style="width:280px;"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#8c8c8c" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg><input type="text" placeholder="搜索事件..."></div>
     <div class="filter-select-light"><span>平台</span></div>
   </div>
  <div class="table-card">
     <table class="sync-table">
       <thead><tr><th style="width:180px;">时间</th><th>事件类型</th><th>仓库</th><th>分支</th><th>触发者</th><th>状态</th><th style="width:120px;text-align:center;">操作</th></tr></thead>
        <tbody><tr v-for="e in filteredEvents" :key="e.id">
          <td style="color:#8c8c8c;">{{ e.time }}</td>
          <td><span class="event-badge" :class="e.type">{{ e.type }}</span></td>
          <td>{{ e.repo }}</td>
          <td>{{ e.branch || '-' }}</td>
          <td>{{ e.sender }}</td>
          <td><span class="status-badge" :class="e.status">{{ e.status==='success'?'成功':e.status==='running'?'处理中':'失败' }}</span></td>
          <td class="action-col"><button class="action-btn view" @click="showDetail(e)">详情</button></td>
        </tr></tbody>
     </table>
   </div>

   <!-- 详情弹窗 -->
   <el-dialog v-model="detailVisible" title="Webhook 事件详情" width="600px" class="detail-modal">
     <div class="detail-content" v-if="currentEvent">
       <div class="detail-row">
         <span class="detail-label">事件ID</span>
         <span class="detail-value">{{ currentEvent.id }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">触发时间</span>
         <span class="detail-value">{{ currentEvent.time }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">事件类型</span>
         <span class="event-badge" :class="currentEvent.type">{{ currentEvent.type }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">仓库</span>
         <span class="detail-value">{{ currentEvent.repo }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">分支</span>
         <span class="detail-value">{{ currentEvent.branch || '-' }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">触发者</span>
         <span class="detail-value">{{ currentEvent.sender }}</span>
       </div>
       <div class="detail-row">
         <span class="detail-label">状态</span>
         <span class="status-badge" :class="currentEvent.status">{{ currentEvent.status==='success'?'成功':currentEvent.status==='running'?'处理中':'失败' }}</span>
       </div>
       <div class="detail-row full">
         <span class="detail-label">请求详情</span>
         <pre class="detail-pre">{{ JSON.stringify(currentEvent.payload || { test: 'payload data' }, null, 2) }}</pre>
       </div>
     </div>
     <template #footer>
       <div class="detail-actions">
         <button class="btn-cancel" @click="detailVisible = false">关闭</button>
         <button class="btn-retry" v-if="currentEvent?.status === 'failed'" @click="handleRetry">重试</button>
       </div>
     </template>
   </el-dialog>
 </div></template>

 <script setup lang="ts">
 import { ref, computed } from 'vue'
 import { ElMessage } from 'element-plus'
 
 const activeTab = ref('all')
 const detailVisible = ref(false)
 const currentEvent = ref<any>(null)
 
 const events = ref([
   { id: 1, time: '2024-05-16 15:30:00', type: 'push', repo: 'owner/repo1', branch: 'main', sender: 'user1', status: 'success' },
   { id: 2, time: '2024-05-16 14:20:00', type: 'merge_request', repo: 'group/project', branch: 'develop', sender: 'user2', status: 'success' },
   { id: 3, time: '2024-05-16 13:00:00', type: 'tag', repo: 'org/repo', branch: '', sender: 'ci-bot', status: 'failed' },
 ])
 
 const filteredEvents = computed(() => {
   if (activeTab.value === 'all') return events.value
   return events.value.filter(e => e.status === activeTab.value)
 })
 
 function showDetail(e: any) {
   currentEvent.value = e
   detailVisible.value = true
 }
 
 function handleRetry() {
   ElMessage.success(`已触发重试: 事件 #${currentEvent.value.id}`)
   detailVisible.value = false
 }
 </script>

<style scoped lang="scss">
 .page-container { background: #f0f2f5; min-height: 100%; }
 .page-header { margin-bottom: 16px; }
 .page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
 .tabs-bar { display: flex; gap: 4px; margin-bottom: 16px; }
 .tab-btn { padding: 6px 16px; border: none; background: transparent; border-radius: 6px; font-size: 13px; color: #595959; cursor: pointer; transition: all 0.2s;
   &:hover { color: #1890ff; }
   &.active { background: #1890ff; color: #fff; }
 }
 .filter-bar-light { display: flex; gap: 12px; }
 .filter-input-light { background: #fff; border: 1px solid #d9d9d9; border-radius: 6px; padding: 0 12px; height: 32px; display: flex; align-items: center; gap: 8px; input { background: transparent; border: none; outline: none; font-size: 13px; color: #262626; width: 100%; } }
 .filter-select-light { background: #fff; border: 1px solid #d9d9d9; border-radius: 6px; padding: 0 12px; height: 32px; display: flex; align-items: center; gap: 8px; font-size: 13px; color: #262626; }
 .table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
 .sync-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
 .event-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block;
   &.push { background: #e6f7ff; color: #1890ff; }
   &.merge_request { background: #f6ffed; color: #52c41a; }
   &.tag { background: #fff7e6; color: #fa8c16; }
   &.issue { background: #fff2f0; color: #ff4d4f; }
 }
 .status-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block;
   &.success { background: #f6ffed; color: #52c41a; }
   &.running { background: #e6f7ff; color: #1890ff; }
   &.failed { background: #fff2f0; color: #ff4d4f; }
 }
  .action-col { display: flex; justify-content: center; }
  .action-btn {
    padding: 4px 8px;
    border-radius: 4px;
    border: none;
    font-size: 12px;
    cursor: pointer;
    transition: all 0.2s;
    
    &.view {
      background: #e6f7ff;
      color: #1890ff;
      &:hover {
        background: #bae7ff;
      }
    }
  }
  .detail-modal :deep(.el-dialog) {
    border-radius: 12px;
  }
  .detail-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .detail-row {
    display: flex;
    align-items: center;
    gap: 12px;
    
    &.full {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
  }
  .detail-label {
    min-width: 80px;
    font-size: 13px;
    color: #8c8c8c;
  }
  .detail-value {
    font-size: 13px;
    color: #262626;
    flex: 1;
  }
  .detail-pre {
    width: 100%;
    max-height: 200px;
    overflow-y: auto;
    background: #f5f5f5;
    border-radius: 6px;
    padding: 12px;
    margin: 0;
    font-size: 12px;
    line-height: 1.6;
    color: #595959;
    box-sizing: border-box;
  }
  .detail-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
  .btn-cancel {
    min-width: 80px;
    height: 32px;
    border-radius: 6px;
    background: #fff;
    border: 1px solid #d9d9d9;
    color: #595959;
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      color: #1890ff;
      border-color: #1890ff;
    }
  }
  .btn-retry {
    min-width: 80px;
    height: 32px;
    border-radius: 6px;
    background: #1890ff;
    border: none;
    color: #fff;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    
    &:hover {
      background: #40a9ff;
    }
  }
  </style>
