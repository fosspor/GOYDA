#!/bin/bash
# Start both Backend and Frontend development servers

set -e

PROJECT_DIR="/Users/fosspor/GOYDA"
VENV_PYTHON="$PROJECT_DIR/venv/bin/python"

echo "🚀 GOYDA - Starting Development Servers"
echo "========================================"

# Check if servers are already running
if lsof -Pi :8000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  Backend server already running on port 8000"
else
    echo "📦 Starting Backend Server..."
    cd "$PROJECT_DIR/backend"
    $VENV_PYTHON manage.py runserver 0.0.0.0:8000 &
    BACKEND_PID=$!
    echo "✓ Backend Server started (PID: $BACKEND_PID)"
    sleep 2
fi

if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  Frontend server already running on port 3000"
else
    echo "📦 Starting Frontend Server..."
    cd "$PROJECT_DIR/frontend"
    npm start &
    FRONTEND_PID=$!
    echo "✓ Frontend Server started (PID: $FRONTEND_PID)"
    sleep 3
fi

echo ""
echo "🎉 Development servers started successfully!"
echo ""
echo "📍 URLs:"
echo "  Frontend:     http://localhost:3000"
echo "  Backend API:  http://localhost:8000"
echo "  Admin Panel:  http://localhost:8000/admin"
echo ""
echo "👤 Admin Credentials:"
echo "  Username: admin"
echo "  Password: admin123456"
echo ""
echo "Press Ctrl+C to stop the servers"
echo ""

# Keep the script running
wait
