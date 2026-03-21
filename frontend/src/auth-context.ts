import { createContext } from 'react'

export type AuthContextValue = {
  token: string | null
  setToken: (value: string | null) => void
}

export const AuthContext = createContext<AuthContextValue | null>(null)
