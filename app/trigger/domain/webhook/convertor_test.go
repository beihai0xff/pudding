package webhook

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
)

func TestWebhookTemplateEntityTOPo(t *testing.T) {
	tests := []struct {
		name      string
		e         *TriggerTemplate
		want      *po.WebhookTriggerTemplate
		wantErr   assert.ErrorAssertionFunc
		assertion assert.ComparisonAssertionFunc
	}{
		{
			name: "normal",
			e: &TriggerTemplate{
				Topic:             "test",
				Payload:           []byte("hello"),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            pb.TriggerStatus_ENABLED,
			},
			want: &po.WebhookTriggerTemplate{
				Topic:             "test",
				Payload:           []byte("hello"),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            pb.TriggerStatus_ENABLED,
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
		{
			name: "Status is empty",
			e: &TriggerTemplate{
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
			},
			want: &po.WebhookTriggerTemplate{
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
	}
	for _, tt := range tests {
		v, err := convEntityTOPo(tt.e)

		tt.wantErr(t, err)
		tt.assertion(t, tt.want, v)
	}
}

func TestWebhookTemplatePoTOEntity(t *testing.T) {
	tests := []struct {
		name      string
		p         *po.WebhookTriggerTemplate
		want      *TriggerTemplate
		wantErr   assert.ErrorAssertionFunc
		assertion assert.ComparisonAssertionFunc
	}{
		{
			name: "normal",
			p: &po.WebhookTriggerTemplate{
				Model:             gorm.Model{ID: 10},
				Topic:             "test",
				Payload:           []byte("hello"),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            pb.TriggerStatus_ENABLED,
			},
			want: &TriggerTemplate{
				ID:                10,
				Topic:             "test",
				Payload:           []byte("hello"),
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
				Status:            pb.TriggerStatus_ENABLED,
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
		{
			name: "Status is empty",
			p: &po.WebhookTriggerTemplate{
				Model:             gorm.Model{ID: 10},
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
			},
			want: &TriggerTemplate{
				ID:                10,
				ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				ExceptedLoopTimes: 1,
				LoopedTimes:       1,
			},
			wantErr:   assert.NoError,
			assertion: assert.Equal,
		},
	}

	for _, tt := range tests {
		v, err := convPoTOEntity(tt.p)

		tt.wantErr(t, err)
		tt.assertion(t, tt.want, v)
	}
}

func TestWebhookTemplateSlicePoTOEntity(t *testing.T) {
	type args struct {
		p []*po.WebhookTriggerTemplate
	}
	tests := []struct {
		name    string
		args    args
		want    []*TriggerTemplate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				p: []*po.WebhookTriggerTemplate{
					{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
							DeletedAt: gorm.DeletedAt{},
						},
						Topic:             "test",
						Payload:           []byte("hello"),
						ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						ExceptedLoopTimes: 1,
						LoopedTimes:       1,
						Status:            1,
					},
					{
						Model: gorm.Model{
							ID:        2,
							CreatedAt: time.Time{},
							UpdatedAt: time.Time{},
							DeletedAt: gorm.DeletedAt{},
						},
						Topic:             "test",
						Payload:           []byte("hello"),
						ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						ExceptedLoopTimes: 1,
						LoopedTimes:       1,
						Status:            1,
					},
				},
			},
			want: []*TriggerTemplate{
				{
					ID:                1,
					Topic:             "test",
					Payload:           []byte("hello"),
					ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            1,
				},
				{
					ID:                2,
					Topic:             "test",
					Payload:           []byte("hello"),
					ExceptedEndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					ExceptedLoopTimes: 1,
					LoopedTimes:       1,
					Status:            1,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convSlicePoTOEntity(tt.args.p)
			if !tt.wantErr(t, err, fmt.Sprintf("CronTemplateSlicePoTOEntity(%v)", tt.args.p)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CronTemplateSlicePoTOEntity(%v)", tt.args.p)
		})
	}
}
