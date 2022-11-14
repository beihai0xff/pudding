package configs

// SchedulerConfig Scheduler Config
type SchedulerConfig struct {
	TimeSliceInterval string `json:"time_slice_interval" yaml:"time_slice_interval" mapstructure:"time_slice_interval"`
	MessageTopic      string `json:"message_topic" yaml:"message_topic" mapstructure:"message_topic"`
	TokenTopic        string `json:"token_topic" yaml:"token_topic" mapstructure:"token_topic"`
}
