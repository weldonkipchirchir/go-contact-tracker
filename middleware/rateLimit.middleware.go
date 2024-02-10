package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	rate       float64    //requests/tokens allowed per second
	burst      int        //burst tokens
	mu         sync.Mutex //access shared resource
	lastUpdate time.Time
	tokens     float64 //track of the number of tokens
}

func NewRateLimiter(rate float64, bust int) *RateLimiter {
	return &RateLimiter{
		rate:  rate,
		burst: bust,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock() //locks the mutex (rl.mu) to ensure thread safety.
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.tokens += elapsed * rl.rate
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastUpdate = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
