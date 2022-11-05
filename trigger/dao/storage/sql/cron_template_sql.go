package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/trigger/dao/convertor"
	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
	"github.com/beihai0xff/pudding/trigger/entity"
	"github.com/beihai0xff/pudding/types"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

type CronTemplate struct {
	db *mysql.Client
}

func NewCronTemplate(db *mysql.Client) *CronTemplate {
	return &CronTemplate{db: db}
}

func (dao *CronTemplate) Insert(ctx context.Context, e *entity.CronTriggerTemplate) error {
	if err := validate.Struct(e); err != nil {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	p, err := convertor.CronTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	if err := dao.db.WithContext(ctx).Create(p).Error; err != nil {
		return err
	}
	e.ID = p.ID

	return nil
}

func (dao *CronTemplate) Update(ctx context.Context, e *entity.CronTriggerTemplate) error {
	p, err := convertor.CronTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	return dao.db.WithContext(ctx).Updates(p).Error
}

func (dao *CronTemplate) FindEnableRecords(ctx context.Context, t time.Time, batchSize int,
	f types.CronTempHandler) error {
	var results []po.CronTriggerTemplate

	// handle function
	fc := func(tx *gorm.DB, batch int) error {
		for i := 0; i < len(results); i++ {
			e, err := convertor.CronTemplatePoTOEntity(&results[i])
			if err != nil {
				return err
			}
			if err := f(e); err != nil {
				return err
			}

			if err := copier.Copy(&results[i], e); err != nil {
				return err
			}
		}

		log.Infof("update %+v records", results)

		tx.Save(&results)

		return nil
	}

	// SELECT xxx FROM `xxxx` FOR UPDATE SKIP LOCKED LIMIT batchSize
	return dao.db.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Where("status = ? AND last_execution_time <= ?", types.TemplateStatusEnable, &t).
		FindInBatches(&results, batchSize, fc).Error

}

func (dao *CronTemplate) FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error) {
	var res po.CronTriggerTemplate

	// SELECT * FROM pudding_cron_trigger_template WHERE id =
	if err := dao.db.WithContext(ctx).Find(&res, id).Error; err != nil {
		return nil, err
	}

	e, err := convertor.CronTemplatePoTOEntity(&res)
	if err != nil {
		return nil, err
	}
	return e, nil

}
