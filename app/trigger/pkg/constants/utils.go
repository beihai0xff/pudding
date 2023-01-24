// Package constants provides the constants of the trigger
package constants

import "time"

// PageQuery is a helper struct to page query
type PageQuery struct {
	// Offset is the offset of the query
	// Offset must be greater than or equal to 0
	// Limit is the limit of the query
	// Limit must be greater than 0
	Offset, Limit int
}

const (

	// DefaultMaximumLoopTimes Maximum Loop Times of Cron Trigger: 1024
	DefaultMaximumLoopTimes = 1 << 10
	// DefaultTemplateActiveDuration is the default active duration of cron template: 30 days
	DefaultTemplateActiveDuration = 30 * 24 * time.Hour
)
