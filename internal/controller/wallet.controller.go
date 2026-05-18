package controller

import (
	"backend-golang/internal/service"
	"backend-golang/pkg"
	"backend-golang/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService *service.WalletService
}

func NewWalletController(walletService *service.WalletService) *WalletController {
	return &WalletController{
		walletService: walletService,
	}
}

func (wc *WalletController) GetDashboard(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	user, err := wc.walletService.GetDashboard(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		utils.SendResponse(ctx, http.StatusInternalServerError, false, "internal server error", nil, err.Error())
		return
	}
	utils.SendResponse(ctx, http.StatusOK, true, "ok", user, nil)
}
