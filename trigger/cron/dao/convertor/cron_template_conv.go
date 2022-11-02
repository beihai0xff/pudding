package convertor

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"

	"github.com/beihai0xff/pudding/trigger/cron/dao/storage/po"
	"github.com/beihai0xff/pudding/trigger/cron/domain/entity"
)

// use a single instance of Validate, it caches struct info
var validate = validator.New()

func CronTemplateEntityTOPo(e *entity.CrontriggerTemplate) (*po.CronTriggerTemplate, error) {
	if err := validate.Struct(e); err != nil {
		return nil, fmt.Errorf("invalid validation error: %w", err)
	}

	p := &po.CronTriggerTemplate{}
	if err := copier.Copy(p, e); err != nil {
		return nil, err
	}
	return p, nil
}
