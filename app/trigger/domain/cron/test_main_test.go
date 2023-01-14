package cron

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/test/mock/api/gen/pudding/scheduler/v1"
)

var testTrigger *Trigger

func TestMain(m *testing.M) {
	// newMySQLServer()
	db := newMySQLClient()
	createTable(db)

	testTrigger = &Trigger{
		schedulerClient: mock.NewMockSchedulerServiceClient(),
		repo:            repo.NewCronTemplate(db),
		wallClock:       clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	code := m.Run()

	dropTable(db)
	os.Exit(code)

}

func newMySQLClient() *mysql.Client {
	c := &configs.MySQLConfig{
		DSN: "root:my-secret-pw@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC",
	}

	return mysql.New(c)
}

func createTable(db *mysql.Client) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&po.CronTriggerTemplate{})
	if err != nil {
		log.Errorf("create table failed: %v", err)
	}
}

func dropTable(db *mysql.Client) {
	err := db.Migrator().DropTable(&po.CronTriggerTemplate{})
	if err != nil {
		log.Errorf("drop test table failed: %v", err)
	}
}
