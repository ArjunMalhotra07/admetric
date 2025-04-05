package server

import (
	"github.com/ArjunMalhotra/config"
	v1 "github.com/ArjunMalhotra/internal/server/v1"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
	"gorm.io/gorm"
)

type Server struct {
	App *http.App
	Web *v1.HttpServer
}

func NewServer(log *logger.Logger, db *gorm.DB, cfg *config.Config) *Server {
	app := http.NewApp(log)
	newHttp := v1.NewHTTP(cfg, app, db, log)

	return &Server{App: app, Web: newHttp}
}
