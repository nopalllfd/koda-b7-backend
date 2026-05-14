package middleware

import (
	"backend-golang/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomMiddleware(ctx *gin.Context) {
	xkodax := ctx.GetHeader("X-Koda-X")

	if xkodax != "aku koda" {
		ctx.AbortWithStatusJSON(http.StatusConflict, dto.Response{
			Success: false,
			Message: "error",
			Data:    nil,
			Error:   "wrong user of header",
		})
	}
}
