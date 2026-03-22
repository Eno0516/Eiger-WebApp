import { Outlet } from 'react-router-dom'

export function UserLayout() {
  return (
    <div style={{ maxWidth: '480px', margin: '0 auto', minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header style={{ padding: '1rem', borderBottom: '1px solid #e5e7eb', background: '#fff', position: 'sticky', top: 0 }}>
        <h1 style={{ margin: 0, fontSize: '1.125rem', fontWeight: 600 }}>Eiger</h1>
      </header>
      <main style={{ flex: 1, padding: '1rem' }}>
        <Outlet />
      </main>
    </div>
  )
}
