//nolint:lll
package po

import (
	"time"

	"gorm.io/gorm"
)

type CronTriggerTemplate struct {
	gorm.Model

	// CronExpr is the cron expression.
	CronExpr string `gorm:"column:cron_expr;type:varchar(255);not null;default:'unknown';comment:'cron表达式'" copier:"must,nopanic"`

	// Topic the message topic
	Topic string `gorm:"column:topic;type:varchar(255);not null;default:'unknown';comment:'消息Topic'" copier:"must,nopanic"`
	// Payload the message payload
	Payload []byte `gorm:"column:payload;type:TEXT;not null;comment:'消息内容'" copier:"must,nopanic"`

	// NextExecutionTime last time to schedule the message
	LastExecutionTime time.Time `gorm:"column:last_execution_time;type:datetime;not null;comment:'上次执行时间'" copier:"must,nopanic"`

	// ExceptedEndTime Excepted Trigger end time, if it is 0, it means that it will not end.
	ExceptedEndTime time.Time `gorm:"column:excepted_end_time;type:datetime;not null;comment:'预期结束时间'" copier:"must,nopanic"`
	// ExceptedLoopTimes except loop times
	ExceptedLoopTimes uint64 `gorm:"column:excepted_loop_times;type:int unsigned;not null;default:0;comment:'预期循环次数'" copier:"must,nopanic"`
	// LoopedTimes already loop times
	LoopedTimes uint64 `gorm:"column:looped_times;type:int unsigned;not null;default:0;comment:'已循环次数'" copier:"must,nopanic"`

	// Status the CronTriggerTemplate status: enable disable offline
	Status int `gorm:"column:status;type:int unsigned;not null;default:0;comment:'CronTriggerTemplate status'" copier:"must,nopanic"`
}

func (CronTriggerTemplate) TableName() string {
	return "pudding_cron_trigger_template"
}
