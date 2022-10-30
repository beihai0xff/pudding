package mysql

import (
	"os"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/configs"
)

// Example of how to implement a MySQL server based on a Engine:
//
// ```
// > mysql --host=127.0.0.1 --port=3306 -u root test -e "SELECT * FROM user"
// +----------+-------------------+-------------------------------+---------------------+
// | name     | email             | phone_numbers                 | created_at          |
// +----------+-------------------+-------------------------------+---------------------+
// | John Doe | john@doe.com      | ["555-555-555"]               |                     |
// | John Doe | johnalt@doe.com   | []                            |                     |
// | Jane Doe | jane@doe.com      | []                            |                     |
// | Evil Bob | evilbob@gmail.com | ["555-666-555","666-666-666"] |                     |
// +----------+-------------------+-------------------------------+---------------------+
// ```
func TestMain(m *testing.M) {
	engine := sqle.NewDefault(
		sql.NewDatabaseProvider(
			createTestDatabase(),
			information_schema.NewInformationSchemaDatabase(),
		))
	engine.Analyzer.Catalog.MySQLDb.AddRootAccount()
	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
	}
	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}
	go s.Start()

	time.Sleep(time.Second)

	defer s.Close()
	os.Exit(m.Run())
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "test"
		tableName = "user"
	)
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Datetime, Nullable: false, Source: tableName},
	}), nil)

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()
	_ = table.Insert(ctx, sql.NewRow("John Doe", "john@doe.com", sql.MustJSON(`["555-555-555"]`), time.Now()))
	_ = table.Insert(ctx, sql.NewRow("John Doe", "johnalt@doe.com", sql.MustJSON(`[]`), time.Now()))
	_ = table.Insert(ctx, sql.NewRow("Jane Doe", "jane@doe.com", sql.MustJSON(`[]`), time.Now()))
	_ = table.Insert(ctx, sql.NewRow("Jane Deo", "janedeo@gmail.com", sql.MustJSON(`["556-565-566", "777-777-777"]`), time.Now()))
	return db
}

func TestNew(t *testing.T) {
	c := configs.MySQLConfig{DSN: "root:@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"}

	client := New(c)

	count := int64(0)
	err := client.Table("user").Count(&count).Error
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(4), count)
}
