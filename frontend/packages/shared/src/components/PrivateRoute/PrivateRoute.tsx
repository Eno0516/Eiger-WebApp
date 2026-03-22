import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import type { UserRole } from '../../types/auth.types'

const ROLE_LEVEL: Record<UserRole, number> = {
  pending: 1,
  member: 2,
  admin: 3,
}

interface PrivateRouteProps {
  /** 最低限必要なロール。未指定の場合はログイン済みであれば通過。 */
  requiredRole?: UserRole
  /** 未認証時のリダイレクト先 */
  redirectTo?: string
}

/**
 * PrivateRoute は認証・ロールを検証するルートガード。
 * 未認証 → redirectTo（デフォルト: /login）にリダイレクト。
 * ロール不足 → /unauthorized にリダイレクト。
 */
export function PrivateRoute({ requiredRole, redirectTo = '/login' }: PrivateRouteProps) {
  const { isAuthenticated, user } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to={redirectTo} replace />
  }

  if (requiredRole && user) {
    const userLevel = ROLE_LEVEL[user.role] ?? 0
    const requiredLevel = ROLE_LEVEL[requiredRole] ?? 0
    if (userLevel < requiredLevel) {
      return <Navigate to="/unauthorized" replace />
    }
  }

  return <Outlet />
}
