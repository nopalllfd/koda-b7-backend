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

func SetupTransactionRoute(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	TransactionRepo := repository.NewTransactionRepo(rc)
	TransactionService := service.NewTransactionService(TransactionRepo, db)
	TransactionController := controller.NewTransactionController(TransactionService)
	trx := app.Group("/transactions")
	{
		trx.POST("/pin", middleware.VerifyMiddleware(rc), TransactionController.CheckPin)
		trx.POST("/topup", middleware.VerifyMiddleware(rc), TransactionController.Topup)
		trx.POST("/transfer", middleware.VerifyMiddleware(rc), TransactionController.Transfer)
		trx.GET("", middleware.VerifyMiddleware(rc), TransactionController.GetAllUserTransaction)
		trx.GET("/payments", TransactionController.GetAllPaymentMethods)
		trx.GET("/chart", middleware.VerifyMiddleware(rc), TransactionController.GetChartData)
		trx.GET("/transfer/receivers", middleware.VerifyMiddleware(rc), TransactionController.GetAllReceiverWithPagination)
	}
}
