package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"
	"github.com/ArjunMalhotra/pkg/logger"
)

const (
	batchSize = 100
)

type ClickService struct {
	clickRepo *repo.ClickRepo
	log       *logger.Logger
	cb        *circuitbreaker.CircuitBreaker
	counters  map[string]*CounterEntry
	kafka     *KafkaService

	batchMutex   sync.Mutex
	counterMutex sync.RWMutex
	currentBatch []model.Click
}

type CounterEntry struct {
	ClickCount int64
	LastUpdate time.Time
}

func NewClickService(clickRepo *repo.ClickRepo, log *logger.Logger, kafka *KafkaService) *ClickService {
	service := &ClickService{
		clickRepo:    clickRepo,
		log:          log,
		cb:           circuitbreaker.NewCircuitBreaker(5, time.Second*30, "click-service"),
		counters:     make(map[string]*CounterEntry),
		kafka:        kafka,
		currentBatch: make([]model.Click, 0, batchSize),
	}

	// Start Kafka consumer
	if err := service.kafka.StartConsumer(service); err != nil {
		log.Logger.Errorf("Failed to start Kafka consumer: %v", err)
	}

	return service
}

func (s *ClickService) ProcessClick(click model.Click) error {
	s.batchMutex.Lock()
	defer s.batchMutex.Unlock()

	s.currentBatch = append(s.currentBatch, click)
	if len(s.currentBatch) >= batchSize {
		return s.processBatch()
	}
	return nil
}

func (s *ClickService) processBatch() error {
	if len(s.currentBatch) == 0 {
		return nil
	}

	if s.cb.IsOpen() {
		// Circuit breaker open, save to Kafka for retry
		for _, click := range s.currentBatch {
			if err := s.kafka.PublishClick(click); err != nil {
				s.log.Logger.Errorf("Failed to republish click to Kafka: %v", err)
			}
		}
		s.currentBatch = s.currentBatch[:0]
		return nil
	}

	// Try database batch insert
	err := s.clickRepo.SaveBatch(s.currentBatch)
	if err != nil {
		s.log.Logger.Errorf("Failed to store batch in database: %v", err)
		s.cb.RecordFailure()
		// Republish to Kafka for retry
		for _, click := range s.currentBatch {
			if err := s.kafka.PublishClick(click); err != nil {
				s.log.Logger.Errorf("Failed to republish click to Kafka: %v", err)
			}
		}
	} else {
		s.cb.RecordSuccess()
		s.counterMutex.Lock()
		for _, click := range s.currentBatch {
			s.updateCounter(click)
		}
		s.counterMutex.Unlock()
	}

	s.currentBatch = s.currentBatch[:0]
	return nil
}

func (s *ClickService) RecordClick(click model.Click) error {
	return s.kafka.PublishClick(click)
}

func (s *ClickService) updateCounter(click model.Click) {
	s.counters["total"].ClickCount += 1
	adID := fmt.Sprintf("ad:%s", click.AdID)
	if entry, exists := s.counters[adID]; !exists {
		s.counters[adID] = &CounterEntry{
			ClickCount: 1,
			LastUpdate: time.Now(),
		}
	} else {
		entry.ClickCount++
		entry.LastUpdate = time.Now()
	}
}

func (s *ClickService) GetClickCount(adID string) int64 {
	s.counterMutex.RLock()
	defer s.counterMutex.RUnlock()
	id := fmt.Sprintf("ad:%s", adID)
	if entry, exists := s.counters[id]; exists {
		return entry.ClickCount
	}
	return 0
}

func (s *ClickService) AdExists(adID string) (bool, error) {
	return s.clickRepo.AdExists(adID)
}
