package webhook

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
	"github.com/beihai0xff/pudding/configs"
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
		repo:      repo.NewWebhookTemplate(db),
		wallClock: clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
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
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&po.WebhookTriggerTemplate{})
	if err != nil {
		log.Errorf("create table failed: %v", err)
	}
}

func dropTable(db *mysql.Client) {
	err := db.Migrator().DropTable(&po.WebhookTriggerTemplate{})
	if err != nil {
		log.Errorf("drop table failed: %v", err)
	}
}
