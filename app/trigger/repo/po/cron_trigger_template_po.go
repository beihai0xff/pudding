// package po is the package of persistent object

//nolint:lll
package po

import (
	"time"

	"gorm.io/gorm"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
)

// CronTriggerTemplate is the po of CronTriggerTemplate
// it define the table name and columns
type CronTriggerTemplate struct {
	gorm.Model

	// CronExpr is the cron expression.
	CronExpr string `gorm:"column:cron_expr;type:varchar(255);not null;default:'unknown';comment:'cron expr'" copier:"must,nopanic"`

	// Topic the message topic
	Topic string `gorm:"column:topic;type:varchar(255);not null;default:'unknown';comment:'message topic'" copier:"must,nopanic"`
	// Payload the message payload
	Payload []byte `gorm:"column:payload;type:TEXT;not null;comment:'message content'" copier:"must,nopanic"`

	// LastExecutionTime last time to schedule the message
	LastExecutionTime time.Time `gorm:"column:last_execution_time;type:TIMESTAMP;not null;comment:'last time to schedule the message'" copier:"must,nopanic"`
	// LoopedTimes already loop times
	LoopedTimes uint64 `gorm:"column:looped_times;type:int unsigned;not null;default:0;comment:'already loop times'" copier:"must,nopanic"`

	// ExceptedEndTime excepted trigger end time, if it is 0, it means that it will not end.
	ExceptedEndTime time.Time `gorm:"column:excepted_end_time;type:TIMESTAMP;not null;comment:'excepted trigger end time, if it is 0, it means that it will not end.'" copier:"must,nopanic"`
	// ExceptedLoopTimes except loop times
	ExceptedLoopTimes uint64 `gorm:"column:excepted_loop_times;type:int unsigned;not null;default:0;comment:'except loop times'" copier:"must,nopanic"`

	// Status the trigger template status: enable->1 disable->2 offline->3 and so on
	Status pb.TriggerStatus `gorm:"column:status;type:int unsigned;not null;default:0;comment:'trigger template status: enable->1 disable->2 offline->3 and so on'" copier:"must,nopanic"`
}

// TableName returns the table name of CronTriggerTemplate
// impl gorm Tabler interface
func (CronTriggerTemplate) TableName() string {
	return "cron_trigger_template"
}
