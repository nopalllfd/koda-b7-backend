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

func SetupUserRoute(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	UserRepo := repository.NewUserRepo(db)
	UserService := service.NewUserService(UserRepo)
	UserController := controller.NewUserController(UserService)

	user := app.Group("/user")
	{
		user.GET("/profile", middleware.VerifyMiddleware(rc), UserController.GetProfile)
		user.PATCH("/profile", middleware.VerifyMiddleware(rc), UserController.EditProfile)
	}
}
