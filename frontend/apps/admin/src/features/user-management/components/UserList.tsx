import { useUsers } from '../hooks/useUsers'
import { UserCard } from './UserCard'

export function UserList() {
  const { data: users, isLoading, isError } = useUsers()

  if (isLoading) return <p>Loading...</p>
  if (isError) return <p style={{ color: '#ef4444' }}>Failed to load users.</p>
  if (!users?.length) return <p style={{ color: '#6b7280' }}>ユーザーがいません。</p>

  return (
    <div style={{ display: 'grid', gap: '0.75rem' }}>
      {users.map((user) => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  )
}
