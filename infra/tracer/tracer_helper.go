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
func (t *TracerHelper) StartSpan(ctx context.Context) (context.Context, trace.Span) {
	return t.templateData.OTELTracer.Start(ctx, t.templateData.RequestNameOtel)
}

// Injeta o contexto OTEL no header da requisição de saída
func (t *TracerHelper) InjectContext(ctx context.Context, req *http.Request) {
	carrier := propagation.HeaderCarrier(req.Header)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}
