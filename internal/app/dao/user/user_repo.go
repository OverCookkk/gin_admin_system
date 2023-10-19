package user

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var UserSet = wire.NewSet(wire.Struct(new(UserRepo), "*"))

type UserRepo struct {
	DB *gorm.DB
}

func (u *UserRepo) getQueryOption(opts ...types.UserQueryOptions) types.UserQueryOptions {
	var opt types.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (u *UserRepo) Query(ctx context.Context, req types.UserQueryReq, opts ...types.UserQueryOptions) (*types.UserQueryResp, error) {
	opt := u.getQueryOption(opts...)

	db := GetUserDB(ctx, u.DB)
	if req.UserName != "" {
		db = db.Where("user_name = ?", req.UserName)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	if len(req.RoleIDs) > 0 { // 通过使用角色去查询  用户
		// 先去user_role表通过role_id查询到它对应的user_id列表
		subQuery := GetUserRoleDB(ctx, db).Select("user_id").Where("role_id in (?)", req.RoleIDs)
		// 再在user表中用user_id列表查询用户的具体信息
		db = db.Where("id in ?", subQuery)
	}
	if req.QueryValue != "" { // 模糊查询用到，模糊查询  用户的名字  和  真实的名字
		v := "%" + req.QueryValue + "%"
		db = db.Where("user_name like ? or real_name like ?", v, v)
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}
	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var userList Users
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, &userList)
	if err != nil {
		return nil, err
	}
	return &types.UserQueryResp{
		Data:       userList.ToTypesUsers(),
		PageResult: *result,
	}, nil
}

func (u *UserRepo) Get(ctx context.Context, id uint64) (*types.User, error) {
	var userItem User
	db := GetUserDB(ctx, u.DB)
	result := db.Where("id=?", id).First(&userItem)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return userItem.ToTypesUser(), nil
}

func (u *UserRepo) Create(ctx context.Context, item types.User) (uint64, error) {
	entityItem := TypesUser(item).ToUser()
	result := GetUserDB(ctx, u.DB).Create(entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return uint64(entityItem.ID), nil
}

func (u *UserRepo) Update(ctx context.Context, id uint64, item types.User) error {
	entityItem := TypesUser(item).ToUser()
	result := GetUserDB(ctx, u.DB).Where("id=?", id).Updates(entityItem)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) Delete(ctx context.Context, id uint64) error {
	result := GetUserDB(ctx, u.DB).Where("id=?", id).Delete(&User{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetUserDB(ctx, u.DB).Where("id=?", id).Update("status", status)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UpdatePassword(ctx context.Context, id uint64, password string) error {
	result := GetUserDB(ctx, u.DB).Where("id=?", id).Update("password", password)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
