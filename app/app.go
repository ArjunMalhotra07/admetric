package app

import (
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
	//! Kafka
	kafkaService, err := services.NewKafkaService(cfg.Kafka.Brokers)
	if err != nil {
		log.Logger.Errorf("Failed to initialize Kafka: %v", err)
		return
	}
	defer kafkaService.Close()
	clickRepo := repo.NewClickRepo(db.DB)
	clickService := services.NewClickService(clickRepo, log, kafkaService)
	adService := services.NewAdService(adRepo, log)
	// http API server based on fiber
	server := server.NewHTTP(cfg, app, log, adService, clickService)
	//! start http server
	go func() {
		err := server.App.Listen(cfg.Http.Host + cfg.Http.Port)
		if err != nil {
			log.Logger.Fatalf("Error trying to listenning on port %s: %v", cfg.Http.Port, err)
		}
	}()
	//! Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Logger.Info("Shutting down server...")
	if err := server.App.Shutdown(); err != nil {
		log.Logger.Fatalf("Server forced to shutdown: %v", err)
	}
}
