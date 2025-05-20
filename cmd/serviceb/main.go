package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/israelalvesmelo/desafio-otel/cmd/config"
	"github.com/israelalvesmelo/desafio-otel/infra/tracer"
	"github.com/israelalvesmelo/desafio-otel/infra/webserver"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/domain/usecase"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/infra/gateway"
	"github.com/israelalvesmelo/desafio-otel/internal/temperature/infra/web/handler"
)

func main() {
	// Load config
	var cfg config.Config
	viperCfg := config.NewViper("env")
	viperCfg.ReadViper(&cfg)

	// Set up signal context
	ctx, cancel := setupSignalContext()
	defer cancel()

	// Initialize telemetry
	shutdown := tracer.SetupTelemetry(ctx, "service_b", cfg.Zipkin.Endpoint)
	defer shutdown()

	// Initialize tracer helper
	tracerHelper := tracer.NewTracerHelper(
		cfg.ServiceB.Host,
		"service_a:all",
		"service_a",
	)

	// Create gateway
	locationGateway := gateway.NewLocationGateway(&cfg.CEP, tracerHelper)
	temperatureGateway := gateway.NewTemperatureGateway(&cfg.Temperature, tracerHelper)

	// Create use case
	useCase := usecase.NewGetTemperatureUseCase(locationGateway, temperatureGateway)

	// Create handler
	handler := handler.NewTemperatureHandler(useCase, tracerHelper)

	// Create server
	server := webserver.NewWebServer(fmt.Sprintf(":%s", cfg.ServiceB.Port))
	server.AddMiddleware(chimiddleware.Logger)

	server.AddHandler("/temperature", handler.GetWeather)

	server.Start()

}

func setupSignalContext() (context.Context, context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, cancel
}
