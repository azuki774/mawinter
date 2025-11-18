import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http'
import { NodeSDK } from '@opentelemetry/sdk-node'
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions'

// @opentelemetry/resources is CommonJS, so import via default export
import * as resources from '@opentelemetry/resources'

let sdk: NodeSDK | null = null

// ユーザー入力の OTLP_SERVER を正規化して HTTP Exporter が理解できる URL に変換
const buildTraceEndpoint = (raw?: string | null): string | null => {
  const trimmed = raw?.trim()
  if (!trimmed) {
    return null
  }

  if (trimmed.includes('://')) {
    try {
      const url = new URL(trimmed)
      if (!url.hostname) {
        return null
      }

      if (!url.pathname || url.pathname === '/') {
        url.pathname = '/v1/traces'
      }

      return url.toString()
    } catch (error) {
      console.error('[telemetry] OTLP_SERVER の値が不正です', error)
      return null
    }
  }

  const normalized = trimmed.replace(/\/+$/, '')
  if (!normalized) {
    return null
  }

  return `http://${normalized}/v1/traces`
}

export default defineNitroPlugin(async (nitroApp) => {
  if (sdk) {
    return
  }

  const endpoint = buildTraceEndpoint(process.env.OTLP_SERVER)
  if (!endpoint) {
    console.info('[telemetry] OTLP_SERVER が未設定のため、SSR のトレース送信は無効です')
    return
  }

  const { resourceFromAttributes, defaultResource } = resources
  const attributes = {
    [SemanticResourceAttributes.SERVICE_NAME]: 'mawinter-api',
    [SemanticResourceAttributes.SERVICE_VERSION]:
      process.env.NUXT_PUBLIC_APP_VERSION || process.env.npm_package_version || 'dev',
  }
  const baseResource = resourceFromAttributes(attributes)
  const resource = defaultResource && typeof defaultResource.merge === 'function'
    ? defaultResource.merge(baseResource)
    : baseResource

  sdk = new NodeSDK({
    resource,
    traceExporter: new OTLPTraceExporter({ url: endpoint }),
    instrumentations: [new HttpInstrumentation()],
  })

  try {
    await sdk.start()
    console.info('[telemetry] OpenTelemetry tracing enabled', endpoint)
  } catch (error) {
    console.error('[telemetry] OpenTelemetry の初期化に失敗しました', error)
    sdk = null
    return
  }

  nitroApp.hooks.hook('close', async () => {
    if (!sdk) {
      return
    }

    try {
      await sdk.shutdown()
      console.info('[telemetry] OpenTelemetry tracing stopped')
    } catch (error) {
      console.error('[telemetry] OpenTelemetry の終了処理に失敗しました', error)
    } finally {
      sdk = null
    }
  })
})
