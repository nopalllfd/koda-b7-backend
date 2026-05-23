package controller

import (
	"backend-golang/internal/dto"
	errs "backend-golang/internal/err"
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"errors"
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
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid request body", nil, err.Error())
		return
	}
	if err := tc.transactionService.CheckPin(ctx.Request.Context(), claims.Id, user.Pin); err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			utils.SendResponse(ctx, http.StatusNotFound, false, "check pin failed", nil, err.Error())
			return
		}
		if errors.Is(err, errs.ErrPINNotSet) || errors.Is(err, errs.ErrInvalidPin) {
			utils.SendResponse(ctx, http.StatusBadRequest, false, "check pin failed", nil, err.Error())
			return
		}
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "check pin failed", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "your pin is valid", nil, nil)
}

func (tc *TransactionController) GetAllUserTransaction(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	transactionsHistory, err := tc.transactionService.GetAllUserTransaction(ctx.Request.Context(), claims.Id)
	if err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, "failed to get all transactions", nil, err.Error())
	}

	utils.SendResponse(ctx, http.StatusOK, true, "success to get all transactions", transactionsHistory, nil)
}

func (tc *TransactionController) Topup(ctx *gin.Context) {
	var input dto.TopupRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, errs.ErrInvalidInput.Error(), nil, err.Error())
		return
	}

	result, err := tc.transactionService.CreateTopup(ctx.Request.Context(), input)
	if err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, errs.ErrInternalServer.Error(), nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "topup successful", result, nil)
}

func (tc *TransactionController) Transfer(ctx *gin.Context) {
	var input dto.TransferRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, errs.ErrInvalidInput.Error(), nil, err.Error())
		return
	}

	result, err := tc.transactionService.CreateTransfer(ctx.Request.Context(), input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, errs.ErrSameWalletTransfer) || errors.Is(err, errs.ErrInsufficientBalance) {
			statusCode = http.StatusBadRequest
		}
		utils.SendResponse(ctx, statusCode, false, err.Error(), nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "transfer successful", result, nil)
}
