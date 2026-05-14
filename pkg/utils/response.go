package utils

import (
	"backend-golang/internal/dto"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, code int, success bool, message string, data any, err any) {
	c.JSON(code, dto.Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   err,
	})
}
