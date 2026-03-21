package handlers

import (
	"strconv"
	"strings"

	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (a *API) ListLocations(c *fiber.Ctx) error {
	q := c.Query("search", "")
	limitStr := strings.TrimSpace(c.Query("limit"))
	offsetStr := strings.TrimSpace(c.Query("offset"))
	if limitStr == "" && offsetStr == "" {
		list, err := a.Store.ListLocations(c.UserContext(), q)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if list == nil {
			list = []store.Location{}
		}
		return c.JSON(list)
	}
	limit := 100
	offset := 0
	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n < 1 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid limit")
		}
		limit = n
	}
	if offsetStr != "" {
		n, err := strconv.Atoi(offsetStr)
		if err != nil || n < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid offset")
		}
		offset = n
	}
	items, total, err := a.Store.ListLocationsPage(c.UserContext(), q, limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
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

type locationPatchBody struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Category    *string   `json:"category,omitempty"`
	Seasons     *[]string `json:"seasons,omitempty"`
	MediaURLs   *[]string `json:"media_urls,omitempty"`
	Lat         *float64  `json:"lat,omitempty"`
	Lng         *float64  `json:"lng,omitempty"`
}

func (a *API) PatchLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var patch locationPatchBody
	if err := c.BodyParser(&patch); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	loc, err := a.Store.GetLocation(c.UserContext(), id)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if patch.Name != nil {
		loc.Name = *patch.Name
	}
	if patch.Description != nil {
		loc.Description = *patch.Description
	}
	if patch.Category != nil {
		loc.Category = *patch.Category
	}
	if patch.Seasons != nil {
		loc.Seasons = *patch.Seasons
	}
	if patch.MediaURLs != nil {
		loc.MediaURLs = *patch.MediaURLs
	}
	if patch.Lat != nil || patch.Lng != nil {
		if patch.Lat == nil || patch.Lng == nil {
			return fiber.NewError(fiber.StatusBadRequest, "lat and lng must be set together")
		}
		loc.Lat, loc.Lng = patch.Lat, patch.Lng
	}
	if loc.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
	}
	if err := a.Store.UpdateLocation(c.UserContext(), loc); err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	out, err := a.Store.GetLocation(c.UserContext(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(out)
}

func (a *API) DeleteLocation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := a.Store.DeleteLocation(c.UserContext(), id); err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
