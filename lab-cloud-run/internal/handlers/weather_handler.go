package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucasafonsokremer/goexpert/lab-cloud-run/internal/services"
)

type WeatherHandler struct {
	cepService     *services.CEPService
	weatherService *services.WeatherService
}

func NewWeatherHandler(cepService *services.CEPService, weatherService *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		cepService:     cepService,
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeatherByCEP(c *gin.Context) {
	cep := c.Param("cep")

	if !h.cepService.ValidateCEP(cep) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	city, err := h.cepService.GetCityFromCEP(cep)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	temp, err := h.weatherService.GetTemperature(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching weather data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, temp)
}
