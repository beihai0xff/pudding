package crond

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/types"
)

type CronD struct {
	s scheduler.Scheduler
}

func (d *CronD) Tracking() {

	t, err := d.getNextTime("")
	if err != nil {
		log.Errorf("failed to get next time, caused by %v", err)
	}
	msg := &types.Message{
		Topic:     "",
		Payload:   nil,
		ReadyTime: t.Unix(),
	}

	_ = d.s.Produce(context.Background(), msg)
}

func (d *CronD) getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(time.Now()), nil
}
