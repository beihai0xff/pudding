package types

import (
	"github.com/beihai0xff/pudding/trigger/entity"
)

const (
	// TemplateStatusEnable is the status of template which is enabled
	TemplateStatusEnable = iota + 1
	// TemplateStatusDisable is the status of template disabled.
	TemplateStatusDisable
	// TemplateStatusMaxTimes the CrontriggerTemplate loop times exceeds the maximum times limit.
	TemplateStatusMaxTimes
	// TemplateStatusMaxAge the CrontriggerTemplate exceeds the maximum age limit.
	TemplateStatusMaxAge
)

type CronTempHandler func(results *entity.CronTriggerTemplate) error
