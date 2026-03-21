package handlers

import (
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (a *API) ListLocations(c *fiber.Ctx) error {
	q := c.Query("search", "")
	list, err := a.Store.ListLocations(c.UserContext(), q)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(list)
}

func (a *API) GetLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	loc, err := a.Store.GetLocation(c.UserContext(), id)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(loc)
}

func (a *API) CreateLocation(c *fiber.Ctx) error {
	var loc store.Location
	if err := c.BodyParser(&loc); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if loc.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name required")
	}
	out, err := a.Store.CreateLocation(c.UserContext(), loc)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(out)
}
