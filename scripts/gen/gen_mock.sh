#!/bin/sh

mockgen -destination=../../test/mock/scheduler_mock.go --package=mock github.com/beihai0xff/pudding/app/scheduler Scheduler
mockgen -destination=../../../test/mock/queue_mock.go --package=mock github.com/beihai0xff/pudding/app/scheduler/broker DelayQueue,RealTimeQueue