package menu

import "gorm.io/gorm"

type MenuAction struct {
	gorm.Model
	MenuID uint64 `gorm:"index;not null;"`              // 菜单ID
	Code   string `gorm:"size:100;"`                    // 动作编号
	Name   string `gorm:"size:100;"`                    // 动作名称
	Menus  []Menu `gorm:"many2many:menu_menu_actions;"` // 菜单关联
}
