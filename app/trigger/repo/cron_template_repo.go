// package repo is a package for the repo layer.
// It contains the repository interfaces and implementations.

package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"gorm.io/gen"
	"gorm.io/gorm/clause"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo/convertor"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/sql"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

type CronTempHandler func(results *entity.CronTriggerTemplate) error

// CronTemplateDAO is the interface for the CronTemplate repository.
type CronTemplateDAO interface {
	// FindByID find a cron template by id
	FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error)
	// PageQuery query cron templates by page
	PageQuery(ctx context.Context, p *entity.PageQuery, status pb.TriggerStatus) (res []*entity.CronTriggerTemplate,
		count int64, err error)

	// Insert create a cron template
	Insert(ctx context.Context, e *entity.CronTriggerTemplate) error
	// UpdateStatus update the status of a cron template
	UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error)
	// BatchHandleRecords batch handle the records which need to be executed
	BatchHandleRecords(ctx context.Context, t time.Time, batchSize int, f CronTempHandler) error
}

type CronTemplate struct{}

func NewCronTemplate(db *mysql.Client) *CronTemplate {
	sql.SetDefault(db.GetDB())
	return &CronTemplate{}
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

func (dao *CronTemplate) PageQuery(ctx context.Context, p *entity.PageQuery, status pb.TriggerStatus) (
	[]*entity.CronTriggerTemplate, int64, error) {

	var res []*po.CronTriggerTemplate
	var count int64
	var err error
	if status > pb.TriggerStatus_UNKNOWN_UNSPECIFIED && status <= pb.TriggerStatus_MAX_AGE {
		res, count, err = sql.CronTriggerTemplate.WithContext(ctx).
			Where(sql.CronTriggerTemplate.Status.Eq(int32(status))).FindByPage(p.Offset, p.Limit)
	} else {
		res, count, err = sql.CronTriggerTemplate.WithContext(ctx).FindByPage(p.Offset, p.Limit)
	}

	if err != nil {
		return nil, 0, err
	}

	e, err := convertor.CronTemplateSlicePoTOEntity(res)
	if err != nil {
		return nil, 0, err
	}
	return e, count, nil

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

func (dao *CronTemplate) UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error) {
	return sql.CronTriggerTemplate.WithContext(ctx).UpdateStatus(ctx, id, status)
}

func (dao *CronTemplate) BatchHandleRecords(ctx context.Context, t time.Time, batchSize int,
	f CronTempHandler) error {
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
			_ = tx.Save(results[i])
			log.Infof("update %+v records", results)
		}

		return nil
	}

	// SELECT xxx FROM `xxxx` FOR UPDATE SKIP LOCKED LIMIT batchSize
	return sql.CronTriggerTemplate.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Where(sql.CronTriggerTemplate.LastExecutionTime.Lte(t)).
		Where(sql.CronTriggerTemplate.Status.Eq(int32(pb.TriggerStatus_ENABLED))).
		FindInBatches(&results, batchSize, fc)

}
