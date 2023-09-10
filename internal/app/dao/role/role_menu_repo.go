package role

import (
	"context"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RoleMenuSet = wire.NewSet(wire.Struct(new(RoleMenuRepo), "*"))

type RoleMenuRepo struct {
	DB *gorm.DB
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
