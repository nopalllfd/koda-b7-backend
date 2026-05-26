package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/middleware"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupTransactionRoute(app *gin.Engine, db *pgxpool.Pool) {
	TransactionRepo := repository.NewTransactionRepo()
	TransactionService := service.NewTransactionService(TransactionRepo, db)
	TransactionController := controller.NewTransactionController(TransactionService)
	trx := app.Group("/transactions")
	{
		trx.POST("/pin", middleware.VerifyMiddleware(db), TransactionController.CheckPin)
		trx.POST("/topup", middleware.VerifyMiddleware(db), TransactionController.Topup)
		trx.POST("/transfer", middleware.VerifyMiddleware(db), TransactionController.Transfer)
		trx.GET("", middleware.VerifyMiddleware(db), TransactionController.GetAllUserTransaction)
		trx.GET("/payments", TransactionController.GetAllPaymentMethods)
		trx.GET("/chart", middleware.VerifyMiddleware(db), TransactionController.GetChartData)
		trx.GET("/transfer/receivers", middleware.VerifyMiddleware(db), TransactionController.GetAllReceiverWithPagination)
	}
}
