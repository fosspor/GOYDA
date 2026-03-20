# GOYDA Architecture Documentation

## System Overview

GOYDA — это веб-приложение для создания персональных маршрутов по Краснодарскому краю с использованием ИИ.

```
┌─────────────────────────────────────────────────────────┐
│                    GOYDA System                         │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────────────────┐          ┌──────────────────┐  │
│  │   React Frontend   │ ◄───────► │  Django Backend  │  │
│  │   (Port 3000)      │   REST    │  (Port 8000)     │  │
│  └────────────────────┘   API     └──────────────────┘  │
│           │                               │              │
│           │                               ├────────────┐ │
│           │                               │            │ │
│           │                        ┌──────▼────────┐   │ │
│           └─────────────────────►  │ PostgreSQL +  │   │ │
│                                    │   PostGIS     │   │ │
│                                    │   (Geo DB)    │   │ │
│                                    └───────────────┘   │ │
│                                                        │ │
│                                         ┌──────────────▼─┤
│                                         │  Yandex LLM    │
│                                         │  (AI Routes)   │
│                                         └────────────────┤
│                                                          │
└─────────────────────────────────────────────────────────┘
```

## Backend Architecture (Django)

### Apps Structure

#### 1. **locations** — Локации и категории
```
locations/
├── models.py
│   ├── Category              # Категории (винодельня, ресторан и т.д.)
│   ├── Season                # Сезоны (весна, лето и т.д.)
│   ├── Location              # Основная модель локации
│   ├── LocationReview        # Отзывы пользователей
│   └── LocationPhoto         # Фото пользователей
├── views.py                  # API endpoints
├── serializers.py            # DRF сериализаторы
└── urls.py
```

**API:**
- `GET /api/locations/` — все локации
- `GET /api/locations/{id}/` — детали
- `GET /api/locations/?category=wine&season=summer` — фильтры

#### 2. **routes** — Маршруты
```
routes/
├── models.py
│   ├── Route               # Персональный маршрут
│   ├── RouteLocation       # Порядок посещений
│   └── RouteInterest       # Интересы пользователя
├── views.py
├── serializers.py
└── urls.py
```

**API:**
- `POST /api/routes/` — создать маршрут
- `GET /api/routes/` — мои маршруты
- `GET /api/routes/{id}/` — детали маршрута

#### 3. **users** — Пользователи
```
users/
├── models.py
│   ├── User                # Кастомный User с профилем
│   └── UserPreference      # Предпочтения для рекомендаций
├── views.py
└── urls.py
```

#### 4. **ai_recommendations** — ИИ рекомендации
```
ai_recommendations/
├── models.py
│   ├── RouteRecommendation # Рекомендованный маршрут
│   └── AIPromptLog         # Логирование AI запросов
├── services.py
│   └── YandexAIService     # Интеграция с Yandex LLM
└── urls.py
```

### Database Schema (Key Models)

```sql
-- Locations
CREATE TABLE locations_location (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    location GEOMETRY(POINT, 4326),  -- PostGIS point
    category_id INT REFERENCES locations_category,
    rating DECIMAL(3,2),
    image VARCHAR(100),
    video_url VARCHAR(200),
    created_at TIMESTAMP
);

-- Routes
CREATE TABLE routes_route (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    user_id INT REFERENCES users_user,
    season VARCHAR(10),
    duration_days INT,
    created_at TIMESTAMP
);

CREATE TABLE routes_routelocation (
    id SERIAL PRIMARY KEY,
    route_id INT REFERENCES routes_route,
    location_id INT REFERENCES locations_location,
    order INT,
    start_time TIME
);

-- Users
CREATE TABLE users_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(150) UNIQUE,
    email VARCHAR(254),
    avatar VARCHAR(100),
    is_business BOOLEAN,
    created_at TIMESTAMP
);

CREATE TABLE users_userpreference (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users_user,
    min_budget INT,
    max_budget INT,
    physical_activity_level INT
);
```

### AI Integration Flow

```
User Request
    ↓
[RouteRecommendationViewSet]
    ↓
Build Prompt with:
  - User interests
  - Budget
  - Season
  - Group info
    ↓
[YandexAIService.generate_route_prompt()]
    ↓
Call Yandex LLM API
    ↓
Parse Response
    ↓
Create RouteRecommendation object
    ↓
Log in AIPromptLog
    ↓
Return to Frontend
```

## Frontend Architecture (React)

### Component Structure

```
src/
├── components/
│   ├── Navigation.js          # Главное меню
│   ├── LocationCard.js        # Карточка локации
│   ├── RouteBuilder.js        # Конструктор маршрутов
│   ├── Map.js                 # Интерактивная карта
│   ├── VirtualTour.js         # 360° тур
│   └── ...
├── pages/
│   ├── Home.js                # Главная страница
│   ├── Discover.js            # Поиск и фильтры
│   ├── SuggestedRoutes.js      # Мои маршруты
│   ├── RouteDetail.js         # Детали маршрута
│   └── Profile.js             # Профиль пользователя
├── services/
│   ├── api.js                 # Axios API client
│   └── auth.js                # Аутентификация
├── store/
│   └── index.js               # Zustand stores
├── hooks/
│   └── ...                    # Кастомные hooks
└── App.js                     # Root component
```

### State Management (Zustand)

```javascript
// Global stores
useUserStore     // Текущий пользователь
useRouteStore    // Маршруты
useLocationStore // Локации и фильтры
```

### Key Pages

#### 1. Home (`/`)
- Hero section с call-to-action
- Фишки сервиса
- Ссылки на discover и routes

#### 2. Discover (`/discover`)
- Список локаций
- Фильтры (категория, сезон, теги)
- Поиск по названию
- Карта с маркерами

#### 3. Suggested Routes (`/suggested-routes`)
- Список AI-рекомендованных маршрутов
- Сохраненные маршруты
- Фильтры по сезону/интересам

#### 4. Route Detail (`/route/:id`)
- Карта с порядком посещений
- Информация о каждой локации
- Галерея фото
- Логистика (время, расстояние)

## API Documentation

### Authentication
```
POST /api/auth/token/
{
  "username": "user@example.com",
  "password": "password123"
}

Response:
{
  "access": "eyJhbGc...",
  "refresh": "eyJhbGc..."
}
```

### Locations
```
GET /api/locations/?category=wine&season=summer&search=маршам

{
  "count": 42,
  "next": "...",
  "results": [
    {
      "id": 1,
      "name": "Усадьба Абрау",
      "category": {
        "id": 1,
        "name": "Винодельня",
        "slug": "winery",
        "icon": "wine-glass"
      },
      "address": "р-н Абрау, ул. Апрельская",
      "latitude": 44.2,
      "longitude": 37.5,
      "image": "https://...",
      "rating": 4.5,
      "tags": ["organic", "tours", "tasting"],
      "best_season": ["spring", "autumn"],
      "amenities": ["parking", "wifi", "cafe"]
    }
  ]
}
```

### Routes
```
GET /api/routes/

{
  "count": 5,
  "results": [
    {
      "id": 1,
      "name": "Винная осень - 3 дня",
      "slug": "wine-autumn-3days",
      "description": "...",
      "season": "autumn",
      "duration_days": 3,
      "locations": [
        {
          "id": 1,
          "name": "Абрау",
          "order": 1,
          "duration_hours": 3
        },
        ...
      ],
      "created_at": "2024-03-20T10:00:00Z"
    }
  ]
}
```

### AI Generate Route
```
POST /api/ai/generate-route/

{
  "interests": ["wine", "gastronomy"],
  "budget": 50000,
  "duration": 3,
  "season": "autumn",
  "group_info": {
    "size": 2,
    "ages": [30, 35],
    "with_children": false
  }
}

{
  "score": 0.92,
  "reasoning": "Perfect combination for autumn wine lovers...",
  "route": {
    "locations": [1, 5, 12, 8],
    "estimated_duration": 18,
    "total_distance_km": 85.5
  }
}
```

## Deployment

### Prerequisites
- Python 3.10
- Node.js 16+
- PostgreSQL 12+ with PostGIS
- Yandex Cloud credentials

### Environment Setup

```bash
# Backend
export DJANGO_SETTINGS_MODULE=config.settings
export SECRET_KEY=your-secret-key
export YANDEX_FOLDER_ID=your-folder-id
export YANDEX_API_KEY=your-api-key
```

### Running Locally

```bash
# Terminal 1: Backend
source venv/bin/activate
cd backend
python manage.py runserver 0.0.0.0:8000

# Terminal 2: Frontend
cd frontend
npm start
```

### Production Deployment

#### Backend (e.g., Heroku)
```dockerfile
FROM python:3.10-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY backend .
CMD ["gunicorn", "config.wsgi:application", "--bind", "0.0.0.0:8000"]
```

#### Frontend (e.g., Vercel/Netlify)
```bash
npm run build
# Deploy dist/ directory
```

## Performance Considerations

1. **Database Indexing**
   - Index on `location.location` for spatial queries
   - Index on `route.user_id`, `location.category_id`

2. **Caching**
   - Cache popular locations and routes
   - Cache Yandex AI responses

3. **API Optimization**
   - Pagination (20 items per page)
   - Select/Prefetch related queries
   - Compression (gzip)

4. **Frontend**
   - Code splitting for routes
   - Image lazy loading
   - Service workers for offline support

## Security

- JWT for authentication
- CORS enabled only for trusted origins
- Environment variables for secrets
- HTTPS in production
- CSRF protection
- Rate limiting on AI endpoints

## Testing

```bash
# Backend
pytest backend/

# Frontend
npm test
```

## Monitoring & Logging

- Django logging to stdout/files
- Frontend error tracking (Sentry)
- Database query logging
- API response time monitoring
