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

// 查询
func (m *MenuRepo) Query(ctx context.Context, req types.MenuQueryReq, opt types.MenuQueryOptions) (*types.MenuQueryResp, error) {
    db := GetMenuDB(ctx, m.DB)
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
