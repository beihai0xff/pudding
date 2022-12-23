// package shutdown provides a graceful shutdown mechanism.

package shutdown

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/resolver"
)

// OptionFunc is a function that can be used to configure a graceful shutdown.
type OptionFunc func(ctx context.Context) error

type ResolverPair struct {
	R         resolver.Resolver
	ServiceID string
}

// ResolverDeregister deregister the service from the resolver.
func ResolverDeregister(pairs ...ResolverPair) OptionFunc {
	return func(ctx context.Context) error {
		for _, p := range pairs {
			if err := p.R.Deregister(p.ServiceID); err != nil {
				log.Errorf("failed to deregister service: %v", err)
				return err
			}
		}
		return nil
	}
}

// HealthcheckServerShutdown shutdown the healthcheck server.
func HealthcheckServerShutdown(healthcheck *health.Server) OptionFunc {
	return func(ctx context.Context) error {
		healthcheck.Shutdown()
		return nil
	}
}

// HTTPServerShutdown shutdown the HTTP server.
func HTTPServerShutdown(httpServer *http.Server) OptionFunc {
	return func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	}
}

// GRPCServerShutdown shutdown the gRPC server.
func GRPCServerShutdown(s *grpc.Server) OptionFunc {
	return func(ctx context.Context) error {
		s.GracefulStop()
		return nil
	}
}

// LogSync flushes any buffered log entries.
func LogSync() OptionFunc {
	return func(ctx context.Context) error {
		log.Sync()
		return nil
	}

}

// GracefulShutdown
// 1. tell the load balancer to stop sending new requests
// 2. stop accepting new HTTP requests and wait for existing HTTP requests to finish
// 3. stop accepting new connections and RPCs and blocks until all the pending RPCs are finished.
// 4. flushing any buffered log entries
func GracefulShutdown(ctx context.Context, opts ...OptionFunc) {
	for _, opt := range opts {
		if err := opt(ctx); err != nil {
			log.Errorf("failed to shutdown gracefully: %v", err)
		}
	}
}
