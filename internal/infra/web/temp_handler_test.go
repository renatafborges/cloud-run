package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeather(t *testing.T) {
	expectedWeather := WeatherAPI{
		Current: Current{
			CelsiusTemperature:    25.0,
			FarhenheitTemperature: 77.0,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedWeather)
	}))
	defer server.Close()

	originalURL := weatherApiURL
	weatherApiURL = server.URL
	defer func() { weatherApiURL = originalURL }()

	location := ViaCEP{Localidade: "City Example"}
	weather, err := GetWeather(location)
	assert.NoError(t, err)
	assert.Equal(t, expectedWeather, weather)
}

func TestGetLocation(t *testing.T) {
	postCode := "12345678"
	expectedLocation := ViaCEP{
		Localidade: "City Example",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedLocation)
	}))

	defer server.Close()

	originalURL := viaCepApiURL
	viaCepApiURL = server.URL + "/"
	defer func() { viaCepApiURL = originalURL }()

	location, err := GetLocation(httptest.NewRecorder(), postCode)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, location)
}

func TestIsValidZipCode(t *testing.T) {
	validZip := "12345678"
	invalidZip := "1234abcd"

	assert.True(t, IsValidPostCode(validZip))
	assert.False(t, IsValidPostCode(invalidZip))
}

func TestConvertCelsiusToFahrenheit(t *testing.T) {
	celsius := 45.7
	expected := "114.3"

	result := ConvertCelsiusToFahrenheit(celsius)

	assert.Equal(t, expected, result)
}

func TestConvertCelsiusToKelvin(t *testing.T) {
	celsius := 35.6
	expected := "308.6"

	result := ConvertCelsiusToKelvin(celsius)

	assert.Equal(t, expected, result)
}
