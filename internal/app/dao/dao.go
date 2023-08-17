package dao

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/dao/menu"
	"gorm.io/gorm"
	"strings"
)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db.AutoMigrate(
		&menu.Menu{},
		&menu.MenuAction{},
	)
}
