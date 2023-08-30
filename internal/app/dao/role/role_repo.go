package role

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleRepo), "*"))

type RoleRepo struct {
	DB *gorm.DB
}

func (r *RoleRepo) getQueryOption(opts ...types.RoleQueryOptions) types.RoleQueryOptions {
	var opt types.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}
func (r *RoleRepo) Query(ctx context.Context, req types.RoleQueryReq, opts ...types.RoleQueryOptions) (*types.RoleQueryResp, error) {
	opt := r.getQueryOption(opts...)

	db := GetRoleDB(ctx, r.DB)
	if len(req.IDs) > 0 {
		db = db.Where("id in (?)", req.IDs)
	}
	if req.Name != "" {
		db = db.Where("name=?", req.Name)
	}
	if req.Status != 0 {
		db = db.Where("status=?", req.Status)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}
	var roleList *Roles
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, roleList)
	if err != nil {
		return nil, err
	}
	return &types.RoleQueryResp{
		Data:       roleList.ToTypesRoles(),
		PageResult: *result,
	}, nil
}
