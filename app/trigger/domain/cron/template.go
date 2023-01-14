// Package cron implemented the cron trigger and handler
// template.go implements the cron template
package cron

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/pkg/constants"
	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/log"
)

var validate = validator.New()

// FindByID find cron template by id
func (t *Trigger) FindByID(ctx context.Context, id uint) (*TriggerTemplate, error) {
	if id <= 0 {
		err := fmt.Errorf("invalid id, id: %d", id)
		log.Errorf("%v", err)
		return nil, err
	}

	p, err := t.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("failed to find webhook template, caused by %v", err)
		return nil, err
	}

	res, err := convPoTOEntity(p)
	if err != nil {
		return nil, fmt.Errorf("convert po to entity failed: %w", err)
	}

	return res, nil
}

// PageQuery page query cron templates
func (t *Trigger) PageQuery(ctx context.Context, p *constants.PageQuery,
	status pb.TriggerStatus) ([]*TriggerTemplate, int64, error) {
	// check params
	if p.Offset < 0 || p.Limit <= 0 {
		err := fmt.Errorf("invalid offset or limit, offset: %d, limit: %d", p.Offset, p.Limit)
		log.Errorf("%v", err)
		return nil, 0, err
	}

	po, count, err := t.repo.PageQuery(ctx, p, status)
	if err != nil {
		log.Errorf("failed to PageQuery cron template, caused by %v", err)
		return nil, 0, err
	}

	res, err := convSlicePoTOEntity(po)
	if err != nil {
		log.Errorf("failed to convert po to entity, caused by %v", err)
	}

	return res, count, nil
}

// Register register a cron template
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

// checkRegisterParams check the params of register cron template
func (t *Trigger) checkRegisterParams(temp *TriggerTemplate) error {
	// 1. check cron expression
	if _, err := cronexpr.Parse(temp.CronExpr); err != nil {
		log.Errorf("Invalid cron expression: %v", err)
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
		temp.ExceptedEndTime = t.wallClock.Now().Add(constants.DefaultTemplateActiveDuration)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = constants.DefaultMaximumLoopTimes
	}

	temp.LastExecutionTime = defaultLastExecutionTime
	// default status is Disable
	temp.Status = pb.TriggerStatus_DISABLED

	return nil
}

// UpdateStatus update cron template status
func (t *Trigger) UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error) {
	// 1. set template status

	// 2. update the template status to db
	rowsAffected, err := t.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Errorf("failed to update cron template, caused by %v", err)
		return 0, err
	}

	return rowsAffected, nil
}
