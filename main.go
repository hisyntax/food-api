package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hisyntax/food-api/controllers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error finding .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/food/create", controllers.CreateFood)
	router.Run(":" + port)

}
