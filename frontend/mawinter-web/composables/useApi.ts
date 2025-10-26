/**
 * API呼び出し用のcomposable
 *
 * runtimeConfigから環境変数を取得し、バックエンドAPIへのリクエストを簡単に行えるようにする
 *
 * 使用例:
 * ```typescript
 * const { getApiUrl, fetchApi } = useApi()
 *
 * // URLを取得
 * const url = getApiUrl('/v3/categories')
 * // => http://localhost:8080/api/v3/categories
 *
 * // データを取得
 * const { data, error } = await fetchApi('/v3/categories')
 * ```
 */
export const useApi = () => {
  const config = useRuntimeConfig()

  /**
   * APIの完全なURLを取得する
   * @param path - APIエンドポイントのパス (例: '/v3/categories', '/v3/record')
   * @returns 完全なURL (例: 'http://localhost:8080/api/v3/categories')
   */
  const getApiUrl = (path: string): string => {
    const baseUrl = config.public.mawinterApi
    const baseEndpoint = config.public.mawinterApiBaseEndpoint

    // パスの先頭のスラッシュを正規化
    const normalizedPath = path.startsWith('/') ? path : `/${path}`

    return `${baseUrl}${baseEndpoint}${normalizedPath}`
  }

  /**
   * APIにリクエストを送信する (useFetchのラッパー)
   * @param path - APIエンドポイントのパス
   * @param options - useFetchのオプション
   * @returns useFetchの戻り値
   */
  const fetchApi = <T = any>(path: string, options?: any) => {
    const url = getApiUrl(path)
    return useFetch<T>(url, options)
  }

  /**
   * APIにリクエストを送信する ($fetchのラッパー)
   * @param path - APIエンドポイントのパス
   * @param options - $fetchのオプション
   * @returns $fetchの戻り値
   */
  const callApi = async <T = any>(path: string, options?: any): Promise<T> => {
    const url = getApiUrl(path)
    return $fetch<T>(url, options)
  }

  return {
    getApiUrl,
    fetchApi,
    callApi
  }
}
