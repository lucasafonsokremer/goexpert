package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CepRequest struct {
	CEP string `json:"cep" binding:"required"`
}

type CepHandler struct {
	weatherServiceURL string
	httpClient        *http.Client
	tracer            trace.Tracer
}

func NewCepHandler() *CepHandler {
	return &CepHandler{
		weatherServiceURL: os.Getenv("WEATHER_SERVICE_URL"),
		httpClient:        &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		tracer:            otel.GetTracerProvider().Tracer("service-a"),
	}
}

func (h *CepHandler) HandleCep(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "handle_cep_request")
	defer span.End()

	// Log TraceID e SpanID do span criado
	spanCtx := span.SpanContext()
	println("[service-a] TraceID:", spanCtx.TraceID().String(), "SpanID:", spanCtx.SpanID().String())

	var request CepRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		span.SetAttributes(attribute.String("error", "invalid_json"))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	// Validate CEP format (8 digits)
	matched, _ := regexp.MatchString(`^\d{8}$`, request.CEP)
	if !matched {
		span.SetAttributes(attribute.String("error", "invalid_cep_format"))
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	span.SetAttributes(
		attribute.String("cep", request.CEP),
	)

	// Create request to Service B with context
	req, _ := http.NewRequestWithContext(ctx, "GET", h.weatherServiceURL+"/weather/"+request.CEP, nil)

	// Forward to Service B
	resp, err := h.httpClient.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("error", "service_b_error"))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error forwarding request"})
		return
	}
	defer resp.Body.Close()

	// Forward the response from Service B
	var responseData interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		span.SetAttributes(attribute.String("error", "decode_error"))
	}

	span.SetAttributes(attribute.Int("status_code", resp.StatusCode))
	c.JSON(resp.StatusCode, responseData)
}
