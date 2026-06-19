package router

import (
	"github.com/nopalllfd/koda-b7-backend/internal/controller"
	"github.com/nopalllfd/koda-b7-backend/internal/middleware"
	"github.com/nopalllfd/koda-b7-backend/internal/repository"
	"github.com/nopalllfd/koda-b7-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetupUserWallet(api *gin.RouterGroup, db *pgxpool.Pool, rc *redis.Client) {
	WalletRepo := repository.NewWalletRepo(db)
	WalletService := service.NewWalletService(WalletRepo)
	WalletController := controller.NewWalletController(WalletService)

	user := api.Group("/wallet")
	user.Use(middleware.VerifyMiddleware(rc))
	{
		user.GET("/dashboard", WalletController.GetDashboard)

	}
}
