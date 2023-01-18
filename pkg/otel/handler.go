// Package otel provides OpenTelemetry utilities.
// handler.go register a http handler to handle OpenTelemetry requests.
package otel

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/beihai0xff/pudding/pkg/log"
)

// RegisterHandler register Swagger UI handler.
func RegisterHandler(gwmux *runtime.ServeMux) {
	if err := gwmux.HandlePath("GET", "/metrics", func(w http.ResponseWriter,
		r *http.Request, pathParams map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		log.Panicf("failed to register metrics handler: %v", err)
	}
}
