<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  LineElement,
  PointElement,
  LineController,
  BarController
} from 'chart.js'

// Chart.jsのコンポーネントを登録
ChartJS.register(
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  LineElement,
  PointElement,
  LineController,
  BarController
)

// 現在の年度を計算 (4月開始の会計年度)
const getCurrentFiscalYear = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = now.getMonth() + 1 // 0-indexed なので +1
  return month >= 4 ? year : year - 1
}

// 利用可能な年度を取得
const availableYears = ref([])
const isLoadingYears = ref(true)
const selectedYear = ref(getCurrentFiscalYear())

const { data: availableData } = await useFetch('/api/v3/record/available', {
  baseURL: useRuntimeConfig().public.mawinterApi,
  server: false
})

// availableDataの変更を監視
watch(availableData, (newData) => {
  if (newData?.fy) {
    // 文字列の配列を数値の配列に変換
    availableYears.value = newData.fy.map(y => parseInt(y))
    isLoadingYears.value = false

    // 選択された年度がリストにない場合は最初の年度を選択
    if (availableYears.value.length > 0 && !availableYears.value.includes(selectedYear.value)) {
      selectedYear.value = availableYears.value[0]
    }
  }
}, { immediate: true })

// サマリーデータ
const summaryData = ref([])
const isLoadingSummary = ref(false)

// サマリーデータを取得
const fetchSummaryData = async () => {
  if (!selectedYear.value) return

  isLoadingSummary.value = true
  try {
    const response = await $fetch(`/api/v3/record/summary/${selectedYear.value}`, {
      baseURL: useRuntimeConfig().public.mawinterApi
    })
    summaryData.value = response || []
  } catch (error) {
    console.error('Failed to fetch summary data:', error)
    summaryData.value = []
  } finally {
    isLoadingSummary.value = false
  }
}

// 年度が変更されたらデータを再取得
watch(selectedYear, () => {
  fetchSummaryData()
})

// 初期データ取得
fetchSummaryData()

// 月次支出推移グラフのデータを生成
const monthlyChartData = computed(() => {
  if (!summaryData.value || summaryData.value.length === 0) {
    return {
      labels: [],
      datasets: []
    }
  }

  // 4月から3月までの月ラベル
  const months = ['4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月', '1月', '2月', '3月']

  // カテゴリ別に集計
  const categoryMap = {}
  const incomeData = Array(12).fill(0)

  summaryData.value.forEach(item => {
    const categoryId = item.category_id
    // price配列は12ヶ月分（インデックス0=4月, 11=3月）
    const priceArray = item.price || []

    if (item.category_type === 'income') {
      // 収入カテゴリ
      priceArray.forEach((price, index) => {
        incomeData[index] += Math.abs(price)
      })
    } else if (item.category_type === 'outgoing') {
      // 支出カテゴリ
      if (!categoryMap[categoryId]) {
        categoryMap[categoryId] = {
          name: item.category_name,
          data: priceArray.map(p => Math.abs(p))
        }
      }
    }
  })

  // カテゴリ別のデータセットを作成
  const datasets = Object.entries(categoryMap).map(([categoryId, data]) => ({
    type: 'bar',
    label: data.name,
    data: data.data,
    backgroundColor: getCategoryColor(parseInt(categoryId))
  }))

  // 収入の折れ線グラフを追加
  datasets.push({
    type: 'line',
    label: '収入合計',
    data: incomeData,
    borderColor: 'rgb(16, 185, 129)',
    backgroundColor: 'rgba(16, 185, 129, 0.1)',
    borderWidth: 2,
    tension: 0.1
  })

  return {
    labels: months,
    datasets
  }
})

// カテゴリ別年間支出グラフのデータを生成
const categoryChartData = computed(() => {
  if (!summaryData.value || summaryData.value.length === 0) {
    return {
      labels: [],
      datasets: []
    }
  }

  // カテゴリ別に年間合計を集計
  const categoryTotals = {}

  summaryData.value.forEach(item => {
    const categoryId = item.category_id
    const price = Math.abs(item.total || 0)

    // 支出カテゴリのみ
    if (item.category_type === 'outgoing') {
      categoryTotals[categoryId] = {
        name: item.category_name,
        total: price
      }
    }
  })

  // カテゴリ名とデータを配列に変換
  const sortedCategories = Object.entries(categoryTotals)
    .sort(([, a], [, b]) => b.total - a.total) // 金額でソート

  const labels = sortedCategories.map(([, data]) => data.name)
  const data = sortedCategories.map(([, data]) => data.total)
  const colors = sortedCategories.map(([categoryId]) => getCategoryColor(parseInt(categoryId)))

  return {
    labels,
    datasets: [{
      label: '年間支出',
      data,
      backgroundColor: colors
    }]
  }
})

// カテゴリIDに基づいて色を生成
function getCategoryColor(categoryId) {
  // カテゴリIDに基づいて色相を決定
  const hue = (categoryId * 137) % 360
  return `hsl(${hue}, 70%, 60%)`
}

// グラフのオプション
const monthlyChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: false
    },
    legend: {
      position: 'bottom',
      labels: {
        boxWidth: 12,
        padding: 12,
        font: { size: 11 }
      }
    }
  },
  scales: {
    x: {
      stacked: true,
      grid: { display: false }
    },
    y: {
      stacked: true,
      ticks: {
        callback: (value) => `${value.toLocaleString()}`
      }
    }
  }
}

const categoryChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: false
    },
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      ticks: {
        callback: (value) => `${value.toLocaleString()}`
      }
    },
    x: {
      grid: { display: false }
    }
  }
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-slate-800 mb-4">
      グラフ
    </h1>

    <div class="flex items-center gap-3 mb-6">
      <label class="text-sm font-semibold text-slate-700" for="year-select">
        年度
      </label>
      <select
        id="year-select"
        v-model="selectedYear"
        :disabled="isLoadingYears"
        class="rounded-md border border-slate-300 px-3 py-1.5 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none disabled:bg-slate-100 disabled:cursor-not-allowed"
      >
        <option
          v-for="year in availableYears"
          :key="year"
          :value="year"
        >
          {{ year }}年度
        </option>
      </select>
    </div>

    <div v-if="isLoadingSummary" class="text-sm text-slate-500 py-12 text-center">
      読み込み中...
    </div>

    <div v-else class="space-y-6">
      <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-4 sm:p-6">
        <h2 class="text-base font-semibold text-slate-800 mb-4">
          月次収支推移
        </h2>
        <div class="h-[400px] relative">
          <Bar
            :data="monthlyChartData"
            :options="monthlyChartOptions"
          />
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-4 sm:p-6">
        <h2 class="text-base font-semibold text-slate-800 mb-4">
          カテゴリ別年間支出
        </h2>
        <div class="h-[400px] relative">
          <Bar
            :data="categoryChartData"
            :options="categoryChartOptions"
          />
        </div>
      </div>
    </div>
  </div>
</template>
