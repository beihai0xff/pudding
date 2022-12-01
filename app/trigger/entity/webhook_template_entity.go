package entity

import (
	"time"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
)

type WebhookTriggerTemplate struct {
	ID uint

	// Topic the message topic
	Topic string `validate:"required"`
	// Payload the message payload
	Payload []byte `validate:"required"`

	// ExceptedEndTime Excepted Trigger end time, if it is 0, it means that it will not end.
	ExceptedEndTime time.Time `json:"excepted_end_time"`
	// ExceptedLoopTimes except loop times
	ExceptedLoopTimes uint64 `json:"excepted_loop_times"`
	// LoopedTimes already loop times
	LoopedTimes uint64 `json:"looped_times"`

	// Status the CronTriggerTemplate status: enable offline
	Status pb.TriggerStatus `json:"status" validate:"required"`
}
