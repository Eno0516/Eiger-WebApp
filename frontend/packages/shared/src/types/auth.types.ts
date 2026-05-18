export interface AuthUser {
  userId: string
  isSystemAdmin: boolean
}

export interface LineAuthPayload {
  idToken: string // LINE が発行した ID トークン
}

export interface AuthResponse {
  token: string       // アプリが発行する JWT
  is_registered: boolean // 初期登録（LIFFフォーム）完了済みか
}
