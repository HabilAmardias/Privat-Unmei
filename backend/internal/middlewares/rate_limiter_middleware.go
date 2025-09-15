package middlewares

import (
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	tokens     int
	capacity   int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	tokensToAdd := int(elapsed.Seconds()) * rl.refillRate
	rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
	rl.lastRefill = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func NewRateLimiter(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

var rateLimitersMutex sync.RWMutex

func RateLimiterMiddleware(rateLimiters map[string]*RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()

		rateLimitersMutex.RLock()
		limiter, exists := rateLimiters[clientIP]
		rateLimitersMutex.RUnlock()

		if !exists {
			rateLimitersMutex.Lock()
			if limiter, exists = rateLimiters[clientIP]; !exists {
				limiter = NewRateLimiter(constants.BurstSize, constants.RequestPerSecond)
				rateLimiters[clientIP] = limiter
			}
			rateLimitersMutex.Unlock()
		}

		if !limiter.Allow() {
			ctx.Error(customerrors.NewError(
				"too many request going on",
				errors.New("too many request going on"),
				customerrors.TooManyRequest,
			))
			ctx.Abort()
		}
		ctx.Next()
	}
}
