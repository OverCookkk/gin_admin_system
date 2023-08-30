package service

import (
	"context"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleSrv), "*"))

type RoleSrv struct {
	// Enforcer               *casbin.SyncedEnforcer
	// TransRepo              *dao.TransRepo
	RoleRepo *role.RoleRepo
	// RoleMenuRepo           *dao.RoleMenuRepo
	// UserRepo               *dao.UserRepo
	// MenuActionResourceRepo *dao.MenuActionResourceRepo
}

func (r RoleSrv) Query(ctx context.Context, req types.RoleQueryReq, opt types.RoleQueryOptions) (*types.RoleQueryResp, error) {
	return r.RoleRepo.Query(ctx, req, opt)
}
