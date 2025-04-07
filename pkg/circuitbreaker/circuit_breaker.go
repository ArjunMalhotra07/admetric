package circuitbreaker

import (
	"log"
	"time"

	"github.com/sony/gobreaker"
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // Number of requests allowed in half-open state
		Interval:    10 * time.Second, // Time window for counting failures
		Timeout:     30 * time.Second, // Time to wait before switching from open to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5 // Trip after 5 consecutive failures
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("%s: %s -> %s\n", name, from, to)
		},
	})
}
