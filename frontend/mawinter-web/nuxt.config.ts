// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: ['@nuxt/eslint'],

  // 開発サーバー設定
  devServer: {
    port: 3000
  },

  // Nitroサーバー設定（プロキシ）
  nitro: {
    devProxy: {
      '/api': {
        target: 'http://localhost:8080/api',
        changeOrigin: true
      }
    }
  },

  // ランタイム設定 (環境変数)
  runtimeConfig: {
    public: {
      // バックエンドAPIのベースURL（プロキシ経由で空文字列）
      mawinterApi: '',

      // APIのベースエンドポイント (固定)
      mawinterApiBaseEndpoint: '/api'
    }
  }
})
