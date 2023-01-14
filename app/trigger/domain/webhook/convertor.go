// Package webhook implemented the webhook trigger and handler
// convertor.go implements the conversion between entity and po
package webhook

import (
	"github.com/jinzhu/copier"

	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
)

func convEntityTOPo(e *TriggerTemplate) (*po.WebhookTriggerTemplate, error) {
	p := &po.WebhookTriggerTemplate{}
	if err := copier.Copy(p, e); err != nil {
		return nil, err
	}

	return p, nil
}

func convPoTOEntity(p *po.WebhookTriggerTemplate) (*TriggerTemplate, error) {
	e := &TriggerTemplate{}
	if err := copier.Copy(e, p); err != nil {
		return nil, err
	}
	return e, nil
}

func convSlicePoTOEntity(p []*po.WebhookTriggerTemplate) ([]*TriggerTemplate, error) {
	var e []*TriggerTemplate
	if err := copier.Copy(&e, p); err != nil {
		return nil, err
	}

	return e, nil
}
