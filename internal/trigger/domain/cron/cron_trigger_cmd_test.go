package cron

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/internal/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

func TestTrigger_Register(t1 *testing.T) {
	type args struct {
		ctx  context.Context
		temp *entity.CronTriggerTemplate
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx: context.Background(),
				temp: &entity.CronTriggerTemplate{
					ID:                0,
					CronExpr:          "*/1 * * * * * *",
					Topic:             "test",
					Payload:           []byte("hello"),
					LoopedTimes:       0,
					ExceptedEndTime:   testTrigger.clock.Now().AddDate(1, 1, 0),
					ExceptedLoopTimes: 10,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.wantErr(t1, testTrigger.Register(tt.args.ctx, tt.args.temp), fmt.Sprintf("Register(%+v, %+v)", tt.args.ctx, tt.args.temp))
			e, err := testTrigger.dao.FindByID(tt.args.ctx, tt.args.temp.ID)
			assert.NoError(t1, err)
			assert.Equal(t1, tt.args.temp, e)
		})
	}
}

func TestTrigger_checkRegisterParams(t1 *testing.T) {
	type args struct {
		temp *entity.CronTriggerTemplate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.CronTriggerTemplate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				&entity.CronTriggerTemplate{
					CronExpr: "*/1 * * * * * *",
					Topic:    "test",
					Payload:  []byte("hello"),
				},
			},
			want: &entity.CronTriggerTemplate{
				CronExpr:          "*/1 * * * * * *",
				Topic:             "test",
				Payload:           []byte("hello"),
				LastExecutionTime: defaultLastExecutionTime,
				ExceptedEndTime:   testTrigger.clock.Now().Add(defaultTemplateActiveDuration),
				ExceptedLoopTimes: defaultMaximumLoopTimes,
				Status:            types.TemplateStatusDisabled,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid cron expression",
			args: args{
				&entity.CronTriggerTemplate{
					CronExpr: "*/70 * * * * * *",
					Topic:    "test",
					Payload:  []byte("hello"),
				},
			},
			want: &entity.CronTriggerTemplate{
				CronExpr: "*/70 * * * * * *",
				Topic:    "test",
				Payload:  []byte("hello"),
			},
			wantErr: assert.Error,
		},
		{
			name: "topic not found",
			args: args{
				&entity.CronTriggerTemplate{
					CronExpr: "*/20 * * * * * *",
					Payload:  []byte("hello"),
				},
			},
			want: &entity.CronTriggerTemplate{
				CronExpr: "*/20 * * * * * *",
				Payload:  []byte("hello"),
			},
			wantErr: assert.Error,
		},
		{
			name: "payload not found",
			args: args{
				&entity.CronTriggerTemplate{
					CronExpr: "*/20 * * * * * *",
					Topic:    "test",
				},
			},
			want: &entity.CronTriggerTemplate{
				CronExpr: "*/20 * * * * * *",
				Topic:    "test",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			err := testTrigger.checkRegisterParams(tt.args.temp)
			tt.wantErr(t1, err, fmt.Errorf("checkRegisterParams got error (%w)", err))
			assert.Equalf(t1, tt.want, tt.args.temp, fmt.Sprintf("checkRegisterParams(%+v)", tt.args.temp))
		})
	}
}

func TestTrigger_UpdateStatus(t1 *testing.T) {
	temp := &entity.CronTriggerTemplate{
		ID:                0,
		CronExpr:          "*/1 * * * * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LoopedTimes:       0,
		ExceptedEndTime:   testTrigger.clock.Now().AddDate(1, 1, 0),
		ExceptedLoopTimes: 10,
	}
	ctx := context.Background()
	testTrigger.Register(ctx, temp)

	type args struct {
		ctx    context.Context
		id     uint
		status int
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx:    nil,
				id:     temp.ID,
				status: types.TemplateStatusEnabled,
			},
			wantErr: assert.NoError,
		},
		{
			name: "normal",
			args: args{
				ctx:    nil,
				id:     temp.ID,
				status: types.TemplateStatusDisabled,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.wantErr(t1, testTrigger.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.status), fmt.Sprintf("UpdateStatus(%+v, %+v, %+v)", tt.args.ctx, tt.args.id, tt.args.status))
			e, _ := testTrigger.dao.FindByID(ctx, tt.args.id)
			assert.Equalf(t1, tt.args.status, e.Status, fmt.Sprintf("get UpdateStatus(%+v)", e))
		})
	}
}
