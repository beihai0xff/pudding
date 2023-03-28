// Package server provides the start and dependency registration of the broker server
package server

import (
	"google.golang.org/grpc"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/app/broker"
	"github.com/beihai0xff/pudding/pkg/grpc/launcher"

	"github.com/beihai0xff/pudding/configs"
)

func withSchedulerService(conf *configs.BrokerConfig) launcher.StartServiceFunc {
	return func(server *grpc.Server, serviceName *string) error {
		// set serviceName
		// use string point to store serviceName, so that we can return it to the caller
		// it is not a good way to do that, but it works
		*serviceName = pb.SchedulerService_ServiceDesc.ServiceName

		// Initialize dependencies
		delay, realtime := newQueue(conf)
		s := broker.New(conf, delay, realtime)
		s.Run()
		handler := broker.NewHandler(s)

		// register scheduler server
		pb.RegisterSchedulerServiceServer(server, handler)

		return nil
	}
}
