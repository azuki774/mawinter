import {
  context,
  propagation,
  trace,
  SpanStatusCode,
  type TextMapSetter,
} from '@opentelemetry/api'

const headerSetter: TextMapSetter<Record<string, string | string[] | undefined>> = {
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
  const tracer = trace.getTracer('mawinter-web')

  return tracer.startActiveSpan(
    `proxy ${event.node.req.method ?? 'HTTP'} ${event.path}`,
    async (span) => {
      try {
        const path = event.path
        const targetUrl = `${config.mawinterApiUrl}${path}`

        const spanContext = trace.setSpan(context.active(), span)
        propagation.inject(
          spanContext,
          event.node.req.headers as Record<string, string | string[] | undefined>,
          headerSetter,
        )

        return await proxyRequest(event, targetUrl)
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
