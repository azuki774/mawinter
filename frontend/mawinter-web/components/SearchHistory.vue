<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// 検索パラメータ
const searchParams = ref({
  yyyymm: '',
  category_id: null,
  num: 20,
  offset: 0
})

// カテゴリ一覧
const { data: categories } = await useFetch('/api/v3/categories', {
  baseURL: useRuntimeConfig().public.mawinterApi,
  server: false // クライアントサイドでのみ取得
})

// 利用可能な年月一覧
const { data: availableData } = await useFetch('/api/v3/record/available', {
  baseURL: useRuntimeConfig().public.mawinterApi,
  server: false // クライアントサイドでのみ取得
})

// 年月の選択肢（降順でソート）
const availableYYYYMM = computed(() => {
  if (!availableData.value?.yyyymm) return []
  // 降順にソート（新しい月が上）
  return [...availableData.value.yyyymm].sort((a, b) => b.localeCompare(a))
})

// 記録一覧と総件数
const records = ref([])
const totalCount = ref(null)
const pending = ref(false)

// データ取得関数
const fetchRecords = async () => {
  try {
    const params = new URLSearchParams()
    if (searchParams.value.yyyymm) {
      params.append('yyyymm', searchParams.value.yyyymm)
    }
    if (searchParams.value.category_id) {
      params.append('category_id', searchParams.value.category_id.toString())
    }
    params.append('num', searchParams.value.num.toString())
    params.append('offset', searchParams.value.offset.toString())

    const config = useRuntimeConfig()
    const url = `${config.public.mawinterApi}/api/v3/record?${params.toString()}`

    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    records.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('fetchRecords エラー:', error)
    records.value = []
  }
}

const fetchCount = async () => {
  try {
    const params = new URLSearchParams()
    if (searchParams.value.yyyymm) {
      params.append('yyyymm', searchParams.value.yyyymm)
    }
    if (searchParams.value.category_id) {
      params.append('category_id', searchParams.value.category_id.toString())
    }

    const config = useRuntimeConfig()
    const url = `${config.public.mawinterApi}/api/v3/record/count?${params.toString()}`

    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    totalCount.value = data
  } catch (error) {
    console.error('fetchCount エラー:', error)
    totalCount.value = { count: 0 }
  }
}

// ページネーション計算
const totalPages = computed(() => {
  if (!totalCount.value) return 0
  return Math.ceil(totalCount.value.count / searchParams.value.num)
})

const currentPage = computed(() => {
  return Math.floor(searchParams.value.offset / searchParams.value.num) + 1
})

// 検索実行
const search = async () => {
  pending.value = true
  try {
    searchParams.value.offset = 0
    await Promise.all([fetchRecords(), fetchCount()])
  } catch (error) {
    console.error('検索エラー:', error)
  } finally {
    pending.value = false
  }
}

// ページ変更
const changePage = async (page) => {
  pending.value = true
  try {
    searchParams.value.offset = (page - 1) * searchParams.value.num
    await fetchRecords()
  } catch (error) {
    console.error('ページ変更エラー:', error)
  } finally {
    pending.value = false
  }
}

// 日時フォーマット
const formatDateTime = (datetime) => {
  return new Date(datetime).toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// カテゴリ名取得
const getCategoryName = (categoryId) => {
  const category = categories.value?.find(c => c.category_id === categoryId)
  return category?.category_name || '不明'
}

// 金額フォーマット
const formatPrice = (price) => {
  return price.toLocaleString('ja-JP') + '円'
}

// 初回データ取得
onMounted(() => {
  // ページロード時にデータを取得
  search()
})
</script>

<template>
  <div class="search-history">
    <div class="search-form">
      <div class="form-row">
        <div class="form-group">
          <label for="yyyymm">年月</label>
          <select
            id="yyyymm"
            v-model="searchParams.yyyymm"
          >
            <option value="">
              全て
            </option>
            <option
              v-for="ym in availableYYYYMM"
              :key="ym"
              :value="ym"
            >
              {{ ym.substring(0, 4) }}年{{ ym.substring(4, 6) }}月
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="category">カテゴリ</label>
          <select
            id="category"
            v-model="searchParams.category_id"
          >
            <option :value="null">
              全て
            </option>
            <option
              v-for="cat in categories"
              :key="cat.category_id"
              :value="cat.category_id"
            >
              {{ cat.category_name }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <label for="num">表示件数</label>
          <select
            id="num"
            v-model.number="searchParams.num"
          >
            <option :value="10">
              10件
            </option>
            <option :value="20">
              20件
            </option>
            <option :value="50">
              50件
            </option>
            <option :value="100">
              100件
            </option>
          </select>
        </div>
      </div>

      <button
        class="search-button"
        :disabled="pending"
        @click="search"
      >
        {{ pending ? '検索中...' : '検索' }}
      </button>
    </div>

    <div
      v-if="pending"
      class="loading"
    >
      読み込み中...
    </div>

    <div v-else-if="records.length > 0" class="results">
      <p class="result-info">
        全 {{ totalCount?.count || 0 }} 件中 {{ searchParams.offset + 1 }} - {{ Math.min(searchParams.offset + searchParams.num, totalCount?.count || 0) }} 件を表示
      </p>

      <table class="records-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>日時</th>
            <th>カテゴリ</th>
            <th>送信元</th>
            <th>金額</th>
            <th>メモ</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="record in records"
            :key="record.id"
          >
            <td>{{ record.id }}</td>
            <td>{{ formatDateTime(record.datetime) }}</td>
            <td>{{ getCategoryName(record.category_id) }}</td>
            <td>{{ record.from || '-' }}</td>
            <td class="price">
              {{ formatPrice(record.price) }}
            </td>
            <td>{{ record.memo || '-' }}</td>
          </tr>
        </tbody>
      </table>

      <div
        v-if="totalPages > 1"
        class="pagination"
      >
        <button
          :disabled="currentPage === 1"
          @click="changePage(currentPage - 1)"
        >
          前へ
        </button>
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        <button
          :disabled="currentPage === totalPages"
          @click="changePage(currentPage + 1)"
        >
          次へ
        </button>
      </div>
    </div>

    <div
      v-else
      class="no-results"
    >
      記録が見つかりませんでした
    </div>
  </div>
</template>

<style scoped>
.search-history {
  margin: 2rem 0;
}

.search-form {
  padding: 1rem;
  background-color: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: bold;
  color: #333;
  font-size: 0.9rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}

.search-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-button:hover:not(:disabled) {
  background-color: #0052a3;
}

.search-button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.loading,
.no-results {
  text-align: center;
  padding: 2rem;
  color: #666;
}

.result-info {
  margin-bottom: 1rem;
  color: #666;
  font-size: 0.9rem;
}

.records-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.records-table th,
.records-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.records-table th {
  background-color: #f8f9fa;
  font-weight: bold;
  color: #333;
}

.records-table tr:hover {
  background-color: #f8f9fa;
}

.records-table td.price {
  text-align: right;
  font-weight: bold;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.pagination button {
  padding: 0.5rem 1rem;
  background-color: #0066cc;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.pagination button:hover:not(:disabled) {
  background-color: #0052a3;
}

.pagination button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.page-info {
  font-weight: bold;
  color: #333;
}

@media (max-width: 768px) {
  .records-table {
    font-size: 0.85rem;
  }

  .records-table th,
  .records-table td {
    padding: 0.5rem;
  }
}
</style>
