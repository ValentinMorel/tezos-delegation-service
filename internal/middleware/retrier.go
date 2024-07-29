package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Retrier struct{}

func NewRetrier() *Retrier {
	return &Retrier{}
}

func (r *Retrier) RetryMiddleware(attempts int, sleep time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		for i := 0; ; i++ {
			c.Next()
			if len(c.Errors) == 0 {
				return
			}

			if i >= (attempts - 1) {
				return
			}

			time.Sleep(sleep)
		}
	}
}
