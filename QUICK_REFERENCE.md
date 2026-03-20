# 🚀 GOYDA Quick Reference

## Start Development

```bash
# Single command to start both servers
cd /Users/fosspor/GOYDA
bash start-dev.sh

# Or manually...
# Terminal 1:
cd /Users/fosspor/GOYDA/backend
source ../venv/bin/activate
python manage.py runserver

# Terminal 2:
cd /Users/fosspor/GOYDA/frontend
npm start
```

## URLs
- **Frontend**: http://localhost:3000
- **Admin**: http://localhost:8000/admin (admin/admin123456)
- **API**: http://localhost:8000/api/

## Django Commands

```bash
cd /Users/fosspor/GOYDA/backend
source ../venv/bin/activate

python manage.py makemigrations        # Create new migrations
python manage.py migrate               # Apply migrations
python manage.py shell                 # Interactive shell
python manage.py createsuperuser       # Create admin user
python manage.py runserver            # Start dev server
python manage.py test locations        # Run app tests
```

## NPM Commands

```bash
cd /Users/fosspor/GOYDA/frontend

npm start                  # Start dev server
npm run build             # Build for production
npm test                  # Run tests
npm eject                 # Expose config (careful!)
npm install <package>     # Add new package
```

## File Locations

| File | Path |
|------|------|
| Django Settings | `/backend/config/settings.py` |
| Environment | `/backend/.env` |
| Database | PostgreSQL `krasnodar_tourism` |
| Models | `/backend/{app}/models.py` |
| Admin | `/backend/{app}/admin.py` |
| React App | `/frontend/src/App.js` |
| Tailwind Config | `/frontend/tailwind.config.js` |

## Database Info

```
Database: krasnodar_tourism
User: fosspor
Host: localhost
Port: 5432
```

Connect with psql:
```bash
psql -d krasnodar_tourism -U fosspor
```

## Project Database Models

### Location ⭐
- name, description, short_description
- category_id (ForeignKey)
- latitude, longitude
- address, region
- image, phone, email, website
- rating, reviews_count
- best_season (ManyToMany)

### Route 🗺️
- name, description, slug
- user_id (ForeignKey)
- season_id (ForeignKey)
- duration_days
- is_public, featured
- locations (ManyToMany via RouteLocation)

### User 👤
- Custom Django User
- email, phone, date_joined
- has UserPreference

### Category 📂
- name, slug, description
- icon, color

### Season 🌤️
- name (spring/summer/autumn/winter)
- description
- start_month, end_month

## Common Development Tasks

### Add New Model
1. Edit `/backend/{app}/models.py`
2. Register in `/backend/{app}/admin.py`
3. Create migration: `python manage.py makemigrations`
4. Apply: `python manage.py migrate`

### Create API Endpoint
1. `views.py`: Create ViewSet or APIView
2. `serializers.py`: Create serializer
3. `urls.py`: Add URL pattern
4. Frontend: Update `services/api.js`

### Add React Component
1. Create file in `frontend/src/components/`
2. Import in needed pages
3. Use Tailwind classes for styling
4. Connect to Zustand store if needed

## Environment Variables

```
SECRET_KEY=django-insecure-...
DEBUG=True
DB_NAME=krasnodar_tourism
DB_USER=fosspor
DB_HOST=localhost
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8000
YANDEX_API_KEY=...  (for AI features)
```

## Useful Links

- Django Docs: https://docs.djangoproject.com/
- DRF Docs: https://www.django-rest-framework.org/
- React Docs: https://react.dev/
- Tailwind CSS: https://tailwindcss.com/
- Zustand: https://github.com/pmndrs/zustand

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| Ctrl+C | Stop development server |
| Cmd+Shift+P | VS Code command palette |
| Cmd+K Cmd+0 | Fold all in editor |
| Cmd+/ | Toggle comment |

## Git Commands

```bash
cd /Users/fosspor/GOYDA

git status                     # Check status
git add <file>                 # Stage changes
git commit -m "message"        # Commit
git push                       # Push to remote
git pull                       # Get latest

# Useful for development:
git diff                       # View changes
git log --oneline             # View commits
git checkout -b feature/name   # Create new branch
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Port already in use | `lsof -i :8000` then `kill <PID>` |
| Database connection error | Check PostgreSQL is running: `brew services list` |
| Module not found | Reinstall deps: `pip install -r requirements.txt` |
| React app won't start | Check Node version, reinstall deps: `npm ci` |
| Migration conflicts | Review migration files, may need manual fix |

## Performance Tips

- Enable Django debug toolbar in development
- Use React DevTools for component inspection
- Monitor network requests in browser DevTools
- Check database query count with Django shell_plus

## Security Reminders

⚠️ **Before Production**:
- [ ] Change SECRET_KEY
- [ ] Set DEBUG=False
- [ ] Update ALLOWED_HOSTS
- [ ] Set secure CORS origins
- [ ] Use environment variables for secrets
- [ ] Enable HTTPS
- [ ] Set SECURE_HSTS_SECONDS
- [ ] Enable CSRF protection
- [ ] Review user permissions

---

**Keep this file handy during development!** 📌
