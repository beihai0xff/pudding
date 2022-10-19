package configs

// DelayQueueConfig DelayQueue Config
type DelayQueueConfig struct {
	PartitionInterval string `json:"partitionInterval" yaml:"partitionInterval"`
}

func GetDelayQueueConfig() *DelayQueueConfig {
	return c.delayQueue
}
