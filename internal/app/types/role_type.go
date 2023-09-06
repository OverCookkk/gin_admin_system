package types

// Role 角色对象
type Role struct {
	// ID        uint64    `json:"id,string"`                             // 唯一标识
	Name      string    `json:"name" binding:"required"`               // 角色名称
	Sequence  int       `json:"sequence"`                              // 排序值
	Memo      string    `json:"memo"`                                  // 备注
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator   uint64    `json:"creator"`                               // 创建者
	RoleMenus RoleMenus `json:"role_menus" binding:"required,gt=0"`    // 角色菜单列表
}

type RoleQueryReq struct {
	PaginationParam
	IDs    []uint64 `form:"-"`      // 唯一标识列表
	Name   string   `form:"-"`      // 角色名称
	Status int      `form:"status"` // 状态(1:启用 2:禁用)
}

// RoleQueryResp 菜单查询响应
type RoleQueryResp struct {
	Data       []Role           `json:"role_list"`
	PageResult PaginationResult `json:"page_result"`
}

// RoleQueryOptions 查询可选参数项
type RoleQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// --------------------角色下管理的菜单------------------------
type RoleMenu struct {
	// ID       uint64 `json:"id,string"`                           // 唯一标识
	RoleID   uint64 `json:"role_id,string" binding:"required"`   // 角色ID，角色有哪些菜单权限，所以角色要与RoleMenus绑定
	MenuID   uint64 `json:"menu_id,string" binding:"required"`   // 菜单ID
	ActionID uint64 `json:"action_id,string" binding:"required"` // 动作ID
}

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu
