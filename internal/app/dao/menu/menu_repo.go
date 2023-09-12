package menu

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuRepo), "*"))

type MenuRepo struct {
	DB *gorm.DB
}

func (m *MenuRepo) getQueryOption(opts ...types.MenuQueryOptions) types.MenuQueryOptions {
	var opt types.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// 查询
func (m *MenuRepo) Query(ctx context.Context, req types.MenuQueryReq, opts ...types.MenuQueryOptions) (*types.MenuQueryResp, error) {
	opt := m.getQueryOption(opts...)

	db := GetMenuDB(ctx, m.DB)
	// 设置查询条件
	if len(req.IDs) > 0 {
		db = db.Where("id in (?)", req.IDs)
	}
	if req.Name != "" {
		db = db.Where("name=?", req.Name)
	}
	if req.PrefixParentPath != "" { // 修改菜单的时候需要查询到它的父路径
		db = db.Where("parent_path like ?", req.PrefixParentPath+"%") // 模糊查询
	}
	if req.IsShow != 0 {
		db = db.Where("is_show=?", req.IsShow)
	}
	if req.Status != 0 {
		db = db.Where("status=?", req.Status)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var menuList Menus
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, menuList)
	if err != nil {
		return nil, err
	}
	return &types.MenuQueryResp{
		Data:       menuList.ToTypesMenus(),
		PageResult: *result,
	}, nil
}

func (m *MenuRepo) Get(ctx context.Context, id uint64) (*types.Menu, error) {
	var menuItem Menu
	db := GetMenuDB(ctx, m.DB)
	result := db.Where("id=?", id).First(&menuItem)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return menuItem.ToTypesMenu(), nil
}

func (m *MenuRepo) Create(ctx context.Context, item types.Menu) (uint64, error) {
	entityItem := TypesMenu(item).ToMenu()
	result := GetMenuDB(ctx, m.DB).Create(&entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return uint64(entityItem.ID), nil
}

func (m *MenuRepo) Update(ctx context.Context, id uint64, item types.Menu) error {
	entityItem := TypesMenu(item).ToMenu()
	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Updates(&entityItem)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (m *MenuRepo) UpdateParentPath(ctx context.Context, id uint64, parentPath string) error {
	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Update("parent_path", parentPath)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (m *MenuRepo) Delete(ctx context.Context, id uint64) error {
	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Delete(&Menu{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (m *MenuRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetMenuDB(ctx, m.DB).Where("id=?", id).Update("status", status)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
