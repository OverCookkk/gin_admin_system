package service

import (
	"context"
	"fmt"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/types"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"strconv"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleSrv), "*"))

type RoleSrv struct {
	Enforcer *casbin.SyncedEnforcer
	// TransRepo              *dao.TransRepo
	RoleRepo     *role.RoleRepo
	RoleMenuRepo *role.RoleMenuRepo
	// UserRepo               *role.UserRepo
	MenuActionResourceRepo *menu.MenuActionResourceRepo
}

func (r *RoleSrv) Query(ctx context.Context, req types.RoleQueryReq, opt types.RoleQueryOptions) (*types.RoleQueryResp, error) {
	return r.RoleRepo.Query(ctx, req, opt)
}

func (r *RoleSrv) Get(ctx context.Context, id uint64) (*types.Role, error) {
	roleItem, err := r.RoleRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if roleItem == nil {
		return nil, errors.New("role item not found")
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

	// 先创建角色表
	roleId, err := r.RoleRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	// 再创建角色菜单表 create
	for _, roleMenuItem := range item.RoleMenus {
		// MenuID、ActionID前端会带过来
		roleMenuItem.RoleID = roleId
		err := r.RoleMenuRepo.Create(ctx, roleMenuItem)
		if err != nil {
			return nil, err
		}
	}

	// 权限控制
	// 通过menu_id查询到menu_action_resource信息
	resources, err := r.MenuActionResourceRepo.Query(ctx, types.MenuActionResourceQueryReq{
		PaginationParam: types.PaginationParam{}, // 默认参数
		MenuIDs:         item.RoleMenus.ToMenuIDs(),
	})
	// 去重。不同的菜单动作下可能有相同的资源
	resourceMap := make(map[string]*types.MenuActionResource)
	for _, v := range resources.Data {
		resourceMap[v.Method+v.Path] = &v
	}
	for _, v := range resourceMap {
		// role_id, path, method
		_, err := r.Enforcer.AddPermissionForUser(strconv.FormatUint(item.ID, 10), v.Path, v.Method)
		if err != nil {
			return nil, err
		}
	}

	return &types.IDResult{ID: roleId}, nil
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

	// 角色表 Create
	roleId, err := r.RoleRepo.Create(ctx, item)
	if err != nil {
		return err
	}

	// 角色菜单表 Create；菜单可能增加或者减少，菜单操作可能增加或者减少，所以更新前要先对比新菜单和老菜单
	addRoleMenus, deleteRoleMenus := r.compareRoleMenus(ctx, oldItem.RoleMenus, item.RoleMenus)
	for _, roleMenuItem := range addRoleMenus {
		// MenuID、ActionID前端会带过来
		roleMenuItem.RoleID = roleId
		err := r.RoleMenuRepo.Create(ctx, roleMenuItem)
		if err != nil {
			return err
		}
	}

	for _, roleMenuItem := range deleteRoleMenus {
		err := r.Delete(ctx, roleMenuItem.RoleID)
		if err != nil {
			return err
		}
	}

	// 先删除这角色全部的权限，重新给该角色逐一增加权限
	r.Enforcer.DeletePermissionsForUser(strconv.FormatUint(item.ID, 10))

	roleMenus, err := r.RoleMenuRepo.Query(ctx, types.RoleMenuQueryReq{
		RoleID: id,
	})
	if err != nil {
		return err
	}
	// 通过menu_id查询到menu_action_resource信息
	resources, err := r.MenuActionResourceRepo.Query(ctx, types.MenuActionResourceQueryReq{
		PaginationParam: types.PaginationParam{}, // 默认参数
		MenuIDs:         types.RoleMenus(roleMenus.Data).ToMenuIDs(),
	})
	// 去重。不同的菜单动作下可能有相同的资源
	resourceMap := make(map[string]*types.MenuActionResource)
	for _, v := range resources.Data {
		resourceMap[v.Method+v.Path] = &v
	}
	for _, v := range resourceMap {
		// role_id, path, method
		_, err := r.Enforcer.AddPermissionForUser(strconv.FormatUint(item.ID, 10), v.Path, v.Method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleSrv) compareRoleMenus(ctx context.Context, oldRoleMenus, newRoleMenus types.RoleMenus) (addList, delList types.RoleMenus) {
	// 先转成map，方便查找
	oldMap := make(map[string]*types.RoleMenu)
	for _, item := range oldRoleMenus {
		// 同一个RoleID下有不同的MenuID和ActionID的组合，所以key如下
		oldMap[fmt.Sprintf("%s-%s", item.MenuID, item.ActionID)] = &item
	}

	newMap := make(map[string]*types.RoleMenu)
	for _, item := range newRoleMenus {
		newMap[fmt.Sprintf("%s-%s", item.MenuID, item.ActionID)] = &item
	}

	for k, v := range newMap {
		if _, ok := oldMap[k]; !ok { // 新的item没找到，说明是新增的
			addList = append(addList, *v)
		} else { // 找到了，就删除oldMap里面的值，这样剩下的就是删除的
			delete(oldMap, k)
		}
	}
	for _, v := range oldMap {
		delList = append(delList, *v)
	}

	return
}

func (r *RoleSrv) Delete(ctx context.Context, id uint64) error {
	// 先查询该角色是否存在
	item, err := r.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		errors.New("角色不存在")
	}

	// todo:查询是否有用户属于这角色，如果有，则不用删除该角色
	// queryUserResult, err := r.UserRepo.Query(ctx, types.UserQueryReq{
	// 	// PaginationParam:  types.PaginationParam{},
	// 	RoleIDs: nil,
	// })
	// if err != nil {
	// 	return err
	// } else if len(queryUserResult.Data) != 0 {
	// 	return errors.New("不允许删除已经存在用户的角色")
	// }

	// 先删除角色菜单
	err = r.RoleMenuRepo.DeleteByRoleID(ctx, id)
	if err != nil {
		return err
	}

	// 再删除角色
	err = r.RoleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// 最后删除该角色整个的权限（DeleteRole）
	r.Enforcer.DeleteRole(strconv.FormatUint(id, 10))
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

	err = r.RoleRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// 禁用就删除该角色整个的权限（DeleteRole），启用就为该角色增加权限
	if status == 1 {
		roleMenus, err := r.RoleMenuRepo.Query(ctx, types.RoleMenuQueryReq{
			RoleID: id,
		})
		if err != nil {
			return err
		}
		// 通过menu_id查询到menu_action_resource信息
		resources, err := r.MenuActionResourceRepo.Query(ctx, types.MenuActionResourceQueryReq{
			PaginationParam: types.PaginationParam{}, // 默认参数
			MenuIDs:         types.RoleMenus(roleMenus.Data).ToMenuIDs(),
		})
		// 去重。不同的菜单动作下可能有相同的资源
		resourceMap := make(map[string]*types.MenuActionResource)
		for _, v := range resources.Data {
			resourceMap[v.Method+v.Path] = &v
		}
		for _, v := range resourceMap {
			// role_id, path, method
			_, err := r.Enforcer.AddPermissionForUser(strconv.FormatUint(item.ID, 10), v.Path, v.Method)
			if err != nil {
				return err
			}
		}
	} else {
		r.Enforcer.DeleteRole(strconv.FormatUint(id, 10))
	}

	return nil
}
