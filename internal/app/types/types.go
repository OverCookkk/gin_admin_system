package types

type OrderDirection int

const (
	OrderByAsc OrderDirection = iota + 1
	OrderByDesc
)

// PaginationResult 分页结果
type PaginationResult struct {
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}

// OrderField 数据库排序参数结构
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

// PaginationParam 分页参数
type PaginationParam struct {
	Pagination bool `form:"-"`
	OnlyCount  bool `form:"-"`
	Current    int  `form:"current,default=1"`
	PageSize   int  `form:"pageSize,default=10" binding:"max=100"`
}

type IDResult struct {
	ID uint64 `json:"id"`
}
