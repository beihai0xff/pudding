// Package: init_mysql_env.go init mysql env for test
package main

import (
	"github.com/beihai0xff/pudding/app/trigger/repo/po"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/test/mock/utils"
)

func main() {
	db := mysql.New(test_utils.TestMySQLConfig)
	if err := createTable(db); err != nil {
		panic(err)
	}
}

func createTable(db *mysql.Client) error {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&po.WebhookTriggerTemplate{}); err != nil {
		log.Errorf("create table failed: %v", err)

		return err
	}

	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&po.CronTriggerTemplate{}); err != nil {
		log.Errorf("create table failed: %v", err)

		return err
	}

	return nil
}

func dropTable(db *mysql.Client) error {
	if err := db.Migrator().DropTable(&po.WebhookTriggerTemplate{}); err != nil {
		log.Errorf("drop test table failed: %v", err)

		return err
	}

	if err := db.Migrator().DropTable(&po.CronTriggerTemplate{}); err != nil {
		log.Errorf("drop test table failed: %v", err)

		return err
	}

	return nil
}
