package router

import (
	"fmt"

	"github.com/nopalllfd/koda-b7-backend/internal/middleware"

	_ "github.com/nopalllfd/koda-b7-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(app *gin.Engine, db *pgxpool.Pool, rc *redis.Client) {
	app.Use(middleware.CORSMiddleware())
	fmt.Println("Swagger route loaded")
	app.Static("/img", "public/img")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	
	api := app.Group("/api")
	{
		SetupAuthRoute(api, db, rc)
		SetupUserRoute(api, db, rc)
		SetupTransactionRoute(api, db, rc)
		SetupUserWallet(api, db, rc)
	}

}
