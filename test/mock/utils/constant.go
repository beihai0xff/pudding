package test_utils

import "github.com/beihai0xff/pudding/configs"

const (
	TestConfigFilePath = "./test/data/config.test.yaml"
)

var TestMySQLConfig = &configs.MySQLConfig{
	DSN: "root:my-secret-pw@tcp(mysql_host:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC",
}
