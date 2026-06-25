package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nopalllfd/koda-b7-backend/pkg"
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter := pkg.GetLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"message": "too many requests",
			})
			return
		}
		c.Next()
	}
}
