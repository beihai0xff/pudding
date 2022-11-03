package mysql

import (
	"os"
	"testing"

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
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	c := &configs.MySQLConfig{
		DSN: "root:my-secret-pw@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		Log: &configs.LogConfig{
			Writers:    []string{"console"},
			Format:     "console",
			Level:      "info",
			CallerSkip: 3,
		},
	}

	client := New(c)

	sqlDB, err := client.DB.DB()
	assert.Equal(t, nil, err)

	// Ping
	assert.Equal(t, nil, sqlDB.Ping())

	assert.NotEqual(t, nil, sqlDB.Stats())
}
