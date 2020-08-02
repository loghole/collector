package handlers

import (
	"net/http"

	"github.com/gadavy/tracing"
)

type TracingMiddleware struct {
	tracer *tracing.Tracer
}

func NewTracingMiddleware(tracer *tracing.Tracer) *TracingMiddleware {
	return &TracingMiddleware{tracer: tracer}
}

func (m *TracingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := m.tracer.NewSpan().WithName(r.URL.String()).ExtractHeaders(r.Header).Build()
		defer span.Finish()

		next.ServeHTTP(w, r.WithContext(span.Context(r.Context())))
	})
}
