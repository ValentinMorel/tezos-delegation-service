package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

type CircuitBreaker struct{}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{}
}

func (mw *CircuitBreaker) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := hystrix.Do("my_command", func() error {
			c.Next()
			if len(c.Errors) > 0 {
				return c.Errors.Last()
			}
			return nil
		}, nil)

		if err != nil {
			c.JSON(500, gin.H{"error": "service unavailable"})
			c.Abort()
		}
	}
}
