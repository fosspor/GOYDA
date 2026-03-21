# Локальный смоук GOYDA

## Автоматика (терминал)

1. Поднимите стек из корня репозитория:

   ```bash
   docker compose up --build
   ```

2. В другом терминале:

   ```bash
   ./scripts/smoke.sh http://127.0.0.1:8080
   ```

   Скрипт проверяет `GET /health` и публичный `GET /api/locations`.

## Браузер

Откройте **http://localhost:8080** (UI и API на одном порту в Docker-образе).

Порядок проверки:

1. **Регистрация** — sidebar → Register, создать пользователя.
2. **Профиль** — Me / PATCH интересы.
3. **Локации** — ListLocations, открыть карточку по ссылке из списка.
4. **AI** — GenerateRoute, при необходимости «Сохранить как маршрут».
5. **Маршруты** — список и деталь сохранённого маршрута.
6. **Создать локацию** — форма CreateLocation (нужен JWT).

## CI на GitHub

После push в `main` смотрите статус:  
[github.com/fosspor/GOYDA/actions](https://github.com/fosspor/GOYDA/actions)

## Без Docker

Соберите фронт и бинарь со вшитым SPA (нужны Go, Node, PostgreSQL):

```bash
./scripts/sync-spa-dist.sh
export DATABASE_URL=postgres://...
export JWT_SECRET=dev-secret-min-8-chars
go run -tags embed ./cmd/server
```

Затем снова `./scripts/smoke.sh http://127.0.0.1:8080`.
