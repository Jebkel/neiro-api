package handlers

import "github.com/gin-gonic/gin"

type successResponse struct {
	Status  bool   `json:"status"`
	Data    *gin.H `json:"data"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Errors  *gin.H `json:"errors,omitempty"`
}

func NewSuccessResponse(c *gin.Context, statusCode int, data *gin.H, message string) {
	c.AbortWithStatusJSON(statusCode, &successResponse{
		Status:  true,
		Data:    data,
		Message: message,
	})
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, errors *gin.H) {
	c.AbortWithStatusJSON(statusCode, &errorResponse{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}
