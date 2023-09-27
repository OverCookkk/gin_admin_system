package menu

import (
	"context"
	"gin_admin_system/internal/app/dao/util"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var MenuActionResourceSet = wire.NewSet(wire.Struct(new(MenuActionResourceRepo), "*"))

type MenuActionResourceRepo struct {
	DB *gorm.DB
}

func (m *MenuActionResourceRepo) getQueryOption(opts ...types.MenuActionResourceQueryOptions) types.MenuActionResourceQueryOptions {
	var opt types.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// 查询
func (m *MenuActionResourceRepo) Query(ctx context.Context, req types.MenuActionResourceQueryReq, opts ...types.MenuActionResourceQueryOptions) (*types.MenuActionResourceQueryResp, error) {
	opt := m.getQueryOption(opts...)

	db := GetMenuActionResourceDB(ctx, m.DB)

	// 设置查询条件
	// 先用menu_id在menuAction表中查，得到action_id再去menuActionResource表中去查
	if req.MenuID != 0 { // 查询一个menu_id
		subDb := GetMenuActionDB(ctx, m.DB)
		subQuery := subDb.Where("menu_id=?", req.MenuID).Select("id") // 查询action_id，返回[]int
		db = db.Where("action_id IN (?)", subQuery)
	}

	if len(req.MenuIDs) != 0 { // 查询多个menu_id
		subDb := GetMenuActionDB(ctx, m.DB)
		subQuery := subDb.Where("menu_id IN (?)", req.MenuIDs).Select("id") // 查询action_id，返回[]int
		db = db.Where("action_id IN (?)", subQuery)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var menuAcitonResourceList MenuAcitonResources
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, &menuAcitonResourceList)
	if err != nil {
		return nil, err
	}
	return &types.MenuActionResourceQueryResp{
		Data:       menuAcitonResourceList.ToTypesMenuActionResources(), // 转为返回给前端的结构体
		PageResult: *result,
	}, nil
}

func (m *MenuActionResourceRepo) Create(ctx context.Context, item types.MenuActionResource) (uint64, error) {
	entityItem := TypesMenuActionResource(item).ToMenuActionResource()
	result := GetMenuActionResourceDB(ctx, m.DB).Create(&entityItem)
	if err := result.Error; err != nil {
		return 0, err
	}
	return uint64(entityItem.ID), nil
}

func (m *MenuActionResourceRepo) DeleteByMenuID(ctx context.Context, menuID uint64) error {
	// 先查出该menuID下的所有actionID，再用actionID删除MenuActionResource
	subQuery := GetMenuActionDB(ctx, m.DB).Where("menu_id=?", menuID).Select("id")
	result := GetMenuActionResourceDB(ctx, m.DB).Where("action_id IN (?)", subQuery).Delete(MenuActionResource{})

	return result.Error
}
