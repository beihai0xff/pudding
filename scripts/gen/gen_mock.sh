#!/bin/sh

mockgen -destination=../../test/mock/scheduler_mock.go -package=mock -source=../../app/scheduler/scheduler.go Scheduler
mockgen -destination=../../test/mock/broker_mock.go -package=mock -source=../../app/scheduler/broker/broker.go DelayBroker
mockgen -destination=../../test/mock/connector_mock.go -package=mock -source=../../app/scheduler/connector/connector.go RealTimeConnector

# mockgen -destination=../../test/mock/api/gen/pudding/scheduler/v1/scheduler_grpc_mock.pb.go -package=mock -source=../../api/gen/pudding/scheduler/v1/scheduler_grpc.pb.go SchedulerServiceClient