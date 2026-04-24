/**
 * 日期格式化工具
 * 使用原生 Intl.DateTimeFormat API，无需额外依赖
 */

/**
 * 格式化日期为本地化字符串
 * @param dateStr - ISO 8601 日期字符串或 Date 对象
 * @param format - 格式类型：'datetime'（默认）, 'date', 'time', 'relative'
 * @returns 格式化后的日期字符串
 */
export function formatDate(dateStr: string | Date | null | undefined, format: 'datetime' | 'date' | 'time' = 'datetime'): string {
  if (!dateStr) return '-'

  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  if (isNaN(date.getTime())) return '-'

  const locale = 'zh-CN'

  switch (format) {
    case 'datetime':
      return new Intl.DateTimeFormat(locale, {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      }).format(date)
    case 'date':
      return new Intl.DateTimeFormat(locale, {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      }).format(date)
    case 'time':
      return new Intl.DateTimeFormat(locale, {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false,
      }).format(date)
    default:
      return date.toISOString()
  }
}

/**
 * 格式化相对时间（如 "3 小时前"）
 * @param dateStr - ISO 8601 日期字符串或 Date 对象
 * @returns 相对时间字符串
 */
export function formatRelativeTime(dateStr: string | Date | null | undefined): string {
  if (!dateStr) return '-'

  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  if (isNaN(date.getTime())) return '-'

  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 60) return '刚刚'
  if (diffMin < 60) return `${diffMin} 分钟前`
  if (diffHour < 24) return `${diffHour} 小时前`
  if (diffDay < 30) return `${diffDay} 天前`

  return formatDate(dateStr, 'date')
}
