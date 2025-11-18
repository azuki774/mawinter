/**
 * OpenTelemetry トレーシング初期化
 *
 * このファイルはアプリケーション起動前に読み込む必要があります:
 * node --require ./tracing.mjs .output/server/index.mjs
 *
 * 環境変数:
 * - OTLP_SERVER: OTLP コレクタのエンドポイント (例: localhost:4318)
 */

import { NodeSDK } from '@opentelemetry/sdk-node'
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { Resource } from '@opentelemetry/resources'
import { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } from '@opentelemetry/semantic-conventions'

// ユーザー入力の OTLP_SERVER を正規化して HTTP Exporter が理解できる URL に変換
function buildTraceEndpoint(raw) {
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
    }
    catch (error) {
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

const endpoint = buildTraceEndpoint(process.env.OTLP_SERVER)

if (!endpoint) {
  console.info('[telemetry] OTLP_SERVER が未設定のため、SSR のトレース送信は無効です')
}
else {
  const sdk = new NodeSDK({
    resource: new Resource({
      [ATTR_SERVICE_NAME]: 'mawinter-web',
      [ATTR_SERVICE_VERSION]: process.env.npm_package_version || 'dev',
    }),
    traceExporter: new OTLPTraceExporter({ url: endpoint }),
    instrumentations: [
      getNodeAutoInstrumentations({
        // undici, http, https などを自動計装
        '@opentelemetry/instrumentation-fs': { enabled: false }, // ファイルシステムは不要
      }),
    ],
  })

  sdk.start()
    .then(() => {
      console.info('[telemetry] OpenTelemetry tracing enabled', endpoint)
    })
    .catch((error) => {
      console.error('[telemetry] OpenTelemetry の初期化に失敗しました', error)
    })

  // プロセス終了時にクリーンアップ
  process.on('SIGTERM', async () => {
    try {
      await sdk.shutdown()
      console.info('[telemetry] OpenTelemetry tracing stopped')
    }
    catch (error) {
      console.error('[telemetry] OpenTelemetry の終了処理に失敗しました', error)
    }
    finally {
      process.exit(0)
    }
  })
}
