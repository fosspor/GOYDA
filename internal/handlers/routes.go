package handlers

import (
	"encoding/json"

	"github.com/fosspor/GOYDA/internal/middleware"
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createRouteBody struct {
	Title   string          `json:"title"`
	Season  string          `json:"season"`
	Payload json.RawMessage `json:"payload"`
}

func (a *API) ListMyRoutes(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	list, err := a.Store.ListRoutesForUser(c.UserContext(), uid)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if list == nil {
		list = []store.Route{}
	}
	return c.JSON(list)
}

func (a *API) CreateRoute(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	var b createRouteBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	r, err := a.Store.CreateRoute(c.UserContext(), uid, b.Title, b.Season, b.Payload)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(r)
}

func (a *API) GetRoute(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	r, err := a.Store.RouteBelongsTo(c.UserContext(), id, uid)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(r)
}

type routePatchBody struct {
	Title   *string          `json:"title,omitempty"`
	Season  *string          `json:"season,omitempty"`
	Payload *json.RawMessage `json:"payload,omitempty"`
}

func (a *API) PatchRoute(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var patch routePatchBody
	if err := c.BodyParser(&patch); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	r, err := a.Store.RouteBelongsTo(c.UserContext(), id, uid)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	title := r.Title
	season := r.Season
	payload := r.Payload
	if patch.Title != nil {
		title = *patch.Title
	}
	if patch.Season != nil {
		season = *patch.Season
	}
	if patch.Payload != nil {
		payload = *patch.Payload
	}
	if err := a.Store.UpdateRoute(c.UserContext(), id, uid, title, season, payload); err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	out, err := a.Store.GetRoute(c.UserContext(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(out)
}

func (a *API) DeleteRoute(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := a.Store.DeleteRouteForUser(c.UserContext(), id, uid); err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
