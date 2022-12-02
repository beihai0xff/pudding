// package repo is a package for the repo layer.
// It contains the repository interfaces and implementations.

package repo

import (
	"context"
	"fmt"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo/convertor"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/sql"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
)

// WebhookTemplate is the interface for the webhook template repository.
type WebhookTemplate interface {
	// FindByID find a cron template by id
	FindByID(ctx context.Context, id uint) (*entity.WebhookTriggerTemplate, error)
	// PageQuery query cron templates by page
	PageQuery(ctx context.Context, p *entity.PageQuery, status pb.TriggerStatus) (res []*entity.WebhookTriggerTemplate,
		count int64, err error)

	// Insert create a cron template
	Insert(ctx context.Context, e *entity.WebhookTriggerTemplate) error
	// UpdateStatus update the status of a cron template
	UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error)
}

// webhookTemplate is the implementation of WebhookTemplate
type webhookTemplate struct{}

// NewWebhookTemplate create a new Webhook template repository
func NewWebhookTemplate(db *mysql.Client) WebhookTemplate {
	sql.SetDefault(db.GetDB())
	return &webhookTemplate{}
}

// FindByID find a Webhook template by id
func (dao *webhookTemplate) FindByID(ctx context.Context, id uint) (*entity.WebhookTriggerTemplate, error) {
	// SELECT * FROM pudding_webhook_trigger_template WHERE id =
	res, err := sql.WebhookTriggerTemplate.WithContext(ctx).FindByID(id)
	if err != nil {
		return nil, err
	}

	e, err := convertor.WebhookTemplatePoTOEntity(res)
	if err != nil {
		return nil, err
	}
	return e, nil

}

// PageQuery query Webhook templates by page
func (dao *webhookTemplate) PageQuery(ctx context.Context, p *entity.PageQuery, status pb.TriggerStatus) (
	[]*entity.WebhookTriggerTemplate, int64, error) {

	var res []*po.WebhookTriggerTemplate
	var count int64
	var err error
	if status > pb.TriggerStatus_UNKNOWN_UNSPECIFIED && status <= pb.TriggerStatus_MAX_AGE {
		res, count, err = sql.WebhookTriggerTemplate.WithContext(ctx).
			Where(sql.WebhookTriggerTemplate.Status.Eq(int32(status))).FindByPage(p.Offset, p.Limit)
	} else {
		res, count, err = sql.WebhookTriggerTemplate.WithContext(ctx).FindByPage(p.Offset, p.Limit)
	}

	if err != nil {
		return nil, 0, err
	}

	e, err := convertor.WebhookTemplateSlicePoTOEntity(res)
	if err != nil {
		return nil, 0, err
	}
	return e, count, nil

}

// Insert create a Webhook template
func (dao *webhookTemplate) Insert(ctx context.Context, e *entity.WebhookTriggerTemplate) error {
	if err := validate.Struct(e); err != nil {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	p, err := convertor.WebhookTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	if err := sql.WebhookTriggerTemplate.WithContext(ctx).Create(p); err != nil {
		return err
	}
	e.ID = p.ID

	return nil
}

// UpdateStatus update the status of a Webhook template
func (dao *webhookTemplate) UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error) {
	return sql.WebhookTriggerTemplate.WithContext(ctx).UpdateStatus(ctx, id, status)
}
