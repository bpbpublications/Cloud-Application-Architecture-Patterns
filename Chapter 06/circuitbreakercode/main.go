package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Define breaker states
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// CircuitBreaker struct
type CircuitBreaker struct {
	state           State
	failureCount    int
	successCount    int
	failureThreshold int
	retryTimeout    time.Duration
	lastFailureTime time.Time
	mutex           sync.Mutex
}

func NewCircuitBreaker(failureThreshold int, retryTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:           Closed,
		failureThreshold: failureThreshold,
		retryTimeout:    retryTimeout,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case Open:
		if time.Since(cb.lastFailureTime) > cb.retryTimeout {
			fmt.Println("Retry timeout expired. Trying Half-Open state...")
			cb.state = HalfOpen
		} else {
			return errors.New("circuit breaker is OPEN")
		}
	case HalfOpen:
		// allow one trial
		err := fn()
		if err != nil {
			cb.trip()
			return errors.New("Half-Open trial failed. Going back to OPEN")
		}
		cb.reset()
		fmt.Println("Half-Open trial successful. Back to CLOSED state.")
		return nil
	}

	// Normal operation in Closed or Half-Open
	err := fn()
	if err != nil {
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.trip()
		}
		return err
	}

	cb.reset()
	return nil
}

func (cb *CircuitBreaker) trip() {
	cb.state = Open
	cb.lastFailureTime = time.Now()
	cb.failureCount = 0
	fmt.Println("Circuit Breaker tripped to OPEN state!")
}

func (cb *CircuitBreaker) reset() {
	cb.state = Closed
	cb.failureCount = 0
	cb.successCount = 0
}

func externalServiceCall() error {
	resp, err := http.Get("http://localhost:8081/downstream")
	if err != nil || resp.StatusCode != 200 {
		return errors.New("service call failed")
	}
	return nil
}

func main() {
	breaker := NewCircuitBreaker(3, 10*time.Second)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := breaker.Call(externalServiceCall)
		if err != nil {
			http.Error(w, "Downstream service unavailable: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pong from upstream service!"))
	})

	fmt.Println("Starting Circuit Breaker Demo on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
