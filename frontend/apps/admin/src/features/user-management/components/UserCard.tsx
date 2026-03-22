import type { User } from '../types/user.types'

interface UserCardProps {
  user: User
}

export function UserCard({ user }: UserCardProps) {
  return (
    <div style={{ padding: '1rem', border: '1px solid #e5e7eb', borderRadius: '0.5rem', background: '#fff' }}>
      <p style={{ margin: '0 0 0.25rem', fontWeight: 600 }}>{user.name}</p>
      <p style={{ margin: 0, fontSize: '0.875rem', color: '#6b7280' }}>{user.email}</p>
    </div>
  )
}
