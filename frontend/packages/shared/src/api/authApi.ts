import { httpClient } from './httpClient'
import type { LineAuthPayload, AuthResponse } from '../types/auth.types'

export const authApi = {
  // LINE ID トークンを送信して JWT を取得する
  lineAuth: (payload: LineAuthPayload) =>
    httpClient.post<AuthResponse>('/auth/line', payload).then((r) => r.data),

  logout: () =>
    httpClient.post<{ message: string }>('/auth/logout').then((r) => r.data),
}
