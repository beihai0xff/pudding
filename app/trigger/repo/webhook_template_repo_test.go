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

func TestWebhookTemplate_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		p   *po.WebhookTriggerTemplate
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
				p: &po.WebhookTriggerTemplate{
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
		err := testWebhookTemplate.Insert(tt.args.ctx, tt.args.p)
		tt.wantErr(t, err)

		res, _ := testWebhookTemplate.FindByID(tt.args.ctx, tt.args.p.ID)
		res.CreatedAt, res.UpdatedAt = tt.args.p.CreatedAt, tt.args.p.UpdatedAt
		tt.assertion(t, tt.args.p, res)
	}
}

func TestWebhookTemplate_Update(t *testing.T) {
	ctx := context.Background()
	p := &po.WebhookTriggerTemplate{
		Topic:             "test",
		Payload:           []byte("hello"),
		DeliverAfter:      10,
		ExceptedEndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		ExceptedLoopTimes: 1,
		LoopedTimes:       1,
		Status:            pb.TriggerStatus_DISABLED,
	}
	_ = testWebhookTemplate.Insert(ctx, p)

	// test set status to enable

	update := &po.WebhookTriggerTemplate{
		Model: gorm.Model{
			ID: p.ID,
		},
		Status: pb.TriggerStatus_ENABLED,
	}
	_, err := testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testWebhookTemplate.FindByID(ctx, p.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_ENABLED)
		p.Status = pb.TriggerStatus_ENABLED
		assert.Equal(t, res.Status, pb.TriggerStatus_ENABLED)
	}

	// test set status to disable
	p.Status, update.Status = pb.TriggerStatus_DISABLED, pb.TriggerStatus_DISABLED
	_, err = testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	if assert.NoError(t, err) {
		res, _ := testWebhookTemplate.FindByID(ctx, p.ID)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
		assert.Equal(t, res.Status, pb.TriggerStatus_DISABLED)
	}

	// test update not exist record
	update = &po.WebhookTriggerTemplate{
		Model: gorm.Model{
			ID: p.ID * 100,
		},
		Status: pb.TriggerStatus_DISABLED,
	}
	_, err = testWebhookTemplate.UpdateStatus(ctx, update.ID, update.Status)
	assert.NoError(t, err)

}
