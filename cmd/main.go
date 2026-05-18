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

	// route
	router.InitRoutes(app, db)
	app.Run("localhost:5000")
}
