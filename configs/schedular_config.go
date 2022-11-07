package configs

// SchedulerConfig DelayQueue Config
type SchedulerConfig struct {
	TimeSliceInterval string `json:"timeSliceInterval" yaml:"timeSliceInterval" mapstructure:"timeSliceInterval"`
}
