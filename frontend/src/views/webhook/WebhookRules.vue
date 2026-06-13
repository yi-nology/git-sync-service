<template><div class="page-container">
   <div class="page-header">
     <h1 class="page-title-light">Webhook 规则管理</h1>
     <div class="header-actions">
       <button class="btn-default-light">批量导出</button>
       <button class="btn-primary-light" @click="openCreate">添加规则</button>
     </div>
   </div>
   <div class="table-card" style="margin-top:24px;">
     <table class="sync-table">
       <thead><tr><th style="width:200px;">规则名称</th><th>Git 平台</th><th>触发事件</th><th>仓库过滤</th><th>触发次数</th><th>状态</th><th style="width:120px;text-align:center;">操作</th></tr></thead>
       <tbody><tr v-for="rule in rules" :key="rule.id">
         <td class="task-name-text">{{ rule.name }}</td>
         <td><span class="badge platform" :class="rule.platform">{{ rule.platform }}</span></td>
         <td><span v-for="e in rule.events" :key="e" class="badge event-tag">{{ e }}</span></td>
         <td><span class="badge filter-tag">{{ rule.filter || '全部' }}</span></td>
         <td>{{ rule.count }}</td>
         <td><span class="status-badge" :class="rule.status">{{ rule.status==='enabled'?'已启用':'已停用' }}</span></td>
         <td class="action-col"><button class="icon-btn" @click="openEdit(rule)">编辑</button><button class="icon-btn danger" @click="handleDelete(rule.id)">删除</button></td>
       </tr></tbody>
     </table>
   </div>

   <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
     <el-form :model="formData" label-width="100px">
       <el-form-item label="规则名称">
         <el-input v-model="formData.name" placeholder="请输入规则名称"/>
       </el-form-item>
       <el-form-item label="Git 平台">
         <el-select v-model="formData.platform" style="width: 100%">
           <el-option label="GitHub" value="GitHub"/>
           <el-option label="GitLab" value="GitLab"/>
           <el-option label="Gitee" value="Gitee"/>
         </el-select>
       </el-form-item>
       <el-form-item label="触发事件">
         <el-checkbox-group v-model="formData.events">
           <el-checkbox label="push"/>
           <el-checkbox label="merge_request"/>
           <el-checkbox label="tag"/>
           <el-checkbox label="issue"/>
         </el-checkbox-group>
       </el-form-item>
       <el-form-item label="仓库过滤">
         <el-input v-model="formData.filter" placeholder="请输入仓库过滤规则"/>
       </el-form-item>
     </el-form>
     <template #footer>
       <el-button @click="closeDialog">取消</el-button>
       <el-button type="primary" @click="handleSubmit">确定</el-button>
     </template>
   </el-dialog>
 </div></template>

 <script setup lang="ts">
 import { ref, reactive } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'

 const rules = ref([{ id: 1, name: '主分支同步触发', platform: 'GitHub', events: ['push'], filter: 'owner/*', count: 156, status: 'enabled' }, { id: 2, name: 'PR 自动同步', platform: 'GitLab', events: ['merge_request'], filter: 'group/*', count: 42, status: 'enabled' }])
 const platformColors = {
   GitHub: { bg: '#e6f7ff', color: '#1890ff' },
   GitLab: { bg: '#f6ffed', color: '#52c41a' },
   Gitee: { bg: '#fff7e6', color: '#fa8c16' }
 }

 const dialogVisible = ref(false)
 const dialogTitle = ref('')
 const formData = reactive<any>({ events: [] })

 function openCreate() {
   dialogTitle.value = '新建Webhook规则'
   Object.assign(formData, { name: '', platform: 'GitHub', events: ['push'], filter: '' })
   dialogVisible.value = true
 }

 function openEdit(rule: any) {
   dialogTitle.value = '编辑Webhook规则'
   Object.assign(formData, rule)
   dialogVisible.value = true
 }

 async function handleDelete(id: number) {
   try {
     await ElMessageBox.confirm('确定要删除该规则吗？', '提示', {
       confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning'
     })
     rules.value = rules.value.filter(r => r.id !== id)
     ElMessage.success('删除成功')
   } catch {}
 }

 function handleSubmit() {
   if (formData.id) {
     const index = rules.value.findIndex(r => r.id === formData.id)
     if (index > -1) {
       rules.value[index] = { ...rules.value[index], ...formData }
     }
     ElMessage.success('更新成功')
   } else {
     rules.value.unshift({ id: Date.now(), ...formData, count: 0, status: 'enabled' })
     ElMessage.success('创建成功')
   }
   dialogVisible.value = false
 }

 function closeDialog() {
   dialogVisible.value = false
   Object.keys(formData).forEach(key => delete formData[key])
 }
 </script>

<style scoped lang="scss">
 .page-container { background: #f0f2f5; min-height: 100%; }
 .page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
 .page-title-light { font-size: 18px; font-weight: 600; color: #1a1a1a; margin: 0; }
 .header-actions { display: flex; gap: 12px; }
 .btn-primary-light { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-radius: 6px; background: #1890ff; border: none; color: #fff; font-size: 13px; font-weight: 500; cursor: pointer; transition: all 0.2s; &:hover { background: #40a9ff; } }
 .btn-default-light { padding: 7px 14px; border-radius: 6px; background: #fff; border: 1px solid #d9d9d9; color: #595959; font-size: 13px; cursor: pointer; transition: all 0.2s; &:hover { color: #1890ff; border-color: #1890ff; } }
 .table-card { background: #fff; border-radius: 8px; border: 1px solid #f0f0f0; overflow: hidden; }
 .sync-table { width: 100%; border-collapse: collapse; th { background: #fafafa; height: 56px; padding: 0 24px; font-size: 13px; font-weight: 500; color: #8c8c8c; text-align: left; border-bottom: 1px solid #f0f0f0; } td { height: 64px; padding: 0 24px; font-size: 13px; color: #262626; border-bottom: 1px solid #f0f0f0; } }
 .badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block; margin-right: 4px; }
 .badge.platform {
   &.GitHub { background: #e6f7ff; color: #1890ff; }
   &.GitLab { background: #f6ffed; color: #52c41a; }
   &.Gitee { background: #fff7e6; color: #fa8c16; }
 }
 .event-tag { background: #f5f5f5; color: #595959; }
 .filter-tag { background: #fff7e6; color: #fa8c16; }
 .status-badge { padding: 4px 12px; border-radius: 4px; font-size: 12px; display: inline-block;
   &.enabled { background: #f6ffed; color: #52c41a; }
   &.disabled { background: #f5f5f5; color: #8c8c8c; }
 }
 .action-col { display: flex; justify-content: center; gap: 8px; }
 .icon-btn { padding: 4px 8px; border-radius: 4px; border: 1px solid #d9d9d9; background: #fff; font-size: 12px; cursor: pointer; color: #595959; &:hover { color: #1890ff; border-color: #1890ff; } &.danger:hover { color: #ff4d4f; border-color: #ff4d4f; } }
 </style>
