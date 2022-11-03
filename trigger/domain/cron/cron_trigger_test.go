package cron

import (
	"testing"

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
		assert.Equal(t, tt.want, test_trigger.formatMessageKey(tt.args.temp))
	}
}
