<template>
  <el-dialog v-model="visible" :title="null" width="600px" :close-on-click-modal="false" :show-close="false" class="edit-task-modal">
    <div class="modal-header">
      <span class="header-title">{{ isEdit ? '编辑同步任务' : '新建同步任务' }}</span>
      <button class="close-btn" @click="handleClose">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#8c8c8c" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"/>
          <line x1="6" y1="6" x2="18" y2="18"/>
        </svg>
      </button>
    </div>
    
    <div class="modal-body">
      <div class="form-item">
        <label class="form-label">任务名称<span class="required">*</span></label>
        <input type="text" v-model="form.name" class="form-input" placeholder="请输入任务名称" />
      </div>
      
      <div class="form-row">
        <div class="form-item half">
          <label class="form-label">源分支</label>
          <input type="text" v-model="form.sourceBranch" class="form-input" placeholder="main" />
        </div>
        <div class="form-item half">
          <label class="form-label">目标分支</label>
          <input type="text" v-model="form.targetBranch" class="form-input" placeholder="release" />
        </div>
      </div>
      
      <div class="form-item">
        <label class="form-label">同步模式</label>
        <div class="mode-group">
          <div class="mode-btn" :class="{ active: form.mode === 'single' }" @click="form.mode = 'single'">
            <div class="radio-dot" :class="{ checked: form.mode === 'single' }"></div>
            <span>单分支</span>
          </div>
          <div class="mode-btn" :class="{ active: form.mode === 'all' }" @click="form.mode = 'all'">
            <div class="radio-dot" :class="{ checked: form.mode === 'all' }"></div>
            <span>全分支</span>
          </div>
          <div class="mode-btn" :class="{ active: form.mode === 'regex' }" @click="form.mode = 'regex'">
            <div class="radio-dot" :class="{ checked: form.mode === 'regex' }"></div>
            <span>正则匹配</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="modal-footer">
      <button class="btn-cancel" @click="handleClose">取消</button>
      <button class="btn-save" @click="handleSave">保存修改</button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, computed, watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
  task?: any
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', data: any): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const form = reactive({
  name: '',
  sourceBranch: '',
  targetBranch: '',
  mode: 'single',
})

const isEdit = computed(() => !!props.task?.id)

watch(() => props.modelValue, (val) => {
  if (val && props.task) {
    Object.assign(form, {
      name: props.task?.name || '',
      sourceBranch: props.task?.sourceBranch || props.task?.source || '',
      targetBranch: props.task?.targetBranch || props.task?.target || '',
      mode: props.task?.mode === '单分支' ? 'single' : props.task?.mode || 'single',
    })
  } else if (val) {
    Object.assign(form, {
      name: '',
      sourceBranch: '',
      targetBranch: '',
      mode: 'single',
    })
  }
})

function handleClose() {
  visible.value = false
}

function handleSave() {
  emit('confirm', { ...form })
  visible.value = false
}
</script>

<style scoped lang="scss">
.edit-task-modal :deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

.edit-task-modal :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
  border-bottom: 1px solid #f0f0f0;
}

.edit-task-modal :deep(.el-dialog__body) {
  padding: 0;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  font-family: 'Inter';
}

.close-btn {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  border: none;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: #f5f5f5;
  }
}

.modal-body {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-row {
  display: flex;
  gap: 16px;
  
  .form-item.half {
    flex: 1;
  }
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  font-family: 'Inter';
}

.required {
  color: #ff4d4f;
  margin-left: 4px;
}

.form-input {
  width: 100%;
  height: 40px;
  padding: 0 12px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  color: #262626;
  background: #fff;
  box-sizing: border-box;
  transition: all 0.2s;
  font-family: 'Inter';
  
  &:focus {
    outline: none;
    border-color: #1890ff;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
  }
  
  &::placeholder {
    color: #bfbfbf;
  }
}

.mode-group {
  display: flex;
  gap: 12px;
}

.mode-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 40px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  cursor: pointer;
  font-size: 13px;
  color: #595959;
  font-family: 'Inter';
  transition: all 0.2s;
  
  &.active {
    background: #e6f7ff;
    border-color: #1890ff;
    color: #1890ff;
  }
}

.radio-dot {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1px solid #d9d9d9;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  
  &.checked {
    background: #1890ff;
    border-color: #1890ff;
    
    &::after {
      content: '';
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #fff;
    }
  }
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px 24px;
  border-top: 1px solid #f0f0f0;
}

.btn-cancel {
  min-width: 100px;
  height: 40px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  color: #595959;
  font-size: 14px;
  font-family: 'Inter';
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    color: #1890ff;
    border-color: #1890ff;
  }
}

.btn-save {
  min-width: 120px;
  height: 40px;
  border-radius: 6px;
  background: #1890ff;
  border: none;
  color: #fff;
  font-size: 14px;
  font-family: 'Inter';
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: #40a9ff;
  }
}
</style>
