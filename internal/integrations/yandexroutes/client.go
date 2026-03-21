package yandexroutes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Route struct {
	DistanceM int     `json:"distance_m"`
	DurationS float64 `json:"duration_s"`
	Polyline  []Point `json:"polyline"`
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) Enabled() bool {
	return c.apiKey != ""
}

func (c *Client) Route(ctx context.Context, from, to Point) (Route, error) {
	if !c.Enabled() {
		return Route{}, fmt.Errorf("routing api key is empty")
	}
	u := url.URL{
		Scheme: "https",
		Host:   "api.routing.yandex.net",
		Path:   "/v2/route",
	}
	q := u.Query()
	q.Set("apikey", c.apiKey)
	q.Set("waypoints", fmt.Sprintf("%.6f,%.6f|%.6f,%.6f", from.Lng, from.Lat, to.Lng, to.Lat))
	q.Set("lang", "ru_RU")
	q.Set("results", "1")
	q.Set("overview", "full")
	q.Set("geometries", "geojson")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Route{}, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return Route{}, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return Route{}, fmt.Errorf("routing status: %d", res.StatusCode)
	}
	var payload struct {
		Routes []struct {
			Distance float64 `json:"distance"`
			Duration float64 `json:"duration"`
			Geometry struct {
				Coordinates [][]float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"routes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return Route{}, err
	}
	if len(payload.Routes) == 0 {
		return Route{}, fmt.Errorf("routing empty routes")
	}
	r := payload.Routes[0]
	out := Route{
		DistanceM: int(r.Distance),
		DurationS: r.Duration,
	}
	for _, c := range r.Geometry.Coordinates {
		if len(c) < 2 {
			continue
		}
		out.Polyline = append(out.Polyline, Point{
			Lat: c[1],
			Lng: c[0],
		})
	}
	if len(out.Polyline) == 0 {
		out.Polyline = []Point{from, to}
	}
	return out, nil
}
