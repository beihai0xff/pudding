// Package test_utils provides some constants for testing.
package test_utils

import "github.com/beihai0xff/pudding/configs"

const (
	// TestConfigFilePath is the path of test config file.
	TestConfigFilePath = "./test/data/config.test.yaml"
)

// TestMySQLConfig is the test config for MySQL.
var TestMySQLConfig = &configs.MySQLConfig{
	DSN: "root:my-secret-pw@tcp(mysql_host:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC",
}
