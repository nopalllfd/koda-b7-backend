package main

import (
	// untuk logging
	// informasi dari protocol (status code)

	"backend-golang/internal/config"
	"backend-golang/internal/router"
	"log"

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
	// init
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := gin.Default()
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	rc, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis connection error. \ncause : %s", err.Error())
	}
	defer rc.Close()

	// route
	router.InitRoutes(app, db, rc)
	app.Run("localhost:5000")
}
