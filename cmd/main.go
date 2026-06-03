package main

import (
	// untuk logging
	// informasi dari protocol (status code)

	"log"

	"github.com/nopalllfd/koda-b7-backend/internal/config"
	"github.com/nopalllfd/koda-b7-backend/internal/router"

	"github.com/gin-gonic/gin" // framework gin gonic
	"github.com/joho/godotenv" //godotenv
)

// @title						Backend Koda 7
// @version						1.0
// @description					Backend created by Koda using Gin

// @license.name				MIT

// @host						localhost:5000
// @BasePath					/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization
func main() {
	_ = godotenv.Load()

	app := gin.Default()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %s", err.Error())
	}
	defer db.Close()

	rc, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis connection error: %s", err.Error())
	}
	defer rc.Close()

	router.InitRoutes(app, db, rc)

	app.Run("0.0.0.0:8080")
}
