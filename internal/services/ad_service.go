package services

import (
	"fmt"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"
	"github.com/ArjunMalhotra/pkg/logger"
)

type AdService struct {
	adRepo *repo.AdRepo
	log    *logger.Logger
	cb     *circuitbreaker.CircuitBreaker
}

func NewAdService(adRepo *repo.AdRepo, log *logger.Logger) *AdService {
	return &AdService{
		adRepo: adRepo,
		log:    log,
		cb:     circuitbreaker.NewCircuitBreaker(5, 30*time.Second, "ad-service"),
	}
}

func (s *AdService) GetAllAds() ([]model.Ad, error) {
	if s.cb.IsOpen() {
		return nil, fmt.Errorf("circuit breaker is open for ad-service")
	}
	ads, err := s.adRepo.FetchAll()
	if err != nil {
		s.cb.RecordFailure()
		return nil, err
	}
	s.cb.RecordSuccess()
	return ads, nil
}

// ParseTimeframe parses a timeframe string in the format "int+h/d" (e.g., "56h", "3d")
func ParseTimeframe(timeframe string) (time.Duration, error) {
	if timeframe == "" {
		return 1 * time.Hour, nil // Default to 1 hour
	}

	var value int
	var unit string

	// Extract numeric value and unit
	for i, c := range timeframe {
		if c >= '0' && c <= '9' {
			continue
		}
		value = 0
		for j := 0; j < i; j++ {
			value = value*10 + int(timeframe[j]-'0')
		}
		unit = timeframe[i:]
		break
	}

	if value == 0 {
		return 1 * time.Hour, nil // Default to 1 hour if parsing fails
	}

	// Convert to duration based on unit
	switch unit {
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 1 * time.Hour, nil // Default to 1 hour for unknown units
	}
}
