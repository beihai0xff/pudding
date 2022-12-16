package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

const (
	// messageKeyFormat is the format of cron trigger message key
	messageKeyFormat = "pudding_cron_trigger_template_%d_%d"

	// defaultMaximumLoopTimes  Maximum Loop Times of Cron Trigger: 1024
	defaultMaximumLoopTimes = 1 << 10
	// defaultTemplateActiveDuration is the default active duration of cron template: 30 days
	defaultTemplateActiveDuration = 30 * 24 * time.Hour
)

var (
	// errCronTemplateTopicNotFound is the error of cron template topic is empty
	errCronTemplateTopicNotFound = errors.New("cron template topic not found")
	// errCronTemplatePayloadNotFound is the error of cron template payload is empty
	errCronTemplatePayloadNotFound = errors.New("cron template topic payload not found")
	// errCronTemplateAlreadyEnabled  = errors.New("cron template already enabled")
	// errCronTemplateAlreadyDisabled = errors.New("cron template already disabled")

	// defaultLastExecutionTime is the default last execution time for New Registered cron
	defaultLastExecutionTime = time.Unix(1, 0).UTC()
)

type Trigger struct {
	schedulerClient scheduler.SchedulerServiceClient
	repo            repo.CronTemplateDAO
	// wallClock is the clock used to get current time
	wallClock clock.Clock
}

func NewTrigger(db *mysql.Client, client scheduler.SchedulerServiceClient) *Trigger {
	return &Trigger{
		schedulerClient: client,
		repo:            repo.NewCronTemplate(db),
		wallClock:       clock.New(),
	}
}

// Run run cron trigger loop to produce delay message
func (t *Trigger) Run() {
	log.Infof("start produce token")

	now := t.wallClock.Now()
	timer := time.NewTimer(time.Until(now) + time.Second)

	// wait for the next second
	<-timer.C

	tick := time.NewTicker(1 * time.Second)
	for {

		now := <-tick.C
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		if err := t.repo.BatchHandleRecords(ctx, now, 100, t.Tracking); err != nil {
			log.Errorf("failed to find enable cron template, caused by %v", err)
		}

		cancel()
	}
}

// Tracking try to produce Cron Trigger Message
func (t *Trigger) Tracking(temp *entity.CronTriggerTemplate) error {
	nextTime, err := t.getNextTime(temp.CronExpr)
	if err != nil {
		log.Errorf("failed to get next time, caused by %v", err)
		return err
	}

	if !t.checkTempShouldRun(temp, nextTime) {
		return nil
	}

	// produce the message
	msg := &scheduler.SendDelayMessageRequest{
		Topic:     temp.Topic,
		Key:       t.formatMessageKey(temp),
		Payload:   temp.Payload,
		DeliverAt: nextTime.Unix(),
	}

	if _, err = t.schedulerClient.SendDelayMessage(context.Background(), msg); err != nil {
		log.Errorf("failed to send DelayMessage, caused by %v", err)
		return err
	}

	temp.LastExecutionTime = nextTime
	temp.LoopedTimes++
	if temp.LoopedTimes == temp.ExceptedLoopTimes {
		log.Infof("cron template [%d] has reached the maximum loop times, "+
			"update status to TemplateStatusMaxTimes", temp.ID)

		temp.Status = pb.TriggerStatus_MAX_TIMES
	}
	log.Infof("cron template [%d] looped times: %d", temp.ID, temp.LoopedTimes)

	return nil
}

// checkTempShouldRun check whether the template should run
func (t *Trigger) checkTempShouldRun(temp *entity.CronTriggerTemplate, nextTime time.Time) bool {
	if temp.LoopedTimes >= temp.ExceptedLoopTimes {
		log.Warnf("cron template [%d] has reached the maximum loop times, but it has been tracked", temp.ID)

		temp.Status = pb.TriggerStatus_MAX_TIMES
		return false
	}

	// 到达取消执行时间
	if nextTime.After(temp.ExceptedEndTime) {
		log.Warnf("cron template [%d] has reached the maximum age, set it to StatusMaxAge", temp.ID)

		temp.Status = pb.TriggerStatus_MAX_AGE
		return false
	}

	return true
}

// getNextTime get the next time of cron expression
func (t *Trigger) getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(t.wallClock.Now()), nil
}

// formatMessageKey get cron trigger the message key
func (t *Trigger) formatMessageKey(temp *entity.CronTriggerTemplate) string {
	return fmt.Sprintf(messageKeyFormat, temp.ID, temp.LoopedTimes)
}
