package configs

// SchedulerConfig DelayQueue Config
type SchedulerConfig struct {
	TimeSliceInterval string `json:"time_slice_interval" yaml:"time_slice_interval" mapstructure:"time_slice_interval"`
}
