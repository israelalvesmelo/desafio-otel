package tracer

import (
	"go.opentelemetry.io/otel/trace"
)

type templateData struct {
	ExternalCallMethod string
	ExternalCallURL    string
	RequestNameOtel    string
	OTELTracer         trace.Tracer
}
