// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  modules: ['@nuxt/eslint'],

  // ランタイム設定 (環境変数)
  runtimeConfig: {
    // サーバーサイドのみで利用可能な設定
    // ここには機密情報を配置できる

    // クライアントサイドで利用可能な公開設定
    public: {
      // バックエンドAPIのベースURL
      // 開発環境: http://localhost:8080
      // 本番環境: 環境変数から取得
      mawinterApi: process.env.NUXT_PUBLIC_MAWINTER_API || 'http://localhost:8080',

      // APIのベースエンドポイント (固定)
      mawinterApiBaseEndpoint: '/api'
    }
  }
})
