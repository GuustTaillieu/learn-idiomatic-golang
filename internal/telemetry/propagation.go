package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

var propagator = propagation.TraceContext{}

// InjectTraceParent writes the current span's W3C traceparent string
func InjectTraceParent(ctx context.Context) string {
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, carrier)
	return carrier["traceparent"]
}

// ExtractTraceContext extracts the traceparent into a base context
func ExtractTraceContext(ctx context.Context, traceParent string) context.Context {
	if traceParent == "" {
		return ctx
	}
	carrier := propagation.MapCarrier{"traceparent": traceParent}
	return propagator.Extract(ctx, carrier)
}
