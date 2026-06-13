<template>
  <div class="sync-history page-container dark">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title" style="color:#fff;font-size:24px;font-weight:700;margin:0;">同步历史</h1>
        <p style="color:#94a3b8;font-size:14px;margin:4px 0 0 0;">查看所有同步执行记录</p>
      </div>
      <div class="header-actions" style="display:flex;gap:12px;">
        <button class="btn-default-dark">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          导出历史
        </button>
        <button class="btn-default-dark" style="color:#ff4d4f;border-color:rgba(255,77,79,0.3);">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="3 6 5 6 21 6"></polyline>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          </svg>
          清空历史
        </button>
      </div>
    </div>

    <div class="stats-row-dark">
      <div class="stat-card-dark primary">
        <div class="stat-header">
          <span class="stat-label">总执行次数</span>
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
          </svg>
        </div>
        <div class="stat-value">1,234</div>
        <div class="stat-change">↑12% 较上月</div>
      </div>
      <div class="stat-card-dark success">
        <div class="stat-header">
          <span class="stat-label">成功</span>
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
            <polyline points="22 4 12 14.01 9 11.01"></polyline>
          </svg>
        </div>
        <div class="stat-value">1,180</div>
        <div class="stat-change">95.6% 成功率</div>
      </div>
      <div class="stat-card-dark error">
        <div class="stat-header">
          <span class="stat-label">失败</span>
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="15" y1="9" x2="9" y2="15"></line>
            <line x1="9" y1="9" x2="15" y2="15"></line>
          </svg>
        </div>
        <div class="stat-value">42</div>
        <div class="stat-change" style="color:#ff4d4f;">+5 较上月</div>
      </div>
      <div class="stat-card-dark warning">
        <div class="stat-header">
          <span class="stat-label">平均耗时</span>
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
        </div>
        <div class="stat-value">45s</div>
        <div class="stat-change">-8% 较上月</div>
      </div>
    </div>

    <div class="filter-bar-dark">
      <div class="filter-item" style="width:360px;">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input type="text" placeholder="搜索历史记录..." v-model="filters.keyword">
      </div>
      <div class="filter-item">
        <span>状态</span>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </div>
      <div class="filter-item">
        <span>任务</span>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </div>
      <div class="filter-item">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
          <line x1="16" y1="2" x2="16" y2="6"></line>
          <line x1="8" y1="2" x2="8" y2="6"></line>
          <line x1="3" y1="10" x2="21" y2="10"></line>
        </svg>
        <span>日期范围</span>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </div>
    </div>

    <div class="history-list">
      <div class="history-card-dark">
        <div class="card-header">
          <div class="header-left">
            <div class="status-icon success">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                <polyline points="22 4 12 14.01 9 11.01"></polyline>
              </svg>
            </div>
            <div class="title-wrap">
              <div class="task-name">frontend-sync</div>
              <div class="task-time">GitHub main → GitLab main</div>
            </div>
          </div>
          <div class="status-badge success">成功</div>
        </div>
        <div class="card-body">
          <div class="sync-flow">
            <div class="endpoint">GitHub / main</div>
            <svg class="arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <div class="endpoint">GitLab / main</div>
            <span class="meta">5 个提交   14:32:15 - 14:32:27</span>
          </div>
          <div class="commits-wrap">
            <div class="commits-label">提交记录</div>
            <div class="commit-item">
              <span class="commit-hash">a1b2c3d</span>
              <span class="commit-msg">fix: 修复登录页面样式问题</span>
              <span class="commit-author">@zhangwei</span>
            </div>
            <div class="commit-item">
              <span class="commit-hash">d4e5f67</span>
              <span class="commit-msg">feat: 新增用户权限管理模块</span>
              <span class="commit-author">@admin</span>
            </div>
          </div>
          <div class="log-toggle">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
            查看执行日志
          </div>
        </div>
      </div>

      <div class="history-card-dark">
        <div class="card-header">
          <div class="header-left">
            <div class="status-icon error">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
            </div>
            <div class="title-wrap">
              <div class="task-name">repo-mirror-all</div>
              <div class="task-time">全分支同步任务</div>
            </div>
          </div>
          <div class="status-badge error">失败</div>
        </div>
        <div class="card-body">
          <div class="sync-flow">
            <div class="endpoint">GitLab</div>
            <svg class="arrow" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"></polyline>
            </svg>
            <div class="endpoint">Gitee</div>
            <span class="meta">8 个分支   6 成功 / 2 失败</span>
          </div>
          <div class="error-wrap">
            <div class="error-header">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="12" y1="8" x2="12" y2="12"></line>
                <line x1="12" y1="16" x2="12.01" y2="16"></line>
              </svg>
              同步错误详情
            </div>
            <div class="error-message">推送到 backup/feature-auth 失败: remote rejected (pre-receive hook declined)
推送到 backup/hotfix-v2 失败: conflict detected on remote</div>
          </div>
          <div class="log-toggle">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="18 15 12 9 6 15"></polyline>
            </svg>
            收起执行日志
          </div>
          <div class="terminal-log">
            <div class="log-line"><span class="log-time">14:30:01</span> [INFO] Starting sync process...</div>
            <div class="log-line"><span class="log-time">14:30:02</span> [INFO] Fetching from origin...</div>
            <div class="log-line"><span class="log-time">14:30:05</span> [INFO] 6 branches fetched successfully</div>
            <div class="log-line"><span class="log-time">14:30:08</span> [ERROR] Push to feature-auth failed: pre-receive hook declined</div>
            <div class="log-line"><span class="log-time">14:30:12</span> [ERROR] Push to hotfix-v2 failed: conflict detected</div>
            <div class="log-line"><span class="log-time">14:30:15</span> [ERROR] Sync completed with errors</div>
          </div>
        </div>
      </div>
    </div>

    <div class="pagination-dark">
      <span>显示 1-10 共 156 条记录</span>
      <div class="page-buttons">
        <button><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"></polyline></svg></button>
        <button class="active">1</button>
        <button>2</button>
        <button>3</button>
        <button><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg></button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'

const filters = reactive({ keyword: '' })
</script>

<style scoped lang="scss">
.sync-history {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;

    .header-actions { display: flex; gap: 12px; }
  }

  .btn-default-dark {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    border-radius: 8px;
    background: #1e293b;
    border: 1px solid #334155;
    color: #ffffff;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #1890ff;
      color: #1890ff;
    }
  }

  .history-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
}
</style>
