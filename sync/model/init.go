package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(driver, dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 慢查询日志:超过 500ms 的查询记录为 Warn 级别
	customLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 500 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      false,
	})
	gormCfg := &gorm.Config{Logger: customLogger}

	switch driver {
	case DriverMySQL:
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)
	case DriverSQLite:
		db, err = gorm.Open(sqlite.Open(dsn), gormCfg)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&Platform{},
		&Repo{},
		&SyncTask{},
		&SyncRun{},
		&SyncRunStep{},
		&WebhookRule{},
		&WebhookRuleTask{},
		&WebhookEvent{},
		&OperationLog{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
