# GOYDA API

Backend API для хакатона: персональные маршруты по Краснодарскому краю. Стек: **Go 1.22**, **Fiber**, **PostgreSQL + PostGIS**, JWT, опционально **Yandex Foundation Models**.

История репозитория до сброса сохранена в ветке **`archive/pre-go-reset`**.

## Быстрый старт (Docker)

```bash
cp .env.example .env
# при необходимости отредактируйте JWT_SECRET и CORS_ORIGINS
docker compose up --build
```

API: `http://localhost:8080`  
Health: `GET /health`

После `docker compose up --build` **интерфейс вшит в образ**: откройте **`http://localhost:8080`** — тот же порт, что и у API (линия **B**: один origin, без CORS между UI и `/api`).

Проверка:

```bash
./scripts/smoke.sh http://127.0.0.1:8080
```

Подробный чеклист (браузер, без Docker): [`docs/local-smoke.md`](docs/local-smoke.md).  
Полная документация API: [`docs/API.md`](docs/API.md).  
CI после изменений в `main`: [Actions](https://github.com/fosspor/GOYDA/actions).

Идеи следующих итераций: [`docs/NEXT.md`](docs/NEXT.md).

В браузере: регистрация → профиль → локации → AI → сохранение маршрута → создание локации.

## Фронтенд (SPA)

В репозитории есть клиент на `Vite + React + TypeScript` в [`frontend`](./frontend). Интерфейс оформлен как справочник эндпоинтов в духе [alt:V NativeDB](https://natives.altv.mp/): тёмная тема, боковая навигация с группами и поиском.

### Режим «только Docker» (UI уже в образе)

```bash
docker compose up --build
```

Откройте `http://localhost:8080`.

### Режим разработки (hot reload, Vite на :3000)

1. Запустите API и БД:

```bash
docker compose up --build
```

2. В отдельном терминале:

```bash
cd frontend
cp .env.example .env
npm ci
npm run dev
```

SPA: `http://localhost:3000`  
В `.env` задайте `VITE_API_URL=http://localhost:8080`.

### Сборка бинаря со вшитым SPA (локально)

```bash
./scripts/sync-spa-dist.sh
go build -tags embed -o bin/server ./cmd/server
```

Без тега `embed` бинарь отдаёт только API; UI тогда через `npm run dev` или отдельный хостинг.

## Локально (без Docker-образа API)

1. Поднимите только БД: `docker compose up -d db`
2. Установите Go 1.22+, затем:

```bash
export DATABASE_URL=postgres://goyda:goyda@localhost:5432/goyda?sslmode=disable
export JWT_SECRET=dev-secret-min-8-chars
go run ./cmd/server
```

## Переменные окружения

См. [`.env.example`](./.env.example).

- `MIGRATIONS_DIR` — папка с SQL (по умолчанию `./migrations`).
- `YANDEX_FOLDER_ID`, `YANDEX_API_KEY` — для живой генерации через Yandex; без них `POST /api/ai/generate-route` отдаёт **mock**-маршрут.
- `YANDEX_WEATHER_API_KEY`, `YANDEX_ROUTING_API_KEY` — для live-погоды и live-маршрутов; если пусто, используются mock/fallback ответы.

## Основные эндпоинты

| Метод | Путь | Описание |
|--------|------|-----------|
| GET | `/health` | Проверка |
| POST | `/api/auth/register` | Регистрация |
| POST | `/api/auth/login` | Вход |
| GET | `/api/me` | Профиль (Bearer JWT) |
| PATCH | `/api/me` | Обновить интересы |
| GET | `/api/locations` | Список, `?search=` |
| GET | `/api/locations/:id` | Детали |
| POST | `/api/locations` | Создать (JWT) |
| GET | `/api/routes` | Мои маршруты (JWT) |
| POST | `/api/routes` | Создать (JWT) |
| GET | `/api/routes/:id` | Детали (JWT, только свои) |
| POST | `/api/ai/generate-route` | Генерация (Yandex или mock) |
| GET | `/api/ai/recommendations` | Подбор по `?season=` |
| GET | `/api/weather/point` | Погода в точке (`lat`,`lng`) |
| POST | `/api/routes/weather-aware` | Маршрут по погоде + сохранение (JWT) |

Контракт подробнее: [`api/openapi.yaml`](./api/openapi.yaml) и [`docs/API.md`](docs/API.md).

## Лицензия

MIT
