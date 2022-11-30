package cron

import (
	"context"
	"fmt"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/pkg/cronexpr"
	"github.com/beihai0xff/pudding/pkg/log"
)

// FindByID page query cron templates
func (t *Trigger) FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error) {
	if id <= 0 {
		err := fmt.Errorf("invalid id, id: %d", id)
		log.Errorf("%v", err)
		return nil, err
	}

	res, err := t.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("failed to insert cron template, caused by %v", err)
		return nil, err
	}

	return res, nil
}

// PageQuery page query cron templates
func (t *Trigger) PageQuery(ctx context.Context, offset, limit int) ([]*entity.CronTriggerTemplate, int64, error) {
	if offset < 0 || limit <= 0 {
		err := fmt.Errorf("invalid offset or limit, offset: %d, limit: %d", offset, limit)
		log.Errorf("%+v", err)
		return nil, 0, err
	}
	res, count, err := t.repo.PageQuery(ctx, offset, limit)
	if err != nil {
		log.Errorf("failed to PageQuery cron template, caused by %v", err)
		return nil, 0, err
	}

	return res, count, nil
}

// Register register a cron template
func (t *Trigger) Register(ctx context.Context, temp *entity.CronTriggerTemplate) error {
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

// checkRegisterParams check the params of register cron template
func (t *Trigger) checkRegisterParams(temp *entity.CronTriggerTemplate) error {
	// 1. check cron expression
	if _, err := cronexpr.Parse(temp.CronExpr); err != nil {
		log.Errorf("Invalid cron expression: %v", err)
		return fmt.Errorf("invalid cron expression: %v", err)
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
		temp.ExceptedEndTime = t.wallClock.Now().Add(defaultTemplateActiveDuration)
	}
	if temp.ExceptedLoopTimes == 0 {
		temp.ExceptedLoopTimes = defaultMaximumLoopTimes
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
