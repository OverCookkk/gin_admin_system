package dao

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/dao/user"
	"gin_admin_system/internal/app/dao/util"
	"github.com/google/wire"
	"gorm.io/gorm"
	"strings"
)

var RepoSet = wire.NewSet(
	util.TransSet,
	menu.MenuSet,
	menu.MenuActionSet,
	menu.MenuActionResourceSet,
	role.RoleSet,
	role.RoleMenuSet,
	user.UserSet,
	user.UserRoleSet,
)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	return db.AutoMigrate(
		&menu.Menu{},
		&menu.MenuAction{},
		&menu.MenuActionResource{},
		&role.Role{},
		&role.RoleMenu{},
		&user.User{},
		&user.UserRole{},
	)
}
