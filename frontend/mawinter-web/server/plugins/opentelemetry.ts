import process from "node:process";
import { lookup as dnsLookup, type LookupFunction } from "node:dns";

import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { NodeSDK } from "@opentelemetry/sdk-node";
import { defineNitroPlugin } from "#imports";

const SERVICE_NAME = "mawinter-web";
const STATE_KEY = "__mawinter_otlp_sdk__";

type OtelGlobal = typeof globalThis & {
  [STATE_KEY]?: NodeSDK;
};

const buildCollectorURL = (server: string) => {
  const normalized = server.trim().replace(/\/+$/, "");
  if (!normalized) {
    throw new Error("OTLP_SERVER is empty");
  }

  if (normalized.startsWith("http://") || normalized.startsWith("https://")) {
    return `${normalized}/v1/traces`;
  }
  return `http://${normalized}/v1/traces`;
};

export default defineNitroPlugin(async () => {
  const endpoint = process.env.OTLP_SERVER;
  if (!endpoint) {
    return;
  }

  const globalState = globalThis as OtelGlobal;
  if (globalState[STATE_KEY]) {
    return;
  }

  let collectorURL: string;
  try {
    collectorURL = buildCollectorURL(endpoint);
  } catch (error) {
    console.error("[otel] OTLP_SERVER が不正です", error);
    return;
  }

  if (!process.env.OTEL_SERVICE_NAME) {
    process.env.OTEL_SERVICE_NAME = SERVICE_NAME;
  }

  const forceIPv4Lookup: LookupFunction = (hostname, options, callback) => {
    if (typeof options === "function") {
      return dnsLookup(hostname, { family: 4, all: false }, options);
    }

    if (typeof options === "number") {
      return dnsLookup(hostname, { family: 4, hints: options, all: false }, callback);
    }

    return dnsLookup(hostname, { ...options, family: 4, all: false }, callback);
  };

  const sdk = new NodeSDK({
    traceExporter: new OTLPTraceExporter({
      url: collectorURL,
      httpAgentOptions: {
        lookup: forceIPv4Lookup,
      },
    }),
    instrumentations: [getNodeAutoInstrumentations()],
  });

  try {
    await sdk.start();
    globalState[STATE_KEY] = sdk;
  } catch (error) {
    console.error("[otel] トレース初期化に失敗しました", error);
    return;
  }

  const shutdown = async () => {
    try {
      await sdk.shutdown();
    } catch (error) {
      console.error("[otel] トレースの終了処理に失敗しました", error);
    }
  };

  process.once("SIGTERM", shutdown);
  process.once("SIGINT", shutdown);
});
