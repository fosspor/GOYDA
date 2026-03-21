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

## Фронтенд (SPA)

В репозитории есть клиент на `Vite + React + TypeScript` в [`frontend`](./frontend). Интерфейс оформлен как справочник эндпоинтов в духе [alt:V NativeDB](https://natives.altv.mp/): тёмная тема, боковая навигация с группами и поиском.

1. Запустите API и БД:

```bash
docker compose up --build
```

2. В отдельном терминале запустите фронтенд:

```bash
cd frontend
cp .env.example .env
npm ci
npm run dev
```

SPA: `http://localhost:3000`  
По умолчанию API берётся из `VITE_API_URL=http://localhost:8080`.

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

Контракт подробнее: [`api/openapi.yaml`](./api/openapi.yaml).

## Лицензия

MIT
