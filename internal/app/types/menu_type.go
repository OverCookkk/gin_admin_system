package types

import "time"

type OrderDirection int

const (
	OrderByAsc OrderDirection = iota + 1
	OrderByDesc
)

// Menu 菜单对象
type Menu struct {
	ID         uint64    `json:"id,string"`                              // 唯一标识
	Name       string    `json:"name" binding:"required"`                // 菜单名称
	Sequence   int       `json:"sequence"`                               // 排序值
	Icon       string    `json:"icon"`                                   // 菜单图标
	Router     string    `json:"router"`                                 // 访问路由
	ParentID   uint64    `json:"parent_id,string"`                       // 父级ID
	ParentPath string    `json:"parent_path"`                            // 父级路径
	IsShow     int       `json:"is_show" binding:"required,max=2,min=1"` // 是否显示(1:显示 2:隐藏)
	Status     int       `json:"status" binding:"required,max=2,min=1"`  // 状态(1:启用 2:禁用)
	Memo       string    `json:"memo"`                                   // 备注
	Creator    uint64    `json:"creator"`                                // 创建者
	CreatedAt  time.Time `json:"created_at"`                             // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                             // 更新时间
	// Actions    MenuActions `json:"actions"`                                // 动作列表
}

// MenuQueryReq 菜单查询请求
type MenuQueryReq struct {
	PaginationParam
	IDs              []uint64 `form:"-"`        // 唯一标识列表
	Name             string   `form:"-"`        // 菜单名称
	ParentID         uint64   `form:"parentID"` // 父级id
	PrefixParentPath string   `form:"-"`        // 父级路径(前缀模糊查询)
	IsShow           int      `form:"isShow"`   // 是否显示(1:显示 2:隐藏)
	Status           int      `form:"status"`   // 状态(1:启用 2:禁用)
}

// MenuQueryResp 菜单查询响应
type MenuQueryResp struct {
	Data       []Menu           `json:"menu_list"`
	PageResult PaginationResult `json:"page_result"`
}

// PaginationResult 分页结果
type PaginationResult struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}

// MenuQueryOptions 菜单查询可选参数，如排序方向等
type MenuQueryOptions struct {
	OrderFields  []*OrderField
	SelectFields []string
}

type OrderField struct {
	Key       string
	Direction OrderDirection // 排序的方向
}

func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// 分页参数
type PaginationParam struct {
	Pagination bool `form:"-"`
	OnlyCount  bool `form:"-"`
	Current    int  `form:"current,default=1"`
	PageSize   int  `form:"pageSize,default=10" binding:"max=100"`
}

type IDResult struct {
	ID uint64 `json:"id"`
}
