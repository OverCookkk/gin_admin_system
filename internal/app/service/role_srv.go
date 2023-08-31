package service

import (
	"context"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"github.com/pkg/errors"
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

func (r *RoleSrv) Query(ctx context.Context, req types.RoleQueryReq, opt types.RoleQueryOptions) (*types.RoleQueryResp, error) {
	return r.RoleRepo.Query(ctx, req, opt)
}

func (r *RoleSrv) Get(ctx context.Context, id uint64) (*types.Role, error) {
	roleItem, err := r.RoleRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return roleItem, nil
}

func (r *RoleSrv) Create(ctx context.Context, item types.Role) (*types.IDResult, error) {
	return nil, nil
}

func (r *RoleSrv) Update(ctx context.Context, id uint64, item types.Role) error {
	return nil
}

func (r *RoleSrv) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (r *RoleSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	// 先查询该id的角色是否存在
	item, err := r.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.New("not found")
	} else if item.Status == status {
		return nil
	}

	return r.RoleRepo.UpdateStatus(ctx, id, status)
}
