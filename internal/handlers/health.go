package handlers

import "github.com/gofiber/fiber/v2"

func (a *API) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "goyda-api"})
}
