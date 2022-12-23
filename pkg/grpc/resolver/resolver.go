package resolver

import (
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/utils"
)

type Resolver interface {
	RegisterGRPC(serviceName, ip string, port int) (string, error)
	RegisterHTTP(path, ip string, port int) (string, error)
	Deregister(serviceID string) error
}

type OptionResolver func() Resolver

type Pair struct {
	Resolver  Resolver
	ServiceID string
}

func GRPCRegistration(serviceName string, port int, opt OptionResolver) *Pair {
	rsv := opt()
	serviceID, err := rsv.RegisterGRPC(serviceName, utils.GetOutBoundIP(), port)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}
	return &Pair{rsv, serviceID}
}

func HTTPRegistration(path string, port int, opt OptionResolver) *Pair {
	rsv := opt()
	serviceID, err := rsv.RegisterHTTP(path, utils.GetOutBoundIP(), port)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}
	return &Pair{rsv, serviceID}
}
