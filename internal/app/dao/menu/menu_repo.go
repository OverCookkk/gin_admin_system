package menu

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuRepo), "*"))

type MenuRepo struct {
	DB *gorm.DB
}

func (m *MenuRepo) Query(ctx context.Context, req types.MenuQueryReq, opt types.MenuQueryOptions) (*types.MenuQueryResp, error) {

	return nil, nil
}
