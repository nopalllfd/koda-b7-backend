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

func SetupAuthRoute(api *gin.RouterGroup, db *pgxpool.Pool, rc *redis.Client) {
	UserRepo := repository.NewUserRepo(db)
	WalletRepo := repository.NewWalletRepo(db)
	AuthRepo := repository.NewAuthRepo(db, rc)
	AuthService := service.NewAuthService(AuthRepo, UserRepo, WalletRepo)
	AuthController := controller.NewAuthController(AuthService)

	auth := api.Group("/auth")
	{
		//login
		auth.POST("/login", AuthController.Login)
		//register
		auth.POST("/register", AuthController.Register)
		auth.DELETE("/logout", middleware.VerifyMiddleware(rc), AuthController.Logout)

		auth.POST("/pin", middleware.VerifyMiddleware(rc), AuthController.SetUserPin)
		auth.POST("/forgot-password", middleware.VerifyMiddleware(rc), AuthController.ForgotPassword)
		auth.PATCH("/pin", middleware.VerifyMiddleware(rc), AuthController.UpdateUserPin)
		auth.PATCH("/reset-password", AuthController.ResetPassword)
		auth.PATCH("/password", middleware.VerifyMiddleware(rc), AuthController.UpdateUserPassword)
		auth.GET("/user/:id", AuthController.GetUserDetail)
	}
}
