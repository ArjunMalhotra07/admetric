package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/ArjunMalhotra/internal/repo"
	"github.com/ArjunMalhotra/pkg/circuitbreaker"
	"github.com/ArjunMalhotra/pkg/logger"
)

type ClickService struct {
	clickRepo    *repo.ClickRepo
	metricsRepo  *repo.MetricsRepo
	log          *logger.Logger
	clickChan    chan model.Click
	cb           *circuitbreaker.CircuitBreaker
	counterMutex sync.RWMutex
	counters     map[string]*CounterEntry
}
type CounterEntry struct {
	value      int64
	lastUpdate time.Time
}

func NewClickService(clickRepo *repo.ClickRepo, metricsRepo *repo.MetricsRepo, log *logger.Logger) *ClickService {
	service := &ClickService{
		clickRepo:   clickRepo,
		metricsRepo: metricsRepo,
		log:         log,
		clickChan:   make(chan model.Click, 10000),
		cb:          circuitbreaker.NewCircuitBreaker(5, time.Second*30, "click-service"),
		counters:    make(map[string]*CounterEntry),
	}
	go service.processClickQueue()
	go service.periodicFlush(40 * time.Minute)
	go service.Prune(1 * time.Hour)
	return service
}

func (s *ClickService) AdExists(adID string) (bool, error) {
	return s.clickRepo.AdExists(adID)
}
func (s *ClickService) RecordClick(click model.Click) error {
	select {
	case s.clickChan <- click:
		return nil
	//! in case channel is full
	default:
		return s.saveToLocalBackup(click)
	}
}
func (s *ClickService) GetClickCount(adID string) int64 {
	s.counterMutex.RLock()
	defer s.counterMutex.RUnlock()
	adKey := fmt.Sprintf("ad:%s", adID)
	return s.counters[adKey].value
}
func (s *ClickService) updateCounters(click model.Click) {
	s.counterMutex.Lock()
	defer s.counterMutex.Unlock()
	s.counters["total"].value += 1
	adKey := fmt.Sprintf("ad:%s", click.AdID)
	s.counters[adKey].value += 1
}
func (s *ClickService) processClickQueue() {
	for click := range s.clickChan {
		//! Try to store in DB with circuit breaker
		if s.cb.IsOpen() {
			s.log.Logger.Warn("Circuit breaker open, saving click to backup", "click_id", click.ID)
			if err := s.saveToLocalBackup(click); err != nil {
				s.log.Logger.Errorf("Failed to save click to backup", "error", err)
			}
			continue
		}
		// Try to store in database
		err := s.clickRepo.Save(click)
		if err != nil {
			s.log.Logger.Error("Failed to store click in database", "error", err)
			s.cb.RecordFailure()
			// Save to local backup
			if backupErr := s.saveToLocalBackup(click); backupErr != nil {
				s.log.Logger.Errorf("Failed to save click to backup", "error", backupErr)
			}
		} else {
			s.cb.RecordSuccess()
			// Update in-memory counters for real-time stats
			s.updateCounters(click)
		}
	}
}
func (s *ClickService) saveToLocalBackup(click model.Click) error {
	// Ensure backup directory exists
	backupDir := "./backup/clicks"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	// Create a backup file with timestamp and click ID
	filename := filepath.Join(backupDir, fmt.Sprintf("%d-%s.json",
		time.Now().UnixNano(), click.ID))
	// Marshal click data
	data, err := json.Marshal(click)
	if err != nil {
		return err
	}
	// Write to file
	return os.WriteFile(filename, data, 0644)
}
func (s *ClickService) Prune(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {

	}
}
func (s *ClickService) RemoveInMemoryAds() {

}

// periodicFlush periodically syncs in-memory counters to the database
func (s *ClickService) periodicFlush(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		s.flushCountersToDatabase()
		s.processPendingClickFiles()
	}
}

// flushCountersToDatabase syncs in-memory counters to database
func (s *ClickService) flushCountersToDatabase() {
	s.counterMutex.RLock()
	adsData := make(map[string]*CounterEntry)
	for k, ad := range s.counters {
		adsData[k] = ad
	}
	s.counterMutex.RUnlock()
	// Update metrics in database
	for key, count := range adsData {
		// Update metrics in a transaction if needed
		s.log.Logger.Info("Flushing counter to database", "key", key, "count", count)
		// s.metricsRepo.UpdateMetric(key, count)
	}
}
func (s *ClickService) processPendingClickFiles() {
	// Skip if circuit breaker is still open
	if s.cb.IsOpen() {
		return
	}
	backupDir := "./backup/clicks"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		if !os.IsNotExist(err) {
			s.log.Logger.Error("Failed to read backup directory", "error", err)
		}
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Read click data from file
		data, err := os.ReadFile(filepath.Join(backupDir, file.Name()))
		if err != nil {
			s.log.Logger.Error("Failed to read backup file", "file", file.Name(), "error", err)
			continue
		}
		// Unmarshal click data
		var click model.Click
		if err := json.Unmarshal(data, &click); err != nil {
			s.log.Logger.Error("Failed to unmarshal click data", "file", file.Name(), "error", err)
			continue
		}
		// Try to store in database
		if err := s.clickRepo.StoreClick(click); err != nil {
			s.log.Logger.Error("Failed to store click from backup", "error", err)
			s.cb.RecordFailure()
			break // Stop processing if DB is having issues
		} else {
			s.cb.RecordSuccess()
			s.updateCounters(click)
			// Delete the backup file
			if err := os.Remove(filepath.Join(backupDir, file.Name())); err != nil {
				s.log.Logger.Errorf("Failed to delete backup file", "file", file.Name(), "error", err)
			}
		}
	}
}
