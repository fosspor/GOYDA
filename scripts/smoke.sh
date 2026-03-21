#!/usr/bin/env sh
# Смоук: API должен быть запущен (docker compose или go run).
# Пример: ./scripts/smoke.sh http://127.0.0.1:8080
set -e
BASE="${1:-http://127.0.0.1:8080}"
printf 'GET %s/health ... ' "$BASE"
curl -sfS "$BASE/health" | head -c 300
echo
printf 'GET %s/api/locations ... ' "$BASE"
curl -sfS "$BASE/api/locations" | head -c 400
echo
echo "OK"
