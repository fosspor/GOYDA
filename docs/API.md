# GOYDA API Documentation

Полная документация по REST API проекта GOYDA.  
Базовый URL по умолчанию: `http://localhost:8080`.

Источник истины для схем: [`api/openapi.yaml`](../api/openapi.yaml).

## Быстрый старт

```bash
cp .env.example .env
docker compose up -d --build
curl -sSf http://127.0.0.1:8080/health
```

## Авторизация

Для защищенных endpoint используется Bearer JWT:

```http
Authorization: Bearer <token>
```

Получение токена:

1. `POST /api/auth/register`
2. или `POST /api/auth/login`

## Формат ошибок

Все ошибки возвращаются в едином формате:

```json
{
  "detail": "error message"
}
```

## Endpoints

### Health

#### `GET /health`

Проверка доступности API.

Пример ответа:

```json
{
  "service": "goyda-api",
  "status": "ok"
}
```

---

### Auth

#### `POST /api/auth/register`

Регистрация пользователя.

Body:

```json
{
  "email": "user@example.com",
  "password": "secret12",
  "display_name": "User",
  "interests": ["wine", "gastro"]
}
```

Ответ: `201 Created`

```json
{
  "token": "<jwt>",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "display_name": "User",
    "interests": ["wine", "gastro"]
  }
}
```

#### `POST /api/auth/login`

Вход по email/паролю.

Body:

```json
{
  "email": "user@example.com",
  "password": "secret12"
}
```

Ответ: `200 OK` (та же схема, что у register).

---

### Profile

#### `GET /api/me` (JWT)

Текущий профиль пользователя.

#### `PATCH /api/me` (JWT)

Обновление интересов.

Body:

```json
{
  "interests": ["eco", "active"]
}
```

---

### Locations

#### `GET /api/locations`

Параметры:

- `search` (опционально)
- `limit` (опционально, `>= 1`)
- `offset` (опционально, `>= 0`)

Важно:

- если `limit/offset` **не переданы**, ответ — массив `Location[]`;
- если `limit` или `offset` переданы, ответ — объект пагинации:

```json
{
  "items": [],
  "total": 0,
  "limit": 20,
  "offset": 0
}
```

#### `GET /api/locations/{id}`

Получить одну локацию.

#### `POST /api/locations` (JWT)

Создать локацию.

Body:

```json
{
  "name": "Станица Тамань",
  "description": "История и кухня",
  "category": "culture",
  "seasons": ["spring", "summer"],
  "lat": 45.2129,
  "lng": 36.7184,
  "media_urls": []
}
```

#### `PATCH /api/locations/{id}` (JWT)

Частичное обновление локации.

Поддерживаются поля:
`name`, `description`, `category`, `seasons`, `media_urls`, `lat`, `lng`.

Если передается геопозиция, `lat` и `lng` должны быть переданы вместе.

#### `DELETE /api/locations/{id}` (JWT)

Удаление локации.  
Ответ: `204 No Content`.

---

### Routes

#### `GET /api/routes` (JWT)

Список маршрутов текущего пользователя.

#### `POST /api/routes` (JWT)

Создание маршрута.

Body:

```json
{
  "title": "Маршрут на выходные",
  "season": "summer",
  "payload": {}
}
```

#### `GET /api/routes/{id}` (JWT)

Получить маршрут (только свой).

#### `PATCH /api/routes/{id}` (JWT)

Частичное обновление маршрута.

Поддерживаются поля: `title`, `season`, `payload`.

#### `DELETE /api/routes/{id}` (JWT)

Удаление маршрута.  
Ответ: `204 No Content`.

---

### AI

#### `POST /api/ai/generate-route`

Генерация маршрута (через Yandex LLM или mock fallback).

Body:

```json
{
  "interests": ["wine", "gastro"],
  "season": "autumn",
  "days": 3,
  "notes": "без длинных переездов"
}
```

Ответ всегда содержит единый контейнер `route`:

```json
{
  "source": "mock",
  "user_id": "uuid-or-null",
  "route": {
    "kind": "mock",
    "title": "Маршрут...",
    "summary": "Черновик...",
    "stops": []
  }
}
```

Для `source = "yandex"` внутри `route` приходит `kind = "llm"`, `text` и при валидном JSON поле `json`.

#### `GET /api/ai/recommendations`

Подбор локаций по сезону.

Параметры:

- `season` (опционально, по умолчанию `summer`)

Ответ:

```json
{
  "season": "summer",
  "items": []
}
```

---

### Weather + Routing

#### `GET /api/weather/point`

Параметры:

- `lat` (обязательно)
- `lng` (обязательно)

Ответ:

```json
{
  "source": "yandex|mock",
  "lat": 45.0,
  "lng": 38.9,
  "temp_c": 23.1,
  "condition": "clear",
  "wind_speed_ms": 3.5
}
```

#### `POST /api/routes/weather-aware` (JWT)

Строит маршрут, оценивает погодный риск и сохраняет его в `routes`.

Body:

```json
{
  "from_location_id": "uuid",
  "to_location_id": "uuid",
  "date": "2026-03-21",
  "avoid_rain": true,
  "max_wind_ms": 10
}
```

Поля `from_location_id` и `to_location_id` можно не передавать — тогда backend выберет точки автоматически из геокодированных локаций.

---

## Полезные curl-примеры

Регистрация:

```bash
curl -sS -X POST http://127.0.0.1:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"secret12","display_name":"demo","interests":[]}'
```

Логин и токен:

```bash
TOKEN="$(curl -sS -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"secret12"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")"
```

Список локаций с пагинацией:

```bash
curl -sS "http://127.0.0.1:8080/api/locations?limit=20&offset=0"
```

AI генерация:

```bash
curl -sS -X POST http://127.0.0.1:8080/api/ai/generate-route \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"interests":["gastro"],"season":"summer","days":2,"notes":"семейный формат"}'
```

Weather-aware:

```bash
curl -sS -X POST http://127.0.0.1:8080/api/routes/weather-aware \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"avoid_rain":true,"max_wind_ms":10}'
```

## Smoke

Быстрая проверка API:

```bash
./scripts/smoke.sh http://127.0.0.1:8080
```
