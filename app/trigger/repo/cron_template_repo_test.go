package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
)

func TestCronTemplate_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		p   *po.CronTriggerTemplate
	}
	tests := []struct {
		name      string
		args      args
		wantErr   assert.ErrorAssertionFunc
		assertion assert.ComparisonAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				ctx: context.Background(),
				p: &po.CronTriggerTemplate{
					CronExpr:          "0 0 0 * * *",
					Topic:             "test",
					Payload:           []byte("hello"),
					LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            pb.TriggerStatus_DISABLED,
				},
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},

		{
			name: "no_CronExpr",
			args: args{
				ctx: context.Background(),
				p: &po.CronTriggerTemplate{
					Topic:             "test",
					Payload:           []byte("hello"),
					LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            pb.TriggerStatus_DISABLED,
				},
			},
			wantErr:   assert.Error,
			assertion: assert.NotEqual,
		},
	}
	for _, tt := range tests {
		err := testCronTemplate.Insert(tt.args.ctx, tt.args.p)
		tt.wantErr(t, err)
		if err != nil {
			continue
		}
		res, _ := testCronTemplate.FindByID(tt.args.ctx, tt.args.p.ID)
		res.CreatedAt, res.UpdatedAt = tt.args.p.CreatedAt, tt.args.p.UpdatedAt
		tt.assertion(t, tt.args.p, res)
	}
}

func TestCronTemplate_Update(t *testing.T) {
	ctx := context.Background()
	e := &po.CronTriggerTemplate{
		CronExpr:          "0 0 0 * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 1,
		LoopedTimes:       1,
		Status:            pb.TriggerStatus_DISABLED,
	}
	_ = testCronTemplate.Insert(ctx, e)

	// test set status to enable

	update := &po.CronTriggerTemplate{
		Model: gorm.Model{
			ID: e.ID,
		},
		Status: pb.TriggerStatus_ENABLED,
	}
	_, err := testCronTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testCronTemplate.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_ENABLED)
	}

	// test set status to disable
	e.Status, update.Status = pb.TriggerStatus_DISABLED, pb.TriggerStatus_DISABLED
	_, err = testCronTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testCronTemplate.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
	}

	// test update not exist record
	update = &po.CronTriggerTemplate{
		Model: gorm.Model{
			ID: e.ID + 100,
		},
		Status: pb.TriggerStatus_DISABLED,
	}
	_, err = testCronTemplate.UpdateStatus(ctx, update.ID, update.Status)
	assert.NoError(t, err)

}

func TestCronTemplate_FindEnableRecords(t *testing.T) {
	ctx := context.Background()
	e := &po.CronTriggerTemplate{
		CronExpr:          "0 0 0 * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 10,
		LoopedTimes:       1,
		Status:            pb.TriggerStatus_ENABLED,
	}
	if err := testCronTemplate.Insert(ctx, e); err != nil {
		t.Fatal(err)
	}

	fc := func(e2 *po.CronTriggerTemplate) error {
		e.CreatedAt, e.UpdatedAt = e2.CreatedAt, e2.UpdatedAt
		assert.Equal(t, e, e2)
		e2.LoopedTimes = 2
		return nil
	}
	// test find enable records
	err := testCronTemplate.BatchHandleRecords(ctx, time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local), 10, fc)
	assert.NoError(t, err)

	e3, err := testCronTemplate.FindByID(ctx, e.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), e3.LoopedTimes)
}
