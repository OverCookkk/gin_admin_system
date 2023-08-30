package dao

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"github.com/google/wire"
	"gorm.io/gorm"
	"strings"
)

var RepoSet = wire.NewSet(
	menu.MenuSet,
	role.RoleSet,
)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db.AutoMigrate(
		&menu.Menu{},
		&menu.MenuAction{},
		&role.Role{},
	)
}
