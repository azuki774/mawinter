# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# pnpmのインストール
RUN corepack enable && corepack prepare pnpm@latest --activate

# 依存関係のインストール（キャッシュ最適化）
COPY mawinter-web/package.json mawinter-web/pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

# アプリケーションコードのコピー
COPY mawinter-web/ ./

# Nuxtアプリケーションのビルド
# 環境変数はビルド時に埋め込まれる（ARGとして受け取り可能）
ARG NUXT_PUBLIC_MAWINTER_API=""
ARG NUXT_PUBLIC_MAWINTER_API_BASE_ENDPOINT="/api"
ENV NUXT_PUBLIC_MAWINTER_API=${NUXT_PUBLIC_MAWINTER_API}
ENV NUXT_PUBLIC_MAWINTER_API_BASE_ENDPOINT=${NUXT_PUBLIC_MAWINTER_API_BASE_ENDPOINT}

RUN pnpm run build

# Production stage
FROM node:22-alpine AS runner

WORKDIR /app

# 非rootユーザーの作成
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nuxt -u 1001

# pnpmのインストール
RUN corepack enable && corepack prepare pnpm@latest --activate

# 必要なファイルのコピー
COPY --from=builder --chown=nuxt:nodejs /app/.output /app/.output
COPY --from=builder --chown=nuxt:nodejs /app/package.json /app/package.json
COPY --from=builder --chown=nuxt:nodejs /app/pnpm-lock.yaml /app/pnpm-lock.yaml
COPY --from=builder --chown=nuxt:nodejs /app/tracing.mjs /app/tracing.mjs

# 本番用の依存関係をインストール（tracing.mjs が必要とする OpenTelemetry パッケージ）
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --prod --frozen-lockfile

# 環境変数の設定（実行時にオーバーライド可能）
ENV NUXT_PUBLIC_MAWINTER_API=""
ENV NUXT_PUBLIC_MAWINTER_API_BASE_ENDPOINT="/api"
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

USER nuxt

# Nuxt 3の本番サーバーを起動（OpenTelemetry トレーシング有効）
# OTLP_SERVER 環境変数が設定されている場合のみトレース送信が有効化されます
CMD ["node", "--require", "./tracing.mjs", "/app/.output/server/index.mjs"]
