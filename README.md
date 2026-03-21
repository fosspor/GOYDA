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

## Если осталась папка `frontend/`

На macOS иногда нельзя удалить `node_modules` из песочницы IDE. Удалите вручную в терминале:

```bash
cd /path/to/GOYDA
chmod -R u+w frontend 2>/dev/null
rm -rf frontend
```

В `.gitignore` добавлено `/frontend/`, чтобы мусор не попадал в коммиты.

## Лицензия

MIT
