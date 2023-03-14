package repo

import (
	"os"
	"testing"

	"github.com/beihai0xff/pudding/pkg/db/mysql"
	test_utils "github.com/beihai0xff/pudding/test/mock/utils"
)

var (
	testCronTemplate    CronTemplateDAO
	testWebhookTemplate WebhookTemplate
)

func TestMain(m *testing.M) {
	db := mysql.New(test_utils.TestMySQLConfig)

	createDao(db)

	code := m.Run()

	os.Exit(code)

}

func createDao(db *mysql.Client) {
	testCronTemplate = NewCronTemplate(db)
	testWebhookTemplate = NewWebhookTemplate(db)
}
