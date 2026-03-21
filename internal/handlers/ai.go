package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fosspor/GOYDA/internal/middleware"
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
)

type generateRouteBody struct {
	Interests []string `json:"interests"`
	Season    string   `json:"season"`
	Days      int      `json:"days"`
	Notes     string   `json:"notes"`
}

func (a *API) GenerateRoute(c *fiber.Ctx) error {
	var b generateRouteBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if b.Days <= 0 {
		b.Days = 3
	}
	if b.Season == "" {
		b.Season = "summer"
	}
	ctx := c.UserContext()
	uid, _ := middleware.UserID(c)

	locs, err := a.Store.ListLocations(ctx, "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	prompt := fmt.Sprintf(
		"Сезон: %s. Дней: %d. Интересы: %s. Заметки: %s. Доступные локации (кратко JSON): %s",
		b.Season, b.Days, strings.Join(b.Interests, ", "), b.Notes, locationsMiniJSON(locs, 40),
	)

	if a.LLM.Enabled() {
		raw, err := a.LLM.CompletionRaw(ctx, prompt)
		if err == nil {
			route := fiber.Map{"kind": "llm", "text": raw}
			if json.Valid([]byte(raw)) {
				route["json"] = json.RawMessage([]byte(raw))
			}
			return c.JSON(fiber.Map{
				"source":  "yandex",
				"user_id": uid,
				"route":   route,
			})
		}
	}

	mock := buildMockRoute(b, locs)
	return c.JSON(fiber.Map{
		"source":  "mock",
		"user_id": uid,
		"route":   mock,
	})
}

func locationsMiniJSON(locs []store.Location, limit int) string {
	type mini struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	var m []mini
	for i, l := range locs {
		if i >= limit {
			break
		}
		m = append(m, mini{ID: l.ID.String(), Name: l.Name, Category: l.Category})
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func buildMockRoute(b generateRouteBody, locs []store.Location) fiber.Map {
	stops := []fiber.Map{}
	max := b.Days * 2
	if max < 2 {
		max = 2
	}
	for i, l := range locs {
		if i >= max {
			break
		}
		stops = append(stops, fiber.Map{
			"location_id": l.ID.String(),
			"name":        l.Name,
			"reason":      fmt.Sprintf("Подходит под сезон %s и интересы: %s", b.Season, strings.Join(b.Interests, ", ")),
		})
	}
	return fiber.Map{
		"kind":    "mock",
		"title":   fmt.Sprintf("Маршрут на %d дн. (%s)", b.Days, b.Season),
		"summary": "Черновик маршрута (mock). Подключите Yandex LLM для генерации текста.",
		"stops":   stops,
	}
}

func (a *API) AIRecommendations(c *fiber.Ctx) error {
	// Простая выдача: топ локаций по сезону из query
	season := c.Query("season", "summer")
	all, err := a.Store.ListLocations(c.UserContext(), "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var out []store.Location
	for _, l := range all {
		if season == "" {
			out = append(out, l)
			continue
		}
		for _, s := range l.Seasons {
			if strings.EqualFold(s, season) {
				out = append(out, l)
				break
			}
		}
	}
	if len(out) > 20 {
		out = out[:20]
	}
	items := out
	if items == nil {
		items = []store.Location{}
	}
	return c.JSON(fiber.Map{"season": season, "items": items})
}
