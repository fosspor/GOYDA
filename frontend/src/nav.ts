export type NavItem = {
  to: string
  /** Короткое имя «как у натива» */
  id: string
  /** Путь API для отображения */
  path: string
  method: 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'
  group: string
  requireAuth?: boolean
}

export const NAV_GROUPS: { id: string; title: string }[] = [
  { id: 'system', title: 'System' },
  { id: 'auth', title: 'Auth' },
  { id: 'profile', title: 'Profile' },
  { id: 'locations', title: 'Locations' },
  { id: 'routes', title: 'Routes' },
  { id: 'ai', title: 'AI' },
]

export const NAV_ITEMS: NavItem[] = [
  { to: '/', id: 'Overview', path: '/', method: 'GET', group: 'system' },
  { to: '/locations', id: 'ListLocations', path: '/api/locations', method: 'GET', group: 'locations' },
  { to: '/locations/lookup', id: 'GetLocation', path: '/api/locations/:id', method: 'GET', group: 'locations' },
  { to: '/recommendations', id: 'AIRecommendations', path: '/api/ai/recommendations', method: 'GET', group: 'ai' },
  { to: '/ai', id: 'GenerateRoute', path: '/api/ai/generate-route', method: 'POST', group: 'ai' },
  { to: '/register', id: 'Register', path: '/api/auth/register', method: 'POST', group: 'auth' },
  { to: '/login', id: 'Login', path: '/api/auth/login', method: 'POST', group: 'auth' },
  { to: '/profile', id: 'Me', path: '/api/me', method: 'GET', group: 'profile', requireAuth: true },
  { to: '/profile', id: 'PatchMe', path: '/api/me', method: 'PATCH', group: 'profile', requireAuth: true },
  { to: '/create-location', id: 'CreateLocation', path: '/api/locations', method: 'POST', group: 'locations', requireAuth: true },
  { to: '/routes', id: 'ListMyRoutes', path: '/api/routes', method: 'GET', group: 'routes', requireAuth: true },
  { to: '/routes', id: 'CreateRoute', path: '/api/routes', method: 'POST', group: 'routes', requireAuth: true },
  { to: '/routes/lookup', id: 'GetRoute', path: '/api/routes/:id', method: 'GET', group: 'routes', requireAuth: true },
]
