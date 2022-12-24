// Package server provides the start and dependency registration of the trigger server
package server

import (
	"google.golang.org/grpc"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/domain/cron"
	"github.com/beihai0xff/pudding/app/trigger/domain/webhook"
)

func startCronTriggerService(server *grpc.Server, serviceName *string) error {
	// set serviceName
	// use string point to store serviceName, so that we can return it to the caller
	// this is not a good way to do this, but it works
	*serviceName = pb.CronTriggerService_ServiceDesc.ServiceName

	cronHandler := cron.NewHandler(cron.NewTrigger(db, brokerClient))
	pb.RegisterCronTriggerServiceServer(server, cronHandler)

	return nil
}
func startWebhookTriggerService(server *grpc.Server, serviceName *string) error {
	// set serviceName
	// use string point to store serviceName, so that we can return it to the caller
	// this is not a good way to do this, but it works
	*serviceName = pb.WebhookTriggerService_ServiceDesc.ServiceName

	// register Trigger server
	webhookHandler := webhook.NewHandler(webhook.NewTrigger(db, brokerClient))
	pb.RegisterWebhookTriggerServiceServer(server, webhookHandler)

	return nil
}
