import type {
  GenerateRouteResponse,
  Location,
  LocationsPage,
  PointWeather,
  RouteItem,
  User,
  WeatherAwareRoute,
} from './types'

/** Не задано в .env → dev на :3000 ходит на :8080. Пустая строка в .env → тот же origin (embed в Docker). */
const raw = import.meta.env.VITE_API_URL
const API_PREFIX =
  raw === undefined ? 'http://localhost:8080'.replace(/\/$/, '') : String(raw).replace(/\/$/, '')

type RequestOptions = {
  method?: string
  token?: string | null
  body?: unknown
}

export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export function isApiError(err: unknown): err is ApiError {
  return err instanceof ApiError
}

export function getApiBaseUrl() {
  return API_PREFIX === '' ? window.location.origin : API_PREFIX
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (opts.token) {
    headers.Authorization = `Bearer ${opts.token}`
  }
  const res = await fetch(`${API_PREFIX}${path}`, {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  })
  if (!res.ok) {
    const data = await res.json().catch(() => ({}))
    throw new ApiError(res.status, data.detail ?? `HTTP ${res.status}`)
  }
  if (res.status === 204) {
    return undefined as T
  }
  const text = await res.text()
  if (!text) {
    return undefined as T
  }
  return JSON.parse(text) as T
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

export async function listLocations(
  search = '',
  page?: { limit: number; offset?: number },
): Promise<Location[] | LocationsPage> {
  const params = new URLSearchParams()
  if (search) params.set('search', search)
  if (page) {
    params.set('limit', String(page.limit))
    params.set('offset', String(page.offset ?? 0))
  }
  const qs = params.toString()
  const path = qs ? `/api/locations?${qs}` : '/api/locations'
  const data = await request<unknown>(path)
  if (page) {
    const p = data as Partial<LocationsPage>
    return {
      items: Array.isArray(p.items) ? p.items : [],
      total: typeof p.total === 'number' ? p.total : 0,
      limit: typeof p.limit === 'number' ? p.limit : page.limit,
      offset: typeof p.offset === 'number' ? p.offset : page.offset ?? 0,
    }
  }
  return Array.isArray(data) ? (data as Location[]) : []
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

export async function patchLocation(
  token: string,
  id: string,
  payload: Partial<Omit<Location, 'id'>>,
): Promise<Location> {
  return request(`/api/locations/${id}`, { method: 'PATCH', token, body: payload })
}

export async function deleteLocation(token: string, id: string): Promise<void> {
  await request(`/api/locations/${id}`, { method: 'DELETE', token })
}

export async function listRoutes(token: string): Promise<RouteItem[]> {
  const data = await request<unknown>('/api/routes', { token })
  return Array.isArray(data) ? (data as RouteItem[]) : []
}

export async function getRoute(token: string, id: string): Promise<RouteItem> {
  return request(`/api/routes/${id}`, { token })
}

export async function createRoute(token: string, payload: { title: string; season: string; payload: unknown }): Promise<RouteItem> {
  return request('/api/routes', { method: 'POST', token, body: payload })
}

export async function patchRoute(
  token: string,
  id: string,
  payload: Partial<{ title: string; season: string; payload: unknown }>,
): Promise<RouteItem> {
  return request(`/api/routes/${id}`, { method: 'PATCH', token, body: payload })
}

export async function deleteRoute(token: string, id: string): Promise<void> {
  await request(`/api/routes/${id}`, { method: 'DELETE', token })
}

export async function aiGenerate(payload: {
  interests: string[]
  season: string
  days: number
  notes: string
}, token?: string | null): Promise<GenerateRouteResponse> {
  return request('/api/ai/generate-route', { method: 'POST', token, body: payload })
}

export async function aiRecommendations(season: string): Promise<{ season: string; items: Location[] }> {
  const q = season ? `?season=${encodeURIComponent(season)}` : ''
  const data = await request<{ season?: string; items?: unknown }>(`/api/ai/recommendations${q}`)
  return {
    season: data.season ?? season,
    items: Array.isArray(data.items) ? (data.items as Location[]) : [],
  }
}

export async function weatherPoint(lat: number, lng: number): Promise<PointWeather> {
  const q = `?lat=${encodeURIComponent(String(lat))}&lng=${encodeURIComponent(String(lng))}`
  return request(`/api/weather/point${q}`)
}

export async function weatherAwareRoute(
  token: string,
  payload: {
    from_location_id?: string
    to_location_id?: string
    date?: string
    avoid_rain?: boolean
    max_wind_ms?: number
  },
): Promise<WeatherAwareRoute> {
  return request('/api/routes/weather-aware', { method: 'POST', token, body: payload })
}
