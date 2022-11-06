package cron

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/internal/trigger/repo"
	"github.com/beihai0xff/pudding/internal/trigger/repo/storage/po"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

var testTrigger *Trigger

func TestMain(m *testing.M) {

	// newMySQLServer()
	db := newMySQLClient()
	createTable(db)

	testTrigger = &Trigger{
		s:     nil,
		dao:   repo.NewCronTemplate(db),
		clock: clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

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