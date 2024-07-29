package middleware

import (
	"github.com/gin-gonic/gin"
)

type Cors struct{}

func NewCors() *Cors {
	return &Cors{}
}

func (r *Cors) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handle actual requests
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}
