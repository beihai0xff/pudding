package cron

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/test/mock/api/gen/pudding/scheduler/v1"
	"github.com/beihai0xff/pudding/test/mock/utils"
)

var testTrigger *Trigger

func TestMain(m *testing.M) {
	// newMySQLServer()
	db := mysql.New(test_utils.TestMySQLConfig)

	testTrigger = &Trigger{
		schedulerClient: mock.NewMockSchedulerServiceClient(),
		repo:            repo.NewCronTemplate(db),
		wallClock:       clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	code := m.Run()
	os.Exit(code)
}
