package role

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(RoleMenu)) // defDB.Model(new(RoleMenu))
}

type RoleMenu struct {
	gorm.Model
	RoleID   uint64 `gorm:"index;not null;"` // 角色ID，角色有哪些菜单权限，所以角色要与RoleMenus绑定
	MenuID   uint64 `gorm:"index;not null;"` // 菜单ID
	ActionID uint64 `gorm:"index;not null;"` // 动作ID
}

func (r RoleMenu) ToTypesRoleMenu() *types.RoleMenu {
	role := &types.RoleMenu{}
	copier.Copy(role, r)
	return role
}

type RoleMenus []*RoleMenu

func (r RoleMenus) ToTypesRoleMenus() []types.RoleMenu {
	list := make([]types.RoleMenu, 0, len(r))
	for _, v := range r {
		list = append(list, *(v.ToTypesRoleMenu()))
	}
	return list
}

type TypesRoleMenu types.RoleMenu

func (r TypesRoleMenu) ToRoleMenu() *RoleMenu {
	roleMenu := &RoleMenu{}
	copier.Copy(roleMenu, r)
	return roleMenu
}
