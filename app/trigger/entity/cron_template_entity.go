package entity

import (
	"time"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
)

type CronTriggerTemplate struct {
	ID uint

	// CronExpr is the cron expression.
	CronExpr string `json:"cron_expr" validate:"required"`
	// Topic the message topic
	Topic string `validate:"required"`
	// Payload the message payload
	Payload []byte `validate:"required"`

	// LastExecutionTime last time to schedule the message
	LastExecutionTime time.Time `json:"last_execution_time"`
	// LoopedTimes already loop times
	LoopedTimes uint64 `json:"looped_times"`

	// ExceptedEndTime Excepted Trigger end time, if it is 0, it means that it will not end.
	ExceptedEndTime time.Time `json:"excepted_end_time"`
	// ExceptedLoopTimes except loop times
	ExceptedLoopTimes uint64 `json:"excepted_loop_times"`

	// Status the CronTriggerTemplate status: enable offline
	Status pb.TriggerStatus `json:"status" validate:"required"`
}
