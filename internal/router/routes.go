package router

import (
	"backend-golang/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRoutes(app *gin.Engine, db *pgxpool.Pool) {
	app.Use(middleware.CORSMiddleware)

	SetupAuthRoute(app, db)
	SetupUserRoute(app, db)
	SetupTransactionRoute(app, db)
	SetupUserWallet(app, db)

}
