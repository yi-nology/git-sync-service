<template>
  <div class="page-container">
    <PageHeader title="仓库管理">
      <template #actions>
        <a-button type="primary" @click="openCreate">
          <template #icon><PlusOutlined /></template>
          添加仓库
        </a-button>
      </template>
    </PageHeader>

    <a-card>
      <template #extra>
        <a-input
          v-model:value="searchText"
          placeholder="搜索仓库..."
          allow-clear
          style="width: 240px"
        >
          <template #prefix><SearchOutlined /></template>
        </a-input>
      </template>

      <a-table
        :columns="columns"
        :data-source="filteredRepos"
        :loading="repoStore.loading"
        row-key="key"
        :scroll="{ x: 900 }"
        :pagination="false"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div class="repo-name-cell">
              <div class="repo-name">{{ record.name }}</div>
              <div class="repo-url">{{ record.clone_url }}</div>
            </div>
          </template>

          <template v-if="column.key === 'platform'">
            <a-tag color="blue">{{ record.platform }}</a-tag>
          </template>

          <template v-if="column.key === 'status'">
            <StatusBadge :status="record.status" />
          </template>

          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button size="small" @click="testConn(record.key)">测试连接</a-button>
              <a-button size="small" @click="openEdit(record)">编辑</a-button>
              <a-popconfirm
                title="确定要删除该仓库吗？"
                ok-text="确定"
                cancel-text="取消"
                @confirm="handleDelete(record.key)"
              >
                <a-button size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>

        <template #emptyText>
          <a-empty description="暂无仓库数据" />
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      :width="500"
      @ok="handleSubmit"
      ok-text="确定"
      cancel-text="取消"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item label="仓库名称" required>
          <a-input v-model:value="formData.name" placeholder="请输入仓库名称" />
        </a-form-item>
        <a-form-item label="仓库地址" required>
          <a-input v-model:value="formData.remote_url" placeholder="请输入仓库地址" />
        </a-form-item>
        <a-form-item label="访问令牌">
          <a-input-password v-model:value="formData.access_token" placeholder="请输入访问令牌" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { useRepoStore } from '@/stores/repo'
import type { Repo } from '@/types'

const repoStore = useRepoStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const editingKey = ref('')
const searchText = ref('')

const formData = reactive({
  name: '',
  remote_url: '',
  access_token: '',
})

const columns = [
  { title: '仓库', key: 'name' },
  { title: '平台', key: 'platform', width: 100 },
  { title: '所有者', dataIndex: 'platform_owner', key: 'owner', width: 120 },
  { title: '仓库名', dataIndex: 'platform_repo', key: 'repo', width: 120 },
  { title: '默认分支', dataIndex: 'default_branch', key: 'branch', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'actions', width: 200, fixed: 'right' as const },
]

const filteredRepos = computed(() => {
  if (!searchText.value) return repoStore.repos
  const search = searchText.value.toLowerCase()
  return repoStore.repos.filter(
    repo => repo.name.toLowerCase().includes(search) ||
            repo.clone_url.toLowerCase().includes(search)
  )
})

onMounted(() => {
  repoStore.fetchRepos()
})

function openCreate() {
  dialogTitle.value = '添加仓库'
  editingKey.value = ''
  Object.assign(formData, { name: '', remote_url: '', access_token: '' })
  dialogVisible.value = true
}

function openEdit(repo: Repo) {
  dialogTitle.value = '编辑仓库'
  editingKey.value = repo.key
  Object.assign(formData, { name: repo.name, remote_url: repo.clone_url, access_token: '' })
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formData.name || !formData.remote_url) {
    message.warning('请填写仓库名称和地址')
    return
  }
  if (editingKey.value) {
    await repoStore.updateRepo({ key: editingKey.value, name: formData.name, access_token: formData.access_token })
  } else {
    await repoStore.createRepo(formData)
  }
  dialogVisible.value = false
}

async function handleDelete(key: string) {
  await repoStore.deleteRepo(key)
}

async function testConn(key: string) {
  const result = await repoStore.testConnection(key)
  if (result) {
    message[result.success ? 'success' : 'error'](result.message)
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

.repo-name-cell {
  .repo-name {
    font-weight: 500;
    color: $text-primary;
  }

  .repo-url {
    font-size: 12px;
    color: $text-secondary;
    margin-top: $spacing-xs;
  }
}
</style>
