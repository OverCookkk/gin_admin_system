package app

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/pkg/mgorm"
	"gorm.io/gorm"
)

func InitGormDB() (*gorm.DB, error) {
	cfg := config.C
	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	}

	db, err := mgorm.NewGormDB(&mgorm.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移创建表
	if cfg.Gorm.EnableAutoMigrate {
		// err = dao.AutoMigrate(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
