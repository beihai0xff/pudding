package main

import (
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/service/scheduler"
	"github.com/beihai0xff/pudding/types"
)

var s *scheduler.Schedule

func main() {
	configs.Init("./config.test.yaml")

	s = scheduler.New(configs.GetSchedulerConfig())
	lock.Init()
	go s.Run()

	for i := 0; i < 100; i++ {
		msg := &types.Message{
			Payload:   []byte("hello"),
			ReadyTime: time.Now().Unix() + 30,
		}
		_ = s.Produce(context.Background(), msg)
	}
	time.Sleep(5 * time.Minute)
	// Exit
	os.Exit(0)
}
