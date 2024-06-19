package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter is the strategy interface for rate limiting
type RateLimiter interface {
	Allow() bool
}

// FixedWindowRateLimiter strategy implementation
type FixedWindowRateLimiter struct {
	requests int
	limit    int
	reset    time.Time
	mu       sync.Mutex
	window   time.Duration
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:  limit,
		window: window,
		reset:  time.Now().Add(window),
	}
}

func (rl *FixedWindowRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if time.Now().After(rl.reset) {
		rl.requests = 0
		rl.reset = time.Now().Add(rl.window)
	}

	if rl.requests < rl.limit {
		rl.requests++
		return true
	}
	return false
}

// SlidingWindowRateLimiter strategy implementation
type SlidingWindowRateLimiter struct {
	requests []time.Time
	limit    int
	mu       sync.Mutex
	window   time.Duration
}

func NewSlidingWindowRateLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *SlidingWindowRateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)
	newRequests := []time.Time{}

	for _, t := range rl.requests {
		if t.After(windowStart) {
			newRequests = append(newRequests, t)
		}
	}

	if len(newRequests) < rl.limit {
		rl.requests = append(newRequests, now)
		return true
	}

	rl.requests = newRequests
	return false
}

// TokenBucketRateLimiter strategy implementation
type TokenBucketRateLimiter struct {
	capacity  int
	tokens    int
	rate      int
	lastCheck time.Time
	mu        sync.Mutex
}

func NewTokenBucketRateLimiter(capacity, rate int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:  capacity,
		tokens:    capacity,
		rate:      rate,
		lastCheck: time.Now(),
	}
}

func (tb *TokenBucketRateLimiter) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastCheck).Seconds()
	tb.lastCheck = now

	tb.tokens += int(elapsed * float64(tb.rate))
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// LeakyBucketRateLimiter strategy implementation
type LeakyBucketRateLimiter struct {
	capacity  int
	queue     int
	rate      int
	lastCheck time.Time
	mu        sync.Mutex
}

func NewLeakyBucketRateLimiter(capacity, rate int) *LeakyBucketRateLimiter {
	return &LeakyBucketRateLimiter{
		capacity:  capacity,
		rate:      rate,
		lastCheck: time.Now(),
	}
}

func (lb *LeakyBucketRateLimiter) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastCheck).Seconds()
	lb.lastCheck = now

	lb.queue -= int(elapsed * float64(lb.rate))
	if lb.queue < 0 {
		lb.queue = 0
	}

	if lb.queue < lb.capacity {
		lb.queue++
		return true
	}

	return false
}

// Context for using rate limiting strategies
type RateLimiterContext struct {
	strategy RateLimiter
}

func NewRateLimiterContext(strategy RateLimiter) *RateLimiterContext {
	return &RateLimiterContext{
		strategy: strategy,
	}
}

func (ctx *RateLimiterContext) Allow() bool {
	return ctx.strategy.Allow()
}

func main() {
	// Example usage
	rl := NewRateLimiterContext(NewFixedWindowRateLimiter(10, time.Minute))
	// rl := NewRateLimiterContext(NewSlidingWindowRateLimiter(10, time.Minute))
	// rl := NewRateLimiterContext(NewTokenBucketRateLimiter(10, 1))
	// rl := NewRateLimiterContext(NewLeakyBucketRateLimiter(10, 1))

	for i := 0; i < 15; i++ {
		if rl.Allow() {
			fmt.Println("Request allowed")
		} else {
			fmt.Println("Rate limit exceeded")
		}
		time.Sleep(2 * time.Second) // simulate time between requests
	}
}
