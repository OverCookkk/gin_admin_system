package user

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var UserRoleSet = wire.NewSet(wire.Struct(new(UserRoleRepo), "*"))

type UserRoleRepo struct {
	DB *gorm.DB
}

func (u *UserRoleRepo) getQueryOption(opts ...types.UserRoleQueryOptions) types.UserRoleQueryOptions {
	var opt types.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (u *UserRoleRepo) Query(ctx context.Context, req types.UserRoleQueryReq, opts ...types.UserRoleQueryOptions) (*types.UserRoleQueryResp, error) {
	opt := u.getQueryOption(opts...)

	db := GetUserDB(ctx, u.DB)
	if req.UserID > 0 {
		db = db.Where("user_id=?", req.UserID)
	}
	// if len(req.UserIDs) > 0 {
	// 	db = db.Where("status in (?)", req.UserIDs)
	// }

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var userRoleList *UserRoles
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, userRoleList)
	if err != nil {
		return nil, err
	}
	return &types.UserRoleQueryResp{
		Data:       userRoleList.ToTypesUserRoles(),
		PageResult: *result,
	}, nil
}

func (u *UserRoleRepo) Get(ctx context.Context, id uint64) (*types.UserRole, error) {
	var userRoleItem UserRole
	db := GetUserRoleDB(ctx, u.DB)
	result := db.Where("id=?", id).First(&userRoleItem)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return userRoleItem.ToTypesUserRole(), nil
}

func (u *UserRoleRepo) Create(ctx context.Context, item types.UserRole) (uint64, error) {
	entityItem := TypesUserRole(item)
	result := GetUserRoleDB(ctx, u.DB).Create(entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return entityItem.ID, nil
}

func (u *UserRoleRepo) Update(ctx context.Context, id uint64, item types.UserRole) error {
	entityItem := TypesUserRole(item).ToUserRole()
	result := GetUserRoleDB(ctx, u.DB).Where("id=?", id).Updates(&entityItem)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRoleRepo) Delete(ctx context.Context, id uint64) error {
	result := GetUserRoleDB(ctx, u.DB).Where("id=?", id).Delete(&UserRole{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRoleRepo) DeleteByUserID(ctx context.Context, userID uint64) error {
	result := GetUserRoleDB(ctx, u.DB).Where("user_id=?", userID).Delete(&UserRole{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
