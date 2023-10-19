package menu

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var MenuActionSet = wire.NewSet(wire.Struct(new(MenuActionRepo), "*"))

type MenuActionRepo struct {
	DB *gorm.DB
}

func (m *MenuActionRepo) getQueryOption(opts ...types.MenuActionQueryOptions) types.MenuActionQueryOptions {
	var opt types.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// 查询
func (m *MenuActionRepo) Query(ctx context.Context, req types.MenuActionQueryReq, opts ...types.MenuActionQueryOptions) (*types.MenuActionQueryResp, error) {
	opt := m.getQueryOption(opts...)

	db := GetMenuActionDB(ctx, m.DB)
	// 设置查询条件
	if req.MenuID != 0 {
		db = db.Where("menu_id=?", req.MenuID)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var menuAcitonList MenuAcitons
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, &menuAcitonList)
	if err != nil {
		return nil, err
	}
	return &types.MenuActionQueryResp{
		Data:       menuAcitonList.ToTypesMenuActions(), // 转为返回给前端的结构体
		PageResult: *result,
	}, nil
}

// func (m *MenuActionRepo) Get(ctx context.Context, id uint64) (*types.Menu, error) {
// 	var menuItem Menu
// 	db := GetMenuDB(ctx, m.DB)
// 	result := db.Where("id=?", id).First(&menuItem)
// 	if err := result.Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return menuItem.ToTypesMenu(), nil
// }

func (m *MenuActionRepo) Create(ctx context.Context, item types.MenuAction) (uint64, error) {
	entityItem := TypesMenuAction(item).ToMenuAction()
	result := GetMenuActionDB(ctx, m.DB).Create(entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return uint64(entityItem.ID), nil
}

//
// func (m *MenuActionRepo) Update(ctx context.Context, id uint64, item types.Menu) error {
// 	entityItem := TypesMenu(item).ToMenu()
// 	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Updates(&entityItem)
// 	if err := result.Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (m *MenuActionRepo) UpdateParentPath(ctx context.Context, id uint64, parentPath string) error {
// 	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Update("parent_path", parentPath)
// 	if err := result.Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (m *MenuActionRepo) DeleteByMenuID(ctx context.Context, menuID uint64) error {
	result := GetMenuActionDB(ctx, m.DB).Where("menu_id=?", menuID).Delete(MenuAction{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// func (m *MenuActionRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
// 	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Update("status", status)
// 	if err := result.Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
