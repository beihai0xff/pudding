package mysql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	test_utils "github.com/beihai0xff/pudding/test/mock/utils"
)

// Example of how to implement a MySQLConfig server based on a Engine:
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
	client := New(test_utils.TestMySQLConfig)

	sqlDB, err := client.DB.DB()
	assert.Equal(t, nil, err)

	// Ping
	assert.Equal(t, nil, sqlDB.Ping())

	assert.NotEqual(t, nil, sqlDB.Stats())
}
