package service

import (
	"context"
	"errors"
	"gin_admin_system/internal/app/dao/user"
	"gin_admin_system/internal/app/types"
	"gin_admin_system/pkg/util/hash"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"strconv"
)

var UserSrvSet = wire.NewSet(wire.Struct(new(UserSrv), "*"))

type UserSrv struct {
	Enforcer *casbin.SyncedEnforcer
	// TransRepo    *dao.TransRepo
	UserRepo     *user.UserRepo
	UserRoleRepo *user.UserRoleRepo
	// RoleRepo     *user.RoleRepo
}

// func (u *UserSrv) QueryShow(ctx context.Context, req types.UserQueryReq, opts ...types.UserQueryOptions) (*types.UserShowQueryResult, error) {
// 	result, err := u.UserRepo.Query(ctx, req, opts...)
// 	if err != nil {
// 		return nil, err
// 	} else if result == nil {
// 		return nil, nil
// 	}
// }

func (u *UserSrv) Get(ctx context.Context, id uint64, opts ...types.UserQueryOptions) (*types.User, error) {
	// 获取用户信息
	userItem, err := u.UserRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if userItem == nil {
		return nil, errors.New("user item not found")
	}

	// 获取用户所拥有的角色信息
	userRoleItem, err := u.UserRoleRepo.Query(ctx, types.UserRoleQueryReq{
		PaginationParam: types.PaginationParam{},
		UserID:          id,
	})
	if err != nil {
		return nil, err
	}
	userItem.UserRoles = userRoleItem.Data

	return userItem, nil
}

func (u *UserSrv) Create(ctx context.Context, item types.User) (*types.IDResult, error) {
	// 先检查用户名是否存在
	err := u.checkUserName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = hash.SHA1String(item.Password)
	// todo: 事务实现

	// 先在user表创建用户信息
	userID, err := u.UserRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	// 再在user_role表创建用户与角色映射信息
	for _, v := range item.UserRoles {
		v.UserID = userID
		_, err := u.UserRoleRepo.Create(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	// 权限控制：给用户逐个添加所拥有的角色
	for _, v := range item.UserRoles {
		u.Enforcer.AddRoleForUser(strconv.FormatUint(v.UserID, 10), strconv.FormatUint(v.RoleID, 10))
	}

	return &types.IDResult{ID: userID}, nil
}

func (u *UserSrv) checkUserName(ctx context.Context, item types.User) error {
	result, err := u.UserRepo.Query(ctx, types.UserQueryReq{
		// PaginationParam: types.PaginationParam{},
		UserName: item.UserName,
	})
	if err != nil {
		return err
	} else if len(result.Data) == 0 {
		return errors.New("角色名称不能重复")
	}
	return nil
}

func (u *UserSrv) Update(ctx context.Context, id uint64, item types.User) error {
	// 更新前先查询用户是否存在
	oldItem, err := u.UserRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("user item not found")
	} else if oldItem.UserName != item.UserName {
		err := u.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = hash.SHA1String(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	addUserRoles, deleteUserRoles := u.compareUserRoles(ctx, oldItem.UserRoles, item.UserRoles)

	// TODO:事务实现
	for _, v := range addUserRoles {
		v.UserID = id
		_, err := u.UserRoleRepo.Create(ctx, v)
		if err != nil {
			return err
		}
	}

	for _, v := range deleteUserRoles {
		err := u.UserRoleRepo.Delete(ctx, v.ID)
		if err != nil {
			return err
		}
	}

	// 权限控制
	for _, v := range addUserRoles {
		u.Enforcer.AddRoleForUser(strconv.FormatUint(id, 10), strconv.FormatUint(v.RoleID, 10))
	}

	for _, v := range deleteUserRoles {
		u.Enforcer.DeleteRoleForUser(strconv.FormatUint(id, 10), strconv.FormatUint(v.RoleID, 10))
	}

	return nil
}

func (u *UserSrv) compareUserRoles(ctx context.Context, oldUserRoles, newUserRoles types.UserRoles) (addList, delList types.UserRoles) {
	// 先转成map，方便查找
	oldMap := make(map[uint64]*types.UserRole)
	for _, item := range oldUserRoles {
		// 同一个UserID下的有不同的RoleID,所以用RoleID作为key
		oldMap[item.RoleID] = &item
	}

	newMap := make(map[uint64]*types.UserRole)
	for _, item := range newUserRoles {
		newMap[item.RoleID] = &item
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

func (u *UserSrv) Delete(ctx context.Context, id uint64) error {
	// 删除前先查询用户是否存在
	item, err := u.UserRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.New("user item not found")
	}

	// todo:事务实现
	err = u.UserRoleRepo.DeleteByUserID(ctx, id)
	if err != nil {
		return err
	}
	err = u.UserRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.Enforcer.DeleteUser(strconv.FormatUint(id, 10))

	return nil
}

func (u *UserSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	// 更新状态前先查询用户是否存在
	item, err := u.UserRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.New("user item not found")
	} else if item.Status == status {
		return nil
	}

	err = u.UserRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	if status == 1 {
		for _, v := range item.UserRoles {
			u.Enforcer.AddRoleForUser(strconv.FormatUint(id, 10), strconv.FormatUint(v.RoleID, 10))
		}
	} else {
		u.Enforcer.DeleteUser(strconv.FormatUint(id, 10))
	}

	return nil
}
