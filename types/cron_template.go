package types

import (
	"github.com/beihai0xff/pudding/trigger/entity"
)

const (
	// TemplateStatusEnabled is the status of template which is enabled
	TemplateStatusEnabled = iota + 1
	// TemplateStatusDisabled is the status of template disabled.
	TemplateStatusDisabled
	// TemplateStatusMaxTimes the CrontriggerTemplate loop times exceeds the maximum times limit.
	TemplateStatusMaxTimes
	// TemplateStatusMaxAge the CrontriggerTemplate exceeds the maximum age limit.
	TemplateStatusMaxAge
)

type CronTempHandler func(results *entity.CronTriggerTemplate) error
