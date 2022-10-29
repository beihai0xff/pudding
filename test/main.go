package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/scheduler"
	"github.com/beihai0xff/pudding/types"
)

var s *scheduler.Schedule

func main() {
	configs.Init("./config.test.yaml")

	s = scheduler.New()
	lock.Init()
	go s.Run()

	for i := 0; i < 100; i++ {
		msg := &types.Message{
			Payload:   []byte("hello"),
			ReadyTime: time.Now().Unix() + 30,
		}
		fmt.Println(msg.ReadyTime)
		_ = s.Produce(context.Background(), msg)
	}
	time.Sleep(5 * time.Minute)
	// Exit
	os.Exit(0)
}
