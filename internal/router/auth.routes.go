package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/middleware"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupAuthRoute(app *gin.Engine, db *pgxpool.Pool) {
	UserRepo := repository.NewUserRepo(db)
	WalletRepo := repository.NewWalletRepo(db)
	AuthRepo := repository.NewAuthRepo(db)
	AuthService := service.NewAuthService(AuthRepo, UserRepo, WalletRepo)
	AuthController := controller.NewAuthController(AuthService)

	auth := app.Group("/auth")
	{
		//login
		auth.POST("/login", AuthController.Login)
		//register
		auth.POST("/register", AuthController.Register)

		auth.POST("/pin", middleware.VerifyMiddleware, AuthController.SetUserPin)

		auth.PATCH("/pin", middleware.VerifyMiddleware, AuthController.UpdateUserPin)
		auth.PATCH("/password", middleware.VerifyMiddleware, AuthController.UpdateUserPassword)
	}
}
