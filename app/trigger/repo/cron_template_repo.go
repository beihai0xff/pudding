package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gorm.io/gen"
	"gorm.io/gorm/clause"

	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo/convertor"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/sql"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/types"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

type CronTemplateDAO interface {
	FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error)

	Insert(ctx context.Context, e *entity.CronTriggerTemplate) error
	UpdateStatus(ctx context.Context, id uint, status int) error

	BatchHandleRecords(ctx context.Context, t time.Time, batchSize int, f types.CronTempHandler) error
}

type CronTemplate struct{}

func NewCronTemplate(db *mysql.Client) *CronTemplate {
	sql.SetDefault(db.GetDB())
	return &CronTemplate{}
}

func (dao *CronTemplate) Insert(ctx context.Context, e *entity.CronTriggerTemplate) error {
	if err := validate.Struct(e); err != nil {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	p, err := convertor.CronTemplateEntityTOPo(e)
	if err != nil {
		return fmt.Errorf("convert entity to po failed: %w", err)
	}
	if err := sql.CronTriggerTemplate.WithContext(ctx).Create(p); err != nil {
		return err
	}
	e.ID = p.ID

	return nil
}

func (dao *CronTemplate) UpdateStatus(ctx context.Context, id uint, status int) error {
	return sql.CronTriggerTemplate.WithContext(ctx).UpdateStatus(ctx, id, status)
}

func (dao *CronTemplate) BatchHandleRecords(ctx context.Context, t time.Time, batchSize int,
	f types.CronTempHandler) error {
	var results []*po.CronTriggerTemplate

	// handle function
	fc := func(tx gen.Dao, batch int) error {
		for i := 0; i < len(results); i++ {
			e, err := convertor.CronTemplatePoTOEntity(results[i])
			if err != nil {
				return err
			}
			if err := f(e); err != nil {
				return err
			}

			if err := copier.Copy(results[i], e); err != nil {
				return err
			}
		}

		log.Infof("update %+v records", results)

		return tx.Save(&results)
	}

	// SELECT xxx FROM `xxxx` FOR UPDATE SKIP LOCKED LIMIT batchSize
	return sql.CronTriggerTemplate.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Where(sql.CronTriggerTemplate.LastExecutionTime.Lte(t)).
		Where(sql.CronTriggerTemplate.Status.Eq(types.TemplateStatusEnabled)).
		FindInBatches(&results, batchSize, fc)

}

func (dao *CronTemplate) FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error) {
	// SELECT * FROM pudding_cron_trigger_template WHERE id =
	res, err := sql.CronTriggerTemplate.WithContext(ctx).FindByID(id)
	if err != nil {
		return nil, err
	}

	e, err := convertor.CronTemplatePoTOEntity(res)
	if err != nil {
		return nil, err
	}
	return e, nil

}
