package circuitbreaker

import (
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	state            State
	failures         int
	lastFailure      time.Time
	mutex            sync.RWMutex
	name             string
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration, name string) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:            StateClosed,
		name:             name,
	}
}
func (cb *CircuitBreaker) IsOpen() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	// If circuit is open, check if it's time to try again
	if cb.state == StateOpen {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// We've waited long enough, move to half-open state
			// Using a read lock so we need to copy and reacquire with write lock
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			defer cb.mutex.Unlock()
			// Double-check state hasn't changed
			if cb.state == StateOpen && time.Since(cb.lastFailure) > cb.resetTimeout {
				cb.state = StateHalfOpen
				return false
			}
			return true
		}
		return true
	}
	return false
}
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.failures++
	cb.lastFailure = time.Now()

	if cb.state == StateHalfOpen || cb.failures >= cb.failureThreshold {
		cb.state = StateOpen // Open the circuit
	}
}
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.state == StateHalfOpen {
		// Reset circuit after a successful operation in half-open state
		cb.failures = 0
		cb.state = StateClosed
	} else if cb.state == StateClosed && cb.failures > 0 {
		// Decrease failure count on success
		cb.failures--
	}
}
