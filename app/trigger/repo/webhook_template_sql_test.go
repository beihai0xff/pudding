package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
)

func TestWebhookTemplate_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		e   *entity.WebhookTriggerTemplate
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
				e: &entity.WebhookTriggerTemplate{
					Topic:             "test",
					Payload:           []byte("hello"),
					DeliverAfter:      10,
					ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            pb.TriggerStatus_DISABLED,
				},
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
	}
	for _, tt := range tests {
		err := testWebhookTemplate.Insert(tt.args.ctx, tt.args.e)
		tt.wantErr(t, err)

		res, _ := testWebhookTemplate.FindByID(tt.args.ctx, tt.args.e.ID)
		tt.assertion(t, tt.args.e, res)
	}
}

func TestWebhookTemplate_Update(t *testing.T) {
	ctx := context.Background()
	e := &entity.WebhookTriggerTemplate{
		Topic:             "test",
		Payload:           []byte("hello"),
		DeliverAfter:      10,
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 1,
		LoopedTimes:       1,
		Status:            pb.TriggerStatus_DISABLED,
	}
	_ = testWebhookTemplate.Insert(ctx, e)

	// test set status to enable

	update := &entity.WebhookTriggerTemplate{
		ID:     e.ID,
		Status: pb.TriggerStatus_ENABLED,
	}
	_, err := testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testWebhookTemplate.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_ENABLED)
		e.Status = pb.TriggerStatus_ENABLED
		assert.Equal(t, res.Status, pb.TriggerStatus_ENABLED)
	}

	// test set status to disable
	e.Status, update.Status = pb.TriggerStatus_DISABLED, pb.TriggerStatus_DISABLED
	_, err = testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testWebhookTemplate.FindByID(ctx, e.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
	}

	// test update not exist record
	update = &entity.WebhookTriggerTemplate{
		ID:     e.ID * 100,
		Status: pb.TriggerStatus_DISABLED,
	}
	_, err = testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	assert.NoError(t, err)

}
