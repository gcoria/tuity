package middleware

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate / 60

	if tokensToAdd > 0 {
		tb.tokens = int(math.Min(float64(tb.capacity), float64(tb.tokens+tokensToAdd)))
		tb.lastRefill = now
	}

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}

	return false
}

type RateLimiter struct {
	buckets map[string]*TokenBucket
	mutex   sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*TokenBucket),
	}
}

func (rl *RateLimiter) GetBucket(key string, capacity, refillRate int) *TokenBucket {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if bucket, exists := rl.buckets[key]; exists {
		return bucket
	}

	bucket := NewTokenBucket(capacity, refillRate)
	rl.buckets[key] = bucket
	return bucket
}

func RateLimitMiddleware(operation string, limit int) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.Next()
			return
		}

		key := fmt.Sprintf("%s:%s", operation, userID)
		bucket := limiter.GetBucket(key, limit, limit)

		if !bucket.Allow() {
			c.JSON(429, gin.H{
				"error":   "validation_error",
				"message": fmt.Sprintf("Rate limit exceeded for %s operation", operation),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func TweetCreateRateLimit(limit int) gin.HandlerFunc {
	return RateLimitMiddleware("tweet_create", limit)
}

func FollowOperationRateLimit(limit int) gin.HandlerFunc {
	return RateLimitMiddleware("follow_ops", limit)
}

func TimelineRequestRateLimit(limit int) gin.HandlerFunc {
	return RateLimitMiddleware("timeline_requests", limit)
}
