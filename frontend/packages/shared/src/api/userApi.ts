import { httpClient } from './httpClient'

export interface MeResponse {
  id: string
  is_registered: boolean
  is_system_admin: boolean
}

export const userApi = {
  // 自分のユーザー情報を取得する（is_registered の確認に使用）
  getMe: () =>
    httpClient.get<MeResponse>('/users/me').then((r) => r.data),

  // display_name を更新する（初回登録・変更時に使用）
  updateDisplayName: (userId: string, displayName: string) =>
    httpClient
      .patch<{ id: string; display_name: string }>(`/users/${userId}/display-name`, { display_name: displayName })
      .then((r) => r.data),
}
