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
	var roleList Roles
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, &roleList)
	if err != nil {
		return nil, err
	}
	return &types.RoleQueryResp{
		Data:       roleList.ToTypesRoles(),
		PageResult: *result,
	}, nil
}

func (r *RoleRepo) Get(ctx context.Context, id uint64) (*types.Role, error) {
	var roleItem Role
	db := GetRoleDB(ctx, r.DB)
	result := db.Where("id=?", id).First(&roleItem)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return roleItem.ToTypesRole(), nil
}

func (r *RoleRepo) Create(ctx context.Context, item types.Role) (uint64, error) {
	entityItem := TypesRole(item).ToRole()
	result := GetRoleDB(ctx, r.DB).Create(entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return uint64(entityItem.ID), nil
}

func (r *RoleRepo) Update(ctx context.Context, id uint64, item types.Role) error {
	entityItem := TypesRole(item).ToRole()
	result := GetRoleDB(ctx, r.DB).Where("id=?", id).Updates(entityItem)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleRepo) UpdateParentPath(ctx context.Context, id uint64, parentPath string) error {
	result := GetRoleDB(ctx, r.DB).Where("id=?", id).Update("parent_path", parentPath)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, id uint64) error {
	result := GetRoleDB(ctx, r.DB).Where("id=?", id).Delete(&Role{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetRoleDB(ctx, r.DB).Where("id=?", id).Update("status", status)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
