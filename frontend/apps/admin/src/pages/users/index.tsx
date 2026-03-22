import { CreateUserForm, UserList } from '@features/user-management'

export function UsersPage() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
      <section>
        <h2 style={{ marginTop: 0 }}>ユーザーを追加</h2>
        <CreateUserForm />
      </section>
      <section>
        <h2 style={{ marginTop: 0 }}>ユーザー一覧</h2>
        <UserList />
      </section>
    </div>
  )
}
