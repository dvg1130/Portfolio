package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens     int
	capacity   int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed.Seconds()) * rl.refillRate
	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.capacity {
			rl.tokens = rl.capacity
		}
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func main() {
	limiter := NewRateLimiter(5, 1) // 5 tokens, refill 1/sec

	for i := 0; i < 8; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: allowed\n", i+1)
		} else {
			fmt.Printf("Request %d: rate limited\n", i+1)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
