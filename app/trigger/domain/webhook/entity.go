// Package webhook implemented the webhook trigger and handler
// entity.go implements the webhook template
package webhook

import (
	"time"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
)

// TriggerTemplate is the webhook template
type TriggerTemplate struct {
	// ID is the id of the template
	ID uint

	// Topic the message topic
	Topic string `validate:"required"`
	// Payload the message payload
	Payload []byte `validate:"required"`
	// Message DeliverAfter time (Seconds)
	DeliverAfter uint64 `json:"deliver_after" validate:"required"`

	// LoopedTimes already loop times
	LoopedTimes uint64 `json:"looped_times"`

	// ExceptedEndTime Excepted Trigger end time, if it is 0, it means that it will not end.
	ExceptedEndTime time.Time `json:"excepted_end_time"`
	// ExceptedLoopTimes except loop times
	ExceptedLoopTimes uint64 `json:"excepted_loop_times"`

	// Status the CronTriggerTemplate status: enable offline
	Status pb.TriggerStatus `json:"status" validate:"required"`
}
