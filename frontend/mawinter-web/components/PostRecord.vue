<script setup lang="ts">
import { ref, computed } from 'vue'

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

// カテゴリをタイプ別にグループ化
const groupedCategories = computed(() => {
  if (!categories.value) return {}

  const groups = {
    income: [],
    outgoing: [],
    saving: [],
    investing: []
  }

  categories.value.forEach((cat) => {
    if (groups[cat.category_type]) {
      groups[cat.category_type].push({
        category_id: cat.category_id,
        category_name: cat.category_name
      })
    }
  })

  return groups
})

// カテゴリタイプの日本語名
const categoryTypeLabels = {
  income: '収入',
  outgoing: '支出',
  saving: '貯金',
  investing: '投資'
}

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

  try {
    // datetime を RFC3339 形式に変換
    const datetimeRFC3339 = new Date(formData.value.datetime).toISOString()

    await $fetch('/api/v3/record', {
      baseURL: useRuntimeConfig().public.mawinterApi,
      method: 'POST',
      body: {
        category_id: formData.value.category_id,
        datetime: datetimeRFC3339,
        from: 'mawinter-web', // システム固定値
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
  <div class="post-record">
    <form @submit.prevent="submitRecord">
      <div class="form-group">
        <label for="category">カテゴリ</label>
        <select
          id="category"
          v-model.number="formData.category_id"
          required
        >
          <optgroup
            v-for="(categoryList, type) in groupedCategories"
            :key="type"
            :label="categoryTypeLabels[type]"
          >
            <option
              v-for="cat in categoryList"
              :key="cat.category_id"
              :value="cat.category_id"
            >
              {{ cat.category_name }}
            </option>
          </optgroup>
        </select>
      </div>

      <div class="form-group">
        <label for="datetime">日時</label>
        <input
          id="datetime"
          v-model="formData.datetime"
          type="datetime-local"
          required
        >
      </div>

      <div class="form-group">
        <label for="price">金額</label>
        <input
          id="price"
          v-model.number="formData.price"
          type="number"
          required
          min="0"
        >
      </div>

      <div class="form-group">
        <label for="memo">メモ</label>
        <input
          id="memo"
          v-model="formData.memo"
          type="text"
          placeholder="メモ（任意）"
        >
      </div>

      <div class="form-actions">
        <button
          type="submit"
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? '送信中...' : '登録' }}
        </button>
      </div>

      <div
        v-if="errorMessage"
        class="message error"
      >
        {{ errorMessage }}
      </div>
      <div
        v-if="successMessage"
        class="message success"
      >
        {{ successMessage }}
      </div>
    </form>
  </div>
</template>

<style scoped>
.post-record {
  max-width: 600px;
  margin: 1rem auto;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #f9f9f9;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.25rem;
  font-weight: bold;
  color: #333;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #0066cc;
}

.form-actions {
  margin-top: 1.5rem;
}

.form-actions button {
  width: 100%;
  padding: 0.75rem;
  background-color: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.2s;
}

.form-actions button:hover:not(:disabled) {
  background-color: #218838;
}

.form-actions button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.message {
  margin-top: 1rem;
  padding: 0.75rem;
  border-radius: 4px;
}

.message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.message.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}
</style>
