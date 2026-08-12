import { message } from 'ant-design-vue'
import { ApiError } from '@/api/http'

/** 统一错误提示:把任意异常转成可读文案后弹出 */
export function notifyError(e: unknown, fallback = '操作失败'): void {
  const text =
    e instanceof ApiError
      ? e.message
      : e instanceof Error
        ? e.message
        : fallback
  message.error(text)
}

export function notifySuccess(text: string): void {
  message.success(text)
}

export function notifyWarning(text: string): void {
  message.warning(text)
}

export function notifyInfo(text: string): void {
  message.info(text)
}
