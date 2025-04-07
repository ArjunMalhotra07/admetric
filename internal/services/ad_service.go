package services

import (
	"log"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"
	"github.com/sony/gobreaker"
)

type AdService struct {
	adRepo *repo.AdRepo
	cb     *gobreaker.CircuitBreaker
}

func NewAdService(adRepo *repo.AdRepo) *AdService {
	return &AdService{
		adRepo: adRepo,
		cb:     circuitbreaker.NewCircuitBreaker("ad-service"),
	}
}

func (s *AdService) GetAllAds() ([]model.Ad, error) {
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.adRepo.FetchAll()
	})
	if err != nil {
		log.Printf("err in cxt breaker -> %v", err)
		return nil, err
	}
	return result.([]model.Ad), nil
}
