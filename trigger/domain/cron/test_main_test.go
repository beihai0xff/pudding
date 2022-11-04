package cron

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
)

var testTrigger *Trigger

func TestMain(m *testing.M) {

	// newMySQLServer()
	db := newMySQLClient()
	testTrigger = &Trigger{
		s:     nil,
		dao:   nil,
		clock: clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	createTable(db)

	code := m.Run()

	dropTable(db)
	os.Exit(code)

}

func newMySQLClient() *mysql.Client {
	c := &configs.MySQLConfig{
		DSN: "root:my-secret-pw@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC",
		Log: &configs.LogConfig{
			Writers:    []string{"console"},
			Format:     "console",
			Level:      "info",
			CallerSkip: 3,
		},
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
		log.Errorf("drop table failed: %v", err)
	}
}
