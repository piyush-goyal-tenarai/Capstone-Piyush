package resilience

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type CBState int

const (
	Closed   CBState = iota
	Open
	HalfOpen
)

func (s CBState) String() string {
	return [...]string{"CLOSED", "OPEN", "HALF-OPEN"}[s]
}

type CircuitBreaker struct {
	mu               sync.Mutex
	name             string
	state            CBState
	failureCount     int
	failureThreshold int
	successCount     int
	lastFailure      time.Time
	cooldown         time.Duration
}

func NewCircuitBreaker(name string, threshold int) *CircuitBreaker {
	return &CircuitBreaker{
		name:             name,
		state:            Closed,
		failureThreshold: threshold,
		cooldown:         10 * time.Second,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	if cb.state == Open {
		if time.Since(cb.lastFailure) > cb.cooldown {
			cb.state = HalfOpen
			cb.successCount = 0
			log.Printf("[CB:%s] OPEN → HALF-OPEN", cb.name)
		} else {
			cb.mu.Unlock()
			return fmt.Errorf("circuit breaker %s is OPEN", cb.name)
		}
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.state == HalfOpen || cb.failureCount >= cb.failureThreshold {
			cb.state = Open
			log.Printf("[CB:%s] → OPEN (failures: %d)", cb.name, cb.failureCount)
		}
		return err
	}

	if cb.state == HalfOpen {
		cb.successCount++
		if cb.successCount >= 2 {
			cb.state = Closed
			cb.failureCount = 0
			log.Printf("[CB:%s] → CLOSED (recovered)", cb.name)
		}
	} else {
		cb.failureCount = 0
	}
	return nil
}

func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state.String()
}
