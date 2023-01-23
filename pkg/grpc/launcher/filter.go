// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
// filter.go provides a filter to handle gRPC-Gateway requests.
package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

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

// WithRedirectToHTTPS returns a GwMuxDecorator that redirects HTTP requests to HTTPS.
func WithRedirectToHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}

var notLogPrefix = []string{"/metrics", "/healthz", "/pudding/broker/swagger"}

// WithRequestLog returns a GwMuxDecorator that logs the request.
func WithRequestLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if lo.ContainsBy(notLogPrefix, func(item string) bool { return strings.HasPrefix(r.RequestURI, item) }) {
			h.ServeHTTP(w, r)
			return
		}

		// x, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		// 	return
		// }
		rspProxy := &responseProxy{ResponseWriter: w}

		h.ServeHTTP(rspProxy, r)

		requestLog(rspProxy, r)
	})
}

func requestLog(p *responseProxy, r *http.Request) {
	realIP := r.Header.Get("X-Forwarded-For")
	if realIP == "" {
		realIP = r.RemoteAddr
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	request := map[string]interface{}{
		"scheme":       scheme,
		"request_line": fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		"remote_addr":  realIP,
		// "request":      string(x),
		"status":   p.GETHTTPStatusCode(),
		"response": p.GetBody(),
	}

	b, _ := json.Marshal(request)
	if p.GETHTTPStatusCode() != http.StatusOK {
		logger.GetGRPCLogger().Error(string(b))
		return
	}
	logger.GetGRPCLogger().Info(string(b))
}
