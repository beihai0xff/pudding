// Package webhook implemented the webhook trigger and handler
// template.go implements the webhook template
package webhook

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/pkg/configs"
	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

var (
	validate = validator.New()

	// errWebhookTemplateTopicNotFound is the error of webhook template topic is empty
	errWebhookTemplateTopicNotFound = errors.New("webhook template topic not found")
	// errWebhookTemplatePayloadNotFound is the error of webhook template payload is empty
	errWebhookTemplatePayloadNotFound = errors.New("webhook template topic payload not found")
)

const webhookURL = "%s/pudding/trigger/webhook/v1/call/%d"

// Trigger is the webhook trigger
type Trigger struct {
	webhookPrefix string

	schedulerClient broker.SchedulerServiceClient
	repo            repo.WebhookTemplate
	// wallClock is the clock used to get current time
	wallClock clock.Clock
}

// NewTrigger create a webhook trigger
func NewTrigger(db *mysql.Client, client broker.SchedulerServiceClient) *Trigger {
	return &Trigger{
		webhookPrefix:   configs.GetWebhookPrefix(),
		schedulerClient: client,
		repo:            repo.NewWebhookTemplate(db),
		wallClock:       clock.New(),
	}
}

// FindByID find webhook template by id
func (t *Trigger) FindByID(ctx context.Context, id uint) (*TriggerTemplate, error) {
	if id <= 0 {
		err := fmt.Errorf("invalid id, id: %d", id)
		log.Errorf("%v", err)
		return nil, err
	}

	po, err := t.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("failed to find webhook template, caused by %v", err)
		return nil, err
	}

	res, err := convPoTOEntity(po)
	if err != nil {
		return nil, fmt.Errorf("convert po to entity failed: %w", err)
	}

	return res, nil
}

// PageQuery page query webhook templates
func (t *Trigger) PageQuery(ctx context.Context, p *entity.PageQuery,
	status pb.TriggerStatus) ([]*TriggerTemplate, int64, error) {
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

	e, err := convSlicePoTOEntity(res)
	if err != nil {
		return nil, 0, fmt.Errorf("convert po to entity failed: %w", err)
	}
	return e, count, nil
}

// Register create a webhook template
func (t *Trigger) Register(ctx context.Context, temp *TriggerTemplate) error {
	// 1. check params
	if err := t.checkRegisterParams(temp); err != nil {
		log.Errorf("failed to check params, caused by %v", err)
		return err
	}
	if err := validate.Struct(temp); err != nil {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	// 2. save the template to db
	p, err := convEntityTOPo(temp)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	if err := t.repo.Insert(ctx, p); err != nil {
		log.Errorf("failed to insert cron template, caused by %v", err)
		return err
	}

	temp.ID = p.ID
	return nil
}

// checkRegisterParams check the params of register webhook template
func (t *Trigger) checkRegisterParams(temp *TriggerTemplate) error {
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

// genWebhookURL generate webhook URL by trigger template id
func (t *Trigger) genWebhookURL(id uint) string {
	return fmt.Sprintf(webhookURL, t.webhookPrefix, id)
}

// Call trigger a webhook by id
func (t *Trigger) Call(ctx context.Context, id uint) (string, error) {
	template, err := t.FindByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to find webhook template, caused by %w", err)
	}
	messageKey := uuid.NewString()
	if _, err := t.schedulerClient.SendDelayMessage(ctx, &broker.SendDelayMessageRequest{
		Topic:        template.Topic,
		Key:          messageKey,
		Payload:      template.Payload,
		DeliverAfter: template.DeliverAfter,
	}); err != nil {
		return "", fmt.Errorf("failed to send delay message, caused by %w", err)
	}

	return messageKey, nil
}
