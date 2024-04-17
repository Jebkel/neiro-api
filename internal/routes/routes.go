package routes

import (
	"github.com/gin-gonic/gin"
	"neiro-api/internal/handlers"
	"neiro-api/internal/handlers/auth"
	"neiro-api/internal/middlewares"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(middlewares.JsonLogMiddleware(), gin.Recovery())
	r.Use(middlewares.RequestID(middlewares.RequestIDOptions{AllowSetting: false}))
	r.Use(middlewares.CORS(middlewares.CORSOptions{}))

	handler := handlers.NewHandler()

	auth.HandlerAuth{Handler: handler}.Init(r.Group("/auth"))

	return r
}
