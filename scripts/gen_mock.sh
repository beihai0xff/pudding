#!/bin/bash

set -e

mockgen -destination=test/mock/app/broker/scheduler_mock.go -package=mock -source=app/broker/scheduler.go Scheduler
mockgen -destination=test/mock/app/broker/storage/storage_mock.go -package=mock -source=app/broker/storage/storage.go DelayStorage
mockgen -destination=test/mock/app/broker/connector/connector_mock.go -package=mock -source=app/broker/connector/connector.go RealTimeConnector

mockgen -destination=test/mock/app/trigger/repo/webhook_template_repo_mock.go -package=mock -source=app/trigger/repo/webhook_template_repo.go WebhookTemplate
mockgen -destination=test/mock/app/trigger/repo/cron_template_repo_mock.go -package=mock -source=app/trigger/repo/cron_template_repo.go CronTemplate
# mockgen -destination=.test/mock/api/gen/pudding/scheduler/v1/scheduler_grpc_mock.pb.go -package=mock -source=../../api/gen/pudding/scheduler/v1/scheduler_grpc.pb.go SchedulerServiceClient