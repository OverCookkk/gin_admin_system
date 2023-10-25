package util

import (
	"context"
	"gin_admin_system/internal/app/contextx"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var TransSet = wire.NewSet(wire.Struct(new(Transaction), "*"))

type Transaction struct {
	DB *gorm.DB
}

// Exec db开始事务（t.DB.Transaction）的封装
func (t *Transaction) Exec(ctx context.Context, fn func(ctx2 context.Context) error) error {
	// 如果上下文中有事务对象，则直接使用
	if _, ok := contextx.GetTrans(ctx); ok {
		return fn(ctx)
	}

	// 上下文中没有事务对象，t.DB.Transaction创建事务对象后，使用SetTrans存进上下文中
	return t.DB.Transaction(func(tx *gorm.DB) error {
		return fn(contextx.SetTrans(ctx, tx))
	})
}
