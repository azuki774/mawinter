<script setup lang="ts">
import { ref } from 'vue'

// イベント定義
const emit = defineEmits(['record-created'])

// フォームデータ
const formData = ref({
  category_id: 200,
  datetime: '',
  type: '',
  price: 0,
  memo: ''
})

// カテゴリ一覧
const { data: categories } = await useFetch('/api/v3/categories', {
  baseURL: useRuntimeConfig().public.mawinterApi,
  server: false // クライアントサイドでのみ取得
})

// 送信中フラグ
const isSubmitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

// 現在日時を YYYY-MM-DDTHH:mm 形式で取得
const getCurrentDateTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

// 初期値として現在日時を設定
formData.value.datetime = getCurrentDateTime()

// フォーム送信
const submitRecord = async () => {
  isSubmitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  // バリデーション: 数値チェック
  if (typeof formData.value.price !== 'number' || isNaN(formData.value.price)) {
    errorMessage.value = '金額は数値を入力してください'
    isSubmitting.value = false
    return
  }

  // バリデーション: 0は許可しない
  if (formData.value.price === 0) {
    errorMessage.value = '金額は0以外の値を入力してください'
    isSubmitting.value = false
    return
  }

  try {
    // datetime を RFC3339 形式に変換
    const datetimeRFC3339 = new Date(formData.value.datetime).toISOString()

    await $fetch('/api/v3/record', {
      baseURL: useRuntimeConfig().public.mawinterApi,
      method: 'POST',
      body: {
        category_id: formData.value.category_id,
        datetime: datetimeRFC3339,
        from: 'mawinter-web',
        type: formData.value.type,
        price: formData.value.price,
        memo: formData.value.memo
      }
    })

    successMessage.value = '記録を登録しました'

    // フォームをリセット
    formData.value.price = 0
    formData.value.memo = ''

    // 親コンポーネントにイベントを通知
    emit('record-created')

  } catch (error) {
    errorMessage.value = `エラーが発生しました: ${error.message || '不明なエラー'}`
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="max-w-xl mx-auto">
    <div class="bg-white rounded-lg shadow-sm border border-slate-200 p-5">
      <form class="space-y-4" @submit.prevent="submitRecord">
        <div>
          <label for="category" class="block text-sm font-medium text-slate-700 mb-1">
            カテゴリ
          </label>
          <select
            id="category"
            v-model.number="formData.category_id"
            required
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          >
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
          <label for="datetime" class="block text-sm font-medium text-slate-700 mb-1">
            日時
          </label>
          <input
            id="datetime"
            v-model="formData.datetime"
            type="datetime-local"
            required
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          >
        </div>

        <div>
          <label for="price" class="block text-sm font-medium text-slate-700 mb-1">
            金額
          </label>
          <input
            id="price"
            v-model.number="formData.price"
            type="number"
            required
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          >
        </div>

        <div>
          <label for="memo" class="block text-sm font-medium text-slate-700 mb-1">
            メモ
          </label>
          <input
            id="memo"
            v-model="formData.memo"
            type="text"
            placeholder="メモ（任意）"
            class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm placeholder:text-slate-400 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
          >
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full rounded-md bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors"
        >
          {{ isSubmitting ? '送信中...' : '登録' }}
        </button>

        <div
          v-if="errorMessage"
          class="rounded-md bg-red-50 border border-red-200 p-3 text-sm text-red-700"
        >
          {{ errorMessage }}
        </div>
        <div
          v-if="successMessage"
          class="rounded-md bg-emerald-50 border border-emerald-200 p-3 text-sm text-emerald-700"
        >
          {{ successMessage }}
        </div>
      </form>
    </div>
  </div>
</template>
