package cron

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/trigger/dao"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

const (
	// messageKeyFormat is the format of cron trigger message key
	messageKeyFormat = "pudding_cron_trigger_template_%d_%d"
)

var (
	// errCronTemplateTopicNotFound is the error of cron template topic is empty
	errCronTemplateTopicNotFound = errors.New("cron template topic not found")
	// errCronTemplatePayloadNotFound is the error of cron template payload is empty
	errCronTemplatePayloadNotFound = errors.New("cron template topic payload not found")
	// errCronTemplateAlreadyEnabled  = errors.New("cron template already enabled")
	// errCronTemplateAlreadyDisabled = errors.New("cron template already disabled")
)

const (
	// defaultMaximumLoopTimes  Maximum Loop Times of Cron Trigger: 1024
	defaultMaximumLoopTimes = 1 << 10
)

type Trigger struct {
	s     scheduler.Scheduler
	dao   dao.CronTemplateDAO
	clock clock.Clock
}

func NewTrigger(db *mysql.Client, s scheduler.Scheduler) *Trigger {
	return &Trigger{
		s:     s,
		dao:   dao.NewCronTemplateDAO(db),
		clock: clock.New(),
	}
}

// Run run cron trigger loop to produce delay message
func (t *Trigger) Run() {
	log.Infof("start produce token")

	now := t.clock.Now()
	timer := time.NewTimer(time.Until(now) + time.Second)

	// wait for the next second
	<-timer.C

	tick := time.NewTicker(1 * time.Second)
	for {

		now := <-tick.C
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		if err := t.dao.FindEnableRecords(ctx, now, 100, t.Tracking); err != nil {
			log.Errorf("failed to find enable cron template, caused by %w", err)
		}

		cancel()
	}
}

// Register register a cron template
func (t *Trigger) Register(ctx context.Context, temp *entity.CronTriggerTemplate) error {
	// 1. check params
	if err := t.checkRegisterParams(temp); err != nil {
		log.Errorf("failed to check params, caused by %w", err)
		return err
	}

	// 2. save the template to db
	if err := t.dao.Insert(ctx, temp); err != nil {
		log.Errorf("failed to insert cron template, caused by %w", err)
		return err
	}

	return nil
}

// checkRegisterParams check the params of register cron template
func (t *Trigger) checkRegisterParams(temp *entity.CronTriggerTemplate) error {
	// 1. check cron expression
	if _, err := cronexpr.Parse(temp.CronExpr); err != nil {
		log.Errorf("Invalid cron expression: %w", err)
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// 2. check topic
	if temp.Topic == "" {
		log.Error(errCronTemplateTopicNotFound.Error())
		return errCronTemplateTopicNotFound
	}

	// 3. check payload
	if len(temp.Payload) == 0 {
		log.Error(errCronTemplatePayloadNotFound.Error())
		return errCronTemplatePayloadNotFound
	}

	// 4. set default value if necessary
	if temp.ExceptedEndTime.IsZero() {
		temp.ExceptedEndTime = t.clock.Now().AddDate(0, 1, 0)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = defaultMaximumLoopTimes
	}
	// default status is Disable
	temp.Status = types.TemplateStatusDisable

	return nil
}

// UpdateStatus update cron template status
func (t *Trigger) UpdateStatus(ctx context.Context, id uint, status int) error {
	// 1. set template status
	temp := &entity.CronTriggerTemplate{
		ID:     id,
		Status: status,
	}

	// 2. update the template status to db
	if err := t.dao.Update(ctx, temp); err != nil {
		log.Errorf("failed to update cron template, caused by %w", err)
		return err
	}

	// 3. tracking the template for first time
	if err := t.Tracking(temp); err != nil {
		log.Errorf("failed to tracking cron template first time, caused by %w", err)
		return err
	}

	return nil
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
	msg := &types.Message{
		Topic:     temp.Topic,
		Key:       t.formatMessageKey(temp),
		Payload:   temp.Payload,
		ReadyTime: nextTime.Unix(),
	}

	if err = t.s.Produce(context.Background(), msg); err != nil {
		log.Errorf("failed to produce message, caused by %v", err)
		return err
	}

	temp.LastExecutionTime = nextTime
	log.Infof("cron template [%d] looped times: %d", temp.ID, temp.LoopedTimes)

	return nil
}

// checkTempShouldRun check whether the template should run
func (t *Trigger) checkTempShouldRun(temp *entity.CronTriggerTemplate, nextTime time.Time) bool {
	temp.LoopedTimes++
	if temp.LoopedTimes > defaultMaximumLoopTimes || temp.LoopedTimes > temp.ExceptedLoopTimes {
		log.Warnf("cron template [%d] has reached the maximum loop times, but it has been tracked", temp.ID)

		temp.Status = types.TemplateStatusMaxTimes
		return false
	}

	// 到达取消执行时间
	if nextTime == temp.LastExecutionTime {
		log.Warnf("cron template [%d] has reached the maximum age, set it to StatusMaxAge", temp.ID)

		temp.Status = types.TemplateStatusMaxAge
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
	return expression.Next(t.clock.Now()), nil
}

// formatMessageKey get cron trigger the message key
func (t *Trigger) formatMessageKey(temp *entity.CronTriggerTemplate) string {
	return fmt.Sprintf(messageKeyFormat, temp.ID, temp.LoopedTimes)
}
