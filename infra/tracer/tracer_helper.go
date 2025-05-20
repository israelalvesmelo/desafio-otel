package tracer

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TracerHelper struct {
	templateData templateData
}

func NewTracerHelper(externalCallUrl, requestNameOtel, serviceName string) *TracerHelper {
	return &TracerHelper{
		templateData: templateData{
			ExternalCallURL: externalCallUrl,
			RequestNameOtel: requestNameOtel,
			OTELTracer:      otel.Tracer(serviceName),
		},
	}
}

// Extrai o contexto OTEL do header HTTP recebido
func (t *TracerHelper) ExtractContext(r *http.Request) context.Context {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// Cria um novo span a partir do contexto
func (t *TracerHelper) StartSpan(ctx context.Context, spanName ...string) (context.Context, trace.Span) {
	name := t.templateData.RequestNameOtel
	if len(spanName) > 0 {
		name = spanName[0]
	}
	return t.templateData.OTELTracer.Start(ctx, name)
}

// Injeta o contexto OTEL no header da requisição de saída
func (t *TracerHelper) InjectContext(ctx context.Context, req *http.Request) {
	carrier := propagation.HeaderCarrier(req.Header)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}
