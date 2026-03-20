#!/bin/bash
# Setup script for Krasnodar Tourism backend

echo "📦 Installing Python dependencies..."
pip install -r requirements.txt

echo "🗄️  Creating database migrations..."
python manage.py makemigrations

echo "📝 Running migrations..."
python manage.py migrate

echo "🔐 Creating superuser..."
python manage.py createsuperuser

echo "📦 Collecting static files..."
python manage.py collectstatic --noinput

echo "✅ Backend setup complete!"
echo ""
echo "To start the development server:"
echo "  python manage.py runserver"
