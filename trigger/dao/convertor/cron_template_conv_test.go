package convertor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

func TestCronTemplateEntityTOPo(t *testing.T) {
	tests := []struct {
		name      string
		e         *entity.CrontriggerTemplate
		want      *po.CronTriggerTemplate
		wantErr   assert.ErrorAssertionFunc
		assertion assert.ComparisonAssertionFunc
	}{
		{
			name: "normal",
			e: &entity.CrontriggerTemplate{
				CronExpr:          "0 0 0 * * *",
				Topic:             "test",
				Payload:           []byte("hello"),
				LastExecutionTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            types.TemplateStatusEnable,
			},
			want: &po.CronTriggerTemplate{
				CronExpr:          "0 0 0 * * *",
				Topic:             "test",
				Payload:           []byte("hello"),
				LastExecutionTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            types.TemplateStatusEnable,
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
		{
			name: "CronExpr is empty",
			e: &entity.CrontriggerTemplate{
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            types.TemplateStatusEnable,
			},
			want:      nil,
			wantErr:   assert.Error,
			assertion: assert.Equal,
		},
		{
			name: "Status is empty",
			e: &entity.CrontriggerTemplate{
				CronExpr:          "0 0 0 * * *",
				LastExecutionTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
			},
			want:      nil,
			wantErr:   assert.Error,
			assertion: assert.Equal,
		},
	}
	for _, tt := range tests {
		v, err := CronTemplateEntityTOPo(tt.e)

		tt.wantErr(t, err)
		tt.assertion(t, tt.want, v)
	}
}
