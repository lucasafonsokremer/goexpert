package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lucasafonsokremer/goexpert/lab-cloud-run/internal/handlers"
	"github.com/lucasafonsokremer/goexpert/lab-cloud-run/internal/services"
)

func main() {
	if os.Getenv("WEATHER_API_KEY") == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	r := gin.Default()

	cepService := services.NewCEPService()
	weatherService := services.NewWeatherService()
	weatherHandler := handlers.NewWeatherHandler(cepService, weatherService)

	r.GET("/weather/:cep", weatherHandler.GetWeatherByCEP)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
