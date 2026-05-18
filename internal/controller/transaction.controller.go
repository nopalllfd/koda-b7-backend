package controller

import (
	"backend-golang/internal/dto"
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (tc *TransactionController) CheckPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var user dto.UserPIN
	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "Internal Server Error", nil, err.Error())
		return
	}
	if err := tc.transactionService.CheckPin(ctx.Request.Context(), claims.Id, user.Pin); err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "error", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "your pin is valid", nil, nil)
}
