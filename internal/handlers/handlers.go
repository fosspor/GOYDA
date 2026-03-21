package handlers

import (
	"time"

	"github.com/fosspor/GOYDA/internal/config"
	"github.com/fosspor/GOYDA/internal/integrations/yandexroutes"
	"github.com/fosspor/GOYDA/internal/integrations/yandexllm"
	"github.com/fosspor/GOYDA/internal/integrations/yandexweather"
	"github.com/fosspor/GOYDA/internal/store"
)

type API struct {
	Cfg    config.Config
	Store  *store.Store
	LLM    *yandexllm.Client
	Weather *yandexweather.Client
	Routes  *yandexroutes.Client
	JWTKey []byte
	JWTTTL time.Duration
}

func NewAPI(cfg config.Config, st *store.Store) *API {
	return &API{
		Cfg:    cfg,
		Store:  st,
		LLM:    yandexllm.New(cfg.YandexFolder, cfg.YandexAPIKey),
		Weather: yandexweather.New(cfg.YandexWeatherKey),
		Routes:  yandexroutes.New(cfg.YandexRoutingKey),
		JWTKey: []byte(cfg.JWTSecret),
		JWTTTL: 7 * 24 * time.Hour,
	}
}
