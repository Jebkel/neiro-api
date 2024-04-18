package routes

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/handlers"
	"neiro-api/internal/handlers/auth"
	"neiro-api/internal/middlewares"
	"neiro-api/internal/redis"
	"neiro-api/pkg/ratelimit"
	"time"
)

func Init() *gin.Engine {
	r := gin.New()

	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.GetRedis(),
		Rate:        time.Second,
		Limit:       20,
	})

	mw := ratelimit.RateLimiter(store, &ratelimit.Options{})

	r.Use(middlewares.JsonLogMiddleware(), gin.Recovery())
	r.Use(middlewares.RequestID(middlewares.RequestIDOptions{AllowSetting: false}))
	r.Use(middlewares.CORS(middlewares.CORSOptions{}))

	handler := handlers.NewHandler()

	authGroup := r.Group("/auth")
	authGroup.Use(mw)
	auth.HandlerAuth{Handler: handler}.Init(authGroup)

	return r
}
