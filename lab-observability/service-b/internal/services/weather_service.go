package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type WeatherService struct {
	client     *http.Client
	apiKey     string
	apiBaseURL string
	tracer     trace.Tracer
}

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type Temperature struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		client:     &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
		apiKey:     os.Getenv("WEATHER_API_KEY"),
		apiBaseURL: "http://api.weatherapi.com/v1",
		tracer:     otel.GetTracerProvider().Tracer("service-b.weather"),
	}
}

func (s *WeatherService) GetTemperature(ctx context.Context, city string) (*Temperature, error) {
	ctx, span := s.tracer.Start(ctx, "get_temperature")
	defer span.End()

	span.SetAttributes(attribute.String("city", city))

	if s.apiKey == "" {
		span.SetAttributes(attribute.String("error", "api_key_not_set"))
		return nil, fmt.Errorf("WEATHER_API_KEY not set")
	}

	// Normalize city name - replace special characters and handle special cases
	city = strings.ReplaceAll(city, "ã", "a")
	city = strings.ReplaceAll(city, "á", "a")
	city = strings.ReplaceAll(city, "â", "a")
	city = strings.ReplaceAll(city, "é", "e")
	city = strings.ReplaceAll(city, "ê", "e")
	city = strings.ReplaceAll(city, "í", "i")
	city = strings.ReplaceAll(city, "ó", "o")
	city = strings.ReplaceAll(city, "ô", "o")
	city = strings.ReplaceAll(city, "ú", "u")

	// URL encode the city name
	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.apiBaseURL, s.apiKey, encodedCity)
	span.SetAttributes(attribute.String("weather.city.encoded", encodedCity))

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", "request_creation_failed"))
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("error", "weather_api_request_failed"))
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("weather.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		span.SetAttributes(attribute.String("error", "weather_api_error"))
		return nil, fmt.Errorf("weather API error: status=%d, response=%v", resp.StatusCode, errResp)
	}

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	celsius := weatherData.Current.TempC
	fahrenheit := celsius*1.8 + 32
	kelvin := celsius + 273.15

	return &Temperature{
		City:       city,
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}, nil
}
