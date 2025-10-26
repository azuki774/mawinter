// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: ['@nuxt/eslint'],

  // 開発サーバー設定
  devServer: {
    port: 3000
  },

  // ランタイム設定 (環境変数)
  runtimeConfig: {
    // サーバーサイドのみで使用される環境変数（ブラウザからはアクセス不可）
    // バックエンドAPIサーバーの実際のURL
    mawinterApiUrl: process.env.MAWINTER_API_URL || 'http://localhost:8080',

    public: {
      // バックエンドAPIのベースURL（プロキシ経由で空文字列）
      mawinterApi: '',

      // APIのベースエンドポイント (固定)
      mawinterApiBaseEndpoint: '/api'
    }
  }
})
