export type UserRole = 'pending' | 'member' | 'admin'

export interface AuthUser {
  userId: string
  email: string
  role: UserRole
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface ResetPasswordPayload {
  token: string
  newPassword: string
}
