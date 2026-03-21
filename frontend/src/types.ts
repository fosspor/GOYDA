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
