package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/internal/server"
	"github.com/ArjunMalhotra/pkg/db"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// config
	cfg := config.NewConfig()
	// logger
	log, _ := logger.NewLogger(cfg)
	app := http.NewApp(log)
	// mysql db
	db, err := db.NewMysqDB(cfg)
	if err != nil {
		log.Logger.Errorf("Failed to load db object")
		os.Exit(1)
	}
	if err := db.Migrate(); err != nil {
		log.Logger.Fatalf("Error trying to migrate: %v", err)
	}
	// http API server based on fiber
	server := server.NewHTTP(cfg, app, log)
	// Register all APP APIs
	server.RegisterRoutes()

	//! start http server
	go func() {
		err := server.App.Listen(cfg.Http.BaseUrl + cfg.Http.Host + cfg.Http.Port)
		if err != nil {
			log.Logger.Fatalf("Error trying to listenning on port %s: %v", cfg.Http.Port, err)
		}
	}()
	//!
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		fmt.Printf("signal.Notify CTRL+C: %v", v)
	case done := <-ctx.Done():
		fmt.Printf("ctx.Done: %v", done)
	}
}
