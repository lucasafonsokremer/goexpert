package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WeatherService struct {
	client     *http.Client
	apiKey     string
	apiBaseURL string
}

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type Temperature struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		client:     &http.Client{},
		apiKey:     os.Getenv("WEATHER_API_KEY"),
		apiBaseURL: "http://api.weatherapi.com/v1",
	}
}

func (s *WeatherService) GetTemperature(city string) (*Temperature, error) {
	if s.apiKey == "" {
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
	fmt.Printf("Requesting weather data for city: %s (encoded: %s)\n", city, encodedCity)

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
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
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     kelvin,
	}, nil
}
