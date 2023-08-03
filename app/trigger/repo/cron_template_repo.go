// Package repo is a package for the repo layer.
// It contains the repository interfaces and implementations.
package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"gorm.io/gen"
	"gorm.io/gorm/clause"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/pkg/constants"
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/sql"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

// CronTempHandler is the handler for cron template
// It will be called when the cron template is ready to be executed
type CronTempHandler func(results *po.CronTriggerTemplate) error

// CronTemplateDAO is the interface for the CronTemplate repository.
type CronTemplateDAO interface {
	// FindByID find a cron template by id
	FindByID(ctx context.Context, id uint) (*po.CronTriggerTemplate, error)
	// PageQuery query cron templates by page
	PageQuery(ctx context.Context, p *constants.PageQuery, status pb.TriggerStatus) (res []*po.CronTriggerTemplate,
		count int64, err error)

	// Insert create a cron template
	Insert(ctx context.Context, e *po.CronTriggerTemplate) error
	// UpdateStatus update the status of a cron template
	UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error)
	// BatchHandleRecords batch handle the records which need to be executed
	BatchHandleRecords(ctx context.Context, t time.Time, batchSize int, f CronTempHandler) error
}

// CronTemplate is the repository for the CronTemplate entity.
type CronTemplate struct{}

// NewCronTemplate returns a CronTemplate repository
func NewCronTemplate(db *mysql.Client) *CronTemplate {
	sql.SetDefault(db.GetDB())
	return &CronTemplate{}
}

// FindByID find a cron template by id
func (dao *CronTemplate) FindByID(ctx context.Context, id uint) (*po.CronTriggerTemplate, error) {
	// SELECT * FROM pudding_cron_trigger_template WHERE id =
	res, err := sql.CronTriggerTemplate.WithContext(ctx).FindByID(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// PageQuery query cron templates by page
func (dao *CronTemplate) PageQuery(ctx context.Context, p *constants.PageQuery, status pb.TriggerStatus) (
	[]*po.CronTriggerTemplate, int64, error) {
	var (
		res   []*po.CronTriggerTemplate
		count int64
		err   error
	)

	if status > pb.TriggerStatus_UNKNOWN_UNSPECIFIED && status <= pb.TriggerStatus_MAX_AGE {
		res, count, err = sql.CronTriggerTemplate.WithContext(ctx).
			Where(sql.CronTriggerTemplate.Status.Eq(int32(status))).FindByPage(p.Offset, p.Limit)
	} else {
		res, count, err = sql.CronTriggerTemplate.WithContext(ctx).FindByPage(p.Offset, p.Limit)
	}

	if err != nil {
		return nil, 0, err
	}

	return res, count, nil
}

// Insert create a cron template
func (dao *CronTemplate) Insert(ctx context.Context, p *po.CronTriggerTemplate) error {
	if p.CronExpr == "" {
		return fmt.Errorf("cron expression is empty")
	}

	return sql.CronTriggerTemplate.WithContext(ctx).Create(p)
}

// UpdateStatus update the status of a cron template
func (dao *CronTemplate) UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (int64, error) {
	return sql.CronTriggerTemplate.WithContext(ctx).UpdateStatus(ctx, id, status)
}

// BatchHandleRecords batch handle the records which need to be executed
func (dao *CronTemplate) BatchHandleRecords(ctx context.Context, t time.Time, batchSize int,
	f CronTempHandler) error {
	var results []*po.CronTriggerTemplate

	// handle function
	fw := func(tx gen.Dao, batch int) error {
		lo.ForEach(results, func(item *po.CronTriggerTemplate, index int) {
			if err := f(item); err != nil {
				log.Errorf("handle cron template [%d] failed: %v", item.ID, err)
				return
			}

			_ = tx.Save(item)
			log.Infof("update records: %+v ", results)
		})

		return nil
	}

	// SELECT xxx FROM `xxxx` FOR UPDATE SKIP LOCKED LIMIT batchSize
	return sql.CronTriggerTemplate.WithContext(ctx).Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Where(sql.CronTriggerTemplate.LastExecutionTime.Lte(t)).
		Where(sql.CronTriggerTemplate.Status.Eq(int32(pb.TriggerStatus_ENABLED))).
		FindInBatches(&results, batchSize, fw)
}
