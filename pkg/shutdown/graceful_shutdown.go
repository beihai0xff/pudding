// package shutdown provides a graceful shutdown mechanism.

package shutdown

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/beihai0xff/pudding/pkg/grpc/resolver"
	"github.com/beihai0xff/pudding/pkg/log"
)

// OptionFunc is a function that can be used to configure a graceful shutdown.
type OptionFunc func(ctx context.Context) error

// ResolverDeregister deregister the service from the resolver.
func ResolverDeregister(pairs ...*resolver.Pair) OptionFunc {
	return func(ctx context.Context) error {
		// reverse the order of the pairs, so that to deregister is in the reverse order of the register
		// to ensure that the service is no longer accepting new requests when deregistering
		for i := range pairs {
			p := pairs[len(pairs)-i-1]
			if err := p.Resolver.Deregister(p.ServiceID); err != nil {
				log.Errorf("failed to deregister service [%s]: %v", p.ServiceID, err)
			} else {
				log.Infof("service [%s] deregistered", p.ServiceID)
			}
		}
		return nil
	}
}

// HealthServerShutdown shutdown the healthcheck server.
func HealthServerShutdown(healthServer *health.Server) OptionFunc {
	return func(ctx context.Context) error {
		healthServer.Shutdown()
		log.Infof("gRPC health server stopped")
		return nil
	}
}

// HTTPServerShutdown shutdown the HTTP server.
func HTTPServerShutdown(httpServer *http.Server) OptionFunc {
	return func(ctx context.Context) error {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Errorf("failed to shutdown HTTP server [%s]: %v", httpServer.Addr, err)
			return err
		}
		log.Infof("HTTP server [%s] stopped", httpServer.Addr)
		return nil
	}
}

// GRPCServerShutdown shutdown the gRPC server.
func GRPCServerShutdown(s *grpc.Server) OptionFunc {
	return func(ctx context.Context) error {
		s.GracefulStop()
		log.Infof("gRPC server stopped")
		return nil
	}
}

// LogSync flushes any buffered log entries.
func LogSync() OptionFunc {
	return func(ctx context.Context) error {
		log.Sync()
		log.Infof("server log flushed")
		return nil
	}
}

// GracefulShutdown
// 1. tell the load balancer this node is offline, and stop sending new requests
// 2. set the healthcheck status to unhealthy
// 3. stop accepting new HTTP requests and wait for existing HTTP requests to finish
// 4. stop accepting new connections and RPCs and blocks until all the pending RPCs are finished.
// 5. flushing any buffered log entries
func GracefulShutdown(ctx context.Context, opts ...OptionFunc) {
	for _, opt := range opts {
		if err := opt(ctx); err != nil {
			log.Errorf("failed to shutdown gracefully: %v", err)
		}
	}
}
