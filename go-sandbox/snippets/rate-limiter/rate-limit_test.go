package main

import (
	"sync"
	"testing"
)

func TestConcurrentAccess(t *testing.T) {
	limiter := NewRateLimiter(5, 2) // 5 tokens max, refill 2/sec
	var wg sync.WaitGroup

	allowedCount := 0
	deniedCount := 0
	mu := sync.Mutex{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			} else {
				mu.Lock()
				deniedCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Allowed: %d, Denied: %d", allowedCount, deniedCount)

	//  expectations to
	if allowedCount+deniedCount != 10 {
		t.Errorf("Expected 10 total attempts, got %d", allowedCount+deniedCount)
	}

	if allowedCount < 1 {
		t.Errorf("Expected at least 1 allowed request")
	}
}
