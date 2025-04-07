package server

import (
	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
)

type HttpServer struct {
	Cfg *config.Config
	App *http.App
	Log *logger.Logger
}

func NewHTTP(cfg *config.Config, app *http.App, log *logger.Logger) *HttpServer {
	return &HttpServer{
		Cfg: cfg,
		App: app,
		Log: log,
	}
}
	