<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

const runtimeConfig = useRuntimeConfig()
const months = ['4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月', '1月', '2月', '3月']
const incomeCategoryOrder = [
  { id: 100, label: '月給' },
  { id: 101, label: 'ボーナス' },
  { id: 110, label: '雑所得' }
]
const outgoingCategoryOrder = [
  { id: 200, label: '家賃' },
  { id: 210, label: '食費' },
  { id: 220, label: '電気代' },
  { id: 221, label: 'ガス代' },
  { id: 222, label: '水道費' },
  { id: 230, label: 'コンピュータリソース' },
  { id: 231, label: '通信費' },
  { id: 240, label: '生活用品' },
  { id: 250, label: '娯楽費' },
  { id: 251, label: '交遊費' },
  { id: 260, label: '書籍・勉強' },
  { id: 270, label: '交通費' },
  { id: 280, label: '衣服等費' },
  { id: 300, label: '保険・税金' },
  { id: 400, label: '医療・衛生' },
  { id: 500, label: '雑費' }
]
const investingCategoryOrder = [
  // { id: 600, label: '家賃用貯金' },
  // { id: 601, label: 'PC用貯金' },
  { id: 700, label: 'NISA入出金' },
  { id: 701, label: 'NISA変動' }
]

const availableYears = ref([])
const selectedYear = ref(null)
const summaryEntries = ref([])
const isLoadingYears = ref(true)
const isLoadingSummary = ref(false)
const yearError = ref('')
const summaryError = ref('')

const normalizeYearList = (list) => {
  if (!Array.isArray(list)) {
    return []
  }

  const years = []
  for (const value of list) {
    const parsed = Number.parseInt(String(value), 10)
    if (!Number.isNaN(parsed)) {
      years.push(parsed)
    }
  }

  return years.sort((a, b) => b - a)
}

const fetchAvailableYears = async () => {
  isLoadingYears.value = true
  yearError.value = ''

  const endpoints = ['/api/v3/available', '/api/v3/record/available']
  let fetched = false

  for (const endpoint of endpoints) {
    try {
      const data = await $fetch(endpoint, {
        baseURL: runtimeConfig.public.mawinterApi
      })
      const years = normalizeYearList(data?.fy)
      if (years.length === 0) {
        continue
      }
      availableYears.value = years
      if (!selectedYear.value || !years.includes(selectedYear.value)) {
        selectedYear.value = years[0]
      }
      fetched = true
      break
    } catch (error) {
      console.error(`Failed to fetch available years from ${endpoint}`, error)
    }
  }

  if (!fetched) {
    availableYears.value = []
    selectedYear.value = null
    yearError.value = '年度情報の取得に失敗しました。'
  }

  isLoadingYears.value = false
}

const fetchSummaryData = async () => {
  if (!selectedYear.value) {
    summaryEntries.value = []
    return
  }

  isLoadingSummary.value = true
  summaryError.value = ''

  try {
    const data = await $fetch(`/api/v3/record/summary/${selectedYear.value}`, {
      baseURL: runtimeConfig.public.mawinterApi
    })
    summaryEntries.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('Failed to fetch summary data', error)
    summaryEntries.value = []
    summaryError.value = 'サマリーの取得に失敗しました。'
  } finally {
    isLoadingSummary.value = false
  }
}

const createRow = (label, theme) => ({
  label,
  monthly: Array(12).fill(0),
  total: 0,
  theme
})

const applyEntry = (target, entry) => {
  for (let i = 0; i < 12; i++) {
    const value = Math.abs(entry.price?.[i] ?? 0)
    target.monthly[i] += value
  }
  target.total += Math.abs(entry.total ?? 0)
}

const groupedSummary = computed(() => {
  if (!summaryEntries.value.length) {
    return {
      total: [],
      income: [],
      outgoing: [],
      investing: []
    }
  }

  const incomeAggregate = createRow('収入合計', 'income')
  const outgoingAggregate = createRow('支出合計', 'outgoing')
  const investingAggregate = createRow('投資合計', 'investing')
  const incomeCategoryMap = new Map()
  const outgoingCategoryMap = new Map()
  const investingCategoryMap = new Map()

  for (const entry of summaryEntries.value) {
    if (entry.category_type === 'income') {
      applyEntry(incomeAggregate, entry)

      if (!incomeCategoryMap.has(entry.category_id)) {
        incomeCategoryMap.set(
          entry.category_id,
          createRow(entry.category_name || `カテゴリ${entry.category_id}`, 'income')
        )
      }
      applyEntry(incomeCategoryMap.get(entry.category_id), entry)
    } else if (entry.category_type === 'outgoing') {
      applyEntry(outgoingAggregate, entry)
      if (!outgoingCategoryMap.has(entry.category_id)) {
        outgoingCategoryMap.set(
          entry.category_id,
          createRow(entry.category_name || `カテゴリ${entry.category_id}`, 'outgoing')
        )
      }
      applyEntry(outgoingCategoryMap.get(entry.category_id), entry)
    } else if (entry.category_type === 'investing' || entry.category_type === 'saving') {
      applyEntry(investingAggregate, entry)
      if (!investingCategoryMap.has(entry.category_id)) {
        investingCategoryMap.set(
          entry.category_id,
          createRow(entry.category_name || `カテゴリ${entry.category_id}`, 'investing')
        )
      }
      applyEntry(investingCategoryMap.get(entry.category_id), entry)
    }
  }

  const orderedIncomeRows = incomeCategoryOrder.map(({ id, label }) => {
    const row = incomeCategoryMap.get(id)
    if (row) {
      row.label = label
      return row
    }
    return createRow(label, 'income')
  })

  const buildOtherRows = (map, order) => Array.from(map.entries())
    .filter(([categoryId]) => !order.some((item) => item.id === categoryId))
    .map(([, row]) => row)

  const otherIncomeRows = buildOtherRows(incomeCategoryMap, incomeCategoryOrder)

  const orderedOutgoingRows = outgoingCategoryOrder.map(({ id, label }) => {
    const row = outgoingCategoryMap.get(id)
    if (row) {
      row.label = label
      return row
    }
    return createRow(label, 'outgoing')
  })
  const otherOutgoingRows = buildOtherRows(outgoingCategoryMap, outgoingCategoryOrder)

  const orderedInvestingRows = investingCategoryOrder.map(({ id, label }) => {
    const row = investingCategoryMap.get(id)
    if (row) {
      row.label = label
      return row
    }
    return createRow(label, 'investing')
  })
  const otherInvestingRows = buildOtherRows(investingCategoryMap, investingCategoryOrder)

  const buildTotalRow = (label, includeInvesting) => {
    const row = createRow(label, 'total')
    for (let i = 0; i < 12; i++) {
      row.monthly[i] = incomeAggregate.monthly[i] - outgoingAggregate.monthly[i]
      if (includeInvesting) {
        row.monthly[i] -= investingAggregate.monthly[i]
      }
    }
    row.total = incomeAggregate.total - outgoingAggregate.total
    if (includeInvesting) {
      row.total -= investingAggregate.total
    }
    return row
  }

  const totalWithInvesting = buildTotalRow('累計', true)
  const totalWithoutInvesting = buildTotalRow('累計（投資除く）', false)

  return {
    total: [totalWithInvesting, totalWithoutInvesting],
    income: [...orderedIncomeRows, ...otherIncomeRows, incomeAggregate],
    outgoing: [...orderedOutgoingRows, ...otherOutgoingRows, outgoingAggregate],
    investing: [...orderedInvestingRows, ...otherInvestingRows, investingAggregate]
  }
})

const hasSummaryData = computed(() => {
  const groups = groupedSummary.value
  return groups.total.length > 0 || groups.income.length > 0 || groups.outgoing.length > 0 || groups.investing.length > 0
})
const isDropdownDisabled = computed(() => isLoadingYears.value || availableYears.value.length === 0)

onMounted(async () => {
  await fetchAvailableYears()
  if (selectedYear.value) {
    await fetchSummaryData()
  }
})

watch(selectedYear, (newYear, oldYear) => {
  if (newYear && newYear !== oldYear) {
    fetchSummaryData()
  }
})
</script>

<template>
  <section class="summary-page">
    <header>
      <h1>サマリー</h1>
      <nav>
        <ul>
          <li>
            <NuxtLink to="/">トップ</NuxtLink>
          </li>
          <li>
            <NuxtLink to="/graph">グラフ</NuxtLink>
          </li>
        </ul>
      </nav>
    </header>

    <section class="controls">
      <label class="control-label" for="year-select">
        表示年度
      </label>
      <select id="year-select" v-model.number="selectedYear" :disabled="isDropdownDisabled">
        <option v-for="year in availableYears" :key="year" :value="year">
          {{ year }}年度
        </option>
      </select>
      <p v-if="isLoadingYears" class="state-text">
        年度を取得中...
      </p>
      <p v-else-if="yearError" class="state-text state-text--error">
        {{ yearError }}
      </p>
    </section>

    <section>
      <p v-if="isLoadingSummary" class="state-text">
        サマリーを読み込み中です...
      </p>
      <p v-else-if="summaryError" class="state-text state-text--error">
        {{ summaryError }}
      </p>
      <p v-else-if="!hasSummaryData" class="state-text">
        表示できるサマリーデータがありません。
      </p>
      <div v-else>
        <div v-for="section in ['total', 'income', 'outgoing', 'investing']" :key="section" class="table-wrapper">
          <h2 class="table-title">
            {{ section === 'total' ? '累計' : section === 'income' ? '収入' : section === 'outgoing' ? '支出' : '投資' }}
          </h2>
          <table class="summary-table">
            <thead>
              <tr>
                <th scope="col" class="label-col">
                  区分
                </th>
                <th v-for="month in months" :key="month" scope="col" class="numeric month-col">
                  {{ month }}
                </th>
                <th scope="col" class="numeric total-col">
                  合計
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in groupedSummary[section]" :key="section + row.label" class="metric-row"
                :class="`metric-row--${row.theme}`">
                <th scope="row" class="label-col">
                  {{ row.label }}
                </th>
                <td v-for="(value, index) in row.monthly" :key="`${section}-${row.label}-${index}`"
                  class="numeric month-col">
                  {{ value.toLocaleString() }}
                </td>
                <td class="numeric total-col">
                  {{ row.total.toLocaleString() }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  </section>
</template>

<style scoped>
.summary-page {
  padding: 1rem;
}

header {
  margin-bottom: 1rem;
}

h1 {
  margin: 0;
}

nav ul {
  list-style: none;
  padding: 0;
  margin: 0.5rem 0 0;
  display: flex;
  gap: 0.5rem;
}

nav a {
  display: inline-block;
  padding: 0.35rem 0.75rem;
  background-color: #0066cc;
  color: #fff;
  text-decoration: none;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: background-color 0.2s;
}

nav a:hover {
  background-color: #0052a3;
}

.controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.control-label {
  font-weight: 600;
}

select {
  padding: 0.4rem 0.6rem;
  border-radius: 4px;
  border: 1px solid #d1d5db;
  min-width: 140px;
}

.state-text {
  margin: 0;
  font-size: 0.9rem;
  color: #4b5563;
}

.state-text--error {
  color: #c62828;
}

.table-wrapper {
  overflow-x: auto;
}

.summary-table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
  background-color: #fff;
  border: 1px solid #e5e7eb;
}

.table-title {
  margin: 1.5rem 0 0.5rem;
  font-size: 1.25rem;
}

.summary-table th,
.summary-table td {
  border: 1px solid #e5e7eb;
  padding: 0.5rem;
}

.summary-table th {
  text-align: left;
  background-color: #f9fafb;
}

.summary-table th.numeric,
.summary-table td.numeric {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.label-col {
  width: 140px;
  white-space: nowrap;
}

.month-col {
  width: 70px;
}

.total-col {
  width: 110px;
}

.metric-row--total {
  background-color: #fff7cc;
}

.metric-row--total:nth-child(2) {
  background-color: #ffedb5;
}

.metric-row--income:not(:last-child),
.metric-row--outgoing,
.metric-row--investing {
  background-color: #fff;
}

.metric-row--income:last-child {
  background-color: #e5f5d9;
}

.metric-row--outgoing:last-child {
  background-color: #fde2e2;
}

.metric-row--investing:last-child {
  background-color: #f0e7ff;
}
</style>
