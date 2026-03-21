package yandexweather

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

type PointWeather struct {
	Source      string  `json:"source"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	TempC       float64 `json:"temp_c"`
	Condition   string  `json:"condition"`
	WindSpeedMS float64 `json:"wind_speed_ms"`
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Enabled() bool {
	return c.apiKey != ""
}

func (c *Client) Point(ctx context.Context, lat, lng float64) (PointWeather, error) {
	if !c.Enabled() {
		return PointWeather{}, fmt.Errorf("weather api key is empty")
	}
	u := url.URL{
		Scheme: "https",
		Host:   "api.weather.yandex.ru",
		Path:   "/v2/forecast",
	}
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%.6f", lat))
	q.Set("lon", fmt.Sprintf("%.6f", lng))
	q.Set("limit", "1")
	q.Set("hours", "false")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return PointWeather{}, err
	}
	req.Header.Set("X-Yandex-Weather-Key", c.apiKey)

	res, err := c.http.Do(req)
	if err != nil {
		return PointWeather{}, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return PointWeather{}, fmt.Errorf("weather status: %d", res.StatusCode)
	}
	var payload struct {
		Fact struct {
			Temp      float64 `json:"temp"`
			Condition string  `json:"condition"`
			WindSpeed float64 `json:"wind_speed"`
		} `json:"fact"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return PointWeather{}, err
	}
	return PointWeather{
		Source:      "yandex",
		Lat:         lat,
		Lng:         lng,
		TempC:       payload.Fact.Temp,
		Condition:   payload.Fact.Condition,
		WindSpeedMS: payload.Fact.WindSpeed,
	}, nil
}
