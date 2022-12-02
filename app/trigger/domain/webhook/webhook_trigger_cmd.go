package webhook

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

var (
	// errWebhookTemplateTopicNotFound is the error of webhook template topic is empty
	errWebhookTemplateTopicNotFound = errors.New("webhook template topic not found")
	// errWebhookTemplatePayloadNotFound is the error of webhook template payload is empty
	errWebhookTemplatePayloadNotFound = errors.New("webhook template topic payload not found")
)

type Trigger struct {
	repo repo.WebhookTemplate
	// wallClock is the clock used to get current time
	wallClock clock.Clock
}

func NewTrigger(db *mysql.Client) *Trigger {
	return &Trigger{
		repo:      repo.NewWebhookTemplate(db),
		wallClock: clock.New(),
	}
}

// FindByID find webhook template by id
func (t *Trigger) FindByID(ctx context.Context, id uint) (*entity.WebhookTriggerTemplate, error) {
	if id <= 0 {
		err := fmt.Errorf("invalid id, id: %d", id)
		log.Errorf("%v", err)
		return nil, err
	}

	res, err := t.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("failed to find webhook template, caused by %v", err)
		return nil, err
	}

	return res, nil
}

// PageQuery page query webhook templates
func (t *Trigger) PageQuery(ctx context.Context, p *entity.PageQuery,
	status pb.TriggerStatus) ([]*entity.WebhookTriggerTemplate, int64, error) {
	// check params
	if p.Offset < 0 || p.Limit <= 0 {
		err := fmt.Errorf("invalid offset or limit, offset: %d, limit: %d", p.Offset, p.Limit)
		log.Errorf("%v", err)
		return nil, 0, err
	}

	res, count, err := t.repo.PageQuery(ctx, p, status)
	if err != nil {
		log.Errorf("failed to PageQuery cron template, caused by %v", err)
		return nil, 0, err
	}

	return res, count, nil
}

// Register create a webhook template
func (t *Trigger) Register(ctx context.Context, temp *entity.WebhookTriggerTemplate) error {
	// 1. check params
	if err := t.checkRegisterParams(temp); err != nil {
		log.Errorf("failed to check params, caused by %v", err)
		return err
	}

	// 2. save the template to db
	if err := t.repo.Insert(ctx, temp); err != nil {
		log.Errorf("failed to insert cron template, caused by %v", err)
		return err
	}

	return nil
}

// checkRegisterParams check the params of register webhook template
func (t *Trigger) checkRegisterParams(temp *entity.WebhookTriggerTemplate) error {
	// 1. check topic
	if temp.Topic == "" {
		log.Error(errWebhookTemplateTopicNotFound.Error())
		return errWebhookTemplateTopicNotFound
	}

	// 2. check payload
	if len(temp.Payload) == 0 {
		log.Error(errWebhookTemplatePayloadNotFound.Error())
		return errWebhookTemplatePayloadNotFound
	}

	// 3. set default value if necessary
	if temp.ExceptedEndTime.IsZero() {
		temp.ExceptedEndTime = t.wallClock.Now().Add(entity.DefaultTemplateActiveDuration)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = entity.DefaultMaximumLoopTimes
	}

	// default status is Disable
	temp.Status = pb.TriggerStatus_DISABLED

	return nil
}

// UpdateStatus update webhook template status
func (t *Trigger) UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error) {
	// 1. update the template status to db
	rowsAffected, err := t.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Errorf("failed to update webhook template, caused by %v", err)
		return 0, err
	}

	return rowsAffected, nil
}
