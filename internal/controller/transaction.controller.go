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

// Check User PIN
//
//	@Summary		Check user PIN
//	@Description	validate user PIN before transaction
//	@Tags			transaction
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UserPIN	true	"pin payload"
//	@Success		200		{object}	dto.RegisterSwaggerResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		401		{object}	dto.ErrorSwaggerResponse
//	@Failure		404		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/transactions/pin [post]
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

// Get All User Transactions
//
//	@Summary		Get all user transactions
//	@Description	get paginated transaction history for logged-in user
//	@Tags			transaction
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int		false	"page number"
//	@Param			limit	query		int		false	"limit per page"
//	@Param			search	query		string	false	"search keyword"
//	@Success		200		{object}	dto.TransactionPaginationResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/transactions [get]
func (tc *TransactionController) GetAllUserTransaction(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var query dto.TransactionQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, "invalid query params", nil, err.Error())
		return
	}

	result, err := tc.transactionService.GetAllUserTransaction(
		ctx.Request.Context(),
		claims.Id,
		query,
	)

	if err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "failed to get all transactions", nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "success to get all transactions", result, nil)
}

// Topup Balance
//
//	@Summary		Topup wallet balance
//	@Description	add balance to user wallet
//	@Tags			transaction
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.TopupRequest	true	"topup payload"
//	@Success		200		{object}	dto.TopupResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/transactions/topup [post]
func (tc *TransactionController) Topup(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var input dto.TopupRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, errs.ErrInvalidInput.Error(), nil, err.Error())
		return
	}

	result, err := tc.transactionService.CreateTopup(ctx.Request.Context(), claims.Id, input)
	if err != nil {
		utils.SendResponse(ctx, http.StatusInternalServerError, false, errs.ErrInternalServer.Error(), nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "topup successful", result, nil)
}

// Transfer Balance
//
//	@Summary		Transfer balance to another wallet
//	@Description	send money between wallets
//	@Tags			transaction
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.TransferRequest	true	"transfer payload"
//	@Success		200		{object}	dto.TransferResponse
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/transactions/transfer [post]
func (tc *TransactionController) Transfer(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	var input dto.TransferRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.SendResponse(ctx, http.StatusBadRequest, false, errs.ErrInvalidInput.Error(), nil, err.Error())
		return
	}

	result, err := tc.transactionService.CreateTransfer(ctx.Request.Context(), input, claims.Id)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if errors.Is(err, errs.ErrSameWalletTransfer) ||
			errors.Is(err, errs.ErrInsufficientBalance) {
			statusCode = http.StatusBadRequest
		}

		utils.SendResponse(ctx, statusCode, false, err.Error(), nil, err.Error())
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "transfer successful", result, nil)
}

// Transfer Balance
//
//	@Summary		get all methods for top up flow
//	@Description	get all payment methods
//	@Tags			transaction
//	@Produce		json
//	@Success		200		{object}	dto.PaymentMethods
//	@Failure		400		{object}	dto.ErrorSwaggerResponse
//	@Failure		500		{object}	dto.ErrorSwaggerResponse
//	@Router			/transactions/payments [get]
func (tc *TransactionController) GetAllPaymentMethods(ctx *gin.Context) {
	res, err := tc.transactionService.GetPaymentMethods(ctx.Request.Context())
	if err != nil {
		utils.SendResponse(
			ctx,
			http.StatusInternalServerError,
			false,
			"failed to get payment methods",
			nil,
			err.Error(),
		)
		return
	}

	utils.SendResponse(
		ctx,
		http.StatusOK,
		true,
		"success get payment methods",
		res,
		nil,
	)
}
