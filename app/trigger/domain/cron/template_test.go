package cron

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
)

func TestTrigger_Register(t1 *testing.T) {
	type args struct {
		ctx  context.Context
		temp *TriggerTemplate
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
				temp: &TriggerTemplate{
					ID:                0,
					CronExpr:          "*/1 * * * * * *",
					Topic:             "test",
					Payload:           []byte("hello"),
					LoopedTimes:       0,
					ExceptedEndTime:   testTrigger.wallClock.Now().AddDate(1, 1, 0),
					ExceptedLoopTimes: 10,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.wantErr(t1, testTrigger.Register(tt.args.ctx, tt.args.temp), fmt.Sprintf("RegisterGRPC(%+v, %+v)", tt.args.ctx, tt.args.temp))
			p, err := testTrigger.repo.FindByID(tt.args.ctx, tt.args.temp.ID)
			assert.NoError(t1, err)
			e, _ := convPoTOEntity(p)
			assert.Equal(t1, tt.args.temp, e)
		})
	}
}

func TestTrigger_checkRegisterParams(t1 *testing.T) {
	type args struct {
		temp *TriggerTemplate
	}
	tests := []struct {
		name    string
		args    args
		want    *TriggerTemplate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				&TriggerTemplate{
					CronExpr: "*/1 * * * * * *",
					Topic:    "test",
					Payload:  []byte("hello"),
				},
			},
			want: &TriggerTemplate{
				CronExpr:          "*/1 * * * * * *",
				Topic:             "test",
				Payload:           []byte("hello"),
				LastExecutionTime: defaultLastExecutionTime,
				ExceptedEndTime:   testTrigger.wallClock.Now().Add(defaultTemplateActiveDuration),
				ExceptedLoopTimes: defaultMaximumLoopTimes,
				Status:            pb.TriggerStatus_DISABLED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid cron expression",
			args: args{
				&TriggerTemplate{
					CronExpr: "*/70 * * * * * *",
					Topic:    "test",
					Payload:  []byte("hello"),
				},
			},
			want: &TriggerTemplate{
				CronExpr: "*/70 * * * * * *",
				Topic:    "test",
				Payload:  []byte("hello"),
			},
			wantErr: assert.Error,
		},
		{
			name: "topic not found",
			args: args{
				&TriggerTemplate{
					CronExpr: "*/20 * * * * * *",
					Payload:  []byte("hello"),
				},
			},
			want: &TriggerTemplate{
				CronExpr: "*/20 * * * * * *",
				Payload:  []byte("hello"),
			},
			wantErr: assert.Error,
		},
		{
			name: "payload not found",
			args: args{
				&TriggerTemplate{
					CronExpr: "*/20 * * * * * *",
					Topic:    "test",
				},
			},
			want: &TriggerTemplate{
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
	temp := &TriggerTemplate{
		ID:                0,
		CronExpr:          "*/1 * * * * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LoopedTimes:       0,
		ExceptedEndTime:   testTrigger.wallClock.Now().AddDate(1, 1, 0),
		ExceptedLoopTimes: 10,
	}
	ctx := context.Background()
	testTrigger.Register(ctx, temp)

	type args struct {
		ctx    context.Context
		id     uint
		status pb.TriggerStatus
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
				status: pb.TriggerStatus_ENABLED,
			},
			wantErr: assert.NoError,
		},
		{
			name: "normal",
			args: args{
				ctx:    nil,
				id:     temp.ID,
				status: pb.TriggerStatus_DISABLED,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, err := testTrigger.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.status)
			tt.wantErr(t1, err, fmt.Sprintf("UpdateStatus(%+v, %+v, %+v)", tt.args.ctx, tt.args.id, tt.args.status))
			e, _ := testTrigger.repo.FindByID(ctx, tt.args.id)
			assert.Equalf(t1, tt.args.status, e.Status, fmt.Sprintf("get UpdateStatus(%+v)", e))
		})
	}
}
