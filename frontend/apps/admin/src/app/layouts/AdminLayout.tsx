import { Outlet, NavLink } from 'react-router-dom'

export function AdminLayout() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <nav style={{ width: '240px', background: '#1e293b', color: '#fff', padding: '1.5rem 1rem', flexShrink: 0 }}>
        <h1 style={{ margin: '0 0 2rem', fontSize: '1.125rem', fontWeight: 700 }}>Eiger Admin</h1>
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
          <li>
            <NavLink
              to="/users"
              style={({ isActive }) => ({
                display: 'block',
                padding: '0.5rem 0.75rem',
                borderRadius: '0.375rem',
                color: isActive ? '#fff' : '#94a3b8',
                background: isActive ? '#334155' : 'transparent',
                textDecoration: 'none',
                fontSize: '0.875rem',
              })}
            >
              ユーザー管理
            </NavLink>
          </li>
        </ul>
      </nav>
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <header style={{ padding: '1rem 2rem', borderBottom: '1px solid #e5e7eb', background: '#fff' }}>
          <p style={{ margin: 0, fontSize: '0.875rem', color: '#6b7280' }}>管理画面</p>
        </header>
        <main style={{ flex: 1, padding: '2rem' }}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}
