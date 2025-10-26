// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  // カスタムルールをここに追加可能
  {
    rules: {
      // コンソールログの使用を警告 (開発中は許可、本番前に修正)
      'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    },
  }
)
