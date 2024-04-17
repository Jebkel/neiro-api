package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestIDOptions struct {
	AllowSetting bool
}

func RequestID(options RequestIDOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string

		if options.AllowSetting {
			requestID = c.Request.Header.Get("X-Request-ID")
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
