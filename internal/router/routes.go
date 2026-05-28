package router

import (
	"backend-golang/internal/middleware"
	"fmt"

	_ "backend-golang/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	app.Use(middleware.CORSMiddleware())
	fmt.Println("Swagger route loaded")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	SetupAuthRoute(app, db, rc)
	SetupUserRoute(app, db, rc)
	SetupTransactionRoute(app, db, rc)
	SetupUserWallet(app, db, rc)

}
