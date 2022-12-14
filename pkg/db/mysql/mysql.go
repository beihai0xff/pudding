package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
)

type Client struct {
	*gorm.DB
}

func New(c *configs.MySQLConfig) *Client {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       c.DSN, // DSN data source name
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		DisableAutomaticPing:                     false,
		Logger:                                   logger.GetGORMLogger(),
	})

	if err != nil {
		log.Fatalf("failed to connect MySQL database: %v", err)
	}

	setConnPool(db)

	return &Client{db}
}

func setConnPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(500)
	// SetMaxIdleConns sets the maximum number of connections in the idle onnection pool.
	sqlDB.SetMaxIdleConns(100)

	// SetConnMaxIdleTime sets the maximum amount of time a connection may be idle
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func (c *Client) GetDB() *gorm.DB {
	return c.DB
}
