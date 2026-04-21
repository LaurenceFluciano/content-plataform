package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LatencyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Printf("[ROUTE PERF] %s %s | Latência: %v", c.Request.Method, c.Request.URL.Path, latency)
	}
}
