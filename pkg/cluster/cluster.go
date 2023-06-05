// Package cluster provides a cluster manager.
// cluster.go contains the Cluster interface and its implementation.
package cluster

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/beihai0xff/pudding/pkg/log"
)

type (
	// Cluster is a cluster manager.
	Cluster interface {
		// Mutex returns a distributed mutex implementation.
		Mutex(name string, ttl time.Duration, opts ...MutexOption) (Mutex, error)
	}

	// cluster is a cluster manager implementation.
	cluster struct {
		client *clientv3.Client
		opts   *clusterOptions
	}

	// clusterOptions contains options for cluster.
	clusterOptions struct {
		// requestTimeout is the timeout for all requests to etcd.
		requestTimeout time.Duration
	}

	// ClusterOption is a function that applies an option to cluster.
	ClusterOption func(*clusterOptions)
)

// New creates a new cluster manager.
func New(urls []string, opts ...ClusterOption) (Cluster, error) {
	client, err := newETCDClient(urls)
	if err != nil {
		return nil, fmt.Errorf("create etcd client failed: %v", err)
	}

	return newCluster(client, opts...), nil
}

func newETCDClient(urls []string) (*clientv3.Client, error) {
	return clientv3.NewFromURLs(urls)
}

func newCluster(client *clientv3.Client, opts ...ClusterOption) *cluster {
	ops := clusterOptions{requestTimeout: defaultRequestTimeout}
	for _, opt := range opts {
		opt(&ops)
	}

	c := cluster{
		client: client,
		opts:   &ops,
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

	log.Infof("session is ready")

	return session, nil
}

func (c *cluster) grantLease(ctx context.Context, ttl int64) (clientv3.LeaseID, error) {
	respGrant, err := c.client.Grant(ctx, ttl)
	if err != nil {
		return 0, fmt.Errorf("grant lease failed: %v", err)
	}

	return respGrant.ID, nil
}