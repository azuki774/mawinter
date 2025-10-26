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
  <section>
    <header>
      <h1>mawinter-web</h1>
      <nav>
        <ul>
          <li>
            <NuxtLink to="/graph">グラフ</NuxtLink>
          </li>
        </ul>
      </nav>
    </header>

    <section>
      <h2>登録</h2>
      <PostRecord @record-created="handleRecordCreated" />
    </section>

    <section>
      <h2>履歴検索</h2>
      <div class="container-sm">
        <SearchHistory :key="searchHistoryKey" />
      </div>
    </section>
  </section>
</template>

<style scoped>
nav ul {
  list-style: none;
  padding: 0;
  margin: 1rem 0;
}

nav li {
  margin: 0.5rem 0;
}

nav a {
  display: inline-block;
  padding: 0.5rem 1rem;
  background-color: #0066cc;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.2s;
}

nav a:hover {
  background-color: #0052a3;
}

.container-sm {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem;
}
</style>
