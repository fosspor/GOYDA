.PHONY: spa-dist server-embed

# Сборка фронта и копирование в каталог для go:embed (нужен Go и Node)
spa-dist:
	cd frontend && npm ci && npm run build
	rm -rf internal/spa/dist && mkdir -p internal/spa/dist
	cp -R frontend/dist/. internal/spa/dist/

# Пример: go build -tags embed -o bin/server ./cmd/server
server-embed: spa-dist
	go build -tags embed -o bin/server ./cmd/server
