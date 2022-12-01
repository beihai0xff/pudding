package convertor

import (
	"github.com/jinzhu/copier"

	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
)

func WebhookTemplateEntityTOPo(e *entity.WebhookTriggerTemplate) (*po.WebhookTriggerTemplate, error) {
	p := &po.WebhookTriggerTemplate{}
	if err := copier.Copy(p, e); err != nil {
		return nil, err
	}

	return p, nil
}

func WebhookTemplatePoTOEntity(p *po.WebhookTriggerTemplate) (*entity.WebhookTriggerTemplate, error) {
	e := &entity.WebhookTriggerTemplate{}
	if err := copier.Copy(e, p); err != nil {
		return nil, err
	}

	return e, nil
}

func WebhookTemplateSlicePoTOEntity(p []*po.WebhookTriggerTemplate) ([]*entity.WebhookTriggerTemplate, error) {
	var e []*entity.WebhookTriggerTemplate
	if err := copier.Copy(&e, p); err != nil {
		return nil, err
	}

	return e, nil
}
