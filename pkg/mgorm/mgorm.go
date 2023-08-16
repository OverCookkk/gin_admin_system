package mgorm

import (
	"database/sql"
	"fmt"
	"gin_admin_system/internal/app/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

func NewGormDB(c *Config) (*gorm.DB, error) {

	switch strings.ToLower(c.DBType) {
	case "mysql":
		// 不存在database则先创建
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/", config.C.MySQL.User, config.C.MySQL.Password, config.C.MySQL.Host)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		defer db.Close()

		query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", config.C.MySQL.DBName)
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               c.DSN,
		DefaultStringSize: 191,
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	return db, nil
}
