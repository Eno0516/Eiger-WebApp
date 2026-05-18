import axios from 'axios'
import { AppError } from '../errors/AppError'

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

// レスポンスインターセプター: エラーを AppError に正規化する
httpClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response) {
      const status = error.response.status
      const message = (error.response.data as { error?: string })?.error ?? error.message

      // 401 は認証切れ → localStorage をクリアしてログインページへ
      if (status === 401) {
        localStorage.removeItem(TOKEN_KEY)
        window.location.href = '/'
      }

      return Promise.reject(new AppError(message, { status }))
    }
    // ネットワーク断など response がない場合
    return Promise.reject(
      new AppError('ネットワークエラーが発生しました。通信状況を確認してください。', {
        status: 0,
        display: 'dialog',
      }),
    )
  },
)
