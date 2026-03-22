import axios from 'axios'

const TOKEN_KEY = 'auth_token'

export const httpClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// リクエストインターセプター: localStorage の JWT を Authorization ヘッダに付与する
httpClient.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// レスポンスインターセプター: エラーメッセージを正規化する
httpClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response) {
      // 401 は認証切れ → localStorage をクリアしてログインページへ
      if (error.response.status === 401) {
        localStorage.removeItem(TOKEN_KEY)
        // TODO: ページリロードまたは router.navigate('/login') を呼ぶ
      }
      const message = (error.response.data as { error?: string })?.error ?? error.message
      return Promise.reject(new Error(message))
    }
    return Promise.reject(error)
  },
)
