import { context, propagation } from '@opentelemetry/api'

/**
 * APIプロキシサーバー
 *
 * すべての /api/* へのリクエストをバックエンドAPIサーバーに転送します。
 * これにより、ブラウザから直接APIサーバーにアクセスすることを防ぎ、
 * セキュリティを向上させます。
 *
 * 例:
 * - ブラウザ → /api/v3/categories
 * - Nuxtサーバー → http://localhost:8080/api/v3/categories
 */
export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()

  // リクエストパス全体を取得（例: /api/v3/categories）
  const path = event.path

  // バックエンドAPIの完全なURLを構築
  const targetUrl = `${config.mawinterApiUrl}${path}`

  // 既存ヘッダーをコピーし、OpenTelemetry のトレース情報を注入
  const headers: Record<string, string> = {}
  for (const [key, value] of Object.entries(event.req.headers)) {
    if (typeof value === 'string') {
      headers[key] = value
    } else if (Array.isArray(value)) {
      headers[key] = value.join(',')
    }
  }
  propagation.inject(context.active(), headers)

  // すべてのHTTPメソッドとカスタムヘッダーを転送
  return proxyRequest(event, targetUrl, { headers })
})
