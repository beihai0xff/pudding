package sql

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/trigger/dao/convertor"
	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

type CronTempHandler func(results *po.CronTriggerTemplate) error

type CronTemplate struct {
	db *mysql.Client
}

func NewCronTemplate(db *mysql.Client) *CronTemplate {
	return &CronTemplate{db: db}
}

func (dao *CronTemplate) Insert(ctx context.Context, e *entity.CrontriggerTemplate) error {
	p, err := convertor.CronTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	return dao.db.WithContext(ctx).Create(p).Error
}

func (dao *CronTemplate) Update(ctx context.Context, e *entity.CrontriggerTemplate) error {
	p, err := convertor.CronTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	return dao.db.WithContext(ctx).Updates(p).Error
}

func (dao *CronTemplate) FindEnableRecords(ctx context.Context, t time.Time, batchSize int, f CronTempHandler) error {
	var results []po.CronTriggerTemplate

	// handle function
	fc := func(tx *gorm.DB, batch int) error {
		for _, t := range results {
			if err := f(&t); err != nil {
				return err
			}
		}

		tx.Save(&results)

		return nil
	}

	// SELECT xxx FROM `xxxx` FOR UPDATE SKIP LOCKED LIMIT batchSize
	return dao.db.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Where("status = ? AND last_execution_time > ", types.TemplateStatusEnable, t).
		FindInBatches(results, batchSize, fc).Error

}
