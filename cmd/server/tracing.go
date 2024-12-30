package main

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// OTELMux is a wrapper around http.ServeMux that adds OpenTelemetry instrumentation.
type OTELMux struct {
	*http.ServeMux
}

func (m *OTELMux) HandleFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	// Configure the "http.route" for the HTTP instrumentation.
	handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
	m.ServeMux.Handle(pattern, handler)
}
