package services

import (
	"fmt"
	"log"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"

	"github.com/sony/gobreaker"
)

type ClickService struct {
	clickRepo   *repo.ClickRepo
	metricsRepo *repo.MetricsRepo
	cb          *gobreaker.CircuitBreaker
}

func NewClickService(clickRepo *repo.ClickRepo, metricsRepo *repo.MetricsRepo) *ClickService {
	return &ClickService{
		clickRepo:   clickRepo,
		metricsRepo: metricsRepo,
		cb:          circuitbreaker.NewCircuitBreaker("click-service"),
	}
}

func (s *ClickService) AdExists(adID string) (bool, error) {
	return s.clickRepo.AdExists(adID)
}

func (s *ClickService) RecordClick(click model.Click) error {
	if click.AdID == "" {
		return fmt.Errorf("ad ID is required")
	}
	if click.Timestamp.IsZero() {
		return fmt.Errorf("invalid timestamp")
	}
	if click.IP == "" {
		return fmt.Errorf("IP address is required")
	}
	if !s.clickRepo.IsValidIP(click.IP) {
		return fmt.Errorf("invalid IP address")
	}
	if !s.clickRepo.IsPlaybackTimeValid(click.PlaybackTime) {
		return fmt.Errorf("Playback time must be greater than zero")
	}
	//! Check if the adID exists before proceeding
	adExists, err := s.clickRepo.AdExists(click.AdID)
	if err != nil {
		log.Printf("Failed to check if ad exists: %v", err)
		return err
	}
	if !adExists {
		log.Printf("Ad with ID %s not found", click.AdID)
		return fmt.Errorf("ad with ID %s not found", click.AdID)
	}
	//! Rate Limiting: Check if the IP has exceeded the allowed number of clicks
	clickCount, err := s.clickRepo.GetClickCountByIP(click.IP)
	if err != nil {
		log.Printf("Failed to check click count for IP: %v", err)
		return err
	}
	if clickCount > 30 { // Example: Allow a maximum of 30 clicks per hour per IP
		log.Printf("Rate limit exceeded for IP %s", click.IP)
		return fmt.Errorf("rate limit exceeded")
	}

	// Wrap database operation with circuit breaker
	_, err = s.cb.Execute(func() (interface{}, error) {
		if err := s.clickRepo.Save(click); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Printf("Failed to record click (circuit breaker): %v", err)
		return err
	}

	// Wrap Redis operation with circuit breaker
	_, err = s.cb.Execute(func() (interface{}, error) {
		if err := s.metricsRepo.IncrementClickCount(click.AdID); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Printf("Failed to update analytics (circuit breaker): %v", err)
		return err
	}

	log.Printf("Click recorded successfully for ad %s from IP %s", click.AdID, click.IP)
	return nil
}

func (s *ClickService) GetClickCount(adID string) (int64, error) {
	// Wrap Redis operation with circuit breaker
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.metricsRepo.GetClickCount(adID)
	})
	if err != nil {
		log.Printf("Failed to get click count (circuit breaker): %v", err)
		return 0, err
	}

	return result.(int64), nil
}
