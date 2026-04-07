package helpers

import (
	"github.com/gin-gonic/gin"
	"seasaloon-backend-go/structs"
)
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
type LoginResponse struct {
	Token string        `json:"token"`
	User  *structs.Users `json:"user"`
}

func Success[T any](c *gin.Context, code int, message string, data T) {
	c.JSON(code, Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response[any]{
		Code:    code,
		Message: message,
	})
}