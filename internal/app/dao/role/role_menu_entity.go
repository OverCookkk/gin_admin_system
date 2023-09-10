package role

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return defDB.Model(new(RoleMenu))
}

type RoleMenu struct {
	gorm.Model
	RoleID   uint64 `gorm:"index;not null;"` // 角色ID，角色有哪些菜单权限，所以角色要与RoleMenus绑定
	MenuID   uint64 `gorm:"index;not null;"` // 菜单ID
	ActionID uint64 `gorm:"index;not null;"` // 动作ID
}

type TypesRoleMenu types.RoleMenu

func (r TypesRoleMenu) ToRoleMenu() *RoleMenu {
	roleMenu := &RoleMenu{}
	copier.Copy(roleMenu, r)
	return roleMenu
}
