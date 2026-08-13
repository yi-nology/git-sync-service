<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo-icon">
          <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#1677FF" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
          </svg>
        </div>
        <h1>Git Sync Service</h1>
        <p>代码同步管理系统</p>
      </div>

      <a-form
        :model="formState"
        name="loginForm"
        @finish="handleLogin"
        class="login-form"
      >
        <a-form-item
          name="apiKey"
          :rules="[{ required: true, message: '请输入 API Key' }]"
        >
          <a-input-password
            v-model:value="formState.apiKey"
            placeholder="请输入 API Key"
            size="large"
            autofocus
          >
            <template #prefix>
              <LockOutlined style="color: #BFBFBF;" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>

      <div class="login-footer">
        <div class="footer-tip">
          <InfoCircleOutlined />
          <span>请向管理员获取 API Key</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { LockOutlined, InfoCircleOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { systemApi } from '@/api'
import { notifySuccess, notifyError } from '@/utils/notify'
import { ApiError } from '@/api/http'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const formState = reactive({
  apiKey: '',
})

const handleLogin = async () => {
  loading.value = true
  try {
    authStore.setApiKey(formState.apiKey)
    // 探测 API Key 有效性(走统一 axios 实例,自动注入 X-API-Key)
    await systemApi.status()
    notifySuccess('登录成功')
    router.push('/dashboard')
  } catch (e) {
    authStore.clearApiKey()
    if (e instanceof ApiError && e.code === 401) {
      // API Key 无效/过期:给出明确中文提示,避免裸露的后端英文错误
      notifyError('API Key 无效或已过期，请检查后重新输入')
    } else {
      // 网络错误等
      notifyError('无法连接到服务器，请检查服务是否已启动')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

.login-card {
  width: 420px;
  max-width: 100%;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 36px;

  .logo-icon {
    margin-bottom: 16px;
  }

  h1 {
    font-size: 26px;
    color: #1a1a1a;
    margin: 0 0 8px 0;
    font-weight: 700;
    letter-spacing: -0.5px;
  }

  p {
    font-size: 14px;
    color: #8C8C8C;
    margin: 0;
  }
}

.login-form {
  margin-bottom: 8px;

  :deep(.ant-input-affix-wrapper) {
    border-radius: 8px;
  }

  :deep(.ant-btn) {
    height: 44px;
    border-radius: 8px;
    font-size: 15px;
    font-weight: 500;
  }
}

.login-footer {
  text-align: center;
  margin-top: 24px;

  .footer-tip {
    font-size: 13px;
    color: #8C8C8C;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    margin-bottom: 12px;
  }

  .footer-hint {
    font-size: 12px;
    color: #BFBFBF;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 4px;

    .ant-tag {
      margin: 0;
      padding: 0 6px;
      border-radius: 4px;
    }
  }
}
</style>
