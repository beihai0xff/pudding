package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/trigger/entity"
)

func TestTrigger_formatMessageKey(t *testing.T) {
	type args struct {
		temp *entity.CronTriggerTemplate
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{&entity.CronTriggerTemplate{ID: 1, LoopedTimes: 10}}, "pudding_cron_trigger_template_1_10"},
		{"normal", args{&entity.CronTriggerTemplate{ID: 100, LoopedTimes: 100}}, "pudding_cron_trigger_template_100_100"},
		{"normal", args{&entity.CronTriggerTemplate{ID: 1000, LoopedTimes: 1000}}, "pudding_cron_trigger_template_1000_1000"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, testTrigger.formatMessageKey(tt.args.temp))
	}
}

func TestTrigger_getNextTime(t1 *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr assert.ErrorAssertionFunc
	}{
		{"every 1 Second", args{"*/1 * * * * * *"}, testTrigger.clock.Now().Add(time.Second), assert.NoError},
		{"every 1 Minute", args{"* */1 * * * *"}, testTrigger.clock.Now().Add(time.Minute), assert.NoError},
		{"every 1 Hour", args{"0 0 * * * * *"}, testTrigger.clock.Now().Add(time.Hour), assert.NoError},
		{"every 1 Day", args{"0 0 */1 * *"}, testTrigger.clock.Now().AddDate(0, 0, 1), assert.NoError},
		{"every 1 Month", args{"0 0 1 * *"}, testTrigger.clock.Now().AddDate(0, 1, 0), assert.NoError},
		{"every 1 Year", args{"0 0 1 1 *"}, testTrigger.clock.Now().AddDate(1, 0, 0), assert.NoError},
		{"error expr", args{"0 0 0 31 */1"}, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), assert.Error},
		{"every five Second", args{"*/5 * * * * * *"}, testTrigger.clock.Now().Add(5 * time.Second), assert.NoError},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := testTrigger.getNextTime(tt.args.expr)
			if !tt.wantErr(t1, err, fmt.Sprintf("getNextTime(%v)", tt.args.expr)) {
				return
			}
			assert.Equalf(t1, tt.want, got, "getNextTime(%v)", tt.args.expr)
		})
	}
}
