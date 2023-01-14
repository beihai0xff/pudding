// Package cron implemented the cron trigger and handler
// convertor.go implements the conversion between entity and po
package cron

import (
	"github.com/jinzhu/copier"

	"github.com/beihai0xff/pudding/app/trigger/repo/po"
)

func convEntityTOPo(e *TriggerTemplate) (*po.CronTriggerTemplate, error) {
	p := &po.CronTriggerTemplate{}
	if err := copier.Copy(p, e); err != nil {
		return nil, err
	}

	return p, nil
}

func convPoTOEntity(p *po.CronTriggerTemplate) (*TriggerTemplate, error) {
	e := &TriggerTemplate{}
	if err := copier.Copy(e, p); err != nil {
		return nil, err
	}

	return e, nil
}

func convSlicePoTOEntity(p []*po.CronTriggerTemplate) ([]*TriggerTemplate, error) {
	var e []*TriggerTemplate
	if err := copier.Copy(&e, p); err != nil {
		return nil, err
	}

	return e, nil
}
