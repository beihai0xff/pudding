// Package resolver provides a resolver interface for service discovery.
package resolver

import (
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/utils"
)

// Resolver is a service discovery resolver interface.
type Resolver interface {
	// RegisterGRPC register GRPC service
	RegisterGRPC(serviceName, ip string, port int) (string, error)
	// RegisterHTTP register HTTP service
	RegisterHTTP(path, ip string, port int) (string, error)
	// Deregister deregister service with serviceID
	Deregister(serviceID string) error
}

// OptionResolver is a resolver option
type OptionResolver func() Resolver

// Pair is a pair of resolver and serviceID
type Pair struct {
	// Resolver is a service discovery resolver interface.
	Resolver Resolver
	// ServiceID is a service unique ID
	ServiceID string
}

// GRPCRegistration register gRPC service with option resolver
func GRPCRegistration(serviceName string, port int, opt OptionResolver) *Pair {
	rsv := opt()

	serviceID, err := rsv.RegisterGRPC(serviceName, utils.GetOutBoundIP(), port)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	return &Pair{rsv, serviceID}
}

// HTTPRegistration register http service with option resolver
func HTTPRegistration(path string, port int, opt OptionResolver) *Pair {
	rsv := opt()
	serviceID, err := rsv.RegisterHTTP(path, utils.GetOutBoundIP(), port)

	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	return &Pair{rsv, serviceID}
}
