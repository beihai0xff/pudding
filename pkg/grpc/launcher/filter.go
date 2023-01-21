// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
// filter.go provides a filter to handle gRPC-Gateway requests.
package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
)

var gwLogger = log.GetLoggerByName(logger.GRPCLoggerName).WithFields("module", "http")

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

var notLogPrefix = []string{"/metrics", "/pudding/broker/swagger"}

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
			"status":   rspProxy.GETHTTPStatus(),
			"response": rspProxy.Body(),
		}

		b, _ := json.Marshal(request)
		if rspProxy.GETHTTPStatus() != http.StatusOK {
			gwLogger.Error(string(b))
			return
		}
		gwLogger.Info(string(b))
	})

}

// responseProxy wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type responseProxy struct {
	http.ResponseWriter
	status int
	body   []byte
	Len    int
}

// WriteHeader writes the HTTP status code of the response.
func (p *responseProxy) WriteHeader(status int) {
	p.status = status
	p.ResponseWriter.WriteHeader(status)
}

func (p *responseProxy) Write(buf []byte) (int, error) {
	if p.status == 0 {
		p.WriteHeader(http.StatusOK)
	}
	n, err := p.ResponseWriter.Write(buf)
	p.body = append(p.body, buf[:n]...)
	p.Len += n
	return n, err
}

func (p *responseProxy) Body() string {
	return string(p.body)
}

// GETHTTPStatus returns the HTTP status code of the response.
func (p *responseProxy) GETHTTPStatus() int {
	return p.status
}
