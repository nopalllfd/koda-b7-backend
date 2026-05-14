package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(app *gin.Engine, db *pgxpool.Pool) {
	AuthRepo := repository.NewAuthRepo(db)
	AuthService := service.NewAuthService(AuthRepo)
	AuthController := controller.NewAuthController(AuthService)
	auth := app.Group("/auth")
	{
		//login
		auth.POST("/login", AuthController.Login)
		//register
		auth.POST("/register", AuthController.Register)
	}
}
