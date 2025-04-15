package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware function that logs all requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// End timer
		end := time.Now()
		latency := end.Sub(start)

		// Get status
		status := c.Writer.Status()

		// Log request
		fmt.Printf("[%s] %d %s %s %s\n",
			end.Format("2006-01-02 15:04:05"),
			status,
			c.Request.Method,
			c.Request.RequestURI,
			latency,
		)
	}
}
