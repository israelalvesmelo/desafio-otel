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
	"github.com/israelalvesmelo/desafio-otel/internal/inputhandle/infra/handler"
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
	shutdown := tracer.SetupTelemetry(ctx, "service_a_orchestration", cfg.Zipkin.Endpoint)
	defer shutdown()

	// Initialize tracer helper
	tracerHelper := tracer.NewTracerHelper(
		cfg.ServiceB.Host,
		"service_a:all",
		"service_a",
	)

	// Create handler
	inputHandler := handler.NewInputHandler(&cfg, tracerHelper)

	// Create server
	server := webserver.NewWebServer(fmt.Sprintf(":%s", cfg.ServiceA.Port))
	server.AddMiddleware(chimiddleware.Logger)

	server.AddHandler("/temperature", inputHandler.PostWeather)

	server.Start()
}

func setupSignalContext() (context.Context, context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, cancel
}
