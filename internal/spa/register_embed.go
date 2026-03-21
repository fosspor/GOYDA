//go:build embed

package spa

import (
	"embed"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//go:embed all:dist
var root embed.FS

// Register отдаёт собранный SPA (Vite dist) и fallback на index.html для React Router.
func Register(app *fiber.App) {
	sub, err := fs.Sub(root, "dist")
	if err != nil {
		return
	}
	app.Get("/*", func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodGet {
			return fiber.ErrMethodNotAllowed
		}
		p := safePath(c.Path())
		data, err := fs.ReadFile(sub, p)
		if err != nil {
			data, err = fs.ReadFile(sub, "index.html")
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}
			c.Type("text/html; charset=utf-8")
			return c.Send(data)
		}
		ct := mime.TypeByExtension(filepath.Ext(p))
		if ct == "" {
			ct = "application/octet-stream"
		}
		c.Type(ct)
		return c.Send(data)
	})
}

func safePath(raw string) string {
	p := strings.TrimPrefix(raw, "/")
	p = filepath.ToSlash(filepath.Clean(p))
	if p == "." || p == "/" {
		return "index.html"
	}
	if strings.Contains(p, "..") {
		return "index.html"
	}
	return p
}
