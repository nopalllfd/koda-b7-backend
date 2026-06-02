package utils

import (
	"github.com/nopalllfd/koda-b7-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, code int, success bool, message string, data any, err any) {
	var errMsg any = nil

	if err != nil {
		errMsg = err
	}
	c.JSON(code, dto.Response{
		Success: success,
		Message: message,
		Data:    data,
		Error:   errMsg,
	})
}
