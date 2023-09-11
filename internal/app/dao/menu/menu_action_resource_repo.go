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
	if req.MenuID != 0 {
		subDb := GetMenuActionDB(ctx, m.DB)
		subQuery := subDb.Where("menu_id=?", req.MenuID).Select("id") // 查询action_id，返回[]int
		db = db.Where("action_id IN (?)", subQuery)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var menuAcitonResourceList MenuAcitonResources
	result, err := util.WrapPageQuery(ctx, db, req.PaginationParam, menuAcitonResourceList)
	if err != nil {
		return nil, err
	}
	return &types.MenuActionResourceQueryResp{
		Data:       menuAcitonResourceList.ToTypesMenuActionResources(), // 转为返回给前端的结构体
		PageResult: *result,
	}, nil
}
