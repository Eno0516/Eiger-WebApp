import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'
import type { AuthUser } from '../types/auth.types'

// JWT の payload を base64 デコードして AuthUser を取得するヘルパー
function parseToken(token: string): AuthUser | null {
  // TODO: token.split('.')[1] を atob でデコードして JSON.parse する
  // TODO: { user_id, email, role } を AuthUser にマップして返す
  // TODO: デコード失敗時は null を返す
  return null
}

const TOKEN_KEY = 'auth_token'

interface AuthContextValue {
  user: AuthUser | null
  token: string | null
  login: (token: string) => void
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY))
  const [user, setUser] = useState<AuthUser | null>(() => {
    const t = localStorage.getItem(TOKEN_KEY)
    return t ? parseToken(t) : null
  })

  const login = useCallback((newToken: string) => {
    localStorage.setItem(TOKEN_KEY, newToken)
    setToken(newToken)
    setUser(parseToken(newToken))
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY)
    setToken(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
