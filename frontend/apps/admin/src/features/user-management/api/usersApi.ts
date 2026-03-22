import { httpClient } from '@eiger/shared/api'
import type { User, CreateUserPayload } from '../types/user.types'

export const usersApi = {
  getAll: () => httpClient.get<User[]>('/api/v1/users').then((r) => r.data),

  getById: (id: string) => httpClient.get<User>(`/api/v1/users/${id}`).then((r) => r.data),

  create: (payload: CreateUserPayload) =>
    httpClient.post<User>('/api/v1/users', payload).then((r) => r.data),
}
