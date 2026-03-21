//go:build !embed

package spa

import "github.com/gofiber/fiber/v2"

// Register — заглушка: SPA не вшито. Для UI запускайте Vite (`npm run dev` в frontend)
// или соберите бинарь с тегом embed (см. Dockerfile, scripts/sync-spa-dist.sh).
func Register(_ *fiber.App) {}
