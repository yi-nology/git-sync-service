<template>
  <div class="page-container">
    <div class="config-header">
      <div class="config-header-left">
        <div class="repo-icon-lg">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
            <polyline points="13 2 13 9 20 9"/>
          </svg>
        </div>
        <div class="repo-info-lg">
          <div class="repo-name-lg">frontend-app</div>
          <div class="repo-url-lg">https://github.com/example/frontend-app</div>
        </div>
        <span class="badge-success">已同步</span>
      </div>
      <div class="header-actions-lg">
        <button class="btn-default-lg">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
          </svg>
          编辑配置
        </button>
        <button class="btn-primary-lg">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          立即同步
        </button>
      </div>
    </div>

    <div class="config-sections">
      <div class="config-card">
        <div class="card-title">同步配置</div>
        <div class="switch-group">
          <div class="switch-item">
            <div class="switch-info">
              <div class="switch-name">启用增量同步</div>
              <div class="switch-desc">只同步变更的分支，提高效率</div>
            </div>
            <div class="switch-btn" :class="{active: config.incremental}" @click="config.incremental=!config.incremental">
              <span class="switch-dot"></span>
            </div>
          </div>
          <div class="switch-item">
            <div class="switch-info">
              <div class="switch-name">启用分布式锁</div>
              <div class="switch-desc">多实例部署时防止并发冲突</div>
            </div>
            <div class="switch-btn" :class="{active: config.distLock}" @click="config.distLock=!config.distLock">
              <span class="switch-dot"></span>
            </div>
          </div>
          <div class="switch-item">
            <div class="switch-info">
              <div class="switch-name">同步所有标签</div>
              <div class="switch-desc">同步 Git tags 到目标仓库</div>
            </div>
            <div class="switch-btn" :class="{active: config.syncTags}" @click="config.syncTags=!config.syncTags">
              <span class="switch-dot"></span>
            </div>
          </div>
          <div class="switch-item">
            <div class="switch-info">
              <div class="switch-name">强制推送</div>
              <div class="switch-desc">使用 --force 覆盖远程分支</div>
            </div>
            <div class="switch-btn" :class="{active: config.forcePush}" @click="config.forcePush=!config.forcePush">
              <span class="switch-dot"></span>
            </div>
          </div>
        </div>
      </div>

      <div class="config-card">
        <div class="card-title">Webhook 配置</div>
        <div class="webhook-info">
          <div class="webhook-url-wrap">
            <label>Webhook 地址</label>
            <div class="url-input-group">
              <input type="text" readonly :value="webhookUrl" class="url-input">
              <button class="copy-btn">复制</button>
            </div>
          </div>
          <div class="webhook-secret">
            <label>密钥 (Secret)</label>
            <input type="text" readonly :value="webhookSecret" class="secret-input">
          </div>
        </div>
        <div class="webhook-events">
          <div class="events-label">触发事件</div>
          <div class="events-tags">
            <span class="event-tag active">push</span>
            <span class="event-tag active">merge_request</span>
            <span class="event-tag">tag</span>
            <span class="event-tag">issue</span>
          </div>
        </div>
      </div>

      <div class="config-card">
        <div class="card-title">分支映射</div>
        <div class="mapping-list">
          <div class="mapping-item">
            <span class="branch-source">main</span>
            <svg class="mapping-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            <span class="branch-target">main</span>
            <span class="sync-status ok">已同步</span>
          </div>
          <div class="mapping-item">
            <span class="branch-source">develop</span>
            <svg class="mapping-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            <span class="branch-target">develop</span>
            <span class="sync-status ok">已同步</span>
          </div>
          <div class="mapping-item">
            <span class="branch-source">feature/*</span>
            <svg class="mapping-arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            <span class="branch-target">backup/*</span>
            <span class="sync-status pending">待同步</span>
          </div>
        </div>
      </div>

      <div class="config-card">
        <div class="card-title">近期同步历史</div>
        <div class="history-mini-list">
          <div class="history-mini-item" v-for="h in history" :key="h.id">
            <svg class="h-icon" :class="h.status" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline v-if="h.status==='success'" points="22 4 12 14.01 9 11.01"/>
              <circle v-if="h.status==='error'" cx="12" cy="12" r="10"/>
              <line v-if="h.status==='error'" x1="15" y1="9" x2="9" y2="15"/>
              <line v-if="h.status==='error'" x1="9" y1="9" x2="15" y2="15"/>
              <polygon v-if="h.status==='running'" points="5 3 19 12 5 21 5 3"/>
            </svg>
            <div class="h-info">
              <div class="h-branch">{{ h.branch }}</div>
              <div class="h-time">{{ h.time }}</div>
            </div>
            <span class="h-status" :class="h.status">{{ h.statusText }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'

const config = reactive({
  incremental: true,
  distLock: true,
  syncTags: false,
  forcePush: false
})

const webhookUrl = 'https://sync.example.com/webhook/github/frontend-app'
const webhookSecret = 'abc123xyz789'

const history = ref([
  { id: 1, branch: 'main', time: '2分钟前', status: 'success', statusText: '成功' },
  { id: 2, branch: 'develop', time: '15分钟前', status: 'running', statusText: '同步中' },
  { id: 3, branch: 'feature/auth', time: '1小时前', status: 'error', statusText: '失败' },
])
</script>

<style scoped lang="scss">
.page-container {
  background: #f0f2f5;
  min-height: 100%;
}

.config-header {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  .config-header-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .repo-icon-lg {
    width: 56px;
    height: 56px;
    border-radius: 8px;
    background: #e6f7ff;
    color: #1890ff;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .repo-info-lg {
    .repo-name-lg {
      font-size: 18px;
      font-weight: 600;
      color: #262626;
    }
    .repo-url-lg {
      font-size: 13px;
      color: #8c8c8c;
      margin-top: 4px;
    }
  }

  .badge-success {
    padding: 5px 12px;
    background: #f6ffed;
    color: #52c41a;
    border-radius: 4px;
    font-size: 13px;
    font-weight: 500;
    margin-left: 8px;
  }

  .header-actions-lg {
    display: flex;
    gap: 12px;
  }

  .btn-default-lg {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 6px;
    background: #fff;
    border: 1px solid #d9d9d9;
    color: #595959;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      color: #1890ff;
      border-color: #1890ff;
    }
  }

  .btn-primary-lg {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 6px;
    background: #1890ff;
    border: none;
    color: #fff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #40a9ff;
    }
  }
}

.config-sections {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-card {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;

  .card-title {
    font-size: 15px;
    font-weight: 600;
    color: #262626;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f0f0;
  }
}

.switch-group {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .switch-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;

    &:not(:last-child) {
      border-bottom: 1px solid #fafafa;
    }

    .switch-info {
      .switch-name {
        font-size: 14px;
        font-weight: 500;
        color: #262626;
      }
      .switch-desc {
        font-size: 12px;
        color: #8c8c8c;
        margin-top: 4px;
      }
    }
  }
}

.switch-btn {
  width: 44px;
  height: 24px;
  border-radius: 12px;
  background: #d9d9d9;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  &.active {
    background: #52c41a;
  }

  .switch-dot {
    position: absolute;
    left: 2px;
    top: 2px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #fff;
    transition: all 0.2s;
  }

  &.active .switch-dot {
    transform: translateX(20px);
  }
}

.webhook-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;

  label {
    display: block;
    font-size: 13px;
    font-weight: 500;
    color: #595959;
    margin-bottom: 8px;
  }

  .url-input-group {
    display: flex;
    gap: 8px;

    .url-input {
      flex: 1;
      height: 36px;
      padding: 0 12px;
      border: 1px solid #d9d9d9;
      border-radius: 6px;
      font-size: 13px;
      color: #8c8c8c;
      background: #fafafa;
    }

    .copy-btn {
      padding: 0 16px;
      height: 36px;
      border-radius: 6px;
      background: #fff;
      border: 1px solid #d9d9d9;
      color: #595959;
      font-size: 13px;
      cursor: pointer;

      &:hover {
        color: #1890ff;
        border-color: #1890ff;
      }
    }
  }

  .secret-input {
    width: 100%;
    height: 36px;
    padding: 0 12px;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    font-size: 13px;
    color: #8c8c8c;
    background: #fafafa;
  }
}

.webhook-events {
  .events-label {
    font-size: 13px;
    font-weight: 500;
    color: #595959;
    margin-bottom: 8px;
  }

  .events-tags {
    display: flex;
    gap: 8px;

    .event-tag {
      padding: 5px 12px;
      background: #f5f5f5;
      border-radius: 4px;
      font-size: 12px;
      color: #8c8c8c;
      cursor: pointer;
      transition: all 0.2s;

      &.active {
        background: #e6f7ff;
        color: #1890ff;
      }
    }
  }
}

.mapping-list {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .mapping-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: #fafafa;
    border-radius: 6px;

    .branch-source, .branch-target {
      padding: 4px 10px;
      background: #fff;
      border: 1px solid #d9d9d9;
      border-radius: 4px;
      font-size: 13px;
      color: #262626;
    }

    .mapping-arrow {
      color: #bfbfbf;
    }

    .sync-status {
      margin-left: auto;
      font-size: 12px;

      &.ok {
        color: #52c41a;
      }

      &.pending {
        color: #faad14;
      }
    }
  }
}

.history-mini-list {
  display: flex;
  flex-direction: column;
  gap: 8px;

  .history-mini-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    background: #fafafa;
    border-radius: 6px;

    .h-icon {
      &.success {
        color: #52c41a;
      }
      &.error {
        color: #ff4d4f;
      }
      &.running {
        color: #1890ff;
      }
    }

    .h-info {
      flex: 1;

      .h-branch {
        font-size: 13px;
        font-weight: 500;
        color: #262626;
      }
      .h-time {
        font-size: 12px;
        color: #8c8c8c;
      }
    }

    .h-status {
      font-size: 12px;

      &.success {
        color: #52c41a;
      }
      &.error {
        color: #ff4d4f;
      }
      &.running {
        color: #1890ff;
      }
    }
  }
}
</style>

