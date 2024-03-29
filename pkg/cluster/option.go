package cluster

import "time"

const (
	defaultRequestTimeout = 10 * time.Second
)

// WithRequestTimeout set the request timeout for etcd requests.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(opt *clusterOptions) {
		opt.requestTimeout = timeout
	}
}
