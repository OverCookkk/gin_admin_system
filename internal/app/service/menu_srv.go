package service

import (
	"context"
	"errors"
	"fmt"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/types"
	"github.com/google/wire"
)

var MenuSrvSet = wire.NewSet(wire.Struct(new(MenuSrv), "*"))

type MenuSrv struct {
	MenuRepo *menu.MenuRepo
}

func (m *MenuSrv) Query(ctx context.Context, req types.MenuQueryReq, opt types.MenuQueryOptions) (*types.MenuQueryResp, error) {
	result, err := m.MenuRepo.Query(ctx, req, opt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MenuSrv) Get(ctx context.Context, id uint64) (*types.Menu, error) {
	menuItem, err := m.MenuRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return menuItem, nil
}

func (m *MenuSrv) checkMenuName(ctx context.Context, item types.Menu) error {
	queryResult, err := m.MenuRepo.Query(ctx, types.MenuQueryReq{
		PaginationParam: types.PaginationParam{},
		Name:            item.Name,
		ParentID:        item.ParentID,
	})
	if err != nil {
		return err
	} else if len(queryResult.Data) == 0 {
		return errors.New("菜单名称不能重复")
	}
	return nil
}

func (m *MenuSrv) getParentPath(ctx context.Context, parentID uint64) (string, error) {
	if parentID == 0 {
		return "", nil
	}
	pItem, err := m.MenuRepo.Get(ctx, parentID)
	if err != nil {
		return "", err
	} else if pItem == nil {
		return "", errors.New("not found parent node")
	}

	return m.joinParentPath(pItem.ParentPath, pItem.ID), nil
}

func (m *MenuSrv) joinParentPath(parent string, id uint64) string {
	if parent != "" {
		parent += "/"
	}

	return fmt.Sprintf("%s%d", parent, id)
}

func (m *MenuSrv) Create(ctx context.Context, item types.Menu) (*types.IDResult, error) {
	// 检查该菜单是否存在
	if err := m.checkMenuName(ctx, item); err != nil {
		return nil, err
	}

	parentPath, err := m.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}
	item.ParentPath = parentPath

	// todo 事务实现 TransRepo.Exec；是否还需要create menu action
	id, err := m.MenuRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return &types.IDResult{ID: id}, nil
}

func (m *MenuSrv) Update(ctx context.Context, id uint64, item types.Menu) error {
	// 如果修改了menu的name，则要检查新的name是否存在
	oldItem, err := m.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("not found")
	} else if item.Name != oldItem.Name {
		if err = m.checkMenuName(ctx, item); err != nil {
			return err
		}
	}

	// 如果修改了该菜单所挂载的父菜单，则获取新的父菜单的路径赋给当前的item
	if item.ParentID != oldItem.ParentID {
		parentPath, err := m.getParentPath(ctx, item.ParentID)
		if err != nil {
			return err
		}
		item.ParentPath = parentPath
	} else {
		item.ParentPath = oldItem.ParentPath
	}

	// todo 事务实现 TransRepo.Exec；是否还需要update menu action
	// 先更新当前菜单下子菜单中的父菜单路径
	err = m.updateChildParentPath(ctx, *oldItem, item)
	if err != nil {
		return err
	}
	// 再更新当前的菜单
	return m.MenuRepo.Update(ctx, id, item)
}

func (m *MenuSrv) updateChildParentPath(ctx context.Context, oldItem, newItem types.Menu) error {
	// 先查询所有父id是old parentID的对象
	if oldItem.ParentID == newItem.ParentID {
		return nil
	}
	// 模糊查询父级路径
	oldParentPath := m.joinParentPath(oldItem.ParentPath, oldItem.ID)
	result, err := m.MenuRepo.Query(ctx, types.MenuQueryReq{
		PrefixParentPath: oldParentPath,
	})
	if err != nil {
		return err
	}

	// 逐个更新子菜单
	newParentPath := m.joinParentPath(newItem.ParentPath, newItem.ID)
	for _, menu := range result.Data {
		err = m.MenuRepo.UpdateParentPath(ctx, menu.ID, newParentPath+menu.ParentPath[len(oldParentPath):])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MenuSrv) Delete(ctx context.Context, id uint64) error {
	// 先查询该id的菜单是否存在
	item, err := m.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.New("not found")
	}

	// 再查询该id的菜单是否有子菜单，如果有，则禁止删除该菜单
	childItems, err := m.MenuRepo.Query(ctx, types.MenuQueryReq{
		ParentID: id,
	})
	if err != nil {
		return err
	} else if len(childItems.Data) != 0 {
		return errors.New("forbid delete")
	}

	// todo 事务实现 TransRepo.Exec；删除MenuActionResource和MenuAction

	// 删除该id的菜单
	return m.MenuRepo.Delete(ctx, id)
}

func (m *MenuSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	// 先查询该id的菜单是否存在
	item, err := m.MenuRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if item == nil {
		return errors.New("not found")
	} else if item.Status == status {
		return nil
	}

	return m.MenuRepo.UpdateStatus(ctx, id, status)
}
