package server

import (
	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/internal/services"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
)

type HttpServer struct {
	Cfg          *config.Config
	App          *http.App
	Log          *logger.Logger
	AdService    *services.AdService
	ClickService *services.ClickService
}

func NewHTTP(cfg *config.Config, app *http.App, log *logger.Logger, adService *services.AdService, clickService *services.ClickService) *HttpServer {
	server := &HttpServer{
		Cfg:          cfg,
		App:          app,
		Log:          log,
		AdService:    adService,
		ClickService: clickService,
	}
	server.RegisterHttpRoutes()
	return server
}
