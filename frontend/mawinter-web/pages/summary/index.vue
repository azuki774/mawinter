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

  const endpoints = ['/api/v3/record/available']
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

const sectionLabels = {
  total: '累計',
  income: '収入',
  outgoing: '支出',
  investing: '投資'
}

const sectionColors = {
  total: 'border-l-amber-400',
  income: 'border-l-emerald-400',
  outgoing: 'border-l-red-400',
  investing: 'border-l-violet-400'
}

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
  <div>
    <h1 class="text-2xl font-bold text-slate-800 mb-4">
      サマリー
    </h1>

    <div class="flex items-center gap-3 mb-6">
      <label class="text-sm font-semibold text-slate-700" for="year-select">
        表示年度
      </label>
      <select
        id="year-select"
        v-model.number="selectedYear"
        :disabled="isDropdownDisabled"
        class="rounded-md border border-slate-300 px-3 py-1.5 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none disabled:bg-slate-100 disabled:cursor-not-allowed min-w-[140px]"
      >
        <option v-for="year in availableYears" :key="year" :value="year">
          {{ year }}年度
        </option>
      </select>
      <p v-if="isLoadingYears" class="text-sm text-slate-500">
        年度を取得中...
      </p>
      <p v-else-if="yearError" class="text-sm text-red-600">
        {{ yearError }}
      </p>
    </div>

    <div>
      <p v-if="isLoadingSummary" class="text-sm text-slate-500 py-8 text-center">
        サマリーを読み込み中です...
      </p>
      <p v-else-if="summaryError" class="text-sm text-red-600 py-8 text-center">
        {{ summaryError }}
      </p>
      <p v-else-if="!hasSummaryData" class="text-sm text-slate-500 py-8 text-center">
        表示できるサマリーデータがありません。
      </p>
      <div v-else class="space-y-6">
        <div
          v-for="section in ['total', 'income', 'outgoing', 'investing']"
          :key="section"
          class="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden"
        >
          <h2
            class="text-base font-semibold text-slate-800 px-4 py-3 bg-slate-50 border-b border-slate-200 border-l-4"
            :class="sectionColors[section]"
          >
            {{ sectionLabels[section] }}
          </h2>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="bg-slate-50">
                  <th scope="col" class="text-left px-3 py-2 font-semibold text-slate-600 whitespace-nowrap sticky left-0 bg-slate-50 min-w-[120px]">
                    区分
                  </th>
                  <th
                    v-for="month in months"
                    :key="month"
                    scope="col"
                    class="text-right px-3 py-2 font-semibold text-slate-600 whitespace-nowrap min-w-[70px]"
                  >
                    {{ month }}
                  </th>
                  <th scope="col" class="text-right px-3 py-2 font-semibold text-slate-600 whitespace-nowrap min-w-[100px]">
                    合計
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="row in groupedSummary[section]"
                  :key="section + row.label"
                  class="hover:bg-slate-50/50 transition-colors"
                  :class="{
                    'bg-amber-50 font-semibold': row.theme === 'total',
                    'bg-emerald-50 font-semibold': row.theme === 'income' && row.label === '収入合計',
                    'bg-red-50 font-semibold': row.theme === 'outgoing' && row.label === '支出合計',
                    'bg-violet-50 font-semibold': row.theme === 'investing' && row.label === '投資合計'
                  }"
                >
                  <th scope="row" class="text-left px-3 py-2 text-slate-700 whitespace-nowrap sticky left-0 bg-inherit">
                    {{ row.label }}
                  </th>
                  <td
                    v-for="(value, index) in row.monthly"
                    :key="`${section}-${row.label}-${index}`"
                    class="text-right px-3 py-2 text-slate-600 tabular-nums"
                  >
                    {{ value.toLocaleString() }}
                  </td>
                  <td class="text-right px-3 py-2 text-slate-800 font-semibold tabular-nums">
                    {{ row.total.toLocaleString() }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
