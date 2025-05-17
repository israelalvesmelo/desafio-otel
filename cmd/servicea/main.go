package main

import (
	"fmt"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/israelalvesmelo/desafio-otel/cmd/config"
	"github.com/israelalvesmelo/desafio-otel/infra/webserver"
	"github.com/israelalvesmelo/desafio-otel/internal/inputhandle/infra/handler"
)

func main() {
	// Load config
	var cfg config.Config
	viperCfg := config.NewViper("env")
	viperCfg.ReadViper(&cfg)

	// Create handler
	inputHandler := handler.NewInputHandler(&cfg)

	// Create server
	server := webserver.NewWebServer(fmt.Sprintf(":%s", cfg.ServiceA.Port))
	server.AddMiddleware(chimiddleware.Logger)

	server.AddHandler("/temperature", inputHandler.PostWeather)

	server.Start()
}
