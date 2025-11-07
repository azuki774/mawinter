<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

const runtimeConfig = useRuntimeConfig()
const months = ['4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月', '1月', '2月', '3月']

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

const summaryRows = computed(() => {
  if (!summaryEntries.value.length) {
    return []
  }

  const incomeRow = createRow('収入', 'income')
  const outgoingRow = createRow('支出', 'outgoing')
  const investingRow = createRow('投資', 'investing')

  const applyEntry = (target, entry) => {
    for (let i = 0; i < 12; i++) {
      const value = Math.abs(entry.price?.[i] ?? 0)
      target.monthly[i] += value
    }
    target.total += Math.abs(entry.total ?? 0)
  }

  for (const entry of summaryEntries.value) {
    if (entry.category_type === 'income') {
      applyEntry(incomeRow, entry)
    } else if (entry.category_type === 'outgoing') {
      applyEntry(outgoingRow, entry)
    } else if (entry.category_type === 'investing' || entry.category_type === 'saving') {
      applyEntry(investingRow, entry)
    }
  }

  const totalRow = createRow('累計', 'total')
  for (let i = 0; i < 12; i++) {
    totalRow.monthly[i] = incomeRow.monthly[i] - outgoingRow.monthly[i] - investingRow.monthly[i]
  }
  totalRow.total = incomeRow.total - outgoingRow.total - investingRow.total

  return [totalRow, incomeRow, outgoingRow, investingRow]
})

const hasSummaryData = computed(() => summaryRows.value.length > 0)
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
      <label
        class="control-label"
        for="year-select"
      >
        表示年度
      </label>
      <select
        id="year-select"
        v-model.number="selectedYear"
        :disabled="isDropdownDisabled"
      >
        <option
          v-for="year in availableYears"
          :key="year"
          :value="year"
        >
          {{ year }}年度
        </option>
      </select>
      <p
        v-if="isLoadingYears"
        class="state-text"
      >
        年度を取得中...
      </p>
      <p
        v-else-if="yearError"
        class="state-text state-text--error"
      >
        {{ yearError }}
      </p>
    </section>

    <section>
      <p
        v-if="isLoadingSummary"
        class="state-text"
      >
        サマリーを読み込み中です...
      </p>
      <p
        v-else-if="summaryError"
        class="state-text state-text--error"
      >
        {{ summaryError }}
      </p>
      <p
        v-else-if="!hasSummaryData"
        class="state-text"
      >
        表示できるサマリーデータがありません。
      </p>
      <div v-else class="table-wrapper">
        <table class="summary-table">
          <thead>
            <tr>
              <th scope="col">区分</th>
              <th
                v-for="month in months"
                :key="month"
                scope="col"
                class="numeric"
              >
                {{ month }}
              </th>
              <th scope="col" class="numeric">
                合計
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in summaryRows"
              :key="row.label"
              class="metric-row"
              :class="`metric-row--${row.theme}`"
            >
              <th scope="row">
                {{ row.label }}
              </th>
              <td
                v-for="(value, index) in row.monthly"
                :key="`${row.label}-${index}`"
                class="numeric"
              >
                {{ value.toLocaleString() }}
              </td>
              <td class="numeric">
                {{ row.total.toLocaleString() }}
              </td>
            </tr>
          </tbody>
        </table>
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
  border-collapse: collapse;
  background-color: #fff;
  border: 1px solid #e5e7eb;
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

.metric-row--total {
  background-color: #fff7cc;
}

.metric-row--income {
  background-color: #e5f5d9;
}

.metric-row--outgoing {
  background-color: #fde2e2;
}

.metric-row--investing {
  background-color: #f0e7ff;
}
</style>
