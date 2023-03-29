// Package server provides the start and dependency registration of the trigger server
package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/domain/cron"
	"github.com/beihai0xff/pudding/app/trigger/domain/webhook"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/grpc/launcher"
	"github.com/beihai0xff/pudding/pkg/log"
)

// withCronTriggerService returns a function that registers the CronTriggerService to the grpc server
func withCronTriggerService(conf *configs.TriggerConfig) launcher.StartServiceFunc {
	return func(server *grpc.Server, serviceName *string) error {
		brokerClient := getBrokerClient(conf)
		// db is the MySQLConfig database connection.
		db := mysql.New(conf.MySQLConfig)

		// set serviceName
		// use string point to store serviceName, so that we can return it to the caller
		// this is not a good way to do this, but it works
		*serviceName = pb.CronTriggerService_ServiceDesc.ServiceName

		cronHandler := cron.NewHandler(cron.NewTrigger(db, brokerClient))
		pb.RegisterCronTriggerServiceServer(server, cronHandler)

		return nil
	}
}

// withWebhookTriggerService returns a StartServiceFunc that registers the WebhookTriggerService to the gRPC server.
func withWebhookTriggerService(conf *configs.TriggerConfig) launcher.StartServiceFunc {
	return func(server *grpc.Server, serviceName *string) error {
		brokerClient := getBrokerClient(conf)
		// db is the MySQLConfig database connection.
		db := mysql.New(conf.MySQLConfig)

		// set serviceName
		// use string point to store serviceName, so that we can return it to the caller
		// this is not a good way to do this, but it works
		*serviceName = pb.WebhookTriggerService_ServiceDesc.ServiceName

		// register Trigger server
		webhookHandler := webhook.NewHandler(webhook.NewTrigger(db, brokerClient, conf.ServerConfig.WebhookPrefix))
		pb.RegisterWebhookTriggerServiceServer(server, webhookHandler)

		return nil
	}
}

func getBrokerClient(conf *configs.TriggerConfig) broker.SchedulerServiceClient {
	// init grpc connection
	conn, err := grpc.Dial(
		conf.ServerConfig.SchedulerConsulURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatalf("grpc Dial err: %v", err)
	}

	// TODO: close conn
	// defer conn.Close()

	// create scheduler service client
	return broker.NewSchedulerServiceClient(conn)
}
