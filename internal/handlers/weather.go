package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/fosspor/GOYDA/internal/integrations/yandexroutes"
	"github.com/fosspor/GOYDA/internal/integrations/yandexweather"
	"github.com/fosspor/GOYDA/internal/middleware"
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (a *API) WeatherPoint(c *fiber.Ctx) error {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lat")
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lng")
	}
	w, err := a.Weather.Point(c.UserContext(), lat, lng)
	if err == nil {
		return c.JSON(w)
	}
	return c.JSON(mockWeather(lat, lng))
}

type weatherAwareBody struct {
	FromLocationID string `json:"from_location_id"`
	ToLocationID   string `json:"to_location_id"`
	Date           string `json:"date"`
	AvoidRain      bool   `json:"avoid_rain"`
	MaxWindMS      float64 `json:"max_wind_ms"`
}

func (a *API) WeatherAwareRoute(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	_ = uid

	var b weatherAwareBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	from, to, err := a.resolveRouteLocations(c, b.FromLocationID, b.ToLocationID)
	if err != nil {
		return err
	}
	fromPoint := toPoint(from)
	toPointVal := toPoint(to)

	route, routeSource := a.routeWithFallback(c, fromPoint, toPointVal)
	wxFrom, wxFromSource := a.weatherWithFallback(c, fromPoint.Lat, fromPoint.Lng)
	wxTo, wxToSource := a.weatherWithFallback(c, toPointVal.Lat, toPointVal.Lng)

	score := weatherRiskScore(wxFrom, wxTo, b.AvoidRain, b.MaxWindMS)
	response := fiber.Map{
		"source": fiber.Map{
			"routing": routeSource,
			"weather_from": wxFromSource,
			"weather_to": wxToSource,
		},
		"date":    normalizeDate(b.Date),
		"from":    from,
		"to":      to,
		"route":   route,
		"weather": fiber.Map{"from": wxFrom, "to": wxTo},
		"score":   score,
		"reasoning": weatherReasoning(score),
	}
	payload, _ := json.Marshal(response)
	created, err := a.Store.CreateRoute(c.UserContext(), uid, fmt.Sprintf("Weather route: %s -> %s", from.Name, to.Name), "auto", payload)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	response["saved_route_id"] = created.ID.String()
	return c.JSON(response)
}

func (a *API) resolveRouteLocations(c *fiber.Ctx, fromID, toID string) (store.Location, store.Location, error) {
	ctx := c.UserContext()
	if fromID != "" && toID != "" {
		fid, err := uuid.Parse(fromID)
		if err != nil {
			return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusBadRequest, "invalid from_location_id")
		}
		tid, err := uuid.Parse(toID)
		if err != nil {
			return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusBadRequest, "invalid to_location_id")
		}
		from, err := a.Store.GetLocation(ctx, fid)
		if err != nil {
			return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusNotFound, "from location not found")
		}
		to, err := a.Store.GetLocation(ctx, tid)
		if err != nil {
			return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusNotFound, "to location not found")
		}
		return from, to, nil
	}
	all, err := a.Store.ListLocations(ctx, "")
	if err != nil {
		return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	geo := make([]store.Location, 0, len(all))
	for _, l := range all {
		if l.Lat != nil && l.Lng != nil {
			geo = append(geo, l)
		}
	}
	if len(geo) < 2 {
		return store.Location{}, store.Location{}, fiber.NewError(fiber.StatusBadRequest, "need at least two geocoded locations")
	}
	return geo[0], geo[1], nil
}

func toPoint(l store.Location) yandexroutes.Point {
	return yandexroutes.Point{Lat: deref(l.Lat), Lng: deref(l.Lng)}
}

func deref(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func (a *API) routeWithFallback(c *fiber.Ctx, from, to yandexroutes.Point) (fiber.Map, string) {
	r, err := a.Routes.Route(c.UserContext(), from, to)
	if err == nil {
		return fiber.Map{
			"distance_m": r.DistanceM,
			"duration_s": r.DurationS,
			"polyline":   r.Polyline,
		}, "yandex"
	}
	distance := haversineKm(from.Lat, from.Lng, to.Lat, to.Lng) * 1000
	duration := distance / 13.8
	return fiber.Map{
		"distance_m": int(distance),
		"duration_s": duration,
		"polyline": []yandexroutes.Point{
			from, to,
		},
	}, "mock"
}

func (a *API) weatherWithFallback(c *fiber.Ctx, lat, lng float64) (yandexweather.PointWeather, string) {
	w, err := a.Weather.Point(c.UserContext(), lat, lng)
	if err == nil {
		return w, "yandex"
	}
	m := mockWeather(lat, lng)
	return m, "mock"
}

func mockWeather(lat, lng float64) yandexweather.PointWeather {
	return yandexweather.PointWeather{
		Source:      "mock",
		Lat:         lat,
		Lng:         lng,
		TempC:       22,
		Condition:   "cloudy",
		WindSpeedMS: 4,
	}
}

func weatherRiskScore(from, to yandexweather.PointWeather, avoidRain bool, maxWind float64) float64 {
	avgTemp := (from.TempC + to.TempC) / 2
	avgWind := (from.WindSpeedMS + to.WindSpeedMS) / 2
	score := 0.0
	if avgTemp < 5 || avgTemp > 30 {
		score += 25
	}
	if maxWind > 0 && avgWind > maxWind {
		score += 30
	}
	if avoidRain {
		if isRainy(from.Condition) {
			score += 20
		}
		if isRainy(to.Condition) {
			score += 20
		}
	}
	return math.Round(score*10) / 10
}

func isRainy(cond string) bool {
	switch cond {
	case "rain", "light-rain", "showers", "wet-snow":
		return true
	default:
		return false
	}
}

func weatherReasoning(score float64) string {
	if score < 20 {
		return "Погодные условия хорошие, маршрут комфортный."
	}
	if score < 50 {
		return "Условия средние, возможны участки с ветром или осадками."
	}
	return "Высокий погодный риск: рекомендуем выбрать другое время или маршрут."
}

func normalizeDate(v string) string {
	if v == "" {
		return time.Now().Format("2006-01-02")
	}
	return v
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371.0
	toRad := func(v float64) float64 { return v * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	return r * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
