#!/bin/bash
# GOYDA Development Environment - Quick Start Script
# Run this to start development servers

echo "🚀 GOYDA Development Environment"
echo "=================================="
echo ""

# Check if virtual environment is activated
if [[ "$VIRTUAL_ENV" == "" ]]; then
    echo "⚠️  Python virtual environment not activated"
    echo "Run: source venv/bin/activate"
    exit 1
fi

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Shutting down..."
    kill $BACKEND_PID 2>/dev/null
    kill $FRONTEND_PID 2>/dev/null
}

trap cleanup EXIT

# Start Backend
echo "📦 Starting Backend (Django)..."
cd backend
python manage.py runserver 8000 &
BACKEND_PID=$!
echo "✓ Backend running on http://localhost:8000"
echo "  Admin: http://localhost:8000/admin"
echo "  API: http://localhost:8000/api"
cd ..

sleep 2

# Start Frontend
echo ""
echo "🎨 Starting Frontend (React)..."
cd frontend
npm start &
FRONTEND_PID=$!
echo "✓ Frontend running on http://localhost:3000"
cd ..

echo ""
echo "=================================="
echo "✅ Development environment ready!"
echo ""
echo "📋 Available URLs:"
echo "  Frontend:       http://localhost:3000"
echo "  Backend API:    http://localhost:8000/api"
echo "  Django Admin:   http://localhost:8000/admin"
echo ""
echo "💡 Tips:"
echo "  - Backend runs hot reload"
echo "  - Frontend has HMR (Hot Module Replacement)"
echo "  - Edit files to see changes instantly"
echo "  - Press Ctrl+C to stop all servers"
echo ""
echo "📚 Documentation:"
echo "  - docs/QUICKSTART.md     - Setup guide"
echo "  - docs/ARCHITECTURE.md   - System design"
echo "  - docs/YANDEX_SETUP.md   - AI integration"
echo "  - docs/FRONTEND_GUIDE.md - React development"
echo ""
echo "=================================="

# Keep script running
wait
