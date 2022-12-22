package resolver

import (
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/utils"
)

type Resolver interface {
	Register(serviceName string, ip string, port int) (string, error)
	Deregister(serviceID string) error
}

type OptionResolver func() Resolver

func GRPCRegistration(serviceName string, port int, opt OptionResolver) (Resolver, string) {
	rsv := opt()
	serviceID, err := rsv.Register(serviceName, utils.GetOutBoundIP(), port)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}
	return rsv, serviceID
}
