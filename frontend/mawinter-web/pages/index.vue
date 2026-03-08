<script setup lang="ts">
import { ref } from 'vue'

// 最新の年度を取得
// eslint-disable-next-line no-unused-vars
const { data: latestYear } = await useFetch('/api/v3/record/available', {
  baseURL: useRuntimeConfig().public.mawinterApi,
  server: false, // クライアントサイドでのみ取得
  transform: (data) => {
    // available API から最新の年度を取得
    // APIレスポンスは {fy: ["2024", "2023"], yyyymm: [...]} 形式
    if (data?.fy && data.fy.length > 0) {
      return data.fy[0] // 最新の年度（配列の最初の要素）
    }
    return null
  }
})

// SearchHistory を再マウントするためのキー
const searchHistoryKey = ref(0)

// 記録登録後のハンドラ
const handleRecordCreated = () => {
  // キーを変更してSearchHistoryを再マウント
  searchHistoryKey.value++
}
</script>

<template>
  <div>
    <section class="mb-8">
      <h2 class="text-lg font-semibold text-slate-800 mb-3">
        登録
      </h2>
      <PostRecord @record-created="handleRecordCreated" />
    </section>

    <section>
      <h2 class="text-lg font-semibold text-slate-800 mb-3">
        履歴検索
      </h2>
      <SearchHistory :key="searchHistoryKey" />
    </section>
  </div>
</template>
