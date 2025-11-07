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
    borderColor: 'rgb(75, 192, 192)',
    backgroundColor: 'rgba(75, 192, 192, 0.2)',
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
      display: true,
      text: '月次収支推移'
    },
    legend: {
      position: 'bottom'
    }
  },
  scales: {
    x: {
      stacked: true
    },
    y: {
      stacked: true,
      ticks: {
        callback: (value) => `¥${value.toLocaleString()}`
      }
    }
  }
}

const categoryChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: 'カテゴリ別年間支出'
    },
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      ticks: {
        callback: (value) => `¥${value.toLocaleString()}`
      }
    }
  }
}
</script>

<template>
  <section>
    <header>
      <h1>グラフ</h1>
      <nav>
        <ul>
          <li>
            <NuxtLink to="/">トップ</NuxtLink>
          </li>
        </ul>
      </nav>
    </header>

    <!-- 年度選択 -->
    <section>
      <div class="year-selector">
        <label for="year-select">年度: </label>
        <select
          id="year-select"
          v-model="selectedYear"
          :disabled="isLoadingYears"
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
    </section>

    <!-- ローディング表示 -->
    <section v-if="isLoadingSummary">
      <p>読み込み中...</p>
    </section>

    <!-- グラフ表示 -->
    <section v-else>
      <!-- 月次支出推移グラフ -->
      <div class="chart-container">
        <h2>月次収支推移</h2>
        <div class="chart-wrapper">
          <Bar
            :data="monthlyChartData"
            :options="monthlyChartOptions"
          />
        </div>
      </div>

      <!-- カテゴリ別支出グラフ -->
      <div class="chart-container">
        <h2>カテゴリ別年間支出</h2>
        <div class="chart-wrapper">
          <Bar
            :data="categoryChartData"
            :options="categoryChartOptions"
          />
        </div>
      </div>
    </section>
  </section>
</template>

<style scoped>
section {
  padding: 1rem;
}

h1 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

h2 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.year-selector {
  margin: 1rem 0;
}

.year-selector label {
  margin-right: 0.5rem;
  font-weight: bold;
}

.year-selector select {
  padding: 0.5rem;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.chart-container {
  margin: 2rem 0;
  padding: 1rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #fff;
}

.chart-wrapper {
  height: 400px;
  position: relative;
}

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
</style>
