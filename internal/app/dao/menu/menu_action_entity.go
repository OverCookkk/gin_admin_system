package menu

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(MenuAction)) // defDB.Model(new(MenuAction))
}

type MenuAction struct {
	gorm.Model
	MenuID uint64 `gorm:"index;not null;"` // 菜单ID：所属的菜单id
	Code   string `gorm:"size:100;"`       // 动作编号：add、edit、del、query等
	Name   string `gorm:"size:100;"`       // 动作名称：新增、编辑、查询、禁用等
	// Menus  []Menu `gorm:"many2many:menu_menu_actions;"` // 菜单关联
}

// 把数据库的对应的MenuActions转换为返回前端的MenuActions
func (m MenuAction) ToTypesMenuAction() *types.MenuAction {
	menuAction := &types.MenuAction{}
	copier.Copy(menuAction, m) // m赋值给menuAction
	return menuAction
}

type MenuAcitons []*MenuAction

func (m MenuAcitons) ToTypesMenuActions() []types.MenuAction {
	list := make([]types.MenuAction, 0, len(m))
	for _, v := range m {
		list = append(list, *(v.ToTypesMenuAction()))
	}
	return list
}

type TypesMenuAction types.MenuAction

// 把返回前端的MenuAction转换为数据库的对应的MenuAction
func (a TypesMenuAction) ToMenuAction() *MenuAction {
	item := &MenuAction{}
	copier.Copy(item, a)
	return item
}
