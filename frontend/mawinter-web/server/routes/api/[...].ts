/**
 * APIプロキシサーバー
 *
 * すべての /api/* へのリクエストをバックエンドAPIサーバーに転送します。
 * これにより、ブラウザから直接APIサーバーにアクセスすることを防ぎ、
 * セキュリティを向上させます。
 *
 * OpenTelemetry UndiciInstrumentation でトレースされるよう、
 * 明示的に fetch() を使用してプロキシを実装しています。
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

  // リクエストヘッダーを準備（一部のヘッダーは除外）
  const headers = new Headers()
  const requestHeaders = getRequestHeaders(event)
  for (const [key, value] of Object.entries(requestHeaders)) {
    // host, connection などのヘッダーは転送しない
    if (value && !['host', 'connection'].includes(key.toLowerCase())) {
      headers.set(key, value)
    }
  }

  // リクエストボディを取得（GET, HEAD以外）
  let body = undefined
  if (event.method && !['GET', 'HEAD'].includes(event.method)) {
    body = await readRawBody(event)
  }

  // バックエンドAPIにリクエストを転送（fetch() を使用してトレース対象に）
  const response = await fetch(targetUrl, {
    method: event.method,
    headers,
    body,
  })

  // レスポンスヘッダーを設定
  response.headers.forEach((value, key) => {
    setResponseHeader(event, key, value)
  })

  // レスポンスステータスを設定
  setResponseStatus(event, response.status)

  // レスポンスボディを返す
  return response.body
})
