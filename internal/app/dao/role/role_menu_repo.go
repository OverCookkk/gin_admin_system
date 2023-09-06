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
