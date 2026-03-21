#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/frontend"
npm ci
npm run build
rm -rf "$ROOT/internal/spa/dist"
mkdir -p "$ROOT/internal/spa/dist"
cp -R "$ROOT/frontend/dist/." "$ROOT/internal/spa/dist/"
echo "OK: Vite dist → internal/spa/dist (соберите Go с -tags embed)"
