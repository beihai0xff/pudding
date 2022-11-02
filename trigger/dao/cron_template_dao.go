package dao

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/trigger/dao/storage/sql"
	"github.com/beihai0xff/pudding/trigger/entity"
)

type CronTemplateDAO interface {
	Insert(ctx context.Context, e *entity.CrontriggerTemplate) error
	Update(ctx context.Context, e *entity.CrontriggerTemplate) error
	FindEnableRecords(ctx context.Context, t time.Time, batchSize int, f sql.CronTempHandler) error
}

func NewCronTemplateDAO(db *mysql.Client) CronTemplateDAO {
	return sql.NewCronTemplate(db)
}
