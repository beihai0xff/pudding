// Package server provides the start and dependency registration of the broker server
package server

import (
	"google.golang.org/grpc"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/app/broker"
	"github.com/beihai0xff/pudding/app/broker/pkg/configs"
)

func startSchedulerService(server *grpc.Server, serviceName *string) error {
	// set serviceName
	// use string point to store serviceName, so that we can return it to the caller
	// this is not a good way to do this, but it works
	*serviceName = pb.SchedulerService_ServiceDesc.ServiceName

	// Initialize dependencies
	schedulerConfig := configs.GetBrokerConfig()
	delay, realtime := newQueue(schedulerConfig)
	s := broker.New(schedulerConfig, delay, realtime)
	s.Run()
	handler := broker.NewHandler(s)

	// register scheduler server
	pb.RegisterSchedulerServiceServer(server, handler)

	return nil
}
