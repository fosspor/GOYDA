# Быстрый старт - GOYDA

## Предварительные требования

- **Python 3.10** (используется в проекте)
- **Node.js 16+** (для React)
- **PostgreSQL 12+** с расширением PostGIS
- **Git**
- Yandex Cloud Account с API ключом

## 1️⃣ Клонирование и начальная настройка

```bash
# Клонировать проект (если еще не сделано)
cd /Users/fosspor/GOYDA

# Активировать виртуальное окружение Python
source venv/bin/activate
```

## 2️⃣ PostgreSQL + PostGIS Setup

### macOS (Homebrew)

```bash
# Установить PostgreSQL + PostGIS
brew install postgresql@15 postgis

# Запустить PostgreSQL
brew services start postgresql@15

# Создать базу данных и подключить PostGIS
createdb krasnodar_tourism
psql krasnodar_tourism -c "CREATE EXTENSION IF NOT EXISTS postgis;"
psql krasnodar_tourism -c "CREATE EXTENSION IF NOT EXISTS postgis_topology;"

# Проверить установку
psql krasnodar_tourism -c "SELECT PostGIS_version();"
```

## 3️⃣ Backend Setup

```bash
# Перейти в папку backend
cd backend

# Установить Python зависимости
pip install -r requirements.txt

# Создать .env файл
cp .env.example .env
```

### Отредактировать `.env`

```env
# Django
SECRET_KEY=django-insecure-dev-your-secret-key-here
DEBUG=True
ALLOWED_HOSTS=localhost,127.0.0.1,localhost:3000

# Database
DB_NAME=krasnodar_tourism
DB_USER=<ваш-postgres-user>
DB_PASSWORD=<ваш-postgres-password>
DB_HOST=localhost
DB_PORT=5432

# Yandex Cloud
YANDEX_FOLDER_ID=<your-folder-id>
YANDEX_API_KEY=<your-api-key>
YANDEX_IAM_TOKEN=<optional-iam-token>

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8000
```

### Инициализация БД

```bash
# Создать миграции
python manage.py makemigrations

# Применить миграции
python manage.py migrate

# Создать суперпользователя (для админ панели)
python manage.py createsuperuser
# Следовать подсказкам для ввода username, email, password

# (Опционально) Загрузить тестовые данные
python manage.py shell
```

### Запустить Backend

```bash
# Terminal 1
source ../venv/bin/activate
cd backend
python manage.py runserver 8000
```

Backend будет доступен на `http://localhost:8000`
Admin панель: `http://localhost:8000/admin`

## 4️⃣ Frontend Setup

```bash
# Перейти в папку frontend
cd frontend

# Установить Node зависимости
npm install

# Запустить dev сервер
npm start
```

Frontend откроется на `http://localhost:3000`

## 5️⃣ Проверка работы

### Backend API

```bash
# Тестировать API
curl http://localhost:8000/api/locations/

# Должно вернуть:
# {"count": 0, "next": null, "results": []}
```

### Frontend

Откройте в браузере `http://localhost:3000` и вы должны увидеть:
- Навигационный бар с GOYDA логотипом
- Hero section с "Начать путешествие"
- Раздел "Наши возможности"

## 🔨 Разработка

### Python/Django

```bash
# Создание новых миграций
cd backend
python manage.py makemigrations <app_name>
python manage.py migrate

# Запуск тестов
pytest

# Код стиль
black .
isort .
flake8 .
```

### React/Frontend

```bash
cd frontend

# Запустить тесты
npm test

# Собрать для production
npm run build

# Анализ бандла
npm run analyze
```

## 🚀 Основные фичи для реализации

### Фаза 1 — MVP
- [ ] Модели и АПИ для локаций (чтение)
- [ ] Список локаций на фронтенде с поиском
- [ ] Интерактивная карта (Leaflet)
- [ ] Базовая аутентификация
- [ ] AI маршрут (запрос к Yandex)

### Фаза 2 — Расширение
- [ ] Сохранение маршрутов в БД
- [ ] 360° галереи для локаций
- [ ] Отзывы и рейтинги
- [ ] Логистика (маршруты, транспорт)
- [ ] Admin панель для добавления локаций

### Фаза 3 — Оптимизация
- [ ] Мобильное приложение (React Native)
- [ ] Оффлайн поддержка
- [ ] Интеграция платежей
- [ ] Push-уведомления

## 🐛 Troubleshooting

### PostgreSQL не запускается
```bash
# Проверить статус
brew services list

# Перезапустить
brew services stop postgresql@15
brew services start postgresql@15
```

### PostGIS не установлен
```bash
# Проверить установку
psql krasnodar_tourism -c "SELECT PostGIS_version();"

# Если ошибка, переустановить
brew reinstall postgis
```

### Django миграции не работают
```bash
# Проверить статус БД
python manage.py dbshell
# SELECT 1; (должно вернуть 1)

# Пересоздать БД
dropdb krasnodar_tourism
createdb krasnodar_tourism
psql krasnodar_tourism -c "CREATE EXTENSION postgis;"
python manage.py migrate
```

### React зависимости не установились
```bash
# Очистить кэш npm
npm cache clean --force

# Переустановить
rm -rf node_modules package-lock.json
npm install
```

## 📖 Документация

- [Backend API Documentation](ARCHITECTURE.md)
- [Django Models Reference](../backend/README.md)
- [React Components Guide](../frontend/README.md)

## 💻 Git Workflow

```bash
# Создать ветку для фичи
git checkout -b feature/location-search

# Сделать изменения
git add .
git commit -m "Add location search"

# Запушить
git push origin feature/location-search

# Создать Pull Request
```

## 🆘 Быстрая помощь

```bash
# Просмотр логов Django
python manage.py runserver --verbosity 2

# Оптимизация БД
python manage.py shell
from django.core.management import call_command
call_command('optimize')

# Проверка конфига
python manage.py check
```

---

🎉 Готово! Теперь можете начинать разработку основных фич!
