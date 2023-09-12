package role

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RoleMenuSet = wire.NewSet(wire.Struct(new(RoleMenuRepo), "*"))

type RoleMenuRepo struct {
	DB *gorm.DB
}

func (r *RoleMenuRepo) getQueryOption(opts ...types.RoleMenuQueryOptions) types.RoleMenuQueryOptions {
	var opt types.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (r *RoleMenuRepo) Create(ctx context.Context, item types.RoleMenu) error {
	entityItem := TypesRoleMenu(item).ToRoleMenu()
	result := GetRoleMenuDB(ctx, r.DB).Create(entityItem)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleMenuRepo) DeleteByRoleID(ctx context.Context, id uint64) error {
	result := GetRoleDB(ctx, r.DB).Where("role_id=?", id).Delete(&RoleMenu{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleMenuRepo) Query(ctx context.Context, req types.RoleMenuQueryReq, opts ...types.RoleMenuQueryOptions) (*types.RoleMenuQueryResp, error) {
	opt := r.getQueryOption(opts...)

	db := GetRoleMenuDB(ctx, r.DB)
	if v := req.RoleID; v > 0 {
		db = db.Where("role_id=?", v)
	}
	// if v := req.RoleIDs; len(v) > 0 {	// TODO:role_menu   roleIDs
	// 	db = db.Where("role_id IN (?)", v)
	// }

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var roleMenuList RoleMenus
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, &roleMenuList)
	if err != nil {
		return nil, err
	}
	qr := &types.RoleMenuQueryResp{
		Data:       roleMenuList.ToTypesRoleMenus(),
		PageResult: *result,
	}

	return qr, nil
}
