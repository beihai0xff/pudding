package cron

import (
	"os"
	"testing"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
)

var test_trigger *Trigger

func TestMain(m *testing.M) {

	// newMySQLServer()
	db := newMySQLClient()
	test_trigger = &Trigger{}
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
