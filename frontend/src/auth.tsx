import { createContext, useContext, useMemo, useState } from 'react'
import type { ReactNode } from 'react'

type AuthContextValue = {
  token: string | null
  setToken: (value: string | null) => void
}

const AuthContext = createContext<AuthContextValue | null>(null)
const STORAGE_KEY = 'goyda_token'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => localStorage.getItem(STORAGE_KEY))
  const setToken = (value: string | null) => {
    if (value) {
      localStorage.setItem(STORAGE_KEY, value)
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
    setTokenState(value)
  }
  const value = useMemo(() => ({ token, setToken }), [token])
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}
