package menu

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return defDB.Model(new(MenuActionResource))
}

type MenuActionResource struct {
	gorm.Model
	ActionID uint64 `gorm:"index;not null;"` // 菜单动作ID
	Method   string `gorm:"size:50;"`        // 资源请求方式(支持正则)
	Path     string `gorm:"size:255;"`       // 资源请求路径（支持/:id匹配）
}

// 把数据库的对应的MenuActions转换为返回前端的MenuActions
func (m MenuActionResource) ToTypesMenuActionResource() *types.MenuActionResource {
	menuActionResource := &types.MenuActionResource{}
	copier.Copy(menuActionResource, m) // m赋值给menuActionResource
	return menuActionResource
}

type MenuAcitonResources []*MenuActionResource

func (m MenuAcitonResources) ToTypesMenuActionResources() []types.MenuActionResource {
	list := make([]types.MenuActionResource, 0, len(m))
	for _, v := range m {
		list = append(list, *(v.ToTypesMenuActionResource()))
	}
	return list
}
