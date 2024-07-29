package middleware

import (
	"tezos-delegation-service/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		l.Log.Info().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Dur("duration(µs)", time.Duration(duration.Microseconds())).
			Msg("Request processed")
	}
}
