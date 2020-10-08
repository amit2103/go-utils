package ratelimiter

import (
	"sync"
	"time"
)

// A rate limiter only with counters and no tickers
type RateLimiter struct {
	maxCount int
	interval time.Duration

	mu       sync.Mutex
	curCount int
	lastTime time.Time
}

//  maxCount param determines the max burst requests
//  interval specifies the window of the burst
func NewRateLimiter(maxCount int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		maxCount: maxCount,
		interval: interval,
	}
}

// True if okay m false if exceeded limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if time.Since(rl.lastTime) < rl.interval {
		if rl.curCount > 0 {
			rl.curCount--
			return true
		}
		return false
	}
	rl.curCount = rl.maxCount - 1
	rl.lastTime = time.Now()
	return true
}
