package user

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(UserRole)) // defDB.Model(new(UserRole))
}

type UserRole struct {
	gorm.Model
	UserID uint64 `gorm:"index;default:0;"` // 用户ID
	RoleID uint64 `gorm:"index;default:0;"` // 角色ID
}

func (u UserRole) ToTypesUserRole() *types.UserRole {
	userRole := &types.UserRole{}
	copier.Copy(userRole, u)
	return userRole
}

type UserRoles []*UserRole

func (u UserRoles) ToTypesUserRoles() []types.UserRole {
	list := make([]types.UserRole, 0, len(u))
	for _, v := range u {
		list = append(list, *(v.ToTypesUserRole()))
	}
	return list
}

type TypesUserRole types.UserRole

// 把返回前端的userRole转换为数据库的对应的userRole
func (a TypesUserRole) ToUserRole() *UserRole {
	userRole := &UserRole{}
	copier.Copy(userRole, a)
	return userRole
}
