package convertor

import (
	"github.com/jinzhu/copier"

	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
)

func CronTemplateEntityTOPo(e *entity.CronTriggerTemplate) (*po.CronTriggerTemplate, error) {
	p := &po.CronTriggerTemplate{}
	if err := copier.Copy(p, e); err != nil {
		return nil, err
	}

	return p, nil
}

func CronTemplatePoTOEntity(p *po.CronTriggerTemplate) (*entity.CronTriggerTemplate, error) {
	e := &entity.CronTriggerTemplate{}
	if err := copier.Copy(e, p); err != nil {
		return nil, err
	}

	return e, nil
}
