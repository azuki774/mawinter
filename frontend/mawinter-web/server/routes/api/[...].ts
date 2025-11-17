import {
  context,
  propagation,
  trace,
  SpanStatusCode,
  type TextMapSetter,
} from '@opentelemetry/api'

const tracer = trace.getTracer('mawinter-web')

const headerSetter: TextMapSetter<Record<string, string>> = {
  set(carrier, key, value) {
    carrier[key.toLowerCase()] = value
  },
}

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
export default defineEventHandler((event) => {
  const config = useRuntimeConfig()

  return tracer.startActiveSpan(
    `proxy ${event.node.req.method ?? 'HTTP'} ${event.path}`,
    async (span) => {
      try {
        const path = event.path
        const targetUrl = `${config.mawinterApiUrl}${path}`

        const headers: Record<string, string> = {}
        for (const [key, value] of Object.entries(event.node.req.headers)) {
          if (typeof value === 'string') {
            headers[key] = value
          } else if (Array.isArray(value)) {
            headers[key] = value.join(',')
          }
        }

        propagation.inject(context.active(), headers, headerSetter)

        return await proxyRequest(event, targetUrl, { headers })
      } catch (error) {
        span.recordException(error as Error)
        span.setStatus({ code: SpanStatusCode.ERROR })
        throw error
      } finally {
        span.end()
      }
    },
  )
})
