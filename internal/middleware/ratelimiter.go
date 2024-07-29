package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	limiter ratelimit.Limiter
}

func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		limiter: ratelimit.New(rate),
	}
}

func (mw *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.limiter.Take()
		c.Next()
	}
}
