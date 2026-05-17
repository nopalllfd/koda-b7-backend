package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/middleware"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *gin.Engine, db *pgxpool.Pool) {
	UserRepo := repository.NewUserRepo(db)
	UserService := service.NewUserService(UserRepo)
	UserController := controller.NewUserController(UserService)
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
		auth.POST("/register/pin", AuthController.AddPin)
	}

	user := app.Group("/user")
	{
		user.GET("/profile", middleware.VerifyMiddleware, UserController.GetProfile)
		user.PUT("/profile", middleware.VerifyMiddleware, UserController.EditProfile)
	}
}
