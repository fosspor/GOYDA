import { useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { AuthContext } from './auth-context'

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
