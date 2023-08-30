package util

import (
	"context"
	"fmt"
	"gin_admin_system/internal/app/types"
	"gorm.io/gorm"
	"strings"
)

func ParseOrder(orderField []*types.OrderField) string {
	orders := make([]string, len(orderField))
	for _, v := range orderField {
		direction := "desc"
		if v.Direction == types.OrderByAsc {
			direction = "asc"
		}
		orders = append(orders, fmt.Sprintf("%s %s", v.Key, direction))
	}
	return strings.Join(orders, ",")
}

// WrapPageQuery 封装通用分页查询数据库方法，out：输出列表，返回分页结果
func WrapPageQuery(ctx context.Context, db *gorm.DB, params types.PaginationParam, out interface{}) (*types.PaginationResult, error) {
	var count int64
	err := db.Count(&count).Error

	if params.PageSize > 0 && params.Current > 0 {
		db = db.Offset(params.Current - 1).Limit(params.PageSize)
	} else if params.PageSize > 0 {
		db = db.Limit(params.PageSize)
	}
	err = db.Find(out).Error
	return &types.PaginationResult{
		Total:       count,
		CurrentPage: params.Current,
		PageSize:    params.PageSize,
	}, err
}
