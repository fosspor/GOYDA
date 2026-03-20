# GOYDA Project Summary

## 📋 Project Status: ✅ MVP Ready

A comprehensive web application for discovering and creating personalized travel routes in Krasnodar Krai using AI.

---

## 🎯 Quick Navigation

- **📖 Documentation**: See [`docs/`](docs/) folder
  - [Quick Start](docs/QUICKSTART.md) - 5 min setup
  - [Architecture](docs/ARCHITECTURE.md) - System design
  - [Frontend Guide](docs/FRONTEND_GUIDE.md) - React development
  - [Yandex Setup](docs/YANDEX_SETUP.md) - AI integration
  - [Hackathon Requirements](docs/HACKATHON_REQUIREMENTS.md) - MVP checklist

- **💻 Code**: Organized in `backend/` and `frontend/` folders

---

## 🚀 Quick Start (5 minutes)

### 1. Activate Python Environment
```bash
source venv/bin/activate
```

### 2. Setup Database
```bash
# Install & start PostgreSQL
brew install postgresql@15
brew services start postgresql@15

# Create database with PostGIS
createdb krasnodar_tourism
psql krasnodar_tourism -c "CREATE EXTENSION postgis;"
```

### 3. Configure Backend
```bash
cd backend
cp .env.example .env
# Edit .env with your settings (especially Yandex API key)
pip install -r requirements.txt
python manage.py migrate
python manage.py createsuperuser
```

### 4. Configure Frontend
```bash
cd ../frontend
npm install
```

### 5. Start Development
```bash
# Terminal 1: Backend
cd backend && python manage.py runserver

# Terminal 2: Frontend  
cd frontend && npm start
```

Open http://localhost:3000 🎉

---

## 🏗️ Architecture Overview

### Backend (Django + DRF)
```
locations/    - Places, attractions, restaurants
routes/       - Travel routes, user interests
users/        - Custom User model, preferences
ai_recommendations/ - Yandex LLM integration
```

### Frontend (React 18)
```
pages/        - Home, Discover, Routes, Profile
components/   - Reusable UI components
services/     - API client with Axios
store/        - Zustand state management
```

### Database (PostgreSQL + PostGIS)
- Geospatial queries for maps
- Optimized indices for performance
- Full relational schema with migrations

### AI (Yandex LLM)
- Route generation based on preferences
- Location descriptions
- Personalized recommendations

---

## ✨ Key Features

### 1. Location Discovery
- Browse and filter attractions
- Geospatial maps with Leaflet
- Ratings, reviews, seasonal info
- Virtual 360° tours support

### 2. AI Route Generator
- Personalized routes based on:
  - User interests (wine, food, nature, culture, etc.)
  - Budget and duration
  - Seasonal availability
  - Group type (family, seniors, freelancers)

### 3. Route Builder
- Interactive route mapping
- Logistics details (transport, lodging, dining)
- Cost estimation
- Time planning

### 4. Personalization
- User profiles with preferences
- Travel history
- Saved routes
- Preference learning

### 5. Responsive Design
- Mobile-first approach
- Works on all devices
- 60+ FPS animations
- Offline support ready

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | React 18, React Router, Tailwind CSS |
| **Backend** | Django 4.2, DRF, PostgreSQL, PostGIS |
| **State** | Zustand (frontend), Django ORM (backend) |
| **Maps** | Leaflet + React-Leaflet |
| **AI/LLM** | Yandex LLM API |
| **Auth** | JWT (SimpleJWT) |
| **Admin** | Django Admin (customized) |
| **Styling** | Tailwind CSS utility-first |

---

## 📊 Project Statistics

```
Files:
  - Backend: 50+ Python files
  - Frontend: 30+ React components
  - Config: 10+ configuration files
  - Docs: 5+ documentation files

Code:
  - Models: 12+ Django models
  - APIs: 10+ endpoints
  - Components: 15+ React components
  - Pages: 5+ full-page views

Database:
  - Tables: locations, routes, users, reviews, etc.
  - Spatial: PostGIS for geospatial queries
  - Indexes: Optimized for performance

Development:
  - Python 3.10 + venv
  - Node.js 16+ with npm
  - PostgreSQL 12+ with PostGIS
```

---

## 🔧 Development Commands

### Backend
```bash
# Migrations
python manage.py makemigrations
python manage.py migrate

# Create superuser
python manage.py createsuperuser

# Run tests
pytest

# Code quality
black . && isort . && flake8 .
```

### Frontend
```bash
# Development
npm start

# Build production
npm run build

# Tests
npm test

# Code quality
npm run lint
```

---

## 🚢 Deployment

### Backend (Django)
```bash
# Production server
gunicorn config.wsgi:application --bind 0.0.0.0:8000

# Environment variables required:
# - SECRET_KEY
# - DATABASE_URL
# - YANDEX_API_KEY
# - ALLOWED_HOSTS
```

### Frontend (React)
```bash
# Build static files
npm run build

# Deploy 'build/' folder to:
# - Vercel (recommended)
# - Netlify
# - AWS S3 / CloudFront
# - Any static hosting
```

---

## 🔐 Security

- ✅ CORS configured for dev/production
- ✅ JWT authentication for APIs
- ✅ Environment variables for secrets
- ✅ Django security middleware enabled
- ✅ HTTPS-ready (configure in production)
- ✅ Rate limiting ready
- ✅ CSRF protection enabled

---

## 🐛 Testing

### Backend Tests
```bash
pytest backend/
pytest backend/locations/tests/  # Specific app
pytest backend/ -v --cov         # With coverage
```

### Frontend Tests
```bash
npm test                    # All tests
npm test -- Discover       # Specific component
npm test -- --coverage     # With coverage
```

---

## 📱 Mobile Support

- Responsive design (mobile-first)
- Touch-optimized UI
- Fast load times
- Service worker ready
- PWA capable

---

## 🎯 Next Phase Features

- [ ] User authentication/registration UI
- [ ] Payment integration
- [ ] Real-time notifications
- [ ] Social sharing
- [ ] Mobile app (React Native)
- [ ] Advanced search/filters
- [ ] Route timeline builder
- [ ] 3D map visualization
- [ ] Offline support (Service Workers)
- [ ] Dark mode

---

## 📚 Important Files

### Configuration
- `backend/.env.example` - Environment template
- `backend/config/settings.py` - Django settings
- `backend/config/urls.py` - URL routing
- `frontend/tailwind.config.js` - Tailwind config

### Documentation
- `README.md` - This file
- `docs/QUICKSTART.md` - Setup guide
- `docs/ARCHITECTURE.md` - System design
- `docs/YANDEX_SETUP.md` - AI setup

### Key Models
- `backend/locations/models.py` - Location schema
- `backend/routes/models.py` - Route schema
- `backend/users/models.py` - User schema
- `backend/ai_recommendations/models.py` - AI models

---

## 💡 Tips for Development

1. **Use Django Shell**
   ```bash
   python manage.py shell
   from locations.models import Location
   Location.objects.all()
   ```

2. **Check API**
   ```bash
   curl http://localhost:8000/api/locations/
   ```

3. **Admin Panel**
   - Visit http://localhost:8000/admin
   - Add locations, categories, etc.
   - Test your changes in real-time

4. **Frontend Hot Reload**
   - Edit React files to see changes instantly
   - No need to restart `npm start`

5. **Database Queries**
   ```python
   # Optimized geospatial query
   from django.contrib.gis.db.models.functions import Distance
   Location.objects.annotate(
       distance=Distance('location', point)
   ).order_by('distance')
   ```

---

## 🆘 Troubleshooting

### "psycopg2 not found"
```bash
pip install psycopg2-binary
```

### "PostGIS extension not found"
```bash
psql krasnodar_tourism -c "CREATE EXTENSION postgis;"
```

### "Port already in use"
```bash
# Find and kill process
lsof -i :8000
kill -9 <PID>
```

### "Module not found in React"
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

---

## 📖 Learning Resources

- [Django Documentation](https://docs.djangoproject.com/)
- [Django REST Framework](https://www.django-rest-framework.org/)
- [React Documentation](https://react.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [PostgreSQL PostGIS](https://postgis.net/)
- [Yandex Cloud](https://cloud.yandex.com/docs)

---

## 👥 Team Notes

- **Hackathon Challenge**: 48-hour MVP development
- **Target**: Personalized travel routes for Krasnodar Krai
- **MVP Status**: Ready for development
- **Tech Lead**: Python 3.10 + React 18 + PostgreSQL

---

## 📝 Git Workflow

```bash
# Create feature branch
git checkout -b feature/location-search

# Make changes
git add .
git commit -m "Add location search"

# Push and create PR
git push origin feature/location-search
```

---

## ✅ Pre-Deployment Checklist

- [ ] All migrations applied
- [ ] Yandex API configured
- [ ] Frontend builds without errors
- [ ] Backend tests passing
- [ ] .env file created and secret keys set
- [ ] CORS configured for production domain
- [ ] DEBUG = False in production
- [ ] Static files collected
- [ ] Database backed up

---

## 🎉 Ready to Build!

You now have:
- ✅ Complete project structure
- ✅ Database schema
- ✅ API endpoints ready to implement
- ✅ Frontend skeleton with routing
- ✅ AI integration scaffolding
- ✅ Comprehensive documentation
- ✅ Admin panel configured

**Time to start developing the actual features! 🚀**

---

**Questions?** Check the docs folder or read inline code comments.

**Happy Coding! 💻**
