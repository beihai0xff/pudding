package configs

// SchedulerConfig DelayQueue Config
type SchedulerConfig struct {
	TimeSliceInterval string `json:"timeSliceInterval" yaml:"timeSliceInterval"`
}

func GetSchedulerConfig() *SchedulerConfig {
	return c.Scheduler
}
