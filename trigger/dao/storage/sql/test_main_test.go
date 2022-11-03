package sql

import (
	"os"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/trigger/dao/storage/po"
)

var cronTemplateSQL *CronTemplate

func TestMain(m *testing.M) {

	// newMySQLServer()
	db := newMySQLClient()
	createTable(db)
	createDao(db)

	code := m.Run()

	dropTable(db)
	os.Exit(code)

}

func newMySQLServer() {
	engine := sqle.NewDefault(
		sql.NewDatabaseProvider(
			memory.NewDatabase("test"),
			information_schema.NewInformationSchemaDatabase(),
		))
	engine.Analyzer.Catalog.MySQLDb.AddRootAccount()
	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Version:  "8.0.23",
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}
	go s.Start()
	time.Sleep(time.Second)

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

func createDao(db *mysql.Client) {
	cronTemplateSQL = NewCronTemplate(db)
}
