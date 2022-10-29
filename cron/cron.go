package cron

import (
	"time"

	"github.com/beihai0xff/pudding/pkg/cronexpr"
)

func getNextTime(expr string) (time.Time, error) {
	expression, err := cronexpr.Parse(expr)
	if err != nil {
		return time.Time{}, err
	}
	return expression.Next(time.Now()), nil
}
