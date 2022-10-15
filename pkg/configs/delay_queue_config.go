package configs

// DelayQueueConfig DelayQueue Config
type DelayQueueConfig struct {
	PartitionInterval int64 `json:"partitionInterval" yaml:"partitionInterval"`
}

func GetDelayQueueConfig() *DelayQueueConfig {
	return c.delayQueue
}
