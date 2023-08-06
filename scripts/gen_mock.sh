#!/bin/bash

set -ex

mockgen -destination=test/mock/app/broker/scheduler.go -package=mock -source=app/broker/scheduler.go Scheduler
mockgen -destination=test/mock/app/broker/storage/storage.go -package=mock -source=app/broker/storage/storage.go DelayStorage
mockgen -destination=test/mock/app/broker/connector/connector.go -package=mock -source=app/broker/connector/connector.go RealTimeConnector

mockgen -destination=test/mock/app/trigger/repo/webhook_template_repo.go -package=mock -source=app/trigger/repo/webhook_template_repo.go WebhookTemplate
mockgen -destination=test/mock/app/trigger/repo/cron_template_repo.go -package=mock -source=app/trigger/repo/cron_template_repo.go CronTemplate

mockgen -destination=test/mock/pkg/mq/kafka/client.go -package=mock -source=pkg/mq/kafka/client.go Client,Consumer

# mock function for scheduler_grpc.pb.go was manually generated
#mockgen -destination=test/mock/api/gen/pudding/broker/v1/scheduler_grpc.pb.go -package=mock -source=api/gen/pudding/broker/v1/scheduler_grpc.pb.go SchedulerServiceClient
