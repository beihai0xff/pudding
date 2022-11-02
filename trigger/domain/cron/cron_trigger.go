package cron

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/trigger/dao"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

const (
	// defaultMaximumLoopTimes  Maximum Loop Times of Cron Trigger: 1024
	defaultMaximumLoopTimes = 1 << 10
)

type CronTrigger struct {
	s   scheduler.Scheduler
	dao dao.CronTemplateDAO
}

func NewCron(db *mysql.Client, s scheduler.Scheduler) *CronTrigger {
	return &CronTrigger{
		s:   s,
		dao: dao.NewCronTemplateDAO(db),
	}
}

func (t *CronTrigger) Tracking(temp *entity.CrontriggerTemplate) error {
	if temp.LoopedTimes > defaultMaximumLoopTimes {
		temp.Status = types.TemplateStatusMaxTimes
	}

	nextTime, err := t.getNextTime(temp.CronExpr)
	if err != nil {
		log.Errorf("failed to get next time, caused by %v", err)
		return err
	}

	// 到达取消执行时间
	if nextTime == temp.LastExecutionTime {
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

	temp.LoopedTimes++

	if temp.LoopedTimes > defaultMaximumLoopTimes {
		temp.Status = types.TemplateStatusMaxTimes
	}

	return nil
}

func (t *CronTrigger) getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(time.Now()), nil
}
