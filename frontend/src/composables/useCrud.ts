import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

export interface CrudOptions<T> {
  name: string
  fetchApi?: () => Promise<T[]>
  createApi?: (data: T) => Promise<T>
  updateApi?: (data: T) => Promise<T>
  deleteApi?: (id: number) => Promise<void>
}

export function useCrud<T extends { id?: number }>(options: CrudOptions<T>) {
  const dataList = ref<T[]>([])
  const loading = ref(false)
  const dialogVisible = ref(false)
  const dialogTitle = ref('')
  const formData = reactive<Partial<T>>({})
  const isEdit = computed(() => !!formData.id)

  async function fetchData() {
    if (!options.fetchApi) return
    loading.value = true
    try {
      dataList.value = await options.fetchApi()
    } catch (e) {
      ElMessage.error('获取数据失败')
    } finally {
      loading.value = false
    }
  }

  function openCreate(defaultData?: Partial<T>) {
    dialogTitle.value = `新建${options.name}`
    Object.assign(formData, { id: undefined, ...defaultData } as Partial<T>)
    dialogVisible.value = true
  }

  function openEdit(row: T) {
    dialogTitle.value = `编辑${options.name}`
    Object.assign(formData, row)
    dialogVisible.value = true
  }

  async function handleDelete(id: number) {
    try {
      await ElMessageBox.confirm(`确定要删除该${options.name}吗？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      if (options.deleteApi) {
        await options.deleteApi(id)
      }
      dataList.value = dataList.value.filter((item: any) => item.id !== id)
      ElMessage.success('删除成功')
    } catch {
      // 用户取消
    }
  }

  async function handleSubmit() {
    try {
      if (isEdit.value && options.updateApi) {
        await options.updateApi(formData as T)
        const index = dataList.value.findIndex((item: any) => item.id === formData.id)
        if (index > -1) {
          dataList.value[index] = formData as T
        }
        ElMessage.success('更新成功')
      } else if (options.createApi) {
        const result = await options.createApi(formData as T)
        dataList.value.push(result)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
    } catch (e) {
      ElMessage.error('操作失败')
    }
  }

  function closeDialog() {
    dialogVisible.value = false
    Object.keys(formData).forEach(key => {
      delete (formData as any)[key]
    })
  }

  return {
    dataList,
    loading,
    dialogVisible,
    dialogTitle,
    formData,
    isEdit,
    fetchData,
    openCreate,
    openEdit,
    handleDelete,
    handleSubmit,
    closeDialog
  }
}
