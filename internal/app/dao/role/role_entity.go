package role

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return defDB.Model(new(Role))
}

type Role struct {
	gorm.Model
	Name     string  `gorm:"size:100;index;default:'';not null;"` // 角色名称
	Sequence int     `gorm:"index;default:0;"`                    // 排序值
	Memo     *string `gorm:"size:1024;"`                          // 备注
	Status   int     `gorm:"index;default:0;"`                    // 状态(1:启用 2:禁用)
	Creator  uint64  `gorm:""`                                    // 创建者
}

// 把[]*Role定义为Roles，这样可以为它定义方法
type Roles []*Role

func (r Role) ToTypesRole() *types.Role {
	role := &types.Role{}
	copier.Copy(role, r)
	return role
}

func (m Roles) ToTypesRoles() []types.Role {
	list := make([]types.Role, 0, len(m))
	for _, v := range m {
		list = append(list, *(v.ToTypesRole()))
	}
	return list
}

type TypesRole types.Role

// 把返回前端的role转换为数据库的对应的role
func (r TypesRole) ToRole() *Role {
	role := &Role{}
	copier.Copy(role, r)
	return role
}
