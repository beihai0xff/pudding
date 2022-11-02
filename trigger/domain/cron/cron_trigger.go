package cron

import (
	"context"
	"errors"
	"time"

	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/trigger/dao"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

var (
	errCronTemplateTopicNotFound   = errors.New("cron template topic not found")
	errCronTemplatePayloadNotFound = errors.New("cron template topic payload not found")
	errCronTemplateAlreadyEnabled  = errors.New("cron template already enabled")
	errCronTemplateAlreadyDisabled = errors.New("cron template already disabled")
)

const (
	// defaultMaximumLoopTimes  Maximum Loop Times of Cron Trigger: 1024
	defaultMaximumLoopTimes = 1 << 10
)

type Trigger struct {
	s   scheduler.Scheduler
	dao dao.CronTemplateDAO
}

func NewTrigger(db *mysql.Client, s scheduler.Scheduler) *Trigger {
	return &Trigger{
		s:   s,
		dao: dao.NewCronTemplateDAO(db),
	}
}

func (t *Trigger) Register(ctx context.Context, temp *entity.CrontriggerTemplate) error {
	// 1. check params
	if err := t.checkParams(temp); err != nil {
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

func (t *Trigger) checkParams(temp *entity.CrontriggerTemplate) error {
	// 1. check cron expression
	if _, err := cronexpr.Parse(temp.CronExpr); err != nil {
		log.Errorf("Invalid cron expression: %w", err)
		return err
	}

	// 2. check topic
	if temp.Topic == "" {
		log.Error("Cron Template Topic can not be empty")
		return errCronTemplateTopicNotFound
	}

	// 3. check payload
	if len(temp.Payload) == 0 {
		log.Error("Cron Template Payload can not be empty")
		return errCronTemplatePayloadNotFound
	}

	// 4. set default value if necessary
	if temp.ExceptedEndTime.IsZero() {
		temp.ExceptedEndTime = time.Now().AddDate(1, 0, 0)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = defaultMaximumLoopTimes
	}
	// default status is Disable
	temp.Status = types.TemplateStatusDisable

	return nil
}

func (t *Trigger) Enable(ctx context.Context, temp *entity.CrontriggerTemplate) error {
	// 1. set status to enable
	if temp.Status == types.TemplateStatusEnable {
		return errCronTemplateAlreadyEnabled
	}
	temp.Status = types.TemplateStatusEnable

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

func (t *Trigger) Disable(ctx context.Context, temp *entity.CrontriggerTemplate) error {
	// 1. set status to disable
	if temp.Status == types.TemplateStatusDisable {
		return errCronTemplateAlreadyDisabled
	}
	temp.Status = types.TemplateStatusDisable

	// 2. update the template status to db
	if err := t.dao.Update(ctx, temp); err != nil {
		log.Errorf("failed to update cron template, caused by %w", err)
		return err
	}

	return nil
}

func (t *Trigger) Tracking(temp *entity.CrontriggerTemplate) error {
	if temp.LoopedTimes > defaultMaximumLoopTimes {
		log.Warnf("cron template [%d] has reached the maximum loop times, but it has been tracked", temp.ID)

		temp.Status = types.TemplateStatusMaxTimes
		return nil
	}

	nextTime, err := t.getNextTime(temp.CronExpr)
	if err != nil {
		log.Errorf("failed to get next time, caused by %v", err)
		return err
	}

	// 到达取消执行时间
	if nextTime == temp.LastExecutionTime {
		log.Infof("cron template [%d] has reached the maximum age, set it to StatusMaxAge", temp.ID)

		temp.Status = types.TemplateStatusMaxAge
		return nil
	}

	msg := &types.Message{
		Topic:     temp.Topic,
		Payload:   temp.Payload,
		ReadyTime: nextTime.Unix(),
	}

	if err = t.s.Produce(context.Background(), msg); err != nil {
		log.Errorf("failed to produce message, caused by %v", err)
		return err
	}

	temp.LastExecutionTime = nextTime
	temp.LoopedTimes++
	log.Debugf("cron template [%d] looped times: %d", temp.ID, temp.LoopedTimes)

	if temp.LoopedTimes > defaultMaximumLoopTimes {
		log.Infof("cron template [%d] has reached the maximum loop times", temp.ID)
		temp.Status = types.TemplateStatusMaxTimes

	}

	return nil
}

func (t *Trigger) getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(time.Now()), nil
}
