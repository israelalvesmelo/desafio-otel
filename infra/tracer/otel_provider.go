package tracer

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func SetupTelemetry(ctx context.Context, providerName, urlExport string) func() {
	shutdown, err := newTracerProvider(providerName, urlExport)
	if err != nil {
		log.Fatalf("failed to initialize telemetry: %s", err.Error())
	}

	return func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("failed shutting down tracer provider: %s", err.Error())
		}
	}
}

type ShutdownFunction func(ctx context.Context) error

func newTracerProvider(serviceName, collectorURL string) (ShutdownFunction, error) {
	ctx := context.Background()
	res, resErr := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if resErr != nil {
		return nil, fmt.Errorf("failed to create reosource %w", resErr)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := zipkin.New(collectorURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zipkin exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}
