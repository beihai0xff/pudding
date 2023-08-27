package cluster

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var testCluster *cluster

func TestMain(m *testing.M) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	testCluster = newCluster(client, WithRequestTimeout(time.Second))

	m.Run()
	testCluster.client.Delete(context.Background(), "/test", clientv3.WithPrefix())
}

func Test_newCluster(t *testing.T) {
	_, err := testCluster.Mutex("test", 10)
	assert.ErrorAs(t, err, &ErrInvalidTTL)
	_, err = testCluster.Mutex("test", 10*time.Second)
	assert.NoError(t, err)
}
