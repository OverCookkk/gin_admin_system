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
	RoleID   uint64 `json:"role_id,string" binding:"required"`   // 角色ID，角色有哪些菜单权限，所以角色要与RoleMenus绑定
	MenuID   uint64 `json:"menu_id,string" binding:"required"`   // 菜单ID
	ActionID uint64 `json:"action_id,string" binding:"required"` // 动作ID
}

type TypesRoleMenu types.RoleMenu

func (r TypesRoleMenu) ToRoleMenu() *RoleMenu {
	roleMenu := &RoleMenu{}
	copier.Copy(roleMenu, r)
	return roleMenu
}