package main

import (
	"job_portal/internal/repository"
	"job_portal/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := repository.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	// Close DB connection
	defer db.Close()
	// loading PORT variable from .env
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	// Initialize routes
	routes.InitRoutes(r, db)
	r.Run(":" + port)
}
