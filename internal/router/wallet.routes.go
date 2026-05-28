package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/middleware"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SetupUserWallet(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	WalletRepo := repository.NewWalletRepo(db)
	WalletService := service.NewWalletService(WalletRepo)
	WalletController := controller.NewWalletController(WalletService)

	user := app.Group("/wallet")
	user.Use(middleware.VerifyMiddleware(rc))
	{
		user.GET("/dashboard", WalletController.GetDashboard)

	}
}
