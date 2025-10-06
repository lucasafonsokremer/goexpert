package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucasafonsokremer/goexpert/lab-cloud-run/internal/services"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	// Log do traceparent recebido para depuração
	traceparent := c.Request.Header.Get("traceparent")
	if traceparent != "" {
		// Use log padrão para não depender de framework
		println("[service-b] traceparent recebido:", traceparent)
	} else {
		println("[service-b] traceparent NÃO recebido")
	}

	ctx, span := otel.GetTracerProvider().Tracer("service-b").Start(c.Request.Context(), "get_weather_by_cep")
	defer span.End()

	// Log TraceID e SpanID do span criado
	spanCtx := span.SpanContext()
	println("[service-b] TraceID:", spanCtx.TraceID().String(), "SpanID:", spanCtx.SpanID().String())

	cep := c.Param("cep")
	span.SetAttributes(attribute.String("cep", cep))

	if !h.cepService.ValidateCEP(cep) {
		span.SetAttributes(attribute.String("error", "invalid_cep"))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	city, err := h.cepService.GetCityFromCEP(ctx, cep)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	temp, err := h.weatherService.GetTemperature(ctx, city)
	if err != nil {
		span.SetAttributes(attribute.String("error", "weather_service_error"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching weather data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, temp)
}
