package webhook

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	test_utils "github.com/beihai0xff/pudding/test/mock/utils"
)

var testTrigger *Trigger

const testHTTPDomain = "https://example.com"

func TestMain(m *testing.M) {
	// newMySQLServer()
	db := mysql.New(test_utils.TestMySQLConfig)
	createTable(db)

	testTrigger = &Trigger{
		webhookPrefix: testHTTPDomain,
		repo:          repo.NewWebhookTemplate(db),
		wallClock:     clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	code := m.Run()

	dropTable(db)
	os.Exit(code)

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
		log.Errorf("drop test table failed: %v", err)
	}
}
