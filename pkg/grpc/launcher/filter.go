// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
// filter.go provides a filter to handle gRPC-Gateway requests.
package launcher

import (
	"encoding/json"
	"net/http"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
)

// GwMuxDecorator is a decorator for http.Handler.
type GwMuxDecorator func(http.Handler) http.Handler

// Handler is a http.Handler that serves gRPC-Gateway requests.
func Handler(h http.Handler, decors ...GwMuxDecorator) http.Handler {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}

// WithRequestLog returns a GwMuxDecorator that logs the request.
func WithRequestLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		request := map[string]interface{}{
			"protocol":    r.Proto,
			"method":      r.Method,
			"uri":         r.RequestURI,
			"remote_addr": r.RemoteAddr,
		}

		b, _ := json.Marshal(request)
		log.GetLoggerByName(logger.GRPCLoggerName).WithFields("module", "http").Info(string(b))
	})

}
