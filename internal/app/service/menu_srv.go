package service

import (
	"context"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
)

var MenuSrvSet = wire.NewSet(wire.Struct(new(MenuSrv), "*"))

type MenuSrv struct {
	MenuRepo *menu.MenuRepo
}

func (m *MenuSrv) Query(ctx context.Context, req types.MenuQueryReq, opt types.MenuQueryOptions) (*types.MenuQueryResp, error) {
	result, err := m.MenuRepo.Query(ctx, req, opt)
	if err != nil {
		return nil, err
	}
	return result, nil
}
