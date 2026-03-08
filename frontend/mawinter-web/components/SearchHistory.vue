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
    totalCount.value = { num: 0 }
  }
}

// ページネーション計算
const totalPages = computed(() => {
  if (!totalCount.value) return 0
  return Math.ceil(totalCount.value.num / searchParams.value.num)
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

// レコード削除
const deleteRecord = async (id) => {
  if (!confirm(`ID: ${id} のレコードを削除しますか？`)) {
    return
  }

  pending.value = true
  try {
    const config = useRuntimeConfig()
    const url = `${config.public.mawinterApi}/api/v3/record/${id}`

    const response = await fetch(url, {
      method: 'DELETE'
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    await Promise.all([fetchRecords(), fetchCount()])
  } catch (error) {
    console.error('deleteRecord エラー:', error)
    alert('削除に失敗しました')
  } finally {
    pending.value = false
  }
}

// 初回データ取得
onMounted(() => {
  // ページロード時にデータを取得
  search()
})
</script>

<template>
  <div>
    <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-4 mb-5">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
        <div>
          <label for="yyyymm" class="block text-sm font-medium text-slate-700 mb-1">
            年月
          </label>
          <select
            id="yyyymm"
            v-model="searchParams.yyyymm"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
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

        <div>
          <label for="category" class="block text-sm font-medium text-slate-700 mb-1">
            カテゴリ
          </label>
          <select
            id="category"
            v-model="searchParams.category_id"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
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

        <div>
          <label for="num" class="block text-sm font-medium text-slate-700 mb-1">
            表示件数
          </label>
          <select
            id="num"
            v-model.number="searchParams.num"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
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
        :disabled="pending"
        class="w-full rounded-md bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors"
        @click="search"
      >
        {{ pending ? '検索中...' : '検索' }}
      </button>
    </div>

    <div
      v-if="pending"
      class="text-sm text-slate-500 py-8 text-center"
    >
      読み込み中...
    </div>

    <div v-else-if="records.length > 0">
      <p class="text-sm text-slate-500 mb-3">
        全 {{ totalCount?.num || 0 }} 件中 {{ searchParams.offset + 1 }} - {{ Math.min(searchParams.offset + searchParams.num, totalCount?.num || 0) }} 件を表示
      </p>

      <div class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200">
                <th class="text-left px-4 py-3 font-semibold text-slate-600">
                  ID
                </th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600 whitespace-nowrap">
                  日時
                </th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600">
                  カテゴリ
                </th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600">
                  送信元
                </th>
                <th class="text-right px-4 py-3 font-semibold text-slate-600">
                  金額
                </th>
                <th class="text-left px-4 py-3 font-semibold text-slate-600">
                  メモ
                </th>
                <th class="text-center px-4 py-3 font-semibold text-slate-600">
                  削除
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="record in records"
                :key="record.id"
                class="hover:bg-slate-50/50 transition-colors"
              >
                <td class="px-4 py-3 text-slate-600">
                  {{ record.id }}
                </td>
                <td class="px-4 py-3 text-slate-600 whitespace-nowrap">
                  {{ formatDateTime(record.datetime) }}
                </td>
                <td class="px-4 py-3 text-slate-700">
                  {{ getCategoryName(record.category_id) }}
                </td>
                <td class="px-4 py-3 text-slate-500">
                  {{ record.from || '-' }}
                </td>
                <td class="px-4 py-3 text-right font-semibold text-slate-800 tabular-nums">
                  {{ formatPrice(record.price) }}
                </td>
                <td class="px-4 py-3 text-slate-600">
                  {{ record.memo || '-' }}
                </td>
                <td class="px-4 py-3 text-center">
                  <button
                    :disabled="pending"
                    class="rounded-md bg-red-500 px-3 py-1 text-xs font-medium text-white hover:bg-red-600 disabled:bg-slate-300 disabled:cursor-not-allowed transition-colors"
                    @click="deleteRecord(record.id)"
                  >
                    削除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div
        v-if="totalPages > 1"
        class="flex items-center justify-center gap-4 mt-4"
      >
        <button
          :disabled="currentPage === 1"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition-colors"
          @click="changePage(currentPage - 1)"
        >
          前へ
        </button>
        <span class="text-sm font-semibold text-slate-700">{{ currentPage }} / {{ totalPages }}</span>
        <button
          :disabled="currentPage === totalPages"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:bg-slate-300 disabled:cursor-not-allowed transition-colors"
          @click="changePage(currentPage + 1)"
        >
          次へ
        </button>
      </div>
    </div>

    <div
      v-else
      class="text-sm text-slate-500 py-8 text-center"
    >
      記録が見つかりませんでした
    </div>
  </div>
</template>
