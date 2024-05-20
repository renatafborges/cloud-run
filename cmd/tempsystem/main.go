package main

import (
	"log/slog"
	"net/http"

	"github.com/renatafborges/cloud-run/configs"
	"github.com/renatafborges/cloud-run/internal/infra/web"
	"github.com/renatafborges/cloud-run/internal/infra/web/webserver"
)

func main() {

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webserver.AddHandler(http.MethodGet, "/temperature/{postcode}", web.GetTemperatureByPostCode)
	slog.Info("starting web server", "port", configs.WebServerPort)

	webserver.Start()
}
