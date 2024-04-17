package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"neiro-api/internal/utils"
	"time"
)

func JsonLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := utils.GetDurationInMilliseconds(start)

		entry := log.WithFields(log.Fields{
			"client_ip":  utils.GetClientIP(ctx),
			"duration":   duration,
			"method":     ctx.Request.Method,
			"path":       ctx.Request.RequestURI,
			"status":     ctx.Writer.Status(),
			"user_id":    utils.GetUserID(ctx),
			"referer":    ctx.Request.Referer(),
			"request_id": ctx.Writer.Header().Get("X-Request-ID"),
		})

		if ctx.Writer.Status() >= 500 {
			entry.Error(ctx.Errors.String())
			return
		}
		entry.Info("")
	}
}
