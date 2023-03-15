package webhook

import (
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/repo"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	test_utils "github.com/beihai0xff/pudding/test/mock/utils"
)

var testTrigger *Trigger

const testHTTPDomain = "http://localhost:8080"

func TestMain(m *testing.M) {
	// newMySQLServer()
	db := mysql.New(test_utils.TestMySQLConfig)

	testTrigger = &Trigger{
		webhookPrefix: testHTTPDomain,
		repo:          repo.NewWebhookTemplate(db),
		wallClock:     clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	code := m.Run()
	os.Exit(code)

}
