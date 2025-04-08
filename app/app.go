package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/internal/server"
	"github.com/ArjunMalhotra/internal/services"
	"github.com/ArjunMalhotra/pkg/db"
	"github.com/ArjunMalhotra/pkg/http"
	"github.com/ArjunMalhotra/pkg/logger"
	"github.com/go-redis/redis"
)

func Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//! config
	cfg := config.NewConfig()
	//! logger
	log, _ := logger.NewLogger(cfg)
	app := http.NewApp(log)
	//! mysql db
	db, err := db.NewMysqDB(cfg)
	if err != nil {
		log.Logger.Errorf("Failed to load db object")
		return
	}
	if err := db.Migrate(); err != nil {
		log.Logger.Fatalf("Error trying to migrate: %v", err)
		return
	}
	//! Count ads
	adRepo := repo.NewAdRepository(db.DB)
	count, err := adRepo.CountAds()
	if err != nil {
		log.Logger.Error("Failed to count adds ", err)
		return
	}
	if count == 0 {
		if err := db.Seed(); err != nil {
			log.Logger.Error(err)
			return
		} else {
			log.Logger.Info("Successfully seeded data")
		}
	}
	//! Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Url,
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Logger.Error("Failed to connect to redis ->  ", err)
		return
	}
	clickRepo := repo.NewClickRepo(db.DB)
	metricsRepo := repo.NewMetricsRepo(redisClient)
	clickService := services.NewClickService(clickRepo, metricsRepo, log)
	adService := services.NewAdService(adRepo)
	// http API server based on fiber
	server := server.NewHTTP(cfg, app, log, adService, clickService)
	// Register all APP APIs
	server.RegisterRoutes()

	//! start http server
	go func() {
		err := server.App.Listen(cfg.Http.Host + cfg.Http.Port)
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
