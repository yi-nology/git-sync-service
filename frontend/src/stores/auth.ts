import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const STORAGE_KEY = 'git-sync-api-key'

export const useAuthStore = defineStore('auth', () => {
  const apiKey = ref<string>(localStorage.getItem(STORAGE_KEY) || '')

  const isAuthenticated = computed(() => !!apiKey.value)

  function setApiKey(key: string) {
    apiKey.value = key
    localStorage.setItem(STORAGE_KEY, key)
  }

  function clearApiKey() {
    apiKey.value = ''
    localStorage.removeItem(STORAGE_KEY)
  }

  function getApiKey() {
    return apiKey.value
  }

  return {
    apiKey,
    isAuthenticated,
    setApiKey,
    clearApiKey,
    getApiKey,
  }
})
