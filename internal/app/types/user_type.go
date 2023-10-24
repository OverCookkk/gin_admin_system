package types

import (
	"context"
	"gin_admin_system/internal/app/config"
	"gin_admin_system/pkg/util/hash"
	"time"
)

// GetRootUser 获取root用户
func GetRootUser() *User {
	user := config.C.Root
	return &User{
		ID:       user.UserID,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: hash.MD5String(user.Password),
	}
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, userID uint64) bool {
	return GetRootUser().ID == userID
}

// User 用户对象
type User struct {
	ID        uint64    `json:"id,string"`                             // 唯一标识
	UserName  string    `json:"user_name" binding:"required"`          // 用户名
	RealName  string    `json:"real_name" binding:"required"`          // 真实姓名
	Password  string    `json:"password"`                              // 密码
	Phone     string    `json:"phone"`                                 // 手机号
	Email     string    `json:"email"`                                 // 邮箱
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator   uint64    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"created_at"`                            // 创建时间
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // 角色授权
}

func (u *User) CleanSecure() *User {
	u.Password = ""
	return u
}

// UserQueryReq 查询请求
type UserQueryReq struct {
	PaginationParam
	UserName   string   `form:"userName"`   // 用户名
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 用户状态(1:启用 2:停用)
	RoleIDs    []uint64 `form:"-"`          // 角色ID列表，查询有用这些role_id的用户
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields  []*OrderField
	SelectFields []string
}

// UserQueryResp 查询结果
type UserQueryResp struct {
	Data       []User           `json:"user_list"`
	PageResult PaginationResult `json:"page_result"`
}

// ----------------------------------------UserRole--------------------------------------

// UserRole 用户角色
type UserRole struct {
	ID     uint64 `json:"id,string"`      // 唯一标识
	UserID uint64 `json:"user_id,string"` // 用户ID
	RoleID uint64 `json:"role_id,string"` // 角色ID
}

// UserRoleQueryReq 查询请求
type UserRoleQueryReq struct {
	PaginationParam
	UserID uint64 // 用户ID
	// UserIDs []uint64 // 用户ID列表
}

// UserRoleQueryOptions 查询可选参数项
type UserRoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserRoleQueryResp 查询结果
type UserRoleQueryResp struct {
	Data       UserRoles        `json:"user_role_list"`
	PageResult PaginationResult `json:"page_result"`
}

// UserRoles 角色菜单列表
type UserRoles []UserRole

// // ToMap 转换为map
// func (a UserRoles) ToMap() map[uint64]*UserRole {
// 	m := make(map[uint64]*UserRole)
// 	for _, item := range a {
// 		m[item.RoleID] = &item
// 	}
// 	return m
// }

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []uint64 {
	list := make([]uint64, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// // ToUserIDMap 转换为用户ID映射
// func (a UserRoles) ToUserIDMap() map[uint64]UserRoles {
// 	m := make(map[uint64]UserRoles)
// 	for _, item := range a {
// 		m[item.UserID] = append(m[item.UserID], item)
// 	}
// 	return m
// }
