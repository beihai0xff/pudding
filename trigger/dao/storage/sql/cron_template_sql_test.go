package sql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

func TestCronTemplate_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		e   *entity.CronTriggerTemplate
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
				e: &entity.CronTriggerTemplate{
					CronExpr:          "0 0 0 * * *",
					Topic:             "test",
					Payload:           []byte("hello"),
					LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            types.TemplateStatusDisable,
				},
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},

		{
			name: "no_CronExpr",
			args: args{
				ctx: context.Background(),
				e: &entity.CronTriggerTemplate{
					Topic:             "test",
					Payload:           []byte("hello"),
					LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            types.TemplateStatusDisable,
				},
			},
			wantErr:   assert.Error,
			assertion: assert.NotEqual,
		},
	}
	for _, tt := range tests {
		err := cronTemplateSQL.Insert(tt.args.ctx, tt.args.e)
		tt.wantErr(t, err)

		res, _ := cronTemplateSQL.FindByID(tt.args.ctx, tt.args.e.ID)
		tt.assertion(t, tt.args.e, res)
	}
}

func TestCronTemplate_Update(t *testing.T) {
	ctx := context.Background()
	e := &entity.CronTriggerTemplate{
		CronExpr:          "0 0 0 * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 1,
		LoopedTimes:       1,
		Status:            types.TemplateStatusDisable,
	}
	_ = cronTemplateSQL.Insert(ctx, e)

	// test set status to enable

	update := &entity.CronTriggerTemplate{
		ID:     e.ID,
		Status: types.TemplateStatusEnable,
	}
	err := cronTemplateSQL.Update(ctx, update)
	if assert.NoError(t, err) {
		res, _ := cronTemplateSQL.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, types.TemplateStatusEnable)
		e.Status = types.TemplateStatusEnable
		assert.Equal(t, res.Status, types.TemplateStatusEnable)
	}

	// test set status to disable
	e.Status, update.Status = types.TemplateStatusDisable, types.TemplateStatusDisable
	err = cronTemplateSQL.Update(ctx, update)
	if assert.NoError(t, err) {
		res, _ := cronTemplateSQL.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, types.TemplateStatusDisable)
		assert.Equal(t, res.Status, types.TemplateStatusDisable)
	}

	// test update not exist record
	update = &entity.CronTriggerTemplate{
		ID:     e.ID * 100,
		Status: types.TemplateStatusDisable,
	}
	err = cronTemplateSQL.Update(ctx, update)
	assert.NoError(t, err)

}

func TestCronTemplate_FindEnableRecords(t *testing.T) {
	ctx := context.Background()
	e := &entity.CronTriggerTemplate{
		CronExpr:          "0 0 0 * * *",
		Topic:             "test",
		Payload:           []byte("hello"),
		LastExecutionTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 10,
		LoopedTimes:       1,
		Status:            types.TemplateStatusEnable,
	}
	if err := cronTemplateSQL.Insert(ctx, e); err != nil {
		t.Fatal(err)
	}

	fc := types.CronTempHandler(
		func(e2 *entity.CronTriggerTemplate) error {
			assert.Equal(t, e, e2)
			e2.LoopedTimes = 2
			return nil
		})
	// test find enable records
	err := cronTemplateSQL.FindEnableRecords(ctx, time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local), 10, fc)
	assert.NoError(t, err)

	e3, err := cronTemplateSQL.FindByID(ctx, e.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), e3.LoopedTimes)
}
