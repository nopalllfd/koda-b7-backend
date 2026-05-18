package router

import (
	"backend-golang/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRoutes(app *gin.Engine, db *pgxpool.Pool) {
	app.Use(middleware.CORSMiddleware)
	api := app.Group("/api")
	{
		SetupAuthRoute(api, app, db)
		SetupUserRoute(api, app, db)
		SetupTransactionRoute(api, app, db)
	}

}
