package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPAddr      string
	DatabaseURL   string
	JWTSecret     string
	CORSOrigins   []string
	YandexFolder  string
	YandexAPIKey  string
}

func Load() (Config, error) {
	c := Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:  strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:    strings.TrimSpace(os.Getenv("JWT_SECRET")),
		YandexFolder: strings.TrimSpace(os.Getenv("YANDEX_FOLDER_ID")),
		YandexAPIKey: strings.TrimSpace(os.Getenv("YANDEX_API_KEY")),
	}
	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if len(c.JWTSecret) < 8 {
		return Config{}, fmt.Errorf("JWT_SECRET is required (min 8 characters)")
	}
	raw := strings.TrimSpace(os.Getenv("CORS_ORIGINS"))
	if raw == "" {
		c.CORSOrigins = []string{"http://localhost:3000"}
	} else {
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				c.CORSOrigins = append(c.CORSOrigins, p)
			}
		}
	}
	return c, nil
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}
