package v1

import (
	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
	"gorm.io/gorm"
)

type HttpServer struct {
	Cfg *config.Config
	App *http.App
	DB  *gorm.DB
	Log *logger.Logger
}

func NewHTTP(cfg *config.Config, app *http.App, db *gorm.DB, log *logger.Logger) *HttpServer {
	return &HttpServer{
		Cfg: cfg,
		App: app,
		DB:  db,
		Log: log,
	}
}
