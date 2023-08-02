// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
// response_proxy.go provides a response proxy to log the response.
package launcher

import "net/http"

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

// Write writes the body of the response.
func (p *responseProxy) Write(buf []byte) (int, error) {
	if p.status == 0 {
		p.WriteHeader(http.StatusOK)
	}

	n, err := p.ResponseWriter.Write(buf)

	p.body = append(p.body, buf[:n]...)
	p.Len += n

	return n, err
}

// GetBody returns the body of the response.
func (p *responseProxy) GetBody() string {
	return string(p.body)
}

// GETHTTPStatusCode returns the HTTP status code of the response.
func (p *responseProxy) GETHTTPStatusCode() int {
	return p.status
}
