import { defineStore } from 'pinia'
import { ref } from 'vue'
import { webhookApi } from '@/api'
import type { WebhookRule, WebhookEvent, CreateRuleRequest, UpdateRuleRequest } from '@/types'
import { ElMessage } from 'element-plus'

export const useWebhookStore = defineStore('webhook', () => {
  const rules = ref<WebhookRule[]>([])
  const events = ref<WebhookEvent[]>([])
  const loading = ref(false)

  async function fetchRules(repoKey: string) {
    loading.value = true
    try {
      const data = await webhookApi.listRules(repoKey)
      rules.value = data.rules || []
    } catch (e: any) {
      ElMessage.error(e.error || '获取规则列表失败')
    } finally {
      loading.value = false
    }
  }

  async function getRule(id: number): Promise<WebhookRule | null> {
    try {
      const data = await webhookApi.getRule(id)
      return data.rule
    } catch (e: any) {
      ElMessage.error(e.error || '获取规则详情失败')
      return null
    }
  }

  async function createRule(req: CreateRuleRequest): Promise<WebhookRule | null> {
    try {
      const data = await webhookApi.createRule(req)
      ElMessage.success('创建规则成功')
      return data.rule
    } catch (e: any) {
      ElMessage.error(e.error || '创建规则失败')
      return null
    }
  }

  async function updateRule(req: UpdateRuleRequest): Promise<WebhookRule | null> {
    try {
      const data = await webhookApi.updateRule(req)
      ElMessage.success('更新规则成功')
      return data.rule
    } catch (e: any) {
      ElMessage.error(e.error || '更新规则失败')
      return null
    }
  }

  async function deleteRule(id: number) {
    try {
      await webhookApi.deleteRule(id)
      ElMessage.success('删除规则成功')
    } catch (e: any) {
      ElMessage.error(e.error || '删除规则失败')
    }
  }

  async function fetchEvents(repoKey: string, limit = 50) {
    loading.value = true
    try {
      const data = await webhookApi.listEvents({ repo_key: repoKey, limit })
      events.value = data.events || []
    } catch (e: any) {
      ElMessage.error(e.error || '获取事件列表失败')
    } finally {
      loading.value = false
    }
  }

  async function retryEvent(id: number) {
    try {
      await webhookApi.retryEvent(id)
      ElMessage.success('事件重试已触发')
    } catch (e: any) {
      ElMessage.error(e.error || '重试事件失败')
    }
  }

  return {
    rules, events, loading,
    fetchRules, getRule, createRule, updateRule, deleteRule,
    fetchEvents, retryEvent,
  }
})
