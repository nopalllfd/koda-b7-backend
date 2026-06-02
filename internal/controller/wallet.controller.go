package controller

import (
	"log"
	"net/http"

	"github.com/nopalllfd/koda-b7-backend/internal/service"
	"github.com/nopalllfd/koda-b7-backend/pkg"
	"github.com/nopalllfd/koda-b7-backend/pkg/utils"

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

// Get Dashboard
//
//	@Summary		Get dashboard
//	@Description	get wallet dashboard data for authenticated user
//	@Tags			wallet
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	dto.DashboardSwaggerResponse
//	@Failure		401	{object}	dto.ErrorSwaggerResponse
//	@Failure		500	{object}	dto.ErrorSwaggerResponse
//	@Router			/wallet/dashboard [get]
func (wc *WalletController) GetDashboard(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	user, err := wc.walletService.GetDashboard(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())

		utils.SendResponse(
			ctx,
			http.StatusInternalServerError,
			false,
			"get dashboard failed",
			nil,
			err.Error(),
		)
		return
	}

	utils.SendResponse(ctx, http.StatusOK, true, "ok", user, nil)
}
