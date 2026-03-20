#!/bin/bash

# GOYDA Project Setup Script
# Complete setup for Django + React development environment

set -e

echo "🚀 GOYDA Setup - Краснодарский край по-новому"
echo "=============================================="
echo ""

# Check Python version
echo "✓ Проверка Python 3.10..."
python3.10 --version || {
    echo "❌ Python 3.10 не найден. Пожалуйста, установите Python 3.10"
    exit 1
}

# Backend setup
echo ""
echo "📦 Backend Setup"
echo "==============="

cd backend

echo "• Проверка виртуального окружения..."
if [ ! -d "../venv" ]; then
    echo "  Создание виртуального окружения..."
    python3.10 -m venv ../venv
fi

echo "• Активирование виртуального окружения..."
source ../venv/bin/activate

echo "• Установка зависимостей..."
pip install --upgrade pip
pip install -r requirements.txt

echo "• Создание .env файла (если его нет)..."
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo "  ⚠️  Отредактируйте .env с вашими параметрами (Yandex API, DB и т.д.)"
fi

echo "• Применение миграций БД..."
python manage.py migrate --noinput || echo "  ℹ️  Миграции требуют БД. Убедитесь, что PostgreSQL запущен."

echo ""
echo "✓ Backend готов! Запуск: python manage.py runserver"

# Frontend setup
echo ""
echo "🎨 Frontend Setup"
echo "================"

cd ../frontend

echo "• Установка Node.js зависимостей..."
if command -v npm &> /dev/null; then
    npm install
else
    echo "❌ npm не найден. Пожалуйста, установите Node.js https://nodejs.org/"
    exit 1
fi

echo ""
echo "✓ Frontend готов! Запуск: npm start"

# Return to project root
cd ..

echo ""
echo "=============================================="
echo "✅ Проект GOYDA инициализирован!"
echo ""
echo "📋 Следующие шаги:"
echo ""
echo "1️⃣  Отредактируйте backend/.env с вашими параметрами:"
echo "   • Yandex Cloud credentials"
echo "   • PostgreSQL credentials"
echo "   • SECRET_KEY"
echo ""
echo "2️⃣  Запустите PostgreSQL + PostGIS:"
echo "   createdb krasnodar_tourism"
echo ""
echo "3️⃣  Запустите сервер разработки (2 терминала):"
echo ""
echo "   Terminal 1 (Backend):"
echo "   $ source venv/bin/activate"
echo "   $ cd backend"
echo "   $ python manage.py runserver"
echo ""
echo "   Terminal 2 (Frontend):"
echo "   $ cd frontend"
echo "   $ npm start"
echo ""
echo "🌐 Приложение будет доступно на http://localhost:3000"
echo "🔨 Django Admin на http://localhost:8000/admin"
echo ""
echo "🎯 Happy coding! 🚀"
echo "=============================================="
