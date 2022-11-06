package repo

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/internal/trigger/entity"
	"github.com/beihai0xff/pudding/internal/trigger/repo/storage/sql"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/types"
)

type CronTemplateDAO interface {
	FindByID(ctx context.Context, id uint) (*entity.CronTriggerTemplate, error)
	Insert(ctx context.Context, e *entity.CronTriggerTemplate) error
	Update(ctx context.Context, e *entity.CronTriggerTemplate) error
	BatchEnabledRecords(ctx context.Context, t time.Time, batchSize int, f types.CronTempHandler) error
}

func NewCronTemplate(db *mysql.Client) CronTemplateDAO {
	return sql.NewCronTemplate(db)
}
