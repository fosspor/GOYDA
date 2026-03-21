import type { Location, RouteItem, User } from './types'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

type RequestOptions = {
  method?: string
  token?: string | null
  body?: unknown
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (opts.token) {
    headers.Authorization = `Bearer ${opts.token}`
  }
  const res = await fetch(`${API_URL}${path}`, {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    throw new Error(data.detail ?? `HTTP ${res.status}`)
  }
  return (await res.json()) as T
}

export async function register(payload: {
  email: string
  password: string
  display_name: string
  interests: string[]
}): Promise<{ token: string; user: User }> {
  return request('/api/auth/register', { method: 'POST', body: payload })
}

export async function login(payload: { email: string; password: string }): Promise<{ token: string; user: User }> {
  return request('/api/auth/login', { method: 'POST', body: payload })
}

export async function me(token: string): Promise<User> {
  return request('/api/me', { token })
}

export async function patchMe(token: string, interests: string[]): Promise<User> {
  return request('/api/me', { method: 'PATCH', token, body: { interests } })
}

export async function listLocations(search = ''): Promise<Location[]> {
  const q = search ? `?search=${encodeURIComponent(search)}` : ''
  return request(`/api/locations${q}`)
}

export async function getLocation(id: string): Promise<Location> {
  return request(`/api/locations/${id}`)
}

export async function createLocation(
  token: string,
  payload: Omit<Location, 'id'>,
): Promise<Location> {
  return request('/api/locations', { method: 'POST', token, body: payload })
}

export async function listRoutes(token: string): Promise<RouteItem[]> {
  return request('/api/routes', { token })
}

export async function getRoute(token: string, id: string): Promise<RouteItem> {
  return request(`/api/routes/${id}`, { token })
}

export async function createRoute(token: string, payload: { title: string; season: string; payload: unknown }): Promise<RouteItem> {
  return request('/api/routes', { method: 'POST', token, body: payload })
}

export async function aiGenerate(payload: {
  interests: string[]
  season: string
  days: number
  notes: string
}, token?: string | null): Promise<Record<string, unknown>> {
  return request('/api/ai/generate-route', { method: 'POST', token, body: payload })
}

export async function aiRecommendations(season: string): Promise<{ season: string; items: Location[] }> {
  const q = season ? `?season=${encodeURIComponent(season)}` : ''
  return request(`/api/ai/recommendations${q}`)
}
