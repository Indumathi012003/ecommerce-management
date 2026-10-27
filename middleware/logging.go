package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		start := time.Now()
		c.Next()

		reqID, _ := c.Get(RequestIDKey)
		errStr := c.Errors.String()
		if errStr == "" {
			errStr = "none"
		}

		log.Printf("request_id=%v method=%s path=%s status=%d duration=%v error=%s",
			reqID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start), errStr,
		)
	}
}
