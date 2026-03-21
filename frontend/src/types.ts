export type User = {
  id: string
  email: string
  display_name: string
  interests: string[]
}

export type Location = {
  id: string
  name: string
  description: string
  category: string
  seasons: string[]
  lat?: number
  lng?: number
  media_urls: string[]
}

export type RouteItem = {
  id: string
  user_id: string
  title: string
  season: string
  payload: unknown
  created_at: string
}

export type GeoPoint = {
  lat: number
  lng: number
}

export type PointWeather = {
  source: string
  lat: number
  lng: number
  temp_c: number
  condition: string
  wind_speed_ms: number
}

/** Ответ POST /api/ai/generate-route (единый контракт) */
export type GenerateRouteResponse = {
  source: 'mock' | 'yandex'
  user_id: string | null
  route: Record<string, unknown>
}

export type LocationsPage = {
  items: Location[]
  total: number
  limit: number
  offset: number
}

export type WeatherAwareRoute = {
  source: {
    routing: string
    weather_from: string
    weather_to: string
  }
  date: string
  from: Location
  to: Location
  route: {
    distance_m: number
    duration_s: number
    polyline: GeoPoint[]
  }
  weather: {
    from: {
      lat: number
      lng: number
      temp_c: number
      condition: string
      wind_speed_ms: number
    }
    to: {
      lat: number
      lng: number
      temp_c: number
      condition: string
      wind_speed_ms: number
    }
  }
  score: number
  reasoning: string
  saved_route_id: string
}
