package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CEPService struct {
	client *http.Client
	tracer trace.Tracer
}

type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       bool   `json:"erro"`
}

func NewCEPService() *CEPService {
	return &CEPService{
		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		tracer: otel.GetTracerProvider().Tracer("service-b.cep"),
	}
}

func (s *CEPService) ValidateCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func (s *CEPService) GetCityFromCEP(ctx context.Context, cep string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "get_city_from_cep")
	defer span.End()

	span.SetAttributes(attribute.String("cep", cep))

	if !s.ValidateCEP(cep) {
		span.SetAttributes(attribute.String("error", "invalid_cep_format"))
		return "", fmt.Errorf("invalid zipcode")
	}

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	span.SetAttributes(attribute.String("viacep.url", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", "request_creation_failed"))
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("error", "viacep_request_failed"))
		return "", fmt.Errorf("can not find zipcode")
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("viacep.status_code", resp.StatusCode))

	var cepData ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepData); err != nil {
		span.SetAttributes(attribute.String("error", "decode_error"))
		return "", fmt.Errorf("can not find zipcode")
	}

	// Check if CEP exists and has valid data
	if cepData.Erro || cepData.Localidade == "" {
		span.SetAttributes(attribute.String("error", "cep_not_found"))
		return "", fmt.Errorf("can not find zipcode")
	}

	return cepData.Localidade, nil
}
