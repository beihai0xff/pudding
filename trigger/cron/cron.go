package cron

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/trigger/cron/domain/entity"
	"github.com/beihai0xff/pudding/types"
)

const (
	defaultMaximumLoopTimes = 1 << 10
)

type Cron struct {
	s scheduler.Scheduler
}

func (c *Cron) Tracking(temp *entity.CrontriggerTemplate) error {
	if temp.LoopedTimes > defaultMaximumLoopTimes {
		temp.Status = types.TemplateStatusMaxTimes
	}

	t, err := c.getNextTime(temp.CronExpr)
	if err != nil {
		log.Errorf("failed to get next time, caused by %v", err)
		return err
	}

	// 到达取消执行时间
	if t == temp.LastExecutionTime {
		temp.Status = types.TemplateStatusMaxAge
		return nil
	}

	msg := &types.Message{
		Topic:     temp.Topic,
		Payload:   temp.Payload,
		ReadyTime: t.Unix(),
	}

	if err = c.s.Produce(context.Background(), msg); err != nil {
		log.Errorf("failed to produce message, caused by %v", err)
		return err
	}

	temp.LoopedTimes++

	if temp.LoopedTimes > defaultMaximumLoopTimes {
		temp.Status = types.TemplateStatusMaxTimes
	}

	return nil
}

func (c *Cron) getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(time.Now()), nil
}
