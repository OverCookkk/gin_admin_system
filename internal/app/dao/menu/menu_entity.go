package menu

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return defDB.Model(new(Menu))
}

// Menu 菜单数据库实体结构
type Menu struct {
	gorm.Model
	Name       string  `gorm:"size:50;index;default:'';not null;"` // 菜单名称
	Icon       *string `gorm:"size:255;"`                          // 菜单图标
	Router     *string `gorm:"size:255;"`                          // 访问路由
	ParentID   *uint64 `gorm:"index;default:0;"`                   // 父级ID
	ParentPath *string `gorm:"size:512;index;default:'';"`         // 父级路径
	IsShow     int     `gorm:"index;default:0;"`                   // 是否显示(1:显示 2:隐藏)
	Status     int     `gorm:"index;default:0;"`                   // 状态(1:启用 2:禁用)
	Sequence   int     `gorm:"index;default:0;"`                   // 排序值
	Memo       *string `gorm:"size:1024;"`                         // 备注
	Creator    uint64  `gorm:""`                                   // 创建人
	// MenuAction MenuAction `gorm:"foreignKey:MyUserID"
	MenuAction []MenuAction `gorm:"many2many:menu_menu_actions;"` // 菜单动作关联
}

// 把[]*Menu定义为Menus，这样可以为它定义方法
type Menus []*Menu

// 把数据库的对应的menu转换为返回前端的menu
func (m Menu) ToTypesMenu() *types.Menu {
	menu := &types.Menu{}
	copier.Copy(menu, m)
	return menu
}

func (m Menus) ToTypesMenus() []types.Menu {
	list := make([]types.Menu, 0, len(m))
	for _, v := range m {
		list = append(list, *(v.ToTypesMenu()))
	}
	return list
}
