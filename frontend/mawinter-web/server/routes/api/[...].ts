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

  // すべてのHTTPメソッド（GET, POST, DELETE等）とヘッダーをそのまま転送
  return proxyRequest(event, targetUrl)
})
