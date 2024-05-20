package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"golang.org/x/sync/errgroup"
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

var Host string
var URL string

func init() {
	Host = os.Getenv("URL_TEMP")
	if Host == "" {
		Host = "localhost"
	}

	URL = "http://" + Host + ":8080/temperature/"
}

func main() {

	eg := errgroup.Group{}

	eg.Go(func() error {
		if !TestValidPostCode() {
			return errors.New("failed to test validPostCode")
		}
		return nil
	})

	eg.Go(func() error {
		if !TestInvalidPostCode() {
			return errors.New("failed to test invalidPostCode")
		}
		return nil
	})

	eg.Go(func() error {
		if !TestNonexistentPostCode() {
			return errors.New("failed to test nonexistentPostCode")
		}
		return nil
	})

	eg.Wait()

	if erro := eg.Wait(); erro != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("failed to execute automated tests", "caller", fmt.Sprintf("%s :%d", file, line), "error", erro)
		os.Exit(1)
	}

	os.Exit(0)
}

func TestValidPostCode() bool {

	client := &http.Client{}

	validPostCode := "04548004"

	req, err := http.NewRequest("GET", URL+validPostCode, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to make new request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to do request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to read from response body", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	var dto TemperatureOutputDTO

	err = json.Unmarshal(result, &dto)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to parse json", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	if len(result) > 0 && dto.Celsius != "0.0" && dto.Fahrenheit != "32" && dto.Kelvin != "273" && resp.StatusCode == 200 {
		return true
	}

	return false
}

func TestInvalidPostCode() bool {

	client := &http.Client{}

	invalidPostCode := "abc00000"

	req, err := http.NewRequest("GET", URL+invalidPostCode, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to make new request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to do request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to read from response body", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	var dto TemperatureOutputDTO

	err = json.Unmarshal(result, &dto)
	if err == nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to parse json", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	if len(result) > 0 && dto.Celsius != "0.0" && dto.Fahrenheit != "32" && dto.Kelvin != "273" && resp.StatusCode != 422 {
		return false
	}

	return true
}

func TestNonexistentPostCode() bool {

	client := &http.Client{}

	notFoundPostCode := "99999999"

	req, err := http.NewRequest("GET", URL+notFoundPostCode, nil)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to make new request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to do request", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to read from response body", "caller", fmt.Sprintf("%s :%d", file, line), "error", err)
		return false
	}

	var dto TemperatureOutputDTO

	err = json.Unmarshal(result, &dto)
	if err == nil {
		_, file, line, _ := runtime.Caller(0)
		slog.Error("unable to parse json", "caller", fmt.Sprintf("%s :%d", file, line), "result", string(result), "error", err)
		return false
	}

	if len(result) > 0 && dto.Celsius != "0.0" && dto.Fahrenheit != "32" && dto.Kelvin != "273" && resp.StatusCode != 404 {
		return false
	}

	return true
}
