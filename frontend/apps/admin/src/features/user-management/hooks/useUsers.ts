import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { usersApi } from '../api/usersApi'
import type { CreateUserPayload } from '../types/user.types'

const QUERY_KEY = 'users'

export function useUsers() {
  return useQuery({
    queryKey: [QUERY_KEY],
    queryFn: usersApi.getAll,
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateUserPayload) => usersApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEY] })
    },
  })
}
