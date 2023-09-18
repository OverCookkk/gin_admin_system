package user

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return defDB.Model(new(User))
}

type User struct {
	gorm.Model
	UserName string  `gorm:"size:64;uniqueIndex;default:'';not null;"` // 用户名
	RealName string  `gorm:"size:64;index;default:'';"`                // 真实姓名
	Password string  `gorm:"size:40;default:'';"`                      // 密码
	Email    *string `gorm:"size:255;"`                                // 邮箱
	Phone    *string `gorm:"size:20;"`                                 // 手机号
	Status   int     `gorm:"index;default:0;"`                         // 状态(1:启用 2:停用)
	Creator  uint64  `gorm:""`                                         // 创建者
}

func (u User) ToTypesUser() *types.User {
	user := &types.User{}
	copier.Copy(user, u)
	return user
}

type Users []*User

func (u Users) ToTypesUsers() []types.User {
	list := make([]types.User, 0, len(u))
	for _, v := range u {
		list = append(list, *(v.ToTypesUser()))
	}
	return list
}

type TypesUser types.User

// 把返回前端的user转换为数据库的对应的user
func (a TypesUser) ToUser() *User {
	user := &User{}
	copier.Copy(user, a)
	return user
}
