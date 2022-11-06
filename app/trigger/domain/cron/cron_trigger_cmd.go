package cron

import (
	"context"
	"fmt"

	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/types"
)

// Register register a cron template
func (t *Trigger) Register(ctx context.Context, temp *entity.CronTriggerTemplate) error {
	// 1. check params
	if err := t.checkRegisterParams(temp); err != nil {
		log.Errorf("failed to check params, caused by %w", err)
		return err
	}

	// 2. save the template to db
	if err := t.dao.Insert(ctx, temp); err != nil {
		log.Errorf("failed to insert cron template, caused by %w", err)
		return err
	}

	return nil
}

// checkRegisterParams check the params of register cron template
func (t *Trigger) checkRegisterParams(temp *entity.CronTriggerTemplate) error {
	// 1. check cron expression
	if _, err := cronexpr.Parse(temp.CronExpr); err != nil {
		log.Errorf("Invalid cron expression: %w", err)
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// 2. check topic
	if temp.Topic == "" {
		log.Error(errCronTemplateTopicNotFound.Error())
		return errCronTemplateTopicNotFound
	}

	// 3. check payload
	if len(temp.Payload) == 0 {
		log.Error(errCronTemplatePayloadNotFound.Error())
		return errCronTemplatePayloadNotFound
	}

	// 4. set default value if necessary
	if temp.ExceptedEndTime.IsZero() {
		temp.ExceptedEndTime = t.clock.Now().Add(defaultTemplateActiveDuration)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = defaultMaximumLoopTimes
	}

	temp.LastExecutionTime = defaultLastExecutionTime
	// default status is Disable
	temp.Status = types.TemplateStatusDisabled

	return nil
}

// UpdateStatus update cron template status
func (t *Trigger) UpdateStatus(ctx context.Context, id uint, status int) error {
	// 1. set template status
	temp := &entity.CronTriggerTemplate{
		ID:     id,
		Status: status,
	}

	// 2. update the template status to db
	if err := t.dao.Update(ctx, temp); err != nil {
		log.Errorf("failed to update cron template, caused by %w", err)
		return err
	}

	return nil
}
