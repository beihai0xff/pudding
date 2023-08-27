package cluster

import (
	"context"
	"errors"

	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
)

var (
	// ErrKeyExists is returned by putNewKV when the key already exists.
	ErrKeyExists = errors.New("key already exists")
	// ErrNoWatcher is returned when the watcher channel is nil.
	ErrNoWatcher = errors.New("no watcher channel")
)

// putNewKV attempts to create the given key, only succeeding if the key did
// not yet exist.
func putNewKV(kv v3.KV, key, val string, leaseID v3.LeaseID) error {
	cmp := v3.Compare(v3.Version(key), "=", 0)
	req := v3.OpPut(key, val, v3.WithLease(leaseID))

	txnResp, err := kv.Txn(context.TODO()).If(cmp).Then(req).Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return ErrKeyExists
	}

	return nil
}

// deleteRevKey deletes a key by revision, returning false if key is missing
func deleteRevKey(kv v3.KV, key string, rev int64) (bool, error) {
	cmp := v3.Compare(v3.ModRevision(key), "=", rev)
	req := v3.OpDelete(key)

	txnResp, err := kv.Txn(context.TODO()).If(cmp).Then(req).Commit()
	if err != nil {
		return false, err
	} else if !txnResp.Succeeded {
		return false, nil
	}

	return true, nil
}

// waitPrefixEvents waits for the given events on the prefix.
func waitPrefixEvents(c *v3.Client, prefix string, rev int64, evs []mvccpb.Event_EventType) (*v3.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wc := c.Watch(ctx, prefix, v3.WithPrefix(), v3.WithRev(rev))
	if wc == nil {
		return nil, ErrNoWatcher
	}

	return waitEvents(wc, evs), nil
}

// waitKeyEvents waits for the given events on the key.
func waitKeyEvents(c *v3.Client, key string, rev int64, evs []mvccpb.Event_EventType) (*v3.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wc := c.Watch(ctx, key, v3.WithRev(rev))

	if wc == nil {
		return nil, ErrNoWatcher
	}

	return waitEvents(wc, evs), nil
}

func waitEvents(wc v3.WatchChan, evs []mvccpb.Event_EventType) *v3.Event {
	i := 0

	for wresp := range wc {
		for _, ev := range wresp.Events {
			if ev.Type == evs[i] {
				i++
				if i == len(evs) {
					return ev
				}
			}
		}
	}

	return nil
}
