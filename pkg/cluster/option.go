package cluster

import "time"

const (
	defaultRequestTimeout = 10 * time.Second
)

func WithRequestTimeout(timeout time.Duration) ClusterOption {
	return func(opt *clusterOptions) {
		opt.requestTimeout = timeout
	}
}
