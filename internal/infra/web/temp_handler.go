package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type TemperatureInputDTO struct {
	PostCode string `json:"cep"`
}

type TemperatureOutputDTO struct {
	Celsius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

type ViaCEP struct {
	Localidade string `json:"localidade"`
}

type WeatherAPI struct {
	Current Current `json:"current"`
}

type Current struct {
	CelsiusTemperature    float64 `json:"temp_c"`
	FarhenheitTemperature float64 `json:"temp_f"`
}

var viaCepApiURL = "http://viacep.com.br/ws/"
var weatherApiURL = "http://api.weatherapi.com/v1/current.json"

const apiKey = "bdb5695826b246b5a47230638241405"

func GetTemperatureByPostCode(w http.ResponseWriter, r *http.Request) {

	postCode := strings.Trim(r.URL.Path, "/temperature/")

	isValidPostCode := IsValidPostCode(postCode)

	if !isValidPostCode {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := GetLocation(w, postCode)

	if err != nil {
		http.Error(w, "cannot find zipCode", http.StatusNotFound)
		return
	}

	weather, err := GetWeather(location)

	if err != nil {
		http.Error(w, "could not get weather", http.StatusInternalServerError)
		return
	}

	formatCelcius := fmt.Sprintf("%.1f", weather.Current.CelsiusTemperature)

	var dto TemperatureOutputDTO = TemperatureOutputDTO{
		Celsius:    formatCelcius,
		Fahrenheit: ConvertCelsiusToFahrenheit(weather.Current.CelsiusTemperature),
		Kelvin:     ConvertCelsiusToKelvin(weather.Current.CelsiusTemperature),
	}

	err = json.NewEncoder(w).Encode(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func IsValidPostCode(postCode string) bool {
	re := regexp.MustCompile(`^\d{5}\d{3}$`)
	return re.MatchString(postCode)
}

func ConvertCelsiusToFahrenheit(celsius float64) string {
	var f = celsius*1.8 + 32
	return fmt.Sprintf("%.1f", f)

}

func ConvertCelsiusToKelvin(celsius float64) string {
	var k = celsius + 273
	return fmt.Sprintf("%.1f", k)
}

func GetLocation(w http.ResponseWriter, postCode string) (ViaCEP, error) {

	resp, err := http.Get(viaCepApiURL + postCode + "/json/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting postcode	: %v\n", err)
		return ViaCEP{}, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %v\n", err)
		return ViaCEP{}, err
	}

	var viaCepData ViaCEP

	err = json.Unmarshal(result, &viaCepData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling response: %v\n", err)
		return ViaCEP{}, err
	}

	if viaCepData.Localidade == "" {
		err = fmt.Errorf("error validating location: %s", viaCepData.Localidade)
		slog.Error("location is empty", "error", err)
		return ViaCEP{}, err
	}

	return viaCepData, nil
}

func GetWeather(v ViaCEP) (WeatherAPI, error) {

	params := map[string]string{
		"key": apiKey,
		"q":   v.Localidade,
		"aqi": "no",
	}

	u, err := url.Parse(weatherApiURL)
	if err != nil {
		slog.Error("error parsing URL", "url", weatherApiURL, "error", err)
		return WeatherAPI{}, err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		slog.Error("error sending request", "query", u.RawQuery, "error", err)
		return WeatherAPI{}, err
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %v\n", err)
		return WeatherAPI{}, err
	}

	var weather WeatherAPI

	err = json.Unmarshal(result, &weather)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling response: %v\n", err)
		return WeatherAPI{}, err
	}

	return weather, nil
}
