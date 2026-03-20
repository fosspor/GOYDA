# GOYDA — Персональные маршруты по Краснодарскому краю

Веб-приложение для открытия малоизвестных, но интересных мест Краснодарского края и создания персональных маршрутов с помощью ИИ.

## 🎯 Основная проблема

Туристический информационный хаос в Краснодарском крае:
- Реклама одних и тех же популярных мест
- Сложно найти нестандартные локации (малые винодельни, фермы, мастерские)
- Неясная логистика: как добраться, где остановиться, чем питаться
- Непонятно, подойдет ли конкретный тур для вашей группы

## 💡 Решение

GOYDA формирует **персональные маршруты** на любой сезон с опорой на:
- **Интересы пользователя** (гастрономия, вино, природа, история, искусство и т.д.)
- **Сезонность** (разные маршруты для весны, лета, осени, зимы)
- **AI-рекомендации** (Yandex LLM подбирает оптимальный маршрут)
- **Дистанционные туры** (360° галерея перед покупкой)

## 🏗️ Техническая архитектура

### Backend
- **Django + Django REST Framework** (Python 3.10)
- **PostgreSQL + PostGIS** для геоданных
- **Yandex LLM** для генерации маршрутов и описаний
- JWT аутентификация, CORS поддержка

### Frontend
- **React 18** с React Router
- **Tailwind CSS** для дизайна
- **Leaflet + React-Leaflet** для карт
- **Zustand** для управления состоянием
- **Swiper** для каруселей изображений

### Структура проекта
```
GOYDA/
├── backend/
│   ├── config/           # Django settings & URLs
│   ├── locations/        # Локации и категории
│   ├── routes/          # Маршруты и интересы
│   ├── users/           # Кастомный User model
│   ├── ai_recommendations/  # Yandex AI интеграция
│   ├── manage.py
│   ├── requirements.txt
│   └── .env.example
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── store/
│   │   └── App.js
│   ├── package.json
│   ├── tailwind.config.js
│   └── postcss.config.js
└── docs/               # Документация
```

## 🚀 Быстрый старт

### Требования
- Python 3.10
- Node.js 16+
- PostgreSQL 12+ с PostGIS
- Yandex Cloud API ключ

### Backend Setup

```bash
# 1. Активировать виртуальное окружение
source venv/bin/activate

# 2. Установить зависимости
cd backend
pip install -r requirements.txt

# 3. Создать .env файл
cp .env.example .env
# Отредактировать .env с вашими параметрами

# 4. Применить миграции
python manage.py migrate

# 5. Создать суперпользователя
python manage.py createsuperuser

# 6. Запустить сервер разработки
python manage.py runserver 8000
```

### Frontend Setup

```bash
# 1. Перейти в папку frontend
cd frontend

# 2. Установить зависимости
npm install

# 3. Запустить dev сервер
npm start
```

Приложение будет доступно на `http://localhost:3000`

## 📋 Основные компоненты

### 1. **Locations** (Локации)
- Модель для хранения достопримечательностей, ресторанов, винодельн и т.д.
- GIS координаты для карт
- Галереи, видео-туры (360°)
- Рейтинги, отзывы
- Сезонность и доступность

### 2. **Routes** (Маршруты)
- Персональные маршруты от сервиса
- История и сохраненные маршруты
- Порядок посещения локаций
- Длительность и логистика

### 3. **AI Recommendations** (ИИ Рекомендации)
- Интеграция с Yandex LLM
- Генерация маршрутов на основе интересов
- Обучение на основе обратной связи пользователя
- Логирование взаимодействий для аналитики

### 4. **Users** (Пользователи)
- Кастомный модель User с профилем
- Интересы и предпочтения
- История путешествий
- Поддержка разных типов пользователей (туристы, бизнес, гиды)

## 🤖 AI/LLM Интеграция

### Yandex LLM использование
1. **Генерация маршрутов** — на основе интересов, бюджета, сезона
2. **Описания локаций** — креативные и привлекательные описания
3. **Персонализация** — рекомендации под группу (семья, пенсионеры, фрилансеры)
4. **Советы по сезону** — что делать в разные времена года

### Пример запроса

```python
from ai_recommendations.services import get_yandex_service

service = get_yandex_service()
prompt = service.generate_route_prompt(
    user_interests=['wine', 'gastronomy'],
    budget=50000,
    duration=3,
    season='autumn',
    group_info={'size': 2, 'ages': [30, 35]}
)
response = service.call_completion_api(prompt)
```

## 📊 API Endpoints

### Locations
```
GET /api/locations/          - Список всех локаций
GET /api/locations/{id}/     - Детали локации
GET /api/locations/?search=  - Поиск
```

### Routes
```
GET /api/routes/             - Мои маршруты
POST /api/routes/            - Создать маршрут
GET /api/routes/{id}/        - Детали маршрута
```

### AI
```
POST /api/ai/generate-route/ - Генерировать маршрут
GET /api/ai/recommendations/ - Получить рекомендации
```

## 🎨 Дизайн & UX

- **Минималистичный, интуитивный интерфейс**
- **Мобильная оптимизация** (mobile-first)
- **Карты с локациями** и интерактивные туры
- **Персонализированная лента** маршрутов
- **360° галереи** для виртуальных посещений

## 📦 Stack Summary

| Компонент | Технология |
|-----------|-----------|
| Backend Framework | Django 4.2 + DRF |
| Base Language | Python 3.10 |
| Database | PostgreSQL + PostGIS |
| Frontend Framework | React 18 |
| Styling | Tailwind CSS |
| State Management | Zustand |
| Mapping | Leaflet + React-Leaflet |
| AI/LLM | Yandex LLM API |
| Authentication | JWT (SimpleJWT) |
| HTTP Client | Axios |

## 🔐 Безопасность

- CORS настроена только для разрешенных источников
- JWT токены для аутентификации
- Environment переменные для конфиденциальных данных
- Django security middleware включен

## 📝 Environment Variables

[backend/.env.example](backend/.env.example):
```env
SECRET_KEY=<your-secret-key>
DEBUG=True
ALLOWED_HOSTS=localhost,127.0.0.1

# Database
DB_NAME=krasnodar_tourism
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432

# Yandex Cloud
YANDEX_FOLDER_ID=<your-folder-id>
YANDEX_API_KEY=<your-api-key>

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

## 🛠️ Разработка

### Создание миграций
```bash
python manage.py makemigrations
python manage.py migrate
```

### Запуск тестов
```bash
pytest
```

### Code Style
```bash
black .
isort .
flake8 .
```

## 🚢 Развертывание

### Backend (Heroku/Vercel)
```bash
# Требуется PostgreSQL + PostGIS
gunicorn config.wsgi:application
```

### Frontend (Vercel/Netlify)
```bash
npm run build
# Deploy dist/ folder
```

## 📚 Документация

Смотрите [docs/](docs/) для детальной архитектуры и API документации.

## 👥 Команда

Разработано для хакатона "Краснодарский край как цифровой продукт"

## 📄 Лицензия

MIT
