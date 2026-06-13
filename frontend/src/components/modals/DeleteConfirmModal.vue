<template>
  <el-dialog
    v-model="visible"
    :title="null"
    width="500px"
    :close-on-click-modal="false"
    :show-close="false"
    class="delete-confirm-modal"
  >
    <div class="modal-header">
      <svg class="warn-icon" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#faad14" stroke-width="2">
        <path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
        <line x1="12" y1="9" x2="12" y2="13"/>
        <line x1="12" y1="17" x2="12.01" y2="17"/>
      </svg>
      <span class="modal-title">确认删除</span>
    </div>
    
    <div class="modal-body">
      <p class="desc-text">您确定要删除以下同步任务吗？</p>
      <p class="task-name">{{ taskName }}</p>
      
      <div class="warn-box">
        <p class="warn-title">⚠️ 删除影响</p>
        <p class="warn-desc">此操作将永久删除任务配置，历史同步记录会保留。删除后无法恢复。</p>
      </div>
      
      <div class="confirm-check" @click="checked = !checked">
        <div class="checkbox" :class="{ checked }">
          <svg v-if="checked" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="3">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </div>
        <span class="check-text">我已确认，此操作无法撤销</span>
      </div>
    </div>
    
    <div class="modal-footer">
      <button class="btn-cancel" @click="handleCancel">取消</button>
      <button class="btn-delete" :disabled="!checked" @click="handleDelete">确认删除</button>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  modelValue: boolean
  taskName: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const checked = ref(false)

function handleCancel() {
  visible.value = false
  emit('cancel')
}

function handleDelete() {
  if (!checked.value) return
  visible.value = false
  emit('confirm')
}
</script>

<style scoped lang="scss">
.delete-confirm-modal :deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

.delete-confirm-modal :deep(.el-dialog__body) {
  padding: 0;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 24px 16px;
}

.warn-icon {
  flex-shrink: 0;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

.modal-body {
  padding: 0 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.desc-text {
  font-size: 14px;
  color: #595959;
  margin: 0;
}

.task-name {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.warn-box {
  background: #fff7e6;
  border-radius: 6px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.warn-title {
  font-size: 13px;
  font-weight: 500;
  color: #fa8c16;
  margin: 0;
}

.warn-desc {
  font-size: 12px;
  color: #8c8c8c;
  margin: 0;
  line-height: 1.5;
}

.confirm-check {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  cursor: pointer;
}

.checkbox {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  background: #fff;
  border: 1px solid #d9d9d9;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.2s;
  
  &.checked {
    background: #1890ff;
    border-color: #1890ff;
  }
}

.check-text {
  font-size: 13px;
  color: #595959;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px 24px;
}

.btn-cancel {
  min-width: 100px;
  height: 40px;
  border-radius: 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  color: #595959;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    color: #1890ff;
    border-color: #1890ff;
  }
}

.btn-delete {
  min-width: 120px;
  height: 40px;
  border-radius: 6px;
  background: #ff4d4f;
  border: none;
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  &:not(:disabled):hover {
    background: #ff7875;
  }
}
</style>
