import { useEffect, useMemo, useState } from 'react'
import type { FormEvent, ReactElement } from 'react'
import {
  Link,
  Navigate,
  NavLink,
  Outlet,
  Route,
  Routes,
  useNavigate,
  useParams,
} from 'react-router-dom'
import {
  aiGenerate,
  aiRecommendations,
  createLocation,
  createRoute,
  getLocation,
  getRoute,
  listLocations,
  listRoutes,
  login,
  me,
  patchMe,
  register,
} from './api'
import { useAuth } from './auth'
import { NAV_GROUPS, NAV_ITEMS } from './nav'
import type { Location, RouteItem, User } from './types'

function parseCSV(value: string): string[] {
  return value.split(',').map((v) => v.trim()).filter(Boolean)
}

function Protected({ children }: { children: ReactElement }) {
  const { token } = useAuth()
  if (!token) return <Navigate to="/login" replace />
  return children
}

function methodClass(m: string) {
  if (m === 'GET') return 'ndb-method ndb-method-get'
  if (m === 'POST') return 'ndb-method ndb-method-post'
  if (m === 'PATCH') return 'ndb-method ndb-method-patch'
  return 'ndb-method'
}

function Layout() {
  const { token, setToken } = useAuth()
  const [q, setQ] = useState('')
  const filtered = useMemo(() => {
    const s = q.trim().toLowerCase()
    return NAV_ITEMS.filter((item) => {
      if (item.requireAuth && !token) return false
      if (!s) return true
      return (
        item.id.toLowerCase().includes(s) ||
        item.path.toLowerCase().includes(s) ||
        item.method.toLowerCase().includes(s) ||
        item.group.toLowerCase().includes(s)
      )
    })
  }, [q, token])

  return (
    <div className="ndb-root">
      <aside className="ndb-sidebar">
        <div className="ndb-brand">
          <Link to="/" className="ndb-brand-link">
            <span className="ndb-brand-title">GOYDA</span>
            <span className="ndb-brand-sep">/</span>
            <span className="ndb-brand-sub">API</span>
          </Link>
          <p className="ndb-brand-tag">справочник эндпоинтов</p>
        </div>
        <input
          className="ndb-search"
          type="search"
          placeholder="Поиск по имени или пути…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          aria-label="Поиск эндпоинтов"
        />
        <div className="ndb-nav-scroll">
          {NAV_GROUPS.map((g) => {
            const items = filtered.filter((i) => i.group === g.id)
            if (items.length === 0) return null
            return (
              <div key={g.id} className="ndb-group">
                <div className="ndb-group-title">{g.title}</div>
                <ul className="ndb-list">
                  {items.map((item) => (
                    <li key={`${item.id}-${item.method}-${item.path}`}>
                      <NavLink to={item.to} className={({ isActive }) => `ndb-item${isActive ? ' ndb-item-active' : ''}`}>
                        <span className={methodClass(item.method)}>{item.method}</span>
                        <span className="ndb-item-body">
                          <span className="ndb-item-id">{item.id}</span>
                          <code className="ndb-item-path">{item.path}</code>
                        </span>
                      </NavLink>
                    </li>
                  ))}
                </ul>
              </div>
            )
          })}
        </div>
        <div className="ndb-sidebar-foot">
          {token ? (
            <button type="button" className="ndb-btn-ghost" onClick={() => setToken(null)}>
              Выйти
            </button>
          ) : (
            <NavLink to="/login" className="ndb-btn-ghost">
              Вход
            </NavLink>
          )}
        </div>
      </aside>
      <div className="ndb-main">
        <Outlet />
      </div>
    </div>
  )
}

function HomePage() {
  const raw = import.meta.env.VITE_API_URL
  const base =
    raw === undefined ? 'http://localhost:8080' : raw === '' ? 'same origin (как у страницы)' : String(raw)
  const publicCount = NAV_ITEMS.filter((i) => !i.requireAuth).length
  const authCount = NAV_ITEMS.filter((i) => i.requireAuth).length
  return (
    <article className="ndb-panel ndb-panel-hero">
      <h1 className="ndb-hero-title">GOYDA / API</h1>
      <p className="ndb-hero-lead">
        Оболочка в духе{' '}
        <a href="https://natives.altv.mp/" target="_blank" rel="noreferrer">
          alt:V NativeDB
        </a>
        : тёмная тема, боковой список и быстрый поиск по эндпоинтам.
      </p>
      <dl className="ndb-stats">
        <div>
          <dt>Всего в каталоге</dt>
          <dd>{NAV_ITEMS.length}</dd>
        </div>
        <div>
          <dt>Публичных</dt>
          <dd>{publicCount}</dd>
        </div>
        <div>
          <dt>С JWT</dt>
          <dd>{authCount}</dd>
        </div>
        <div>
          <dt>База API</dt>
          <dd>
            <code>{base}</code>
          </dd>
        </div>
        <div>
          <dt>Health</dt>
          <dd>
            <code>GET /health</code>
          </dd>
        </div>
      </dl>
      <p className="ndb-muted">Пункты с JWT в боковой панели скрыты, пока вы не вошли.</p>
    </article>
  )
}

function LocationLookupPage() {
  const [id, setId] = useState('')
  const navigate = useNavigate()
  const onSubmit = (e: FormEvent) => {
    e.preventDefault()
    const v = id.trim()
    if (v) navigate(`/locations/${v}`)
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> GetLocation
      </h2>
      <p className="ndb-muted">
        <code>/api/locations/:id</code>
      </p>
      <form onSubmit={onSubmit} className="ndb-form">
        <label className="ndb-label">UUID локации</label>
        <input value={id} onChange={(e) => setId(e.target.value)} placeholder="00000000-0000-0000-0000-000000000000" />
        <button type="submit">Открыть</button>
      </form>
    </section>
  )
}

function RouteLookupPage() {
  const [id, setId] = useState('')
  const navigate = useNavigate()
  const onSubmit = (e: FormEvent) => {
    e.preventDefault()
    const v = id.trim()
    if (v) navigate(`/routes/${v}`)
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> GetRoute
      </h2>
      <p className="ndb-muted">
        <code>/api/routes/:id</code> — только свои маршруты (JWT).
      </p>
      <form onSubmit={onSubmit} className="ndb-form">
        <label className="ndb-label">UUID маршрута</label>
        <input value={id} onChange={(e) => setId(e.target.value)} placeholder="uuid" />
        <button type="submit">Открыть</button>
      </form>
    </section>
  )
}

function LocationsPage() {
  const [items, setItems] = useState<Location[]>([])
  const [search, setSearch] = useState('')
  const [error, setError] = useState('')
  const load = async (q = '') => {
    try {
      setError('')
      setItems(await listLocations(q))
    } catch (e) {
      setError((e as Error).message)
    }
  }
  useEffect(() => { void load() }, [])
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> ListLocations
      </h2>
      <p className="ndb-muted">
        <code>/api/locations</code>
      </p>
      <form className="ndb-form" onSubmit={(e) => { e.preventDefault(); void load(search) }}>
        <input value={search} onChange={(e) => setSearch(e.target.value)} placeholder="поиск" />
        <button type="submit">Искать</button>
      </form>
      {error && <p className="error">{error}</p>}
      <ul>{items.map((l) => <li key={l.id}><Link to={`/locations/${l.id}`}>{l.name}</Link> ({l.category})</li>)}</ul>
    </section>
  )
}

function LocationDetailsPage() {
  const { id = '' } = useParams()
  const [item, setItem] = useState<Location | null>(null)
  const [error, setError] = useState('')
  useEffect(() => {
    void getLocation(id).then(setItem).catch((e) => setError((e as Error).message))
  }, [id])
  if (error) return <p className="error ndb-panel ndb-panel-pad">{error}</p>
  if (!item) return <p className="ndb-panel ndb-panel-pad ndb-muted">Загрузка...</p>
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">{item.name}</h2>
      <p>{item.description}</p>
      <p>Сезоны: {item.seasons.join(', ') || 'нет'}</p>
      <p>Координаты: {item.lat ?? '-'}, {item.lng ?? '-'}</p>
    </section>
  )
}

function RecommendationsPage() {
  const [season, setSeason] = useState('summer')
  const [items, setItems] = useState<Location[]>([])
  const [error, setError] = useState('')
  const load = async () => {
    try {
      setError('')
      const data = await aiRecommendations(season)
      setItems(data.items)
    } catch (e) {
      setError((e as Error).message)
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> AIRecommendations
      </h2>
      <p className="ndb-muted">
        <code>/api/ai/recommendations?season=</code>
      </p>
      <form className="ndb-form" onSubmit={(e) => { e.preventDefault(); void load() }}>
        <input value={season} onChange={(e) => setSeason(e.target.value)} placeholder="season" />
        <button type="submit">Получить</button>
      </form>
      {error && <p className="error">{error}</p>}
      <ul>{items.map((l) => <li key={l.id}>{l.name}</li>)}</ul>
    </section>
  )
}

function LoginPage() {
  const { setToken } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const data = await login({ email, password })
      setToken(data.token)
      navigate('/profile')
    } catch (err) {
      setError((err as Error).message)
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('POST')}>POST</span> Login
      </h2>
      <p className="ndb-muted">
        <code>/api/auth/login</code>
      </p>
      <form className="ndb-form" onSubmit={onSubmit}>
        <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email" />
        <input value={password} onChange={(e) => setPassword(e.target.value)} placeholder="password" type="password" />
        <button type="submit">Войти</button>
      </form>
      <p><Link to="/register">Нет аккаунта? Регистрация</Link></p>
      {error && <p className="error">{error}</p>}
    </section>
  )
}

function RegisterPage() {
  const { setToken } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [interests, setInterests] = useState('')
  const [error, setError] = useState('')
  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    try {
      const data = await register({ email, password, display_name: displayName, interests: parseCSV(interests) })
      setToken(data.token)
      navigate('/profile')
    } catch (err) {
      setError((err as Error).message)
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('POST')}>POST</span> Register
      </h2>
      <p className="ndb-muted">
        <code>/api/auth/register</code>
      </p>
      <form className="ndb-form" onSubmit={onSubmit}>
        <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email" />
        <input value={password} onChange={(e) => setPassword(e.target.value)} placeholder="password" type="password" />
        <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} placeholder="display_name" />
        <input value={interests} onChange={(e) => setInterests(e.target.value)} placeholder="interests: wine, culture" />
        <button type="submit">Создать аккаунт</button>
      </form>
      {error && <p className="error">{error}</p>}
    </section>
  )
}

function ProfilePage() {
  const { token } = useAuth()
  const [user, setUser] = useState<User | null>(null)
  const [interests, setInterests] = useState('')
  const [error, setError] = useState('')
  useEffect(() => {
    if (!token) return
    void me(token)
      .then((u) => {
        setUser(u)
        setInterests(u.interests.join(', '))
      })
      .catch((e) => setError((e as Error).message))
  }, [token])
  const onSave = async (e: FormEvent) => {
    e.preventDefault()
    if (!token) return
    try {
      const updated = await patchMe(token, parseCSV(interests))
      setUser(updated)
    } catch (err) {
      setError((err as Error).message)
    }
  }
  if (!user) return <p className="ndb-panel ndb-panel-pad">{error || 'Загрузка профиля...'}</p>
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> / <span className={methodClass('PATCH')}>PATCH</span> Me
      </h2>
      <p className="ndb-muted">
        <code>/api/me</code>
      </p>
      <p>{user.email}</p>
      <form className="ndb-form" onSubmit={onSave}>
        <input value={interests} onChange={(e) => setInterests(e.target.value)} />
        <button type="submit">Сохранить интересы</button>
      </form>
      {error && <p className="error">{error}</p>}
    </section>
  )
}

function RoutesPage() {
  const { token } = useAuth()
  const [routes, setRoutes] = useState<RouteItem[]>([])
  const [title, setTitle] = useState('')
  const [season, setSeason] = useState('')
  const [payload, setPayload] = useState('{}')
  const [error, setError] = useState('')
  const load = async () => {
    if (!token) return
    setRoutes(await listRoutes(token))
  }
  useEffect(() => { void load() }, [token])
  const onCreate = async (e: FormEvent) => {
    e.preventDefault()
    if (!token) return
    try {
      await createRoute(token, { title, season, payload: JSON.parse(payload) })
      setTitle('')
      setSeason('')
      setPayload('{}')
      await load()
    } catch (err) {
      setError((err as Error).message)
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('GET')}>GET</span> / <span className={methodClass('POST')}>POST</span> Routes
      </h2>
      <p className="ndb-muted">
        <code>/api/routes</code>
      </p>
      <form className="ndb-form" onSubmit={onCreate}>
        <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="title" />
        <input value={season} onChange={(e) => setSeason(e.target.value)} placeholder="season" />
        <textarea value={payload} onChange={(e) => setPayload(e.target.value)} rows={5} />
        <button type="submit">Создать маршрут</button>
      </form>
      {error && <p className="error">{error}</p>}
      <ul>{routes.map((r) => <li key={r.id}><Link to={`/routes/${r.id}`}>{r.title}</Link></li>)}</ul>
    </section>
  )
}

function RouteDetailsPage() {
  const { token } = useAuth()
  const { id = '' } = useParams()
  const [route, setRoute] = useState<RouteItem | null>(null)
  const [error, setError] = useState('')
  useEffect(() => {
    if (!token) return
    void getRoute(token, id).then(setRoute).catch((e) => setError((e as Error).message))
  }, [token, id])
  if (error) return <p className="error ndb-panel ndb-panel-pad">{error}</p>
  if (!route) return <p className="ndb-panel ndb-panel-pad ndb-muted">Загрузка...</p>
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">{route.title}</h2>
      <pre>{JSON.stringify(route.payload, null, 2)}</pre>
    </section>
  )
}

function CreateLocationPage() {
  const { token } = useAuth()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [category, setCategory] = useState('')
  const [seasons, setSeasons] = useState('')
  const [lat, setLat] = useState('')
  const [lng, setLng] = useState('')
  const [mediaUrls, setMediaUrls] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!token) return
    try {
      const location = await createLocation(token, {
        name,
        description,
        category,
        seasons: parseCSV(seasons),
        lat: lat ? Number(lat) : undefined,
        lng: lng ? Number(lng) : undefined,
        media_urls: parseCSV(mediaUrls),
      })
      setMessage(`Создано: ${location.id}`)
      setError('')
    } catch (err) {
      setError((err as Error).message)
      setMessage('')
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('POST')}>POST</span> CreateLocation
      </h2>
      <p className="ndb-muted">
        <code>/api/locations</code> (JWT)
      </p>
      <form className="ndb-form" onSubmit={onSubmit}>
        <input value={name} onChange={(e) => setName(e.target.value)} placeholder="name" />
        <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="description" />
        <input value={category} onChange={(e) => setCategory(e.target.value)} placeholder="category" />
        <input value={seasons} onChange={(e) => setSeasons(e.target.value)} placeholder="seasons: summer,autumn" />
        <input value={lat} onChange={(e) => setLat(e.target.value)} placeholder="lat (optional)" />
        <input value={lng} onChange={(e) => setLng(e.target.value)} placeholder="lng (optional)" />
        <input value={mediaUrls} onChange={(e) => setMediaUrls(e.target.value)} placeholder="media_urls csv" />
        <button type="submit">Создать</button>
      </form>
      {message && <p>{message}</p>}
      {error && <p className="error">{error}</p>}
    </section>
  )
}

function AIPage() {
  const { token } = useAuth()
  const [interests, setInterests] = useState('')
  const [season, setSeason] = useState('summer')
  const [days, setDays] = useState(3)
  const [notes, setNotes] = useState('')
  const [result, setResult] = useState<Record<string, unknown> | null>(null)
  const [error, setError] = useState('')
  const [saveMsg, setSaveMsg] = useState('')
  const routePayload = useMemo(() => (result?.route ?? result?.raw ?? {}), [result])
  const onGenerate = async (e: FormEvent) => {
    e.preventDefault()
    try {
      setError('')
      setResult(await aiGenerate({ interests: parseCSV(interests), season, days, notes }, token))
    } catch (err) {
      setError((err as Error).message)
    }
  }
  const onSave = async () => {
    if (!token || !result) return
    try {
      setSaveMsg('')
      await createRoute(token, {
        title: `AI route ${season}`,
        season,
        payload: routePayload,
      })
      setSaveMsg('Маршрут сохранён в /api/routes')
    } catch (err) {
      setError((err as Error).message)
    }
  }
  return (
    <section className="ndb-panel">
      <h2 className="ndb-panel-title">
        <span className={methodClass('POST')}>POST</span> GenerateRoute
      </h2>
      <p className="ndb-muted">
        <code>/api/ai/generate-route</code> — Yandex или mock.
      </p>
      <form className="ndb-form" onSubmit={onGenerate}>
        <input value={interests} onChange={(e) => setInterests(e.target.value)} placeholder="interests csv" />
        <input value={season} onChange={(e) => setSeason(e.target.value)} placeholder="season" />
        <input value={days} onChange={(e) => setDays(Number(e.target.value))} type="number" min={1} />
        <textarea value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="notes" />
        <button type="submit">Сгенерировать</button>
      </form>
      {result && <button onClick={() => void onSave()} disabled={!token}>Сохранить как маршрут</button>}
      {saveMsg && <p className="ndb-ok">{saveMsg}</p>}
      {error && <p className="error">{error}</p>}
      {result && <pre>{JSON.stringify(result, null, 2)}</pre>}
    </section>
  )
}

function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/locations" element={<LocationsPage />} />
        <Route path="/locations/lookup" element={<LocationLookupPage />} />
        <Route path="/locations/:id" element={<LocationDetailsPage />} />
        <Route path="/recommendations" element={<RecommendationsPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/profile" element={<Protected><ProfilePage /></Protected>} />
        <Route path="/routes" element={<Protected><RoutesPage /></Protected>} />
        <Route path="/routes/lookup" element={<Protected><RouteLookupPage /></Protected>} />
        <Route path="/routes/:id" element={<Protected><RouteDetailsPage /></Protected>} />
        <Route path="/create-location" element={<Protected><CreateLocationPage /></Protected>} />
        <Route path="/ai" element={<AIPage />} />
      </Route>
    </Routes>
  )
}

export default App
