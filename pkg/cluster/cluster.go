// Package cluster provides a cluster manager.
// cluster.go contains the Cluster interface and its implementation.
package cluster

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/beihai0xff/pudding/pkg/clock"
)

type (
	// Cluster is a cluster manager.
	Cluster interface {
		// Mutex returns a distributed mutex implementation.
		Mutex(name string, ttl time.Duration, opts ...MutexOption) (Mutex, error)
		// Queue returns a distributed queue implementation.
		Queue(topic string) (Queue, error)

		// WallClock returns the wall clock time
		WallClock() time.Time
	}

	// cluster is a cluster manager implementation.
	cluster struct {
		client *clientv3.Client
		opts   *clusterOptions

		// wallClock wall wallClock time
		wallClock clock.Clock
	}

	// clusterOptions contains options for cluster.
	clusterOptions struct {
		// requestTimeout is the timeout for all requests to etcd.
		requestTimeout time.Duration
	}

	// Option is a function that applies an option to cluster.
	Option func(*clusterOptions)
)

// New creates a new cluster manager.
func New(urls []string, opts ...Option) (Cluster, error) {
	client, err := newETCDClient(urls)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed: %v", err)
	}

	return newCluster(client, opts...), nil
}

func newETCDClient(urls []string) (*clientv3.Client, error) {
	return clientv3.NewFromURLs(urls)
}

func newCluster(client *clientv3.Client, opts ...Option) *cluster {
	ops := clusterOptions{requestTimeout: defaultRequestTimeout}
	for _, opt := range opts {
		opt(&ops)
	}

	c := cluster{
		client:    client,
		opts:      &ops,
		wallClock: clock.New(),
	}

	return &c
}

func (c *cluster) getSession(ctx context.Context, ttl int64) (*concurrency.Session, error) {
	leaseID, err := c.grantLease(ctx, ttl)
	if err != nil {
		return nil, err
	}

	session, err := concurrency.NewSession(c.client, concurrency.WithLease(leaseID), concurrency.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("create session failed: %v", err)
	}

	return session, nil
}

func (c *cluster) WallClock() time.Time {
	return c.wallClock.Now()
}

func (c *cluster) grantLease(ctx context.Context, ttl int64) (clientv3.LeaseID, error) {
	respGrant, err := c.client.Grant(ctx, ttl)
	if err != nil {
		return 0, fmt.Errorf("grant lease failed: %v", err)
	}

	return respGrant.ID, nil
}
