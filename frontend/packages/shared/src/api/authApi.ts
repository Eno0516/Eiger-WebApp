import { httpClient } from './httpClient'
import type { LoginPayload, RegisterPayload, ResetPasswordPayload } from '../types/auth.types'

export const authApi = {
  register: (payload: RegisterPayload) =>
    httpClient.post<{ message: string }>('/auth/register', payload).then((r) => r.data),

  login: (payload: LoginPayload) =>
    httpClient.post<{ token: string }>('/auth/login', payload).then((r) => r.data),

  verifyEmail: (token: string) =>
    httpClient.get<{ message: string }>(`/auth/verify-email?token=${token}`).then((r) => r.data),

  resendVerification: (email: string) =>
    httpClient.post<{ message: string }>('/auth/resend-verification', { email }).then((r) => r.data),

  forgotPassword: (email: string) =>
    httpClient.post<{ message: string }>('/auth/forgot-password', { email }).then((r) => r.data),

  resetPassword: (payload: ResetPasswordPayload) =>
    httpClient.post<{ message: string }>('/auth/reset-password', payload).then((r) => r.data),
}
