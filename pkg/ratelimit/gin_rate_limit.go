package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Код взят с https://github.com/JGLTechnologies/gin-rate-limit

type Info struct {
	Limit         uint
	RateLimited   bool
	ResetTime     time.Time
	RemainingHits uint
}

type Store interface {
	Limit(key string, c *gin.Context) Info
}

type Options struct {
	ErrorHandler   func(*gin.Context, Info)
	KeyFunc        func(*gin.Context) string
	BeforeResponse func(c *gin.Context, info Info)
}

func RateLimiter(s Store, options *Options) gin.HandlerFunc {
	if options == nil {
		options = &Options{}
	}
	if options.ErrorHandler == nil {
		options.ErrorHandler = func(c *gin.Context, info Info) {
			c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", info.Limit))
			c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", info.ResetTime.Unix()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
				"errors":  gin.H{},
			})
		}
	}
	if options.BeforeResponse == nil {
		options.BeforeResponse = func(c *gin.Context, info Info) {
			c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", info.Limit))
			c.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%v", info.RemainingHits))
			c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", info.ResetTime.Unix()))
		}
	}
	if options.KeyFunc == nil {
		options.KeyFunc = func(c *gin.Context) string {
			return c.ClientIP() + c.FullPath()
		}
	}
	return func(c *gin.Context) {
		key := options.KeyFunc(c)
		info := s.Limit(key, c)
		options.BeforeResponse(c, info)
		if c.IsAborted() {
			return
		}
		if info.RateLimited {
			options.ErrorHandler(c, info)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
