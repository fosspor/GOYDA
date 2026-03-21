package handlers

import (
	"time"

	"github.com/fosspor/GOYDA/internal/config"
	"github.com/fosspor/GOYDA/internal/integrations/yandexllm"
	"github.com/fosspor/GOYDA/internal/store"
)

type API struct {
	Cfg    config.Config
	Store  *store.Store
	LLM    *yandexllm.Client
	JWTKey []byte
	JWTTTL time.Duration
}

func NewAPI(cfg config.Config, st *store.Store) *API {
	return &API{
		Cfg:    cfg,
		Store:  st,
		LLM:    yandexllm.New(cfg.YandexFolder, cfg.YandexAPIKey),
		JWTKey: []byte(cfg.JWTSecret),
		JWTTTL: 7 * 24 * time.Hour,
	}
}
