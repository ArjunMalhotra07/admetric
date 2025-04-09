package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"
	"github.com/ArjunMalhotra/pkg/logger"
	"gorm.io/gorm"
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
		// Update total clicks in DB for each ad
		for _, click := range s.currentBatch {
			if err := s.clickRepo.UpdateAdTotalClicks(click.AdID, 1); err != nil {
				s.log.Logger.Errorf("Failed to update total clicks for ad %s: %v", click.AdID, err)
			}
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

func (s *ClickService) GetClickCount(adID string) (int64, error) {
	s.counterMutex.RLock()
	defer s.counterMutex.RUnlock()

	// First try to get from in-memory counter
	id := fmt.Sprintf("ad:%s", adID)
	if entry, exists := s.counters[id]; exists {
		return entry.ClickCount, nil
	}

	// If not in memory, get from database
	totalClicks, err := s.clickRepo.GetAdTotalClicks(adID)
	if err != nil {
		return 0, err
	}

	// Update in-memory counter for future requests
	s.counters[id] = &CounterEntry{
		ClickCount: int64(totalClicks),
	}

	return int64(totalClicks), nil
}

func (s *ClickService) AdExists(adID string) (bool, error) {
	return s.clickRepo.AdExists(adID)
}

// ParseTimeFrame parses a timeframe string in the format "int+unit" (e.g., "37m", "7h", "3d")
func (s *ClickService) ParseTimeFrame(timeFrame string) (time.Duration, error) {
	if timeFrame == "" {
		return 1 * time.Hour, nil // Default to 1 hour
	}

	var value int
	var unit string

	// Extract numeric value and unit
	for i, c := range timeFrame {
		if c >= '0' && c <= '9' {
			continue
		}
		value = 0
		for j := 0; j < i; j++ {
			value = value*10 + int(timeFrame[j]-'0')
		}
		unit = timeFrame[i:]
		break
	}

	if value == 0 {
		return 0, fmt.Errorf("invalid time frame format: %s", timeFrame)
	}

	// Convert to duration based on unit
	switch unit {
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid time unit: %s", unit)
	}
}

// GetClickCountByTimeFrame gets the click count for an ad within a specific time frame
func (s *ClickService) GetClickCountByTimeFrame(adID string, timeFrame string) (int64, error) {
	// First check if ad exists
	exists, err := s.AdExists(adID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, gorm.ErrRecordNotFound
	}

	// Parse the time frame
	duration, err := s.ParseTimeFrame(timeFrame)
	if err != nil {
		return 0, err
	}

	// Get click count from database
	count, err := s.clickRepo.GetClickCountByTimeFrame(adID, duration)
	if err != nil {
		return 0, err
	}

	return count, nil
}
