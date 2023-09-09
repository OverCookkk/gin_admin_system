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
	RoleRepo     *role.RoleRepo
	RoleMenuRepo *role.RoleMenuRepo
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
	// 先检查角色名是否存在
	err := r.checkRoleName(ctx, item)
	if err != nil {
		return nil, err
	}

	// TODO: 事务实现

	// 角色表  create
	roleId, err := r.RoleRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	// 角色菜单表 create
	for _, roleMenuItem := range item.RoleMenus {
		// MenuID、ActionID前端会带过来
		roleMenuItem.RoleID = roleId
		err := r.RoleMenuRepo.Create(ctx, *roleMenuItem)
		if err != nil {
			return nil, err
		}
	}

	// TODO：权限控制

	return nil, nil
}

func (r *RoleSrv) checkRoleName(ctx context.Context, item types.Role) error {
	result, err := r.RoleRepo.Query(ctx, types.RoleQueryReq{
		// PaginationParam: types.PaginationParam{},
		Name: item.Name,
	})
	if err != nil {
		return err
	} else if len(result.Data) == 0 {
		return errors.New("角色名称不能重复")
	}
	return nil
}

func (r *RoleSrv) Update(ctx context.Context, id uint64, item types.Role) error {
	oldItem, err := r.RoleRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("not found")
	} else if oldItem.Name != item.Name {
		err = r.checkRoleName(ctx, item)
		if err != nil {
			return errors.New("角色名称不能重复")
		}
	}

	// TODO: 事务实现

	// 角色表
	roleId, err := r.RoleRepo.Create(ctx, item)
	if err != nil {
		return err
	}

	// 角色菜单表，菜单可能增加或者减少，菜单操作可能增加或者减少，所以要先对比新菜单和老菜单
	for _, roleMenuItem := range item.RoleMenus {
		// MenuID、ActionID前端会带过来
		roleMenuItem.RoleID = roleId
		err := r.RoleMenuRepo.Create(ctx, *roleMenuItem)
		if err != nil {
			return err
		}
	}
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

	// todo 权限控制

	return r.RoleRepo.UpdateStatus(ctx, id, status)
}
