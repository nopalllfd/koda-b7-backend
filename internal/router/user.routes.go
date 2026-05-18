package router

import (
	"backend-golang/internal/controller"
	"backend-golang/internal/middleware"
	"backend-golang/internal/repository"
	"backend-golang/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupUserRoute(r *gin.RouterGroup, app *gin.Engine, db *pgxpool.Pool) {
	UserRepo := repository.NewUserRepo(db)
	UserService := service.NewUserService(UserRepo)
	UserController := controller.NewUserController(UserService)

	user := app.Group("/user")
	user.Use(middleware.VerifyMiddleware)
	{
		user.GET("/profile", UserController.GetProfile)
		user.PUT("/profile", UserController.EditProfile)
	}
}
